package logging

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/metrics"
)

// AccessLog emits one structured slog line per HTTP request. Replaces Fiber's
// default text logger so request_id, latency, and status are queryable.
//
// Includes low-cardinality route + status_class so Vercel log drains can derive
// RED metrics without scraping /metrics.
//
// Lines are level Info for 1xx-3xx, Warn for 4xx, and Error for 5xx.
func AccessLog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		latencyMs := time.Since(start).Milliseconds()
		status := c.Response().StatusCode()
		if err != nil {
			var fe *fiber.Error
			if errors.As(err, &fe) && fe.Code > 0 {
				status = fe.Code
			}
		}

		level := slog.LevelInfo
		switch {
		case status >= 500:
			level = slog.LevelError
		case status >= 400:
			level = slog.LevelWarn
		}

		attrs := []any{
			"status", status,
			"status_class", metrics.StatusClass(status),
			"route", metrics.RouteLabel(c, err),
			"latency_ms", latencyMs,
			"bytes", len(c.Response().Body()),
			"ip", c.IP(),
		}
		if ua := c.Get(fiber.HeaderUserAgent); ua != "" {
			attrs = append(attrs, "user_agent", ua)
		}

		FromCtx(c).Log(c.UserContext(), level, "http_request", attrs...)
		return err
	}
}
