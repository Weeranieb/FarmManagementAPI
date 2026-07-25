//go:build cgo

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	apperrors "github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

type PondHandlerTestSuite struct {
	suite.Suite
	pondService *mocks.MockPondService
	pondHandler PondHandler
}

func (s *PondHandlerTestSuite) SetupTest() {
	s.pondService = mocks.NewMockPondService(s.T())
	s.pondHandler = NewPondHandler(s.pondService)
}

func (s *PondHandlerTestSuite) TearDownTest() {
	s.pondService.ExpectedCalls = nil
}

func TestPondHandlerSuite(t *testing.T) {
	suite.Run(t, new(PondHandlerTestSuite))
}

func (s *PondHandlerTestSuite) TestAddPonds_Success() {
	// GIVEN — valid create request; super admin context; service mock returns nil
	createReq := &dto.CreatePondsRequest{
		FarmId: 1,
		Ponds:  []dto.CreatePondItem{{Name: "Pond 1"}, {Name: "Pond 2"}},
	}
	username := "admin"
	s.pondService.On("CreatePonds", mock.Anything, *createReq).Return(nil)
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"userLevel": 3,
	}))
	app.Post("/api/v1/pond", s.pondHandler.AddPonds)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond is sent
	resp, err := app.Test(req)

	// THEN — 200 OK; service was called
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestAddPonds_NonSuperAdmin_ReturnsPermissionDenied() {
	// GIVEN — valid create request; user is not super admin (userLevel 1)
	createReq := &dto.CreatePondsRequest{
		FarmId: 1,
		Ponds:  []dto.CreatePondItem{{Name: "Pond 1"}},
	}
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "user",
		"userLevel": 1,
	}))
	app.Post("/api/v1/pond", s.pondHandler.AddPonds)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond is sent
	resp, err := app.Test(req)

	// THEN — 403 (ErrAuthPermissionDenied — not client-admin or above)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		assert.Equal(s.T(), "500024", errObj["code"])
	}
}

func (s *PondHandlerTestSuite) TestAddPonds_IsSuperAdminError() {
	// GIVEN — valid create request; no user context (empty locals)
	createReq := &dto.CreatePondsRequest{
		FarmId: 1,
		Ponds:  []dto.CreatePondItem{{Name: "Pond 1"}},
	}
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{}))
	app.Post("/api/v1/pond", s.pondHandler.AddPonds)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond is sent
	resp, err := app.Test(req)

	// THEN — 403 (ErrAuthPermissionDenied: IsClientAdminOrAbove fails on missing userLevel)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *PondHandlerTestSuite) TestGetPond_Success() {
	// GIVEN — pond id 1; service returns pond response
	pondId := 1
	expectedResponse := &dto.PondResponse{
		Id:     pondId,
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
	}

	s.pondService.On("Get", mock.Anything, pondId).Return(expectedResponse, nil)

	app := newTestApp()
	app.Get("/api/v1/pond/:id", s.pondHandler.GetPond)

	req := httptest.NewRequest("GET", "/api/v1/pond/1", nil)

	// WHEN — GET /api/v1/pond/1 is sent
	resp, err := app.Test(req)

	// THEN — 200 OK; service was called
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestGetPondList_Success() {
	// GIVEN — farmId 1; service returns list of ponds
	farmId := 1
	expectedResponse := []*dto.PondResponse{
		{Id: 1, FarmId: farmId, Name: "Pond 1", Status: "active"},
		{Id: 2, FarmId: farmId, Name: "Pond 2", Status: "active"},
	}

	s.pondService.On("GetList", mock.Anything, farmId).Return(expectedResponse, nil)

	app := newTestApp()
	app.Get("/api/v1/pond", s.pondHandler.GetPondList)

	req := httptest.NewRequest("GET", "/api/v1/pond?farmId=1", nil)

	// WHEN — GET /api/v1/pond?farmId=1 is sent
	resp, err := app.Test(req)

	// THEN — 200 OK; service was called
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

// fillPondApp returns a Fiber app with FillPond route and optional username in context
func (s *PondHandlerTestSuite) fillPondApp(username string) *fiber.App {
	app := newTestApp()
	locals := map[string]any{}
	if username != "" {
		locals["username"] = username
	}
	app.Use(setLocalsMiddleware(locals))
	app.Post("/api/v1/pond/:pondId/fill", s.pondHandler.FillPond)
	return app
}

func (s *PondHandlerTestSuite) TestFillPond_InvalidPondID_ReturnsValidationError() {
	// GIVEN — body with invalid pond id "abc"
	s.pondService.On("FillPond", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*dto.PondFillResponse)(nil), errors.New("")).Maybe()
	body := []byte(`{"fishType":"nil","amount":100,"fishWeight":0.5,"pricePerUnit":10.5,"activityDate":"2024-01-15"}`)
	app := s.fillPondApp("user")
	req := httptest.NewRequest("POST", "/api/v1/pond/abc/fill", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/abc/fill is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 422 (ErrValidationFailed: pondId param failed strconv.Atoi)
	require.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)

	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"], "expected error for invalid pond ID")
	errObj, ok := result["error"].(map[string]any)
	require.True(s.T(), ok)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestFillPond_MissingUsername_ReturnsAuthError() {
	// GIVEN — valid body; no username in context
	s.pondService.On("FillPond", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*dto.PondFillResponse)(nil), errors.New("")).Maybe()
	body := []byte(`{"fishType":"nil","amount":100,"fishWeight":0.5,"pricePerUnit":10.5,"activityDate":"2024-01-15"}`)
	app := s.fillPondApp("")
	req := httptest.NewRequest("POST", "/api/v1/pond/1/fill", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/1/fill is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 401 (ErrAuthTokenInvalid: username missing in context)
	require.Equal(s.T(), fiber.StatusUnauthorized, resp.StatusCode)

	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.True(s.T(), result["error"] != nil, "expected auth error when username missing")
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		assert.Equal(s.T(), "500022", errObj["code"])
	}
}

