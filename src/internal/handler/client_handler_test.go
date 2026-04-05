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
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
)

type ClientHandlerTestSuite struct {
	suite.Suite
	clientService *mocks.MockClientService
	clientHandler ClientHandler
}

func (s *ClientHandlerTestSuite) SetupTest() {
	s.clientService = mocks.NewMockClientService(s.T())
	s.clientHandler = NewClientHandler(s.clientService)
}

func (s *ClientHandlerTestSuite) TearDownTest() {
	s.clientService.ExpectedCalls = nil
}

func TestClientHandlerSuite(t *testing.T) {
	suite.Run(t, new(ClientHandlerTestSuite))
}

// --- AddClient ---

func (s *ClientHandlerTestSuite) TestAddClient_Success() {
	// GIVEN — valid CreateClientRequest; service returns success
	createReq := &dto.CreateClientRequest{
		Name:          "Acme Corp",
		OwnerName:     "John",
		ContactNumber: "0812345678",
	}
	expectedResponse := &dto.ClientResponse{
		Id:                      1,
		Name:                    createReq.Name,
		OwnerName:               createReq.OwnerName,
		ContactNumber:           createReq.ContactNumber,
		IsActive:                true,
		IsTouristFishingEnabled: false,
	}
	username := "admin"

	s.clientService.On("Create", mock.Anything, *createReq, username).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username": username,
	}))
	app.Post("/api/v1/client", s.clientHandler.AddClient)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/client is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestAddClient_InvalidBody() {
	// GIVEN — malformed JSON body
	s.clientService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.ClientResponse)(nil), errors.New("")).Maybe()
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "admin"}))
	app.Post("/api/v1/client", s.clientHandler.AddClient)

	req := httptest.NewRequest("POST", "/api/v1/client", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST with invalid JSON is sent
	resp, err := app.Test(req)

	// THEN — error or message in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["message"] != nil)
}

func (s *ClientHandlerTestSuite) TestAddClient_ValidationFailed() {
	// GIVEN — body with empty required name
	s.clientService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.ClientResponse)(nil), errors.New("")).Maybe()
	createReq := map[string]any{
		"name":          "",
		"ownerName":     "John",
		"contactNumber": "0812345678",
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "admin"}))
	app.Post("/api/v1/client", s.clientHandler.AddClient)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST with invalid body is sent
	resp, err := app.Test(req)

	// THEN — error or message in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["message"] != nil)
}

func (s *ClientHandlerTestSuite) TestAddClient_MissingUsername() {
	// GIVEN — valid body; no username in context
	createReq := &dto.CreateClientRequest{
		Name:          "Acme",
		OwnerName:     "John",
		ContactNumber: "0812345678",
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{}))
	app.Post("/api/v1/client", s.clientHandler.AddClient)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/client is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestAddClient_ServiceError() {
	// GIVEN — valid body; service returns error
	createReq := &dto.CreateClientRequest{
		Name:          "Acme",
		OwnerName:     "John",
		ContactNumber: "0812345678",
	}
	username := "admin"
	svcErr := errors.New("client already exists")
	s.clientService.On("Create", mock.Anything, *createReq, username).Return((*dto.ClientResponse)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": username}))
	app.Post("/api/v1/client", s.clientHandler.AddClient)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/client is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.clientService.AssertExpectations(s.T())
}

// --- GetClient ---

func (s *ClientHandlerTestSuite) TestGetClient_Success() {
	// GIVEN — clientId in context; service returns client
	clientId := 1
	expectedResponse := &dto.ClientResponse{
		Id:                      1,
		Name:                    "Acme Corp",
		IsTouristFishingEnabled: false,
	}
	s.clientService.On("Get", 1).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/client/:id", s.clientHandler.GetClient)

	req := httptest.NewRequest("GET", "/api/v1/client/1", nil)

	// WHEN — GET /api/v1/client/1 is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestGetClient_InvalidId() {
	// GIVEN — invalid id "not-a-number" in path
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Get("/api/v1/client/:id", s.clientHandler.GetClient)

	req := httptest.NewRequest("GET", "/api/v1/client/not-a-number", nil)

	// WHEN — GET /api/v1/client/not-a-number is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestGetClient_AccessDenied() {
	// GIVEN — user client 1; request for client 2
	app := fiber.New()
	app.Use(userContextFromRequest)
	app.Get("/api/v1/client/:id", s.clientHandler.GetClient)

	req := httptest.NewRequest("GET", "/api/v1/client/2", nil)
	// User belongs to client 1, requesting client 2
	req = req.WithContext(withUserContext("user", 1, 1))

	// WHEN — GET /api/v1/client/2 is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestGetClient_ServiceError() {
	// GIVEN — service returns error
	clientId := 1
	svcErr := errors.New("not found")
	s.clientService.On("Get", 1).Return((*dto.ClientResponse)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/client/:id", s.clientHandler.GetClient)

	req := httptest.NewRequest("GET", "/api/v1/client/1", nil)

	// WHEN — GET /api/v1/client/1 is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.clientService.AssertExpectations(s.T())
}

// --- GetClientList ---

func (s *ClientHandlerTestSuite) TestGetClientList_Success() {
	// GIVEN — super admin; service returns dropdown
	expectedDropdown := []*dto.DropdownItem{
		{Key: 1, Value: "Client A"},
		{Key: 2, Value: "Client B"},
	}
	s.clientService.On("GetClientDropdown").Return(expectedDropdown, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3, // super admin
	}))
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	// WHEN — GET /api/v1/client/list is sent
	resp, err := app.Test(req)

	// THEN — 200, result true, data present
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	assert.NotNil(s.T(), result["data"])
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestGetClientList_NotSuperAdmin() {
	// GIVEN — userLevel 1 (not super admin)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1, // normal user
	}))
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	// WHEN — GET /api/v1/client/list is sent
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

