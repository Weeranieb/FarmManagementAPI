package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// RequestID seeds a per-request ID on every inbound request. The ID is:
//
//   - read from the X-Request-ID header if the caller supplied one
//     (useful for client-side correlation),
//   - otherwise freshly generated as a UUID-v4-style string,
//   - echoed back to the client via the X-Request-ID response header,
//   - stashed in the request context for the logger and error envelope to pick
//     up via logging.RequestIDFrom (backed by requestid.FromContext in v3).
func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: fiber.HeaderXRequestID,
	})
}
