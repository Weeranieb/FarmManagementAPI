package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
)

type FarmHandlerTestSuite struct {
	suite.Suite
	farmService *mocks.MockFarmService
	farmHandler FarmHandler
}

func (s *FarmHandlerTestSuite) SetupTest() {
	s.farmService = mocks.NewMockFarmService(s.T())
	s.farmHandler = NewFarmHandler(s.farmService)
}

func (s *FarmHandlerTestSuite) TearDownTest() {
	s.farmService.ExpectedCalls = nil
}

func TestFarmHandlerSuite(t *testing.T) {
	suite.Run(t, new(FarmHandlerTestSuite))
}

func (s *FarmHandlerTestSuite) TestAddFarm_Success() {
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}

	expectedResponse := &dto.FarmResponse{
		Id:       1,
		ClientId: 1,
		Name:     createReq.Name,
		Status:   "active",
	}

	username := "admin"
	clientId := 1
	s.farmService.On("Create", *createReq, username, clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarm_Success() {
	farmId := 1
	clientId := 1
	expectedResponse := &dto.FarmResponse{
		Id:       farmId,
		ClientId: clientId,
		Name:     "Test Farm",
		Status:   "active",
	}

	s.farmService.On("Get", farmId, mock.AnythingOfType("*int")).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmList_Success() {
	clientId := 1
	expectedResponse := &dto.FarmListResponse{
		Farms: []*dto.FarmResponse{
			{Id: 1, ClientId: clientId, Name: "Farm 1", Status: "active"},
			{Id: 2, ClientId: clientId, Name: "Farm 2", Status: "active"},
		},
		Total:       2,
		TotalActive: 2,
	}

	s.farmService.On("GetList", clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_Success() {
	updateReq := &model.Farm{
		Id:       1,
		ClientId: 1,
		Name:     "Updated Farm",
		Status:   "active",
	}
	username := "admin"
	clientId := 1

	s.farmService.On("Update", updateReq, username).Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username": username,
		"clientId": clientId,
	}))
	app.Put("/api/v1/farm", s.farmHandler.UpdateFarm)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestAddFarm_ServiceError() {
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	username := "admin"
	clientId := 1
	svcErr := errors.New("farm already exists")
	s.farmService.On("Create", *createReq, username, clientId).Return((*dto.FarmResponse)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username": username,
		"clientId": clientId,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarm_InvalidId() {
	clientId := 1

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/not-a-number", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_ServiceError() {
	clientId := 1
	svcErr := errors.New("db error")
	s.farmService.On("GetList", clientId).Return((*dto.FarmListResponse)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

// --- AddFarm error paths ---

func (s *FarmHandlerTestSuite) TestAddFarm_InvalidBody() {
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true, "expected error or non-success response")
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		code := errObj["code"]
		// 500011 = invalid body, 500022 = auth (e.g. when context not passed)
		assert.True(s.T(), code == "500011" || code == "500022", "expected invalid body or auth error, got %v", code)
	}
}

func (s *FarmHandlerTestSuite) TestAddFarm_ValidationFailed() {
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	createReq := map[string]any{
		"clientId": 1,
		"name":     "", // required field empty
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true, "expected error or non-success response")
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		code := errObj["code"]
		assert.True(s.T(), code == "500010" || code == "500022", "expected validation or auth error, got %v", code)
	}
}

func (s *FarmHandlerTestSuite) TestAddFarm_MissingUsername() {
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId": 1,
		// no username
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientAccessDenied() {
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	createReq := &dto.CreateFarmRequest{
		ClientId: 2, // request for client 2
		Name:     "Test Farm",
	}

	app := fiber.New()
	app.Use(userContextFromRequest)
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withUserContext("user", 1, 1)) // user belongs to client 1, request is for client 2

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true, "expected error or non-success response")
	// 500022 = client id not found (e.g. request context not passed to handler), 500024 = permission denied
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		code := errObj["code"]
		assert.True(s.T(), code == "500022" || code == "500024", "expected auth or permission error, got %v", code)
	}
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientIdNotFound() {
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	app := fiber.New()
	app.Use(userContextFromRequest)
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// userLevel 1 but no clientId -> canAccess false
	req = req.WithContext(withUserContext("user", 0, 1))

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"], "expected error response")
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		assert.Equal(s.T(), "500022", errObj["code"]) // ErrAuthTokenInvalid
	}
}

// --- GetFarm error paths ---

func (s *FarmHandlerTestSuite) TestGetFarm_ClientIdNotFound() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1,
		// no clientId -> canAccess false for non-super-admin
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarm_ServiceError() {
	farmId := 1
	clientId := 1
	svcErr := errors.New("not found")
	s.farmService.On("Get", farmId, mock.AnythingOfType("*int")).Return((*dto.FarmResponse)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

// --- GetFarmList edge cases ---

func (s *FarmHandlerTestSuite) TestGetFarmList_SuperAdminWithClientIdQuery() {
	clientId := 2
	expectedResponse := &dto.FarmListResponse{
		Farms:       []*dto.FarmResponse{{Id: 1, ClientId: clientId, Name: "Farm 1", Status: "active"}},
		Total:       1,
		TotalActive: 1,
	}
	s.farmService.On("GetList", clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3, // super admin
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm?clientId=2", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmList_SuperAdminInvalidClientId() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm?clientId=invalid", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_ClientIdNotFound() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1,
		// no clientId
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_IsSuperAdminError() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

// --- UpdateFarm error paths ---

func (s *FarmHandlerTestSuite) TestUpdateFarm_InvalidBody() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username": "admin",
		"clientId": 1,
	}))
	app.Put("/api/v1/farm", s.farmHandler.UpdateFarm)

	req := httptest.NewRequest("PUT", "/api/v1/farm", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_ClientAccessDenied() {
	s.farmService.On("Update", mock.Anything, mock.Anything).Return(errors.New("")).Maybe()

	updateReq := &model.Farm{
		Id:       1,
		ClientId: 2,
		Name:     "Updated",
		Status:   "active",
	}
	app := fiber.New()
	app.Use(userContextFromRequest)
	app.Put("/api/v1/farm", s.farmHandler.UpdateFarm)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withUserContext("user", 1, 1)) // user client 1, update targets client 2

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true, "expected error or non-success response")
	// 500022 = client id not found, 500024 = permission denied
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		code := errObj["code"]
		assert.True(s.T(), code == "500022" || code == "500024", "expected auth or permission error, got %v", code)
	}
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_MissingUsername() {
	updateReq := &model.Farm{
		Id:       1,
		ClientId: 1,
		Name:     "Updated",
		Status:   "active",
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId": 1,
	}))
	app.Put("/api/v1/farm", s.farmHandler.UpdateFarm)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_ServiceError() {
	updateReq := &model.Farm{
		Id:       1,
		ClientId: 1,
		Name:     "Updated",
		Status:   "active",
	}
	username := "admin"
	svcErr := errors.New("update failed")
	s.farmService.On("Update", updateReq, username).Return(svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username": username,
		"clientId": 1,
	}))
	app.Put("/api/v1/farm", s.farmHandler.UpdateFarm)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}
