package http_test

import (
	"bytes"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	appmw "github.com/weeranieb/boonmafarm-backend/src/internal/middleware"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	httputil "github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
)

// End-to-end shape tests for the error envelope. These exercise the full
// stack: request_id middleware → handler → http.Error/NewError → Fiber
// ErrorHandler safety net, asserting on the wire bytes a client would see.

type registerDTO struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,password"`
}

func newApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: appmw.ErrorHandler})
	app.Use(appmw.RequestID())
	return app
}

func TestNewError_ValidationFailure_Returns422_WithFields(t *testing.T) {
	app := newApp()
	app.Post("/x", func(c *fiber.Ctx) error {
		var body registerDTO
		if err := c.BodyParser(&body); err != nil {
			return httputil.NewError(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Wrap(err))
		}
		if err := utils.ValidateStruct(&body); err != nil {
			return httputil.NewError(c, errors.ErrValidationFailed.Code, errors.ErrValidationFailed.Wrap(err))
		}
		return httputil.Success(c, nil)
	})

	req := httptest.NewRequest(stdhttp.MethodPost, "/x", bytes.NewBufferString(`{"username":"ab","password":"short"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusUnprocessableEntity, resp.StatusCode)

	var body httputil.ErrorResponseModel
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "500010", body.Code)
	require.NotEmpty(t, body.RequestID, "request_id must be present on every error response")
	require.NotEmpty(t, body.Fields, "validation errors must include per-field details")

	got := map[string]string{}
	for _, f := range body.Fields {
		got[f.Field] = f.Tag
	}
	require.Equal(t, "min", got["username"], "username should fail on min, not on required (length 2)")
	require.Equal(t, "password", got["password"], "password should fail on the custom password rule")
	require.NotEmpty(t, body.Details, "details summarizes fields[] for back-compat clients")
}

func TestNewError_BadJSONBody_Returns400(t *testing.T) {
	app := newApp()
	app.Post("/x", func(c *fiber.Ctx) error {
		var body registerDTO
		if err := c.BodyParser(&body); err != nil {
			return httputil.NewError(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Wrap(err))
		}
		return httputil.Success(c, nil)
	})

	req := httptest.NewRequest(stdhttp.MethodPost, "/x", bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusBadRequest, resp.StatusCode)

	var body httputil.ErrorResponseModel
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "500011", body.Code)
	require.NotEmpty(t, body.RequestID)
	require.Empty(t, body.Fields)
	require.NotEmpty(t, body.Details, "parser error message should land in details for client-safe code")
}

func TestNewError_NotFoundCode_Returns404(t *testing.T) {
	app := newApp()
	app.Get("/x", func(c *fiber.Ctx) error {
		return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrFarmNotFound)
	})

	resp, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/x", nil))
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusNotFound, resp.StatusCode)

	var body httputil.ErrorResponseModel
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "500040", body.Code)
}

func TestNewError_AlreadyExistsCode_Returns409(t *testing.T) {
	app := newApp()
	app.Post("/x", func(c *fiber.Ctx) error {
		return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrUserEmailAlreadyExists)
	})

	resp, err := app.Test(httptest.NewRequest(stdhttp.MethodPost, "/x", nil))
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusConflict, resp.StatusCode)
}

func TestNewError_AuthInvalidCredentials_Returns401(t *testing.T) {
	app := newApp()
	app.Post("/x", func(c *fiber.Ctx) error {
		return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrAuthInvalidCredentials)
	})

	resp, err := app.Test(httptest.NewRequest(stdhttp.MethodPost, "/x", nil))
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusUnauthorized, resp.StatusCode)
}

func TestNewError_PermissionDenied_Returns403(t *testing.T) {
	app := newApp()
	app.Get("/x", func(c *fiber.Ctx) error {
		return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrAuthPermissionDenied)
	})

	resp, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/x", nil))
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusForbidden, resp.StatusCode)
}

func TestNewError_InternalCode_DoesNotLeakWrappedDetails(t *testing.T) {
	// 500001 (ErrGeneric) is NOT client-safe — wrapped err must stay on server.
	app := newApp()
	app.Get("/x", func(c *fiber.Ctx) error {
		return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrGeneric.Wrap(
			validator.ValidationErrors{}, // arbitrary inner err
		))
	})

	resp, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/x", nil))
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusInternalServerError, resp.StatusCode)

	var body httputil.ErrorResponseModel
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "500001", body.Code)
	require.Empty(t, body.Details, "internal-code details must not be exposed to the client")
	require.Empty(t, body.Fields)
	require.NotEmpty(t, body.RequestID, "request_id must still be present so users can quote it when reporting failures")
}

func TestPanicInHandler_ReturnsGeneric500_NoStackInBody(t *testing.T) {
	// Wire recover.New + ErrorHandler as production does.
	app := fiber.New(fiber.Config{ErrorHandler: appmw.ErrorHandler})
	app.Use(appmw.RequestID())
	app.Use(recover.New())
	app.Get("/boom", func(c *fiber.Ctx) error {
		panic("simulated handler bug")
	})

	resp, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/boom", nil))
	require.NoError(t, err)
	require.Equal(t, stdhttp.StatusInternalServerError, resp.StatusCode)

	var body httputil.ErrorResponseModel
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "500001", body.Code)
	require.NotContains(t, body.Message, "simulated handler bug", "panic message must NOT leak to the client")
	require.Empty(t, body.Details, "internal-code details must not be exposed")
	require.NotEmpty(t, body.RequestID)
}

func TestRequestID_FromClientHeader_IsEchoed(t *testing.T) {
	app := newApp()
	app.Get("/x", func(c *fiber.Ctx) error {
		return httputil.NewError(c, errors.ErrGeneric.Code, errors.ErrFarmNotFound)
	})

	req := httptest.NewRequest(stdhttp.MethodGet, "/x", nil)
	req.Header.Set(fiber.HeaderXRequestID, "test-correlation-123")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, "test-correlation-123", resp.Header.Get(fiber.HeaderXRequestID))

	var body httputil.ErrorResponseModel
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "test-correlation-123", body.RequestID, "client-supplied request_id should appear in body")
}
