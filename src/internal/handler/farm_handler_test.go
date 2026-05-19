//go:build cgo

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
	app := newTestApp()
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
		Status:   "maintenance",
		Summary:  dto.FarmDetailSummary{TotalPonds: 0, ActivePonds: 0},
		Ponds:    []dto.FarmDetailPondItem{},
	}
	s.farmService.On("Get", farmId, mock.AnythingOfType("*int")).Return(expectedResponse, nil)
	app := newTestApp()
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
	app := newTestApp()
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
		{Id: 2, ClientId: clientId, Name: "Delta Farm", Status: "maintenance", Ponds: []dto.FarmDetailPondItem{}},
	}
	s.farmService.On("GetHierarchy", clientId).Return(expectedList, nil)
	app := newTestApp()
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
	app := newTestApp()
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
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/hierarchy", s.farmHandler.GetFarmHierarchy)
	req := httptest.NewRequest("GET", "/api/v1/farm/hierarchy", nil)

	// WHEN — GET /api/v1/farm/hierarchy is sent
	resp, err := app.Test(req)

	// THEN — 500 (service returned an unknown error, mapped via ErrGeneric)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusInternalServerError, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmHierarchy_ClientIdNotFound() {
	// GIVEN — userLevel 1 and no clientId in context
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 1}))
	app.Get("/api/v1/farm/hierarchy", s.farmHandler.GetFarmHierarchy)
	req := httptest.NewRequest("GET", "/api/v1/farm/hierarchy", nil)

	// WHEN — GET /api/v1/farm/hierarchy is sent
	resp, err := app.Test(req)

	// THEN — 401 (ErrAuthTokenInvalid: non-super-admin missing clientId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnauthorized, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_Success() {
	// GIVEN — valid update body; super admin; service returns nil
	updateReq := dto.UpdateFarmRequest{Id: 1, Name: "Updated Farm"}
	username := "admin"
	s.farmService.On("Update", mock.Anything, updateReq).Return(nil)
	app := newTestApp()
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

	app := newTestApp()
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

	// THEN — 500 (service returned an unknown error, mapped via ErrGeneric)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusInternalServerError, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarm_InvalidId() {
	// GIVEN — invalid id "not-a-number" in path
	clientId := 1

	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/not-a-number", nil)

	// WHEN — GET /api/v1/farm/not-a-number is sent
	resp, err := app.Test(req)

	// THEN — 422 (ErrValidationFailed: id param failed strconv.Atoi)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_ServiceError() {
	// GIVEN — service returns error
	clientId := 1
	svcErr := errors.New("db error")
	s.farmService.On("GetList", clientId).Return((*dto.FarmListResponse)(nil), svcErr)

	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	// WHEN — GET /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 500 (service returned an unknown error, mapped via ErrGeneric)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusInternalServerError, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}

// --- AddFarm error paths ---

func (s *FarmHandlerTestSuite) TestAddFarm_InvalidBody() {
	// GIVEN — malformed JSON body
	app := newTestApp()
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

	// THEN — 400 (ErrInvalidRequestBody). validateAndParse short-circuits via the
	// sentinel.
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true, "expected error or non-success response")
}

func (s *FarmHandlerTestSuite) TestAddFarm_ValidationFailed() {
	// GIVEN — body with empty required name
	createReq := map[string]any{
		"clientId": 1,
		"name":     "", // required field empty
	}
	app := newTestApp()
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

	// THEN — 422 (ErrValidationFailed). validateAndParse short-circuits via the
	// sentinel.
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true, "expected error or non-success response")
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientAdmin_Allowed() {
	// GIVEN — client-admin (level 2) targeting their own client (1).
	// Should now be allowed (was super-admin-only before).
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	newFarm := &dto.FarmResponse{Id: 99, ClientId: 1, Name: "Test Farm"}
	s.farmService.On("Create", mock.Anything, *createReq, createReq.ClientId).Return(newFarm, nil)
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "clientAdmin",
		"clientId":  1,
		"userLevel": 2,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(req)

	// THEN — 200 + farm returned; service was called.
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientAdmin_WrongClient_Denied() {
	// GIVEN — client-admin for client 1, but request targets client 2.
	// validateClientAccess returns a sentinel error that short-circuits the
	// handler with a 403 (ErrAuthPermissionDenied).
	createReq := &dto.CreateFarmRequest{
		ClientId: 2,
		Name:     "Test Farm",
	}
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "clientAdmin",
		"clientId":  1,
		"userLevel": 2,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN
	resp, err := app.Test(req)

	// THEN — 403 (ErrAuthPermissionDenied: cross-client check short-circuits).
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true)
}