func (s *ClientHandlerTestSuite) TestGetClientList_IsSuperAdminError() {
	// GIVEN — no userLevel in context
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{})) // no userLevel
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	// WHEN — GET /api/v1/client/list is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestGetClientList_ServiceError() {
	// GIVEN — service returns error
	svcErr := errors.New("db error")
	s.clientService.On("GetClientDropdown").Return(([]*dto.DropdownItem)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 3}))
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	// WHEN — GET /api/v1/client/list is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestGetClientList_EmptyDropdown() {
	// GIVEN — service returns empty dropdown
	expectedDropdown := []*dto.DropdownItem{}
	s.clientService.On("GetClientDropdown").Return(expectedDropdown, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 3}))
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	// WHEN — GET /api/v1/client/list is sent
	resp, err := app.Test(req)

	// THEN — 200 and result true
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	s.clientService.AssertExpectations(s.T())
}

// --- UpdateClient ---

func (s *ClientHandlerTestSuite) TestUpdateClient_Success() {
	// GIVEN — valid update body; service returns nil
	isActive := true
	updateReq := dto.UpdateClientRequest{
		Id:            1,
		Name:          "Updated Name",
		OwnerName:     "Jane",
		ContactNumber: "0898765432",
		IsActive:      &isActive,
	}
	username := "admin"
	s.clientService.On("Update", mock.Anything, updateReq, username).Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Put("/api/v1/client", s.clientHandler.UpdateClient)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/client is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestUpdateClient_InvalidBody() {
	// GIVEN — non-JSON body
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Put("/api/v1/client", s.clientHandler.UpdateClient)

	req := httptest.NewRequest("PUT", "/api/v1/client", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT with invalid body is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestUpdateClient_AccessDenied() {
	// GIVEN — user client 1; update body for client 2
	isActive := true
	updateReq := dto.UpdateClientRequest{
		Id:       2, // update client 2
		Name:     "Updated",
		IsActive: &isActive,
	}
	app := fiber.New()
	app.Use(userContextFromRequest)
	app.Put("/api/v1/client", s.clientHandler.UpdateClient)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withUserContext("user", 1, 1))

	// WHEN — PUT with client 2 body is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestUpdateClient_MissingUsername() {
	// GIVEN — valid body; no username in context
	isActive := true
	updateReq := dto.UpdateClientRequest{
		Id:       1,
		Name:     "Updated",
		IsActive: &isActive,
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Put("/api/v1/client", s.clientHandler.UpdateClient)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/client is sent
	resp, err := app.Test(req)

	// THEN — error in response
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestUpdateClient_ServiceError() {
	// GIVEN — valid body; service returns error
	isActive := true
	updateReq := dto.UpdateClientRequest{
		Id:       1,
		Name:     "Updated",
		IsActive: &isActive,
	}
	username := "admin"
	svcErr := errors.New("client not found")
	s.clientService.On("Update", mock.Anything, updateReq, username).Return(svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Put("/api/v1/client", s.clientHandler.UpdateClient)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/client is sent
	resp, err := app.Test(req)

	// THEN — 200 with message in body
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"])
	s.clientService.AssertExpectations(s.T())
}
