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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
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
	// GIVEN — valid CreateFarmRequest; super admin context; service returns success
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	expectedResponse := &dto.FarmResponse{
		Id:       1,
		ClientId: 1,
		Name:     createReq.Name,
		Status:   "maintenance",
	}
	username := "admin"
	clientId := 1
	s.farmService.On("Create", mock.Anything, *createReq, clientId).Return(expectedResponse, nil)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  clientId,
		"userLevel": 3,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarm_Success() {
	// GIVEN — farm id 1; service returns detail
	farmId := 1
	clientId := 1
	expectedResponse := &dto.FarmDetailResponse{
		Id:       farmId,
		ClientId: clientId,
		Name:     "Test Farm",
		Status:   "active",
		Summary:  dto.FarmDetailSummary{TotalPonds: 0, ActivePonds: 0},
		Ponds:    []dto.FarmDetailPondItem{},
	}
	s.farmService.On("Get", farmId, mock.AnythingOfType("*int")).Return(expectedResponse, nil)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)
	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	// WHEN — GET /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmList_Success() {
	// GIVEN — clientId in context; service returns list
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

	// WHEN — GET /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmHierarchy_Success() {
	// GIVEN — clientId in context; service returns hierarchy
	clientId := 1
	expectedList := []*dto.FarmHierarchyItem{
		{Id: 1, ClientId: clientId, Name: "River Farm", Status: "active", Ponds: []dto.FarmDetailPondItem{{Id: 1, Name: "Pond A1", Status: "active"}}},
		{Id: 2, ClientId: clientId, Name: "Delta Farm", Status: "active", Ponds: []dto.FarmDetailPondItem{}},
	}
	s.farmService.On("GetHierarchy", clientId).Return(expectedList, nil)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/hierarchy", s.farmHandler.GetFarmHierarchy)
	req := httptest.NewRequest("GET", "/api/v1/farm/hierarchy", nil)

	// WHEN — GET /api/v1/farm/hierarchy is sent
	resp, err := app.Test(req)

	// THEN — 200, result true, data present
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	assert.NotNil(s.T(), result["data"])
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmHierarchy_Success_SuperAdminWithClientId() {
	// GIVEN — super admin with clientId query param; service returns hierarchy
	clientId := 2
	expectedList := []*dto.FarmHierarchyItem{
		{Id: 1, ClientId: clientId, Name: "Farm X", Status: "active", Ponds: []dto.FarmDetailPondItem{}},
	}
	s.farmService.On("GetHierarchy", clientId).Return(expectedList, nil)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 3}))
	app.Get("/api/v1/farm/hierarchy", s.farmHandler.GetFarmHierarchy)
	req := httptest.NewRequest("GET", "/api/v1/farm/hierarchy?clientId=2", nil)

	// WHEN — GET with clientId=2 is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmHierarchy_ServiceError() {
	// GIVEN — service returns error
	clientId := 1
	svcErr := errors.New("db error")
	s.farmService.On("GetHierarchy", clientId).Return(([]*dto.FarmHierarchyItem)(nil), svcErr)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/hierarchy", s.farmHandler.GetFarmHierarchy)
	req := httptest.NewRequest("GET", "/api/v1/farm/hierarchy", nil)

	// WHEN — GET /api/v1/farm/hierarchy is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmHierarchy_ClientIdNotFound() {
	// GIVEN — userLevel 1 and no clientId in context
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 1}))
	app.Get("/api/v1/farm/hierarchy", s.farmHandler.GetFarmHierarchy)
	req := httptest.NewRequest("GET", "/api/v1/farm/hierarchy", nil)

	// WHEN — GET /api/v1/farm/hierarchy is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_Success() {
	// GIVEN — valid update body; super admin; service returns nil
	updateReq := dto.UpdateFarmRequest{Id: 1, Name: "Updated Farm"}
	username := "admin"
	s.farmService.On("Update", mock.Anything, updateReq).Return(nil)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  1,
		"userLevel": 3,
	}))
	app.Put("/api/v1/farm/:id", s.farmHandler.UpdateFarm)
	body, _ := json.Marshal(dto.UpdateFarmBody{Name: "Updated Farm"})
	req := httptest.NewRequest("PUT", "/api/v1/farm/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestAddFarm_ServiceError() {
	// GIVEN — valid body; service returns error
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	username := "admin"
	clientId := 1
	svcErr := errors.New("farm already exists")
	s.farmService.On("Create", mock.Anything, *createReq, clientId).Return((*dto.FarmResponse)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  clientId,
		"userLevel": 3, // super admin only
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarm_InvalidId() {
	// GIVEN — invalid id "not-a-number" in path
	clientId := 1

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/not-a-number", nil)

	// WHEN — GET /api/v1/farm/not-a-number is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_ServiceError() {
	// GIVEN — service returns error
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

	// WHEN — GET /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

// --- AddFarm error paths ---

func (s *FarmHandlerTestSuite) TestAddFarm_InvalidBody() {
	// GIVEN — malformed JSON body
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 3, // super admin only
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST with invalid JSON is sent
	resp, err := app.Test(req)

	// THEN — error or non-success response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true, "expected error or non-success response")
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		code := errObj["code"]
		assert.True(s.T(), code == "500011" || code == "500022", "expected invalid body or auth error, got %v", code)
	}
}

func (s *FarmHandlerTestSuite) TestAddFarm_ValidationFailed() {
	// GIVEN — body with empty required name
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	createReq := map[string]any{
		"clientId": 1,
		"name":     "", // required field empty
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 3, // super admin only
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST with invalid body is sent
	resp, err := app.Test(req)

	// THEN — error or non-success response
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

func (s *FarmHandlerTestSuite) TestAddFarm_NotSuperAdmin() {
	// GIVEN — valid body; userLevel 1 (not super admin)
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 1, // not super admin
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — error 500024 (permission denied)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	require.NotNil(s.T(), result["error"], "expected error in response when not super admin")
	errObj, ok := result["error"].(map[string]any)
	require.True(s.T(), ok)
	require.NotNil(s.T(), errObj["code"])
	assert.Equal(s.T(), "500024", errObj["code"])
}

func (s *FarmHandlerTestSuite) TestAddFarm_MissingUsername() {
	// GIVEN — valid body; no username in context
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

	// WHEN — POST /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientAccessDenied() {
	// GIVEN — super admin with client 1; request body has clientId 2
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
	req = req.WithContext(withUserContext("user", 1, 3))

	// WHEN — POST with clientId 2 is sent
	resp, err := app.Test(req)

	// THEN — error (auth or permission)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true)
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		code := errObj["code"]
		assert.True(s.T(), code == "500022" || code == "500024", "expected auth or permission error, got %v", code)
	}
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientIdNotFound() {
	// GIVEN — non–super admin with no clientId in context
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
	// Not super admin and no clientId -> permission denied (500024) before client access check
	req = req.WithContext(withUserContext("user", 0, 1))

	// WHEN — POST /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — error 500024
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		assert.Equal(s.T(), "500024", errObj["code"])
	}
}

// --- GetFarm error paths ---

func (s *FarmHandlerTestSuite) TestGetFarm_ClientIdNotFound() {
	// GIVEN — userLevel 1 and no clientId
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1,
		// no clientId -> canAccess false for non-super-admin
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	// WHEN — GET /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarm_ServiceError() {
	// GIVEN — service returns error
	farmId := 1
	clientId := 1
	svcErr := errors.New("not found")
	s.farmService.On("Get", farmId, mock.AnythingOfType("*int")).Return((*dto.FarmDetailResponse)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	// WHEN — GET /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

// --- GetFarmList edge cases ---

func (s *FarmHandlerTestSuite) TestGetFarmList_SuperAdminWithClientIdQuery() {
	// GIVEN — super admin; clientId=2 in query; service returns list
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

	// WHEN — GET /api/v1/farm?clientId=2 is sent
	resp, err := app.Test(req)

	// THEN — 200 and result true
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmList_SuperAdminInvalidClientId() {
	// GIVEN — super admin; clientId=invalid in query
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm?clientId=invalid", nil)

	// WHEN — GET with invalid clientId is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_ClientIdNotFound() {
	// GIVEN — userLevel 1 and no clientId
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	// WHEN — GET /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_IsSuperAdminError() {
	// GIVEN — empty locals (no userLevel)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)
	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	// WHEN — GET /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

// --- UpdateFarm error paths ---

func (s *FarmHandlerTestSuite) TestUpdateFarm_InvalidBody() {
	// GIVEN — non-JSON body
	s.farmService.On("Update", mock.Anything, mock.Anything).Return(errors.New("")).Maybe()
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 3, // super admin only
	}))
	app.Put("/api/v1/farm/:id", s.farmHandler.UpdateFarm)

	req := httptest.NewRequest("PUT", "/api/v1/farm/1", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT with invalid body is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["code"] != nil)
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_NotSuperAdmin() {
	// GIVEN — userLevel 1 (not super admin)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 1, // not super admin
	}))
	app.Put("/api/v1/farm/:id", s.farmHandler.UpdateFarm)

	body, _ := json.Marshal(dto.UpdateFarmBody{Name: "Updated"})
	req := httptest.NewRequest("PUT", "/api/v1/farm/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — error 500024
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	require.NotNil(s.T(), result["error"])
	errObj, ok := result["error"].(map[string]any)
	require.True(s.T(), ok && errObj["code"] != nil)
	assert.Equal(s.T(), "500024", errObj["code"])
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_NoUsernameInContext() {
	// GIVEN — super admin with no username; service returns nil
	updateReq := dto.UpdateFarmRequest{Id: 1, Name: "Updated"}
	s.farmService.On("Update", mock.Anything, updateReq).Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  1,
		"userLevel": 3, // super admin; no username in context
	}))
	app.Put("/api/v1/farm/:id", s.farmHandler.UpdateFarm)

	body, _ := json.Marshal(dto.UpdateFarmBody{Name: "Updated"})
	req := httptest.NewRequest("PUT", "/api/v1/farm/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — success
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"] == true)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_ServiceError() {
	// GIVEN — valid body; service returns error
	updateReq := dto.UpdateFarmRequest{Id: 1, Name: "Updated"}
	username := "admin"
	svcErr := errors.New("update failed")
	s.farmService.On("Update", mock.Anything, updateReq).Return(svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  1,
		"userLevel": 3, // super admin only
	}))
	app.Put("/api/v1/farm/:id", s.farmHandler.UpdateFarm)

	body, _ := json.Marshal(dto.UpdateFarmBody{Name: "Updated"})
	req := httptest.NewRequest("PUT", "/api/v1/farm/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}