func (s *FarmHandlerTestSuite) TestAddFarm_NotSuperAdmin() {
	// GIVEN — valid body; userLevel 1 (regular user — below client-admin too)
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	app := newTestApp()
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

	// THEN — 403 (ErrAuthPermissionDenied)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	require.NotNil(s.T(), result["error"], "expected error in response when not super admin")
	errObj, ok := result["error"].(map[string]any)
	require.True(s.T(), ok)
	require.NotNil(s.T(), errObj["code"])
	assert.Equal(s.T(), "500024", errObj["code"])
}

func (s *FarmHandlerTestSuite) TestAddFarm_MissingUsername() {
	// GIVEN — valid body; no username and no userLevel in context
	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId": 1,
		// no username, no userLevel
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 403 (ErrAuthPermissionDenied: IsClientAdminOrAbove fails on missing userLevel)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientAccessDenied() {
	// GIVEN — request body has clientId 2. The userContextFromRequest middleware
	// in this suite copies the fasthttp context (not req.Context()) into
	// UserContext, so userLevel/clientId are NOT propagated; IsClientAdminOrAbove
	// fails on the missing level and writes a 403 (ErrAuthPermissionDenied).
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	createReq := &dto.CreateFarmRequest{
		ClientId: 2, // request for client 2
		Name:     "Test Farm",
	}

	app := newTestApp()
	app.Use(userContextFromRequest)
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withUserContext("user", 1, 3))

	// WHEN — POST with clientId 2 is sent
	resp, err := app.Test(req)

	// THEN — 403 (ErrAuthPermissionDenied: IsClientAdminOrAbove fails on missing userLevel)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["result"] != true)
}

func (s *FarmHandlerTestSuite) TestAddFarm_ClientIdNotFound() {
	// GIVEN — non–super admin with no clientId in context
	s.farmService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.FarmResponse)(nil), errors.New("")).Maybe()

	createReq := &dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	app := newTestApp()
	app.Use(userContextFromRequest)
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Not super admin and no clientId -> permission denied (500024) before client access check
	req = req.WithContext(withUserContext("user", 0, 1))

	// WHEN — POST /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 403 (ErrAuthPermissionDenied: user is below client-admin)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
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
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1,
		// no clientId -> canAccess false for non-super-admin
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	// WHEN — GET /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — 401 (ErrAuthTokenInvalid: non-super-admin missing clientId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnauthorized, resp.StatusCode)
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

	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	// WHEN — GET /api/v1/farm/1 is sent
	resp, err := app.Test(req)

	// THEN — 500 (service returned an unknown error, mapped via ErrGeneric)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusInternalServerError, resp.StatusCode)
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

	app := newTestApp()
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
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm?clientId=invalid", nil)

	// WHEN — GET with invalid clientId is sent
	resp, err := app.Test(req)

	// THEN — 422 (ErrValidationFailed: clientId param failed strconv.Atoi)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnprocessableEntity, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_ClientIdNotFound() {
	// GIVEN — userLevel 1 and no clientId
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	// WHEN — GET /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 401 (ErrAuthTokenInvalid: non-super-admin missing clientId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnauthorized, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *FarmHandlerTestSuite) TestGetFarmList_IsSuperAdminError() {
	// GIVEN — empty locals (no userLevel)
	app := newTestApp()
	app.Use(setLocalsMiddleware(map[string]any{}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)
	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	// WHEN — GET /api/v1/farm is sent
	resp, err := app.Test(req)

	// THEN — 401 (ErrAuthTokenInvalid: IsSuperAdmin failed on missing userLevel)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusUnauthorized, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

// --- UpdateFarm error paths ---

func (s *FarmHandlerTestSuite) TestUpdateFarm_InvalidBody() {
	// GIVEN — non-JSON body
	app := newTestApp()
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

	// THEN — 400 (ErrInvalidRequestBody). validateAndParse short-circuits via the
	// sentinel.
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["code"] != nil)
}

func (s *FarmHandlerTestSuite) TestUpdateFarm_NotSuperAdmin() {
	// GIVEN — userLevel 1 (not super admin)
	app := newTestApp()
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

	// THEN — 403 (ErrAuthPermissionDenied: non-super-admin)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
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

	app := newTestApp()
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

	app := newTestApp()
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

	// THEN — 500 (service returned an unknown error, mapped via ErrGeneric)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusInternalServerError, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.farmService.AssertExpectations(s.T())
}
