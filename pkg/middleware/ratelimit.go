package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// RateLimiterConfig defines the rate limiting policy for an endpoint
type RateLimiterConfig struct {
	Action      string
	Limit       int
	Window      time.Duration
	RedisClient *redis.Client
	KeyFunc     func(c echo.Context) string
}

// MemoryRateLimiter fallback storage
type memoryRateEntry struct {
	count     int
	expiresAt time.Time
}

var (
	memoryStore sync.Map
	cleanOnce   sync.Once
)

func initMemoryCleanup() {
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for range ticker.C {
			now := time.Now()
			memoryStore.Range(func(key, value any) bool {
				if entry, ok := value.(memoryRateEntry); ok {
					if now.After(entry.expiresAt) {
						memoryStore.Delete(key)
					}
				}
				return true
			})
		}
	}()
}

// SafeExtractIP extracts the client IP without blindly trusting untrusted headers
func SafeExtractIP(c echo.Context) string {
	// Check remote address first
	remoteIP, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		remoteIP = c.Request().RemoteAddr
	}

	// If remote IP is a loopback or local proxy, check X-Forwarded-For carefully
	ip := net.ParseIP(remoteIP)
	if ip != nil && (ip.IsLoopback() || ip.IsPrivate()) {
		xff := c.Request().Header.Get("X-Forwarded-For")
		if xff != "" {
			parts := strings.Split(xff, ",")
			clientIP := strings.TrimSpace(parts[0])
			if parsed := net.ParseIP(clientIP); parsed != nil {
				return clientIP
			}
		}
		xri := c.Request().Header.Get("X-Real-IP")
		if xri != "" {
			if parsed := net.ParseIP(strings.TrimSpace(xri)); parsed != nil {
				return strings.TrimSpace(xri)
			}
		}
	}

	if remoteIP == "" {
		return "127.0.0.1"
	}
	return remoteIP
}

// RateLimit creates an Echo middleware enforcing the specified rate limit policy
func RateLimit(cfg RateLimiterConfig) echo.MiddlewareFunc {
	cleanOnce.Do(initMemoryCleanup)

	if cfg.Limit <= 0 {
		cfg.Limit = 10
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = func(c echo.Context) string {
			return SafeExtractIP(c)
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			keySuffix := cfg.KeyFunc(c)
			key := fmt.Sprintf("rl:%s:%s", cfg.Action, keySuffix)

			var count int64
			var ttl time.Duration
			usedRedis := false

			// 1. Try Redis limiter first if Redis client is available
			if cfg.RedisClient != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				pipe := cfg.RedisClient.Pipeline()
				incrCmd := pipe.Incr(ctx, key)
				ttlCmd := pipe.TTL(ctx, key)
				_, err := pipe.Exec(ctx)
				cancel()

				if err == nil {
					usedRedis = true
					count = incrCmd.Val()
					ttl = ttlCmd.Val()
					// If key is new (no ttl set or negative), set expiry
					if ttl <= 0 {
						_ = cfg.RedisClient.Expire(context.Background(), key, cfg.Window).Err()
						ttl = cfg.Window
					}
				}
			}

			// 2. Fallback to in-memory limiter if Redis is unavailable or errored
			if !usedRedis {
				now := time.Now()
				rawVal, _ := memoryStore.Load(key)
				var entry memoryRateEntry
				if rawVal != nil {
					entry = rawVal.(memoryRateEntry)
				}

				if entry.expiresAt.IsZero() || now.After(entry.expiresAt) {
					entry = memoryRateEntry{
						count:     1,
						expiresAt: now.Add(cfg.Window),
					}
				} else {
					entry.count++
				}
				memoryStore.Store(key, entry)
				count = int64(entry.count)
				ttl = entry.expiresAt.Sub(now)
				if ttl < 0 {
					ttl = 0
				}
			}

			// Calculate headers
			remaining := int64(cfg.Limit) - count
			if remaining < 0 {
				remaining = 0
			}
			retryAfterSec := int(ttl.Seconds())
			if retryAfterSec <= 0 {
				retryAfterSec = 1
			}

			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Limit))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

			if count > int64(cfg.Limit) {
				c.Response().Header().Set("Retry-After", fmt.Sprintf("%d", retryAfterSec))

				// Return appropriate response for HTMX or HTML
				if c.Request().Header.Get("HX-Request") == "true" {
					c.Response().WriteHeader(http.StatusTooManyRequests)
					return c.HTML(http.StatusTooManyRequests, fmt.Sprintf(`
						<div class="alert alert-error" style="background: #fee2e2; border: 1.5px solid #fca5a5; color: #991b1b; padding: 1rem 1.25rem; border-radius: 12px; font-size: 0.92rem; font-weight: 600; text-align: center; margin: 1rem 0;">
							Too many attempts. For your security, please wait %d seconds before trying again.
						</div>
					`, retryAfterSec))
				}

				if strings.Contains(c.Request().Header.Get("Accept"), "application/json") {
					return c.JSON(http.StatusTooManyRequests, map[string]any{
						"error":       "Too many requests. Please try again later.",
						"retry_after": retryAfterSec,
					})
				}

				return c.HTML(http.StatusTooManyRequests, fmt.Sprintf(`
					<!DOCTYPE html>
					<html lang="en">
					<head>
						<meta charset="UTF-8">
						<title>Too Many Requests — Bite100</title>
						<meta name="viewport" content="width=device-width, initial-scale=1.0">
						<style>
							body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f8fafc; color: #1e293b; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; padding: 1rem; }
							.card { background: #ffffff; border: 1px solid #e2e8f0; border-radius: 16px; padding: 2.5rem 2rem; max-width: 440px; text-align: center; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.05); }
							h1 { font-size: 1.5rem; color: #b91c1c; margin-bottom: 0.75rem; }
							p { color: #64748b; font-size: 0.95rem; line-height: 1.5; margin-bottom: 1.5rem; }
							.btn { display: inline-block; background: #ea580c; color: white; padding: 0.65rem 1.5rem; border-radius: 8px; text-decoration: none; font-weight: 600; }
						</style>
					</head>
					<body>
						<div class="card">
							<h1>Too Many Requests (429)</h1>
							<p>For security, request rates are restricted. Please wait <strong>%d seconds</strong> before attempting this action again.</p>
							<a href="javascript:location.reload()" class="btn">Try Again</a>
						</div>
					</body>
					</html>
				`, retryAfterSec))
			}

			return next(c)
		}
	}
}
