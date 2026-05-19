package middleware

import (
	stderrors "errors"

	"github.com/gofiber/fiber/v2"

	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	httputil "github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
)

// ErrorHandler is the Fiber-level safety net. The 362 existing call sites of
// http.Error / http.NewError write responses themselves and return nil, so
// this handler does NOT run for them. It fires for:
//
//  1. errors returned from middleware before any route handler ran
//     (e.g. body-size limit, custom auth that bubbles an error),
//  2. panics caught by recover.New() and converted into errors,
//  3. any handler that returns an error instead of calling http.Error/NewError.
//
// Goal: produce exactly the same error envelope as http.NewError so clients
// never have to handle two shapes.
func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// http.Error / http.NewError already wrote the response and returned a
	// sentinel so handler chains short-circuit. Nothing more to do here —
	// re-writing would clobber the validation/field details.
	if stderrors.Is(err, httputil.ErrResponseWritten) {
		return nil
	}

	// Map Fiber's built-in *fiber.Error (e.g. 404 route not matched, body limit)
	// onto our envelope using a generic internal code; the HTTP status comes
	// straight from the *fiber.Error so 404 stays 404, 413 stays 413, etc.
	var fe *fiber.Error
	if stderrors.As(err, &fe) {
		return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrGeneric.Wrap(err))
	}

	// AppError chain — let NewError pick the right code + status.
	var appErr *errors.AppError
	if stderrors.As(err, &appErr) {
		return httputil.NewError(c, appErr.GetCode(), err)
	}

	// Anything else — treat as a server-side failure. NewError will log the
	// full chain and emit a generic message (details suppressed because
	// 500001 is not client-safe).
	return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrGeneric.Wrap(err))
}
