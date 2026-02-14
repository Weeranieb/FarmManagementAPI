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
	createReq := &dto.CreateClientRequest{
		Name:          "Acme Corp",
		OwnerName:     "John",
		ContactNumber: "0812345678",
	}
	expectedResponse := &dto.ClientResponse{
		Id:            1,
		Name:          createReq.Name,
		OwnerName:     createReq.OwnerName,
		ContactNumber: createReq.ContactNumber,
		IsActive:      true,
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestAddClient_InvalidBody() {
	s.clientService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return((*dto.ClientResponse)(nil), errors.New("")).Maybe()
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "admin"}))
	app.Post("/api/v1/client", s.clientHandler.AddClient)

	req := httptest.NewRequest("POST", "/api/v1/client", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["message"] != nil, "expected error or message in response")
}

func (s *ClientHandlerTestSuite) TestAddClient_ValidationFailed() {
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["error"] != nil || result["message"] != nil, "expected error or message in response")
}

func (s *ClientHandlerTestSuite) TestAddClient_MissingUsername() {
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestAddClient_ServiceError() {
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"], "service error should return message")
	s.clientService.AssertExpectations(s.T())
}

// --- GetClient ---

func (s *ClientHandlerTestSuite) TestGetClient_Success() {
	clientId := 1
	expectedResponse := &dto.ClientResponse{
		Id:   1,
		Name: "Acme Corp",
	}
	s.clientService.On("Get", 1).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  clientId,
		"userLevel": 1,
	}))
	app.Get("/api/v1/client/:id", s.clientHandler.GetClient)

	req := httptest.NewRequest("GET", "/api/v1/client/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestGetClient_InvalidId() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Get("/api/v1/client/:id", s.clientHandler.GetClient)

	req := httptest.NewRequest("GET", "/api/v1/client/not-a-number", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestGetClient_AccessDenied() {
	app := fiber.New()
	app.Use(userContextFromRequest)
	app.Get("/api/v1/client/:id", s.clientHandler.GetClient)

	req := httptest.NewRequest("GET", "/api/v1/client/2", nil)
	// User belongs to client 1, requesting client 2
	req = req.WithContext(withUserContext("user", 1, 1))

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestGetClient_ServiceError() {
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"], "service error should return message")
	s.clientService.AssertExpectations(s.T())
}

// --- GetClientList ---

func (s *ClientHandlerTestSuite) TestGetClientList_Success() {
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	assert.NotNil(s.T(), result["data"])
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestGetClientList_NotSuperAdmin() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 1, // normal user
	}))
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		assert.Equal(s.T(), "500024", errObj["code"]) // ErrAuthPermissionDenied
	}
}

func (s *ClientHandlerTestSuite) TestGetClientList_IsSuperAdminError() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{})) // no userLevel
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestGetClientList_ServiceError() {
	svcErr := errors.New("db error")
	s.clientService.On("GetClientDropdown").Return(([]*dto.DropdownItem)(nil), svcErr)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 3}))
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"], "service error should return message")
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestGetClientList_EmptyDropdown() {
	expectedDropdown := []*dto.DropdownItem{}
	s.clientService.On("GetClientDropdown").Return(expectedDropdown, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 3}))
	app.Get("/api/v1/client/list", s.clientHandler.GetClientList)

	req := httptest.NewRequest("GET", "/api/v1/client/list", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(s.T(), result["result"].(bool))
	s.clientService.AssertExpectations(s.T())
}

// --- UpdateClient ---

func (s *ClientHandlerTestSuite) TestUpdateClient_Success() {
	updateReq := &model.Client{
		Id:            1,
		Name:          "Updated Name",
		OwnerName:     "Jane",
		ContactNumber: "0898765432",
		IsActive:      true,
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.clientService.AssertExpectations(s.T())
}

func (s *ClientHandlerTestSuite) TestUpdateClient_InvalidBody() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"clientId":  1,
		"userLevel": 1,
	}))
	app.Put("/api/v1/client", s.clientHandler.UpdateClient)

	req := httptest.NewRequest("PUT", "/api/v1/client", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestUpdateClient_AccessDenied() {
	updateReq := &model.Client{
		Id:    2, // update client 2
		Name:  "Updated",
		IsActive: true,
	}
	app := fiber.New()
	app.Use(userContextFromRequest)
	app.Put("/api/v1/client", s.clientHandler.UpdateClient)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/client", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withUserContext("user", 1, 1)) // user client 1, update client 2

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestUpdateClient_MissingUsername() {
	updateReq := &model.Client{
		Id:    1,
		Name:  "Updated",
		IsActive: true,
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
}

func (s *ClientHandlerTestSuite) TestUpdateClient_ServiceError() {
	updateReq := &model.Client{
		Id:    1,
		Name:  "Updated",
		IsActive: true,
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

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(s.T(), result["message"], "service error should return message")
	s.clientService.AssertExpectations(s.T())
}
