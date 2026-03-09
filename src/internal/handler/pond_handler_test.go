package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
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
		Names:  []string{"Pond 1", "Pond 2"},
	}
	username := "admin"
	s.pondService.On("CreatePonds", mock.Anything, *createReq, username).Return(nil)
	app := fiber.New()
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
		Names:  []string{"Pond 1"},
	}
	app := fiber.New()
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

	// THEN — 200 with error body; code 500024 (ErrAuthPermissionDenied)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
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
		Names:  []string{"Pond 1"},
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{}))
	app.Post("/api/v1/pond", s.pondHandler.AddPonds)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond is sent
	resp, err := app.Test(req)

	// THEN — 200 with error body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
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

	s.pondService.On("Get", pondId).Return(expectedResponse, nil)

	app := fiber.New()
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

	s.pondService.On("GetList", farmId).Return(expectedResponse, nil)

	app := fiber.New()
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
	app := fiber.New()
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
	body := []byte(`{"fishType":"nil","amount":100,"pricePerUnit":10.5,"activityDate":"2024-01-15"}`)
	app := s.fillPondApp("user")
	req := httptest.NewRequest("POST", "/api/v1/pond/abc/fill", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/abc/fill is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fiber.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"], "expected error for invalid pond ID")
	errObj, ok := result["error"].(map[string]any)
	// THEN — error with code 500010
	require.True(s.T(), ok)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestFillPond_MissingUsername_ReturnsAuthError() {
	// GIVEN — valid body; no username in context
	s.pondService.On("FillPond", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*dto.PondFillResponse)(nil), errors.New("")).Maybe()
	body := []byte(`{"fishType":"nil","amount":100,"pricePerUnit":10.5,"activityDate":"2024-01-15"}`)
	app := s.fillPondApp("")
	req := httptest.NewRequest("POST", "/api/v1/pond/1/fill", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/1/fill is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fiber.StatusOK, resp.StatusCode)

	// THEN — auth error with code 500022
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
	body := []byte(`{"fishType":"nil","amount":100,"pricePerUnit":10.5,"activityDate":"2024-01-15"}`)
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
	app := fiber.New()
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
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"activityDate":"2024-06-01"}`)
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
	require.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — error with code 500010
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"], "expected error for invalid pond ID")
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestMovePond_MissingUsername_ReturnsAuthError() {
	// GIVEN — valid body; no username in context
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"activityDate":"2024-06-01"}`)
	app := s.movePondApp("")
	req := httptest.NewRequest("POST", "/api/v1/pond/1/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	// WHEN — POST /api/v1/pond/1/move is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — auth error with code 500022
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
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"activityDate":"2024-06-01"}`)
	req := httptest.NewRequest("POST", "/api/v1/pond/999/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond/999/move is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — response code 500070
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
	body := []byte(`{"toPondId":2,"fishType":"nil","amount":50,"activityDate":"2024-06-01"}`)
	req := httptest.NewRequest("POST", "/api/v1/pond/1/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/pond/1/move is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — response code 500074
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), "500074", result["code"])
	s.pondService.AssertExpectations(s.T())
}

// updatePondApp returns a Fiber app with PUT /pond/:id and optional username in context.
func (s *PondHandlerTestSuite) updatePondApp(username string) *fiber.App {
	app := fiber.New()
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
	}, username).Return(nil)
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
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — error with code 500010
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"])
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestUpdatePond_MissingUsername() {
	// GIVEN — valid body; no username in context
	body := []byte(`{"name":"Pond"}`)
	app := s.updatePondApp("")
	req := httptest.NewRequest("PUT", "/api/v1/pond/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/pond/1 is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — error with code 500022
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500022", errObj["code"])
}

func (s *PondHandlerTestSuite) TestUpdatePond_ServiceError() {
	// GIVEN — pond 999; service returns ErrPondNotFound
	username := "user"
	s.pondService.On("Update", mock.Anything, mock.AnythingOfType("dto.UpdatePondRequest"), username).Return(apperrors.ErrPondNotFound)
	app := s.updatePondApp(username)
	req := httptest.NewRequest("PUT", "/api/v1/pond/999", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/pond/999 is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — response code 500070
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), "500070", result["code"])
	s.pondService.AssertExpectations(s.T())
}

// deletePondApp returns a Fiber app with DELETE /pond/:id and optional username in context.
func (s *PondHandlerTestSuite) deletePondApp(username string) *fiber.App {
	app := fiber.New()
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
	s.pondService.On("Delete", pondId, username).Return(nil)
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
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — error with code 500010
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	require.NotNil(s.T(), result["error"])
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500010", errObj["code"])
}

func (s *PondHandlerTestSuite) TestDeletePond_MissingUsername() {
	// GIVEN — no username in context
	app := s.deletePondApp("")
	req := httptest.NewRequest("DELETE", "/api/v1/pond/1", nil)

	// WHEN — DELETE /api/v1/pond/1 is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — error with code 500022
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	errObj := result["error"].(map[string]any)
	assert.Equal(s.T(), "500022", errObj["code"])
}

func (s *PondHandlerTestSuite) TestDeletePond_ServiceError() {
	// GIVEN — pond 1; service returns ErrGeneric
	username := "user"
	s.pondService.On("Delete", 1, username).Return(apperrors.ErrGeneric)
	app := s.deletePondApp(username)
	req := httptest.NewRequest("DELETE", "/api/v1/pond/1", nil)

	// WHEN — DELETE /api/v1/pond/1 is sent
	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	// THEN — response has message
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
			PricePerUnit: decimal.NewFromFloat(10.5),
			ActivityDate: "2024-01-15",
		}
		// WHEN — ValidateStruct is called
		err := utils.ValidateStruct(req)
		// THEN — validation passes
		require.NoError(t, err)
	})
}
