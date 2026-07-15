package logging

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	// TraceIDKey is the Locals key for a W3C trace_id when traceparent is present.
	TraceIDKey = "trace_id"
)

// TraceParent parses an inbound W3C traceparent header (if valid) and stashes
// trace_id for slog. No OpenTelemetry SDK — zero cold-start cost when absent.
func TraceParent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if tid := parseTraceID(c.Get("traceparent")); tid != "" {
			c.Locals(TraceIDKey, tid)
		}
		return c.Next()
	}
}

// TraceIDFrom returns the trace_id from Locals, or "".
func TraceIDFrom(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Locals(TraceIDKey).(string); ok {
		return v
	}
	return ""
}

// parseTraceID extracts the 32-hex trace-id from version-traceid-spanid-flags.
func parseTraceID(header string) string {
	parts := strings.Split(strings.TrimSpace(header), "-")
	if len(parts) != 4 {
		return ""
	}
	// version 00 is current; accept any 2-hex version to stay forward-compatible.
	if len(parts[0]) != 2 || len(parts[1]) != 32 || len(parts[2]) != 16 || len(parts[3]) != 2 {
		return ""
	}
	for _, c := range parts[1] {
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		case c >= 'A' && c <= 'F':
		default:
			return ""
		}
	}
	return strings.ToLower(parts[1])
}
