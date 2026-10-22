package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/pkg/services"
)

type HealthHandler struct {
	c *services.Container
}

func NewHealthHandler(c *services.Container) *HealthHandler {
	return &HealthHandler{c: c}
}

// Live confirms that the application process is running and accepting HTTP requests
func (h *HealthHandler) Live(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "bite100-food-discovery",
	})
}

// Ready verifies that critical backing services (PostgreSQL and Redis) are responsive
func (h *HealthHandler) Ready(ctx echo.Context) error {
	checkCtx, cancel := context.WithTimeout(ctx.Request().Context(), 2*time.Second)
	defer cancel()

	isReady := true
	details := map[string]string{
		"database": "healthy",
		"redis":    "healthy",
	}

	// 1. Check PostgreSQL Database connection
	if h.c.Database == nil {
		isReady = false
		details["database"] = "uninitialized"
	} else if err := h.c.Database.PingContext(checkCtx); err != nil {
		isReady = false
		details["database"] = "unreachable"
	}

	// 2. Check Redis Cache connection
	if h.c.Cache != nil && h.c.Cache.Client != nil {
		if err := h.c.Cache.Client.Ping(checkCtx).Err(); err != nil {
			isReady = false
			details["redis"] = "unreachable"
		}
	} else {
		details["redis"] = "disabled"
	}

	statusCode := http.StatusOK
	statusStr := "ready"
	if !isReady {
		statusCode = http.StatusServiceUnavailable
		statusStr = "degraded"
	}

	return ctx.JSON(statusCode, map[string]any{
		"status":    statusStr,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    details,
	})
}
