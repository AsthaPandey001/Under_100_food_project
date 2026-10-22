package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSafeExtractIP(t *testing.T) {
	e := echo.New()

	t.Run("Untrusted public remote addr ignores X-Forwarded-For", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "203.0.113.10:12345"
		req.Header.Set("X-Forwarded-For", "198.51.100.99, 70.41.3.18")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		ip := SafeExtractIP(c)
		assert.Equal(t, "203.0.113.10", ip)
	})

	t.Run("Loopback proxy respects valid X-Forwarded-For client IP", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		req.Header.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		ip := SafeExtractIP(c)
		assert.Equal(t, "203.0.113.195", ip)
	})

	t.Run("Private proxy respects valid X-Forwarded-For client IP", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.1:8080"
		req.Header.Set("X-Forwarded-For", "198.51.100.55")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		ip := SafeExtractIP(c)
		assert.Equal(t, "198.51.100.55", ip)
	})
}

func TestRateLimiter_InMemoryFallback(t *testing.T) {
	e := echo.New()

	cfg := RateLimiterConfig{
		Action:      "test_fallback",
		Limit:       3,
		Window:      time.Minute,
		RedisClient: nil, // Forces in-memory fallback
	}

	handler := RateLimit(cfg)(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Make 3 allowed requests
	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.RemoteAddr = "192.0.2.1:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "3", rec.Header().Get("X-RateLimit-Limit"))
	}

	// 4th request must be rate limited with HTTP 429
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "192.0.2.1:1234"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("Retry-After"))
	assert.Equal(t, "0", rec.Header().Get("X-RateLimit-Remaining"))

	// Another client IP should still be allowed
	req2 := httptest.NewRequest(http.MethodPost, "/test", nil)
	req2.RemoteAddr = "192.0.2.2:1234"
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	err = handler(c2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec2.Code)
}

func TestRateLimiter_ConcurrentRequests(t *testing.T) {
	e := echo.New()

	limit := 10
	cfg := RateLimiterConfig{
		Action:      "test_concurrent",
		Limit:       limit,
		Window:      time.Minute,
		RedisClient: nil,
	}

	handler := RateLimit(cfg)(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	var wg sync.WaitGroup
	totalRequests := 25
	statusCodes := make([]int, totalRequests)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/test-conc", nil)
			req.RemoteAddr = "198.51.100.10:5432"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = handler(c)
			statusCodes[idx] = rec.Code
		}(i)
	}

	wg.Wait()

	okCount := 0
	rateLimitedCount := 0
	for _, code := range statusCodes {
		if code == http.StatusOK {
			okCount++
		} else if code == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	assert.Equal(t, limit, okCount, "Exactly Limit should succeed")
	assert.Equal(t, totalRequests-limit, rateLimitedCount, "Burst exceeding limit should receive 429")
}