func (s *PondHandlerTestSuite) TestFillPond_Success() {
	// GIVEN — valid fill body; username; service returns success response
	pondId := 1
	username := "admin"
	body := []byte(`{"fishType":"nil","amount":100,"fishWeight":0.5,"pricePerUnit":10.5,"activityDate":"2024-01-15"}`)
	expectedResponse := &dto.PondFillResponse{ActivityId: 1, ActivePondId: 1}
	s.pondService.On("FillPond", mock.Anything, pondId, mock.Anything, username).Return(expectedResponse, nil)
	app := s.fillPondApp(username)
	req := httptest.NewRequest("POST", "/api/v1/pond/1/fill", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/1/fill is sent
	resp, err := app.Test(req)

	// THEN — 200 OK; result true and data present
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	assert.NotNil(s.T(), result["data"])
	s.pondService.AssertExpectations(s.T())
}

// movePondApp returns a Fiber app with POST /pond/:pondId/move and optional username in context.
func (s *PondHandlerTestSuite) movePondApp(username string) *fiber.App {
	app := newTestApp()
	locals := map[string]any{}
	if username != "" {
		locals["username"] = username
	}
	app.Use(setLocalsMiddleware(locals))
	app.Post("/api/v1/pond/:pondId/move", s.pondHandler.MovePond)
	return app
}

func (s *PondHandlerTestSuite) TestMovePond_Success() {
	// GIVEN — valid move body; username; service returns success response
	sourcePondId := 1
	username := "admin"
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"fishWeight":"0.5","pricePerUnit":"10.5","activityDate":"2024-06-01"}`)
	expectedResponse := &dto.PondMoveResponse{ActivityId: 1, ActivePondId: 10, ToActivePondId: 20}
	s.pondService.On("MovePond", mock.Anything, sourcePondId, mock.Anything, username).Return(expectedResponse, nil)
	app := s.movePondApp(username)
	req := httptest.NewRequest("POST", "/api/v1/pond/1/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/1/move is sent
	resp, err := app.Test(req)

	// THEN — 200 OK; result true and data present
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	assert.NotNil(s.T(), result["data"])
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestMovePond_InvalidPondID_ReturnsValidationError() {
	// GIVEN — body with invalid pond id "abc"
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"activityDate":"2024-06-01"}`)
	app := s.movePondApp("user")
	req := httptest.NewRequest("POST", "/api/v1/pond/abc/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/abc/move is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 422 (ErrValidationFailed: pondId param failed strconv.Atoi)
	require.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"], "expected error for invalid pond ID")
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestMovePond_MissingUsername_ReturnsAuthError() {
	// GIVEN — valid body; no username in context
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"fishWeight":"0.5","pricePerUnit":"10.5","activityDate":"2024-06-01"}`)
	app := s.movePondApp("")
	req := httptest.NewRequest("POST", "/api/v1/pond/1/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/1/move is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 401 (ErrAuthTokenInvalid: username missing in context)
	require.Equal(s.T(), fiber.StatusUnauthorized, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		assert.Equal(s.T(), "500022", errObj["code"])
	}
}

func (s *PondHandlerTestSuite) TestMovePond_ServiceError_ErrPondNotFound() {
	// GIVEN — pond 999; service returns ErrPondNotFound
	username := "user"
	s.pondService.On("MovePond", mock.Anything, 999, mock.Anything, username).Return((*dto.PondMoveResponse)(nil), apperrors.ErrPondNotFound)
	app := s.movePondApp(username)
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"fishWeight":"0.5","pricePerUnit":"10.5","activityDate":"2024-06-01"}`)
	req := httptest.NewRequest("POST", "/api/v1/pond/999/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond/999/move is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 404 (ErrPondNotFound → code 500070)
	assert.Equal(s.T(), fiber.StatusNotFound, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), "500070", result["code"])
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestMovePond_ServiceError_ErrPondSourceNotActive() {
	// GIVEN — pond 1; service returns ErrPondSourceNotActive
	username := "user"
	s.pondService.On("MovePond", mock.Anything, 1, mock.Anything, username).Return((*dto.PondMoveResponse)(nil), apperrors.ErrPondSourceNotActive)
	app := s.movePondApp(username)
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"fishWeight":"0.5","pricePerUnit":"10.5","activityDate":"2024-06-01"}`)
	req := httptest.NewRequest("POST", "/api/v1/pond/1/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond/1/move is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 422 (ErrPondSourceNotActive → code 500074)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), "500074", result["code"])
	s.pondService.AssertExpectations(s.T())
}

// updatePondApp returns a Fiber app with PUT /pond/:id and optional username in context.
func (s *PondHandlerTestSuite) updatePondApp(username string) *fiber.App {
	app := newTestApp()
	locals := map[string]any{}
	if username != "" {
		locals["username"] = username
	}
	app.Use(setLocalsMiddleware(locals))
	app.Put("/api/v1/pond/:id", s.pondHandler.UpdatePond)
	return app
}

func (s *PondHandlerTestSuite) TestUpdatePond_Success() {
	// GIVEN — valid body; service returns nil
	pondId := 1
	username := "admin"
	body := dto.UpdatePondBody{Name: "Updated Pond", Status: "active"}
	s.pondService.On("Update", mock.Anything, dto.UpdatePondRequest{
		Id: pondId, FarmId: body.FarmId, Name: body.Name, Status: body.Status,
	}).Return(nil)
	app := s.updatePondApp(username)
	reqBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/pond/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/pond/1 is sent
	resp, err := app.Test(req)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — result is true
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestUpdatePond_InvalidPondID() {
	// GIVEN — invalid pond id "abc"
	body := []byte(`{"name":"Pond"}`)
	app := s.updatePondApp("user")
	req := httptest.NewRequest("PUT", "/api/v1/pond/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/pond/abc is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 422 (ErrValidationFailed: pondId param failed strconv.Atoi)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"])
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestUpdatePond_ServiceError() {
	// GIVEN — pond 999; service returns ErrPondNotFound
	username := "user"
	s.pondService.On("Update", mock.Anything, mock.AnythingOfType("dto.UpdatePondRequest")).Return(apperrors.ErrPondNotFound)
	app := s.updatePondApp(username)
	req := httptest.NewRequest("PUT", "/api/v1/pond/999", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/pond/999 is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 404 (ErrPondNotFound → code 500070)
	assert.Equal(s.T(), fiber.StatusNotFound, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), "500070", result["code"])
	s.pondService.AssertExpectations(s.T())
}

// deletePondApp returns a Fiber app with DELETE /pond/:id and optional username in context.
func (s *PondHandlerTestSuite) deletePondApp(username string) *fiber.App {
	app := newTestApp()
	locals := map[string]any{}
	if username != "" {
		locals["username"] = username
	}
	app.Use(setLocalsMiddleware(locals))
	app.Delete("/api/v1/pond/:id", s.pondHandler.DeletePond)
	return app
}

func (s *PondHandlerTestSuite) TestDeletePond_Success() {
	// GIVEN — pond 1; service returns nil
	pondId := 1
	username := "admin"
	s.pondService.On("Delete", mock.Anything, pondId).Return(nil)
	app := s.deletePondApp(username)
	req := httptest.NewRequest("DELETE", "/api/v1/pond/1", nil)

	// WHEN — DELETE /api/v1/pond/1 is sent
	resp, err := app.Test(req)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — result is true
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestDeletePond_InvalidPondID() {
	// GIVEN — invalid pond id "not-a-number"
	app := s.deletePondApp("user")
	req := httptest.NewRequest("DELETE", "/api/v1/pond/not-a-number", nil)

	// WHEN — DELETE /api/v1/pond/not-a-number is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 422 (ErrValidationFailed: pondId param failed strconv.Atoi)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"])
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestDeletePond_ServiceError() {
	// GIVEN — pond 1; service returns ErrGeneric
	username := "user"
	s.pondService.On("Delete", mock.Anything, 1).Return(apperrors.ErrGeneric)
	app := s.deletePondApp(username)
	req := httptest.NewRequest("DELETE", "/api/v1/pond/1", nil)

	// WHEN — DELETE /api/v1/pond/1 is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	// THEN — 500 (ErrGeneric → code 500001)
	assert.Equal(s.T(), fiber.StatusInternalServerError, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotEmpty(s.T(), result["message"])
	s.pondService.AssertExpectations(s.T())
}

// TestPondFillRequest_Validation ensures PondFillRequest validation rejects invalid input (handler uses validateAndParse).
func TestPondFillRequest_Validation(t *testing.T) {
	t.Run("missing required fields", func(t *testing.T) {
		// GIVEN — empty PondFillRequest
		req := &dto.PondFillRequest{}
		// WHEN — ValidateStruct is called
		err := utils.ValidateStruct(req)
		// THEN — validation fails
		require.Error(t, err)
	})
	t.Run("amount less than 1", func(t *testing.T) {
		// GIVEN — request with Amount 0
		req := &dto.PondFillRequest{
			FishType:     "nil",
			Amount:       0,
			PricePerUnit: decimal.NewFromFloat(10.5),
			ActivityDate: "2024-01-15",
		}
		// WHEN — ValidateStruct is called
		err := utils.ValidateStruct(req)
		// THEN — validation fails
		require.Error(t, err)
	})
	t.Run("pricePerUnit zero", func(t *testing.T) {
		// GIVEN — request with PricePerUnit zero
		req := &dto.PondFillRequest{
			FishType:     "nil",
			Amount:       100,
			PricePerUnit: decimal.Zero,
			ActivityDate: "2024-01-15",
		}
		// WHEN — ValidateStruct is called
		err := utils.ValidateStruct(req)
		// THEN — validation fails
		require.Error(t, err)
	})
	t.Run("fishWeight zero when provided", func(t *testing.T) {
		// GIVEN — request with FishWeight zero
		req := &dto.PondFillRequest{
			FishType:     "nil",
			Amount:       100,
			FishWeight:   decimal.Zero,
			PricePerUnit: decimal.NewFromFloat(10.5),
			ActivityDate: "2024-01-15",
		}
		// WHEN — ValidateStruct is called
		err := utils.ValidateStruct(req)
		// THEN — validation fails
		require.Error(t, err)
	})
	t.Run("valid request", func(t *testing.T) {
		// GIVEN — request with required fields and valid values
		req := &dto.PondFillRequest{
			FishType:     "nil",
			Amount:       100,
			FishWeight:   decimal.NewFromFloat(0.5),
			PricePerUnit: decimal.NewFromFloat(10.5),
			ActivityDate: "2024-01-15",
		}
		// WHEN — ValidateStruct is called
		err := utils.ValidateStruct(req)
		// THEN — validation passes
		require.NoError(t, err)
	})
}

// TestPondMoveRequest_Validation ensures PondMoveRequest validation rejects
// payloads that would book a meaningless move (amount × weight × price = 0).
// Mirrors the fill-request suite — the move DTO previously used a looser
// validator (omitempty,decimal_gte0) on FishWeight, allowing zero-weight
// rows to be persisted. See the comment on PondMoveRequest for the
// accounting rationale.
func TestPondMoveRequest_Validation(t *testing.T) {
	validBase := func() *dto.PondMoveRequest {
		return &dto.PondMoveRequest{
			ToPondId:     2,
			FishType:     "nil",
			Amount:       100,
			FishWeight:   decimal.NewFromFloat(0.5),
			PricePerUnit: decimal.NewFromFloat(70),
			ActivityDate: "2024-01-15",
		}
	}

	t.Run("missing required fields", func(t *testing.T) {
		err := utils.ValidateStruct(&dto.PondMoveRequest{})
		require.Error(t, err)
	})
	t.Run("amount less than 1", func(t *testing.T) {
		req := validBase()
		req.Amount = 0
		require.Error(t, utils.ValidateStruct(req))
	})
	t.Run("fishWeight zero", func(t *testing.T) {
		req := validBase()
		req.FishWeight = decimal.Zero
		require.Error(t, utils.ValidateStruct(req))
	})
	t.Run("fishWeight missing (decimal default)", func(t *testing.T) {
		req := validBase()
		req.FishWeight = decimal.Decimal{}
		require.Error(t, utils.ValidateStruct(req))
	})
	t.Run("pricePerUnit zero", func(t *testing.T) {
		req := validBase()
		req.PricePerUnit = decimal.Zero
		require.Error(t, utils.ValidateStruct(req))
	})
	t.Run("valid request", func(t *testing.T) {
		require.NoError(t, utils.ValidateStruct(validBase()))
	})
}

// --- AddPonds permissions (relaxed to client-admin-or-above) ---

func (s *PondHandlerTestSuite) TestAddPonds_ClientAdmin_Allowed() {
	// GIVEN — client admin (level 2); service returns nil.
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Ponds:  []dto.CreatePondItem{{Name: "P1"}},
	}
	s.pondService.On("CreatePonds", mock.Anything, req).Return(nil)
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "clientAdmin",
		"clientId":  1,
		"userLevel": 2,
	}))
	app.Post("/pond", s.pondHandler.AddPonds)
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/pond", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(httpReq)

	// THEN — 200; per-client scoping is enforced in pondService.CreatePonds
	// (covered by service tests), so the handler delegates without rejecting.
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

// --- BulkImportFarmPond handler ---

func (s *PondHandlerTestSuite) TestBulkImportFarmPond_Success() {
	// GIVEN — super-admin, valid payload, service returns a response.
	clientId := 1
	reqBody := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "Farm A", Ponds: []dto.BulkImportPondItem{{Name: "P1"}}},
		},
	}
	svcResp := &dto.BulkImportFarmPondResponse{
		FarmsCreated: 1,
		PondsCreated: 1,
		Farms: []dto.BulkImportFarmResult{
			{Name: "Farm A", IsNew: true, PondsCreated: 1},
		},
	}
	s.pondService.On("BulkImportFarmPond", mock.Anything, clientId, reqBody).Return(svcResp, nil)
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"userLevel": 3,
	}))
	app.Post("/pond/bulk-import/:clientId", s.pondHandler.BulkImportFarmPond)
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest("POST", "/pond/bulk-import/"+strconv.Itoa(clientId), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(httpReq)

	// THEN — 200 with result.data carrying the service response.
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(s.T(), true, result["result"])
	data, _ := result["data"].(map[string]any)
	require.NotNil(s.T(), data)
	assert.Equal(s.T(), float64(1), data["farmsCreated"])
	assert.Equal(s.T(), float64(1), data["pondsCreated"])
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestBulkImportFarmPond_InvalidClientIdParam() {
	// GIVEN — :clientId is not numeric.
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"userLevel": 3,
	}))
	app.Post("/pond/bulk-import/:clientId", s.pondHandler.BulkImportFarmPond)
	httpReq := httptest.NewRequest("POST", "/pond/bulk-import/abc", bytes.NewBufferString("{}"))
	httpReq.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(httpReq)

	// THEN — 422 (ErrValidationFailed: clientId param failed strconv.Atoi)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	errObj, _ := result["error"].(map[string]any)
	require.NotNil(s.T(), errObj)
	assert.Equal(s.T(), strconv.Itoa(apperrors.ErrValidationFailed.Code), errObj["code"])
	s.pondService.AssertNotCalled(s.T(), "BulkImportFarmPond", mock.Anything, mock.Anything, mock.Anything)
}

func (s *PondHandlerTestSuite) TestBulkImportFarmPond_NonAdmin_Rejected() {
	// GIVEN — regular user (level 1) — not even client-admin.
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "regular",
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Post("/pond/bulk-import/:clientId", s.pondHandler.BulkImportFarmPond)
	body, _ := json.Marshal(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "F", Ponds: []dto.BulkImportPondItem{{Name: "P"}}},
		},
	})
	httpReq := httptest.NewRequest("POST", "/pond/bulk-import/1", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(httpReq)

	// THEN — 403 (ErrAuthPermissionDenied: below client-admin); service untouched.
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	errObj, _ := result["error"].(map[string]any)
	require.NotNil(s.T(), errObj)
	assert.Equal(s.T(), strconv.Itoa(apperrors.ErrAuthPermissionDenied.Code), errObj["code"])
	s.pondService.AssertNotCalled(s.T(), "BulkImportFarmPond", mock.Anything, mock.Anything, mock.Anything)
}

func (s *PondHandlerTestSuite) TestBulkImportFarmPond_ClientAdminWrongClient_Rejected() {
	// GIVEN — client-admin (level 2) for client 1 trying to target client 2.
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 2,
	}))
	app.Post("/pond/bulk-import/:clientId", s.pondHandler.BulkImportFarmPond)
	body, _ := json.Marshal(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "F", Ponds: []dto.BulkImportPondItem{{Name: "P"}}},
		},
	})
	httpReq := httptest.NewRequest("POST", "/pond/bulk-import/2", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(httpReq)

	// THEN — 403 (ErrAuthPermissionDenied: cross-client access denied); service untouched.
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	errObj, _ := result["error"].(map[string]any)
	require.NotNil(s.T(), errObj)
	assert.Equal(s.T(), strconv.Itoa(apperrors.ErrAuthPermissionDenied.Code), errObj["code"])
	s.pondService.AssertNotCalled(s.T(), "BulkImportFarmPond", mock.Anything, mock.Anything, mock.Anything)
}

func (s *PondHandlerTestSuite) TestBulkImportFarmPond_ClientAdminOwnClient_Allowed() {
	// GIVEN — client-admin for client 1 targeting their own client 1.
	clientId := 1
	reqBody := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "F", Ponds: []dto.BulkImportPondItem{{Name: "P"}}},
		},
	}
	svcResp := &dto.BulkImportFarmPondResponse{FarmsCreated: 1, PondsCreated: 1}
	s.pondService.On("BulkImportFarmPond", mock.Anything, clientId, reqBody).Return(svcResp, nil)
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  clientId,
		"userLevel": 2,
	}))
	app.Post("/pond/bulk-import/:clientId", s.pondHandler.BulkImportFarmPond)
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest("POST", "/pond/bulk-import/"+strconv.Itoa(clientId), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(httpReq)

	// THEN — 200; delegation happened.
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

// (Struct-tag validation for nested ponds is exercised by the existing
// TestAddFarm_ValidationFailed pattern and by service-level
// validateBulkImportRequest tests; not duplicated at the handler level.)

func (s *PondHandlerTestSuite) TestBulkImportFarmPond_ServiceError() {
	// GIVEN — admin + valid request, but service returns an error.
	clientId := 1
	reqBody := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "F", Ponds: []dto.BulkImportPondItem{{Name: "P"}}},
		},
	}
	s.pondService.On("BulkImportFarmPond", mock.Anything, clientId, reqBody).Return(
		nil, apperrors.ErrValidationFailed.Wrap(errors.New("duplicate pond")))
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"userLevel": 3,
	}))
	app.Post("/pond/bulk-import/:clientId", s.pondHandler.BulkImportFarmPond)
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequest("POST", "/pond/bulk-import/"+strconv.Itoa(clientId), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(httpReq)

	// THEN — 422 (ErrValidationFailed → code 500010); the AppError's code surfaces
	// via NewError (not the generic default code) so the frontend can map it.
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(s.T(), strconv.Itoa(apperrors.ErrValidationFailed.Code), result["code"])
	s.pondService.AssertExpectations(s.T())
}
