package http

import (
	stderrors "errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
)

type ErrorResponseModel struct {
	Code    string `json:"code" example:"100001"`
	Message string `json:"message" example:"User already exists"`
}

func (err *ErrorResponseModel) Error(c *fiber.Ctx, code string, errMessage string) {
	err.Code = code
	err.Message = fmt.Sprint(errMessage)
}

type ResponseModel struct {
	Data   any  `json:"data,omitempty"`
	Result bool `json:"result"`
	Error  any  `json:"error,omitempty"`
}

// NewError sends an error response (returns ErrorResponseModel directly, not wrapped in ResponseModel)
// Extracts error code from AppError if available (including wrapped errors), otherwise uses defaultCode
func NewError(c *fiber.Ctx, defaultCode int, err error) error {
	var code int
	var message string

	var appErr *errors.AppError
	if stderrors.As(err, &appErr) {
		code = appErr.GetCode()
		message = appErr.GetMessage()
	} else {
		code = defaultCode
		message = err.Error()
	}

	er := ErrorResponseModel{
		Code:    fmt.Sprintf("%d", code),
		Message: message,
	}
	return c.Status(fiber.StatusOK).JSON(er)
}

// Success sends a successful response with optional data
func Success(c *fiber.Ctx, data any) error {
	response := ResponseModel{
		Result: true,
		Data:   data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// SuccessWithoutData sends a successful response without data
func SuccessWithoutData(c *fiber.Ctx) error {
	response := ResponseModel{
		Result: true,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// Error sends an error response wrapped in ResponseModel
// Accepts either numeric code (int) or string code for backward compatibility
func Error(c *fiber.Ctx, code any, message string) error {
	var codeStr string
	switch v := code.(type) {
	case int:
		codeStr = fmt.Sprintf("%d", v)
	case string:
		codeStr = v
	default:
		codeStr = fmt.Sprintf("%v", code)
	}

	errRes := ErrorResponseModel{
		Code:    codeStr,
		Message: message,
	}
	response := ResponseModel{
		Error: errRes,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
