package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/logging"
)

// RequestID seeds a per-request ID on every inbound request. The ID is:
//
//   - read from the X-Request-ID header if the caller supplied one
//     (useful for client-side correlation),
//   - otherwise freshly generated as a UUID-v4-style string,
//   - echoed back to the client via the X-Request-ID response header,
//   - stashed on c.Locals(logging.RequestIDKey) for the logger and
//     error envelope to pick up.
func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header:     fiber.HeaderXRequestID,
		ContextKey: logging.RequestIDKey,
	})
}
