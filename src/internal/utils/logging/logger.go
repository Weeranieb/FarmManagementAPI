// Package logging provides structured request logging via std log/slog.
//
// Observability loop (local + Vercel share UseHTTP):
//
//	RequestID → TraceParent → Metrics → AccessLog → recover
//
// Ops quick ref:
//   - Grep failures: request_id (also X-Request-ID / error JSON).
//   - Liveness: GET /health   Readiness (DB): GET /ready
//   - Metrics: GET /metrics (Prometheus text). On Vercel cold starts wipe
//     process memory — use http_request logs (route, status_class, latency_ms).
//
// Usage rules:
//  1. Do not use log.Printf / fmt.Println on the request path.
//  2. Client-facing errors must go through http.Error / http.NewError
//     so they emit an http_error line with request context.
//  3. Extra request-path logs use FromCtx(c).Info/Warn/Error("snake_event", attrs...).
//  4. Never log bodies, Authorization, cookies, or query strings that may hold tokens.
package logging

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/metrics"
)

// RequestIDKey is the fiber.Ctx Locals key under which the per-request ID is stored.
// UseHTTP / middleware.RequestID seed it; FromCtx and the error envelope read it.
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

// UseHTTP mounts RequestID → TraceParent → Metrics → AccessLog → recover.
// Call after fiber.New(... ErrorHandler: middleware.ErrorHandler) so local
// and serverless bootstraps share one observability stack.
func UseHTTP(app *fiber.App) {
	app.Use(requestid.New(requestid.Config{
		Header:     fiber.HeaderXRequestID,
		ContextKey: RequestIDKey,
	}))
	app.Use(TraceParent())
	app.Use(metrics.Middleware())
	app.Use(AccessLog())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			metrics.IncPanic()
			FromCtx(c).Error("panic recovered",
				"panic", e,
				"stack", string(debug.Stack()),
			)
		},
	}))
}

// FromCtx returns a slog.Logger pre-bound with request-scoped attrs (request_id,
// method, path, user_id, client_id). Callers should use this in handlers /
// services so every log line carries the request context for grep correlation.
// path is c.Path() only — never the raw query string.
func FromCtx(c *fiber.Ctx) *slog.Logger {
	if c == nil {
		return slog.Default()
	}

	attrs := []any{
		"method", c.Method(),
		"path", c.Path(),
	}

	if rid := RequestIDFrom(c); rid != "" {
		attrs = append(attrs, "request_id", rid)
	}
	if tid := TraceIDFrom(c); tid != "" {
		attrs = append(attrs, "trace_id", tid)
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
