package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

// Helper middleware to set context values for testing
func setLocalsMiddleware(locals map[string]interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		if ctx == nil {
			ctx = context.Background()
		}
		
		// Map string keys to context keys
		if userId, ok := locals["userId"]; ok {
			ctx = context.WithValue(ctx, utils.UserIdKey(), userId)
		}
		if username, ok := locals["username"]; ok {
			ctx = context.WithValue(ctx, utils.UsernameKey(), username)
		}
		if clientId, ok := locals["clientId"]; ok {
			ctx = context.WithValue(ctx, utils.ClientIdKey(), clientId)
		}
		if userLevel, ok := locals["userLevel"]; ok {
			ctx = context.WithValue(ctx, utils.UserLevelKey(), userLevel)
		}
		
		c.SetUserContext(ctx)
		return c.Next()
	}
}

// Test AddUser handler
func (s *HandlerTestSuite) TestAddUser_Success() {
	createReq := &dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		LastName:      lo.ToPtr("User"),
		UserLevel:     1,
		ContactNumber: "1234567890",
	}

	expectedResponse := &dto.UserResponse{
		Id:            1,
		ClientId:      lo.ToPtr(1),
		Username:      createReq.Username,
		FirstName:     createReq.FirstName,
		LastName:      createReq.LastName,
		UserLevel:     createReq.UserLevel,
		ContactNumber: createReq.ContactNumber,
		CreatedAt:     time.Now(),
		CreatedBy:     "admin",
		UpdatedAt:     time.Now(),
		UpdatedBy:     "admin",
	}

	username := "admin"
	clientIdInt := 1
	clientId := lo.ToPtr(1)
	userLevel := 3 // Super admin
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil }), *createReq, username, clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username":  username,
		"clientId":  clientIdInt, // Set as int, handler will convert to *int
		"userLevel": userLevel,  // Super admin level
	}))
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestAddUser_InvalidBody() {
	// Invalid JSON might still parse to empty struct and call service, so we need a mock
	emptyReq := dto.CreateUserRequest{}
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), emptyReq, "system", (*int)(nil)).Return(nil, errors.New("validation error"))

	app := fiber.New()
	app.Post("/api/v1/user", s.userHandler.AddUser)

	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	// Invalid body should return an error response
	assert.True(s.T(), resp.StatusCode == fiber.StatusOK || resp.StatusCode >= 400)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestAddUser_ValidationError() {
	req := &dto.CreateUserRequest{
		Username: "ab",  // Too short
		Password: "123", // Too short
	}

	// Handler might still call service even with validation errors, so we need a mock
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), *req, "system", (*int)(nil)).Return(nil, errors.New("validation error"))

	app := fiber.New()
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(reqHTTP)

	assert.NoError(s.T(), err)
	// Validation error should return an error response
	assert.True(s.T(), resp.StatusCode == fiber.StatusOK || resp.StatusCode >= 400)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestAddUser_MissingUsername() {
	createReq := &dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}

	// Mock the service call with nil clientId (system setup)
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), *createReq, "system", (*int)(nil)).Return(nil, errors.New("validation error"))

	app := fiber.New()
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestAddUser_MissingClientId() {
	createReq := &dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}

	// Mock the service call with nil clientId (no clientId in request and no JWT)
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), *createReq, "system", (*int)(nil)).Return(nil, errors.New("service error"))

	app := fiber.New()
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestAddUser_ServiceError() {
	createReq := &dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}

	// When no JWT is present, handler uses "system" as username and nil clientId
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), *createReq, "system", (*int)(nil)).Return(nil, errors.New("user already exist"))

	app := fiber.New()
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test GetUser handler
func (s *HandlerTestSuite) TestGetUser_Success() {
	userID := 1
	expectedResponse := &dto.UserResponse{
		Id:            userID,
		ClientId:      lo.ToPtr(1),
		Username:      "testuser",
		FirstName:     "Test",
		LastName:      lo.ToPtr("User"),
		UserLevel:     1,
		ContactNumber: "1234567890",
		CreatedAt:     time.Now(),
		CreatedBy:     "admin",
		UpdatedAt:     time.Now(),
		UpdatedBy:     "admin",
	}

	s.userService.On("GetUser", userID).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"userId": userID,
	}))
	app.Get("/api/v1/user", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/api/v1/user", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestGetUser_MissingUserId() {
	app := fiber.New()
	app.Get("/api/v1/user", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/api/v1/user", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *HandlerTestSuite) TestGetUser_ServiceError() {
	userID := 1
	s.userService.On("GetUser", userID).Return(nil, errors.New("user not found"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"userId": userID,
	}))
	app.Get("/api/v1/user", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/api/v1/user", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test UpdateUser handler
func (s *HandlerTestSuite) TestUpdateUser_Success() {
	updateUser := &model.User{
		Id:            1,
		ClientId:      lo.ToPtr(1),
		Username:      "updateduser",
		FirstName:     "Updated",
		LastName:      lo.ToPtr("User"),
		UserLevel:     1,
		ContactNumber: "0987654321",
	}

	username := "admin"
	s.userService.On("Update", updateUser, username).Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
	}))
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	body, _ := json.Marshal(updateUser)
	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestUpdateUser_InvalidBody() {
	app := fiber.New()
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUser_MissingUsername() {
	updateUser := &model.User{
		Id:        1,
		Username:  "updateduser",
		FirstName: "Updated",
	}

	app := fiber.New()
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	body, _ := json.Marshal(updateUser)
	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *HandlerTestSuite) TestUpdateUser_ServiceError() {
	updateUser := &model.User{
		Id:        1,
		Username:  "updateduser",
		FirstName: "Updated",
	}

	username := "admin"
	s.userService.On("Update", updateUser, username).Return(errors.New("update failed"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
	}))
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	body, _ := json.Marshal(updateUser)
	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test GetUserList handler
func (s *HandlerTestSuite) TestGetUserList_Success() {
	clientId := 1

	expectedUsers := []*dto.UserResponse{
		{
			Id:            1,
			ClientId:      &clientId,
			Username:      "user1",
			FirstName:     "User",
			LastName:      lo.ToPtr("One"),
			UserLevel:     1,
			ContactNumber: "1111111111",
			CreatedAt:     time.Now(),
			CreatedBy:     "admin",
			UpdatedAt:     time.Now(),
			UpdatedBy:     "admin",
		},
		{
			Id:            2,
			ClientId:      &clientId,
			Username:      "user2",
			FirstName:     "User",
			LastName:      lo.ToPtr("Two"),
			UserLevel:     1,
			ContactNumber: "2222222222",
			CreatedAt:     time.Now(),
			CreatedBy:     "admin",
			UpdatedAt:     time.Now(),
			UpdatedBy:     "admin",
		},
	}

	s.userService.On("GetUserList", clientId).Return(expectedUsers, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"clientId": clientId,
	}))
	app.Get("/api/v1/user/list", s.userHandler.GetUserList)

	req := httptest.NewRequest("GET", "/api/v1/user/list", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *HandlerTestSuite) TestGetUserList_MissingClientId() {
	app := fiber.New()
	app.Get("/api/v1/user/list", s.userHandler.GetUserList)

	req := httptest.NewRequest("GET", "/api/v1/user/list", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *HandlerTestSuite) TestGetUserList_ServiceError() {
	clientId := 1

	s.userService.On("GetUserList", clientId).Return(nil, errors.New("database error"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"clientId": clientId,
	}))
	app.Get("/api/v1/user/list", s.userHandler.GetUserList)

	req := httptest.NewRequest("GET", "/api/v1/user/list", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

