package http

import (
	stderrors "errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/logging"
)

// ErrResponseWritten is returned by Error/NewError so the caller's
// `if err != nil { return err }` chain actually short-circuits. Without it,
// c.Status(...).JSON(body) returns nil on success and the handler keeps
// running past helpers like validateAndParse, letting downstream code
// overwrite a 422 with whatever it does next. The Fiber ErrorHandler
// recognizes this sentinel and skips re-writing the response.
var ErrResponseWritten = stderrors.New("http: response written")

// ErrorResponseModel is the wire shape every error response — direct or wrapped
// — eventually serializes to. Fields are optional so the same struct also
// represents older shapes without breaking existing clients.
type ErrorResponseModel struct {
	Code      string       `json:"code" example:"500010"`
	Message   string       `json:"message" example:"Validation failed"`
	Details   string       `json:"details,omitempty" example:"freshFeedCollectionId is required"`
	RequestID string       `json:"request_id,omitempty" example:"01HXYZ123ABC"`
	Fields    []FieldError `json:"fields,omitempty"`
}

// isClientSafeCode reports whether an error code's wrapped detail can be
// exposed to clients. Validation (500010–500019) and domain rule codes
// (≥ 500020) wrap human-readable messages by convention. Internal/database
// codes (500001–500009, 500050–500059) wrap raw infra errors and must stay
// server-side.
func isClientSafeCode(code int) bool {
	if code >= 500001 && code <= 500009 {
		return false
	}
	if code >= 500050 && code <= 500059 {
		return false
	}
	return code >= 500010
}

type ResponseModel struct {
	Data   any  `json:"data,omitempty"`
	Result bool `json:"result"`
	Error  any  `json:"error,omitempty"`
}

// NewError sends an error response (returns ErrorResponseModel directly, not
// wrapped in ResponseModel). Extracts the AppError code/message/wrapped err
// from anywhere in the error chain, maps the code to a real HTTP status,
// surfaces field-level validation details when present, and logs the full
// chain server-side with the request ID for correlation.
func NewError(c *fiber.Ctx, defaultCode int, err error) error {
	code, message, wrapped := unwrap(err, defaultCode)
	body := buildErrorBody(c, code, message, wrapped)
	status := errors.HTTPStatusFor(code)

	logErr(c, status, code, err)
	if jsonErr := c.Status(status).JSON(body); jsonErr != nil {
		return jsonErr
	}
	return ErrResponseWritten
}

// Success sends a successful response with optional data.
func Success(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(ResponseModel{
		Result: true,
		Data:   data,
	})
}

// SuccessWithoutData sends a successful response without data.
func SuccessWithoutData(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(ResponseModel{
		Result: true,
	})
}

// Error sends an error response wrapped in ResponseModel. Accepts either a
// numeric code (int) or a string code for back-compat with existing handlers.
// Uses the same HTTP-status / logging path as NewError so the wire status
// reflects the real outcome.
func Error(c *fiber.Ctx, code any, message string) error {
	numericCode, codeStr := normalizeCode(code)
	status := errors.HTTPStatusFor(numericCode)

	errRes := ErrorResponseModel{
		Code:      codeStr,
		Message:   message,
		RequestID: logging.RequestIDFrom(c),
	}

	logErr(c, status, numericCode, stderrors.New(message))
	if jsonErr := c.Status(status).JSON(ResponseModel{Error: errRes}); jsonErr != nil {
		return jsonErr
	}
	return ErrResponseWritten
}

// --- internal helpers --------------------------------------------------------

func unwrap(err error, defaultCode int) (code int, message string, wrapped error) {
	var appErr *errors.AppError
	if stderrors.As(err, &appErr) {
		return appErr.GetCode(), appErr.GetMessage(), appErr.Err
	}
	return defaultCode, err.Error(), nil
}

func buildErrorBody(c *fiber.Ctx, code int, message string, wrapped error) ErrorResponseModel {
	body := ErrorResponseModel{
		Code:      fmt.Sprintf("%d", code),
		Message:   message,
		RequestID: logging.RequestIDFrom(c),
	}

	if wrapped == nil || !isClientSafeCode(code) {
		return body
	}

	if fields := BuildFieldErrors(wrapped); len(fields) > 0 {
		body.Fields = fields
		body.Details = SummarizeFieldErrors(fields)
		return body
	}

	body.Details = wrapped.Error()
	return body
}

func normalizeCode(code any) (int, string) {
	switch v := code.(type) {
	case int:
		return v, fmt.Sprintf("%d", v)
	case string:
		var n int
		_, _ = fmt.Sscanf(v, "%d", &n)
		return n, v
	default:
		s := fmt.Sprintf("%v", code)
		var n int
		_, _ = fmt.Sscanf(s, "%d", &n)
		return n, s
	}
}

func logErr(c *fiber.Ctx, status, code int, err error) {
	level := slog.LevelWarn
	if status >= 500 {
		level = slog.LevelError
	}
	logging.FromCtx(c).Log(c.UserContext(), level, "http_error",
		"status", status,
		"code", code,
		"err", err.Error(),
	)
}
