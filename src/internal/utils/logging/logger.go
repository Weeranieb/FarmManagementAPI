package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
)

// RequestIDKey is the fiber.Ctx Locals key under which the per-request ID is stored.
// The request_id middleware seeds it; logging.FromCtx and the error envelope read it.
const RequestIDKey = "request_id"

// Init configures the default slog logger.
//
//   - production / staging  → JSON, level from cfg
//   - everything else (dev) → human-readable text, level from cfg
//
// Returns the logger so callers may store it; slog.SetDefault is called as a
// side effect so any package that just uses slog.Info(...) also benefits.
func Init(env string, level string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     parseLevel(level),
		AddSource: !isProd(env),
	}

	var handler slog.Handler
	if isProd(env) {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// FromCtx returns a slog.Logger pre-bound with request-scoped attrs (request_id,
// method, path, user_id, client_id). Callers should use this in handlers /
// services so every log line carries the request context for grep correlation.
func FromCtx(c *fiber.Ctx) *slog.Logger {
	if c == nil {
		return slog.Default()
	}

	attrs := []any{
		"method", c.Method(),
		"path", c.OriginalURL(),
	}

	if rid := RequestIDFrom(c); rid != "" {
		attrs = append(attrs, "request_id", rid)
	}

	uctx := c.UserContext()
	if uctx == nil {
		uctx = context.Background()
	}
	if uid, ok := uctx.Value(constants.UserIDKey).(int); ok {
		attrs = append(attrs, "user_id", uid)
	}
	if cid, ok := uctx.Value(constants.ClientIDKey).(int); ok {
		attrs = append(attrs, "client_id", cid)
	}

	return slog.Default().With(attrs...)
}

// RequestIDFrom returns the request ID stashed by the request_id middleware,
// or "" if none is set yet (e.g. before middleware runs).
func RequestIDFrom(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Locals(RequestIDKey).(string); ok {
		return v
	}
	return ""
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func isProd(env string) bool {
	e := strings.ToLower(strings.TrimSpace(env))
	return e == "production" || e == "prod" || e == "staging"
}
