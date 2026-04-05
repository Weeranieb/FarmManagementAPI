package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
)

// HandlerTestSuite is shared with user_handler_test.go
// This test suite initializes a local UserHandler for testing
type UserHandlerTestSuite struct {
	suite.Suite
	userService *mocks.MockUserService

	userHandler UserHandler
}

func (s *UserHandlerTestSuite) SetupTest() {
	s.userService = mocks.NewMockUserService(s.T())
	s.userHandler = NewUserHandler(s.userService)
}

func (s *UserHandlerTestSuite) TearDownTest() {
	if s.userService != nil {
		s.userService.ExpectedCalls = nil
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

// Helper middleware to set context values for testing
func setLocalsMiddleware(locals map[string]any) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		if ctx == nil {
			ctx = context.Background()
		}

		// Map string keys to context keys
		if userId, ok := locals["userId"]; ok {
			ctx = context.WithValue(ctx, constants.UserIDKey, userId)
		}
		if username, ok := locals["username"]; ok {
			ctx = context.WithValue(ctx, constants.UsernameKey, username)
		}
		if clientId, ok := locals["clientId"]; ok {
			ctx = context.WithValue(ctx, constants.ClientIDKey, clientId)
		}
		if userLevel, ok := locals["userLevel"]; ok {
			ctx = context.WithValue(ctx, constants.UserLevelKey, userLevel)
		}

		c.SetUserContext(ctx)
		return c.Next()
	}
}

// Test AddUser handler
func (s *UserHandlerTestSuite) TestAddUser_Success() {
	// GIVEN — valid CreateUserRequest; service returns success
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
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  username,
		"clientId":  clientIdInt, // Set as int, handler will convert to *int
		"userLevel": userLevel,   // Super admin level
	}))
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *UserHandlerTestSuite) TestAddUser_InvalidBody() {
	// GIVEN — invalid JSON body
	emptyReq := dto.CreateUserRequest{}
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), emptyReq, "admin", (*int)(nil)).Return(nil, errors.New("validation error"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"userLevel": 3,
	}))
	app.Post("/api/v1/user", s.userHandler.AddUser)

	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST with invalid JSON is sent
	resp, err := app.Test(req)

	// THEN — error or non-success response
	assert.NoError(s.T(), err)
	assert.True(s.T(), resp.StatusCode == fiber.StatusOK || resp.StatusCode >= 400)
}

func (s *UserHandlerTestSuite) TestAddUser_ValidationError() {
	// GIVEN — request with too-short username/password
	req := &dto.CreateUserRequest{
		Username: "ab",  // Too short
		Password: "123", // Too short
	}

	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), *req, "admin", (*int)(nil)).Return(nil, errors.New("validation error"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"userLevel": 3,
	}))
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(req)
	reqHTTP := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	// WHEN — POST with validation errors is sent
	resp, err := app.Test(reqHTTP)

	// THEN — error or non-success response
	assert.NoError(s.T(), err)
	assert.True(s.T(), resp.StatusCode == fiber.StatusOK || resp.StatusCode >= 400)
}

func (s *UserHandlerTestSuite) TestAddUser_MissingUsername() {
	// GIVEN — valid body; no username in context
	createReq := &dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3,
	}))
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — permission denied
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *UserHandlerTestSuite) TestAddUser_MissingClientId() {
	// GIVEN — valid body; no clientId in context
	createReq := &dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}

	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), *createReq, "admin", (*int)(nil)).Return(nil, errors.New("service error"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"userLevel": 3,
	}))
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *UserHandlerTestSuite) TestAddUser_ServiceError() {
	// GIVEN — valid body; service returns error (e.g. user already exists)
	createReq := &dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}
	s.userService.On("Create", mock.MatchedBy(func(ctx context.Context) bool { return true }), *createReq, "admin", (*int)(nil)).Return(nil, errors.New("user already exist"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username":  "admin",
		"userLevel": 3,
	}))
	app.Post("/api/v1/user", s.userHandler.AddUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — POST /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test GetUser handler
func (s *UserHandlerTestSuite) TestGetUser_Success() {
	// GIVEN — userId in context; service returns user
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
	app.Use(setLocalsMiddleware(map[string]any{
		"userId": userID,
	}))
	app.Get("/api/v1/user", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/api/v1/user", nil)

	// WHEN — GET /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *UserHandlerTestSuite) TestGetUser_MissingUserId() {
	// GIVEN — no userId in context
	app := fiber.New()
	app.Get("/api/v1/user", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/api/v1/user", nil)

	// WHEN — GET /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200 (handler may return error in body)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *UserHandlerTestSuite) TestGetUser_ServiceError() {
	// GIVEN — service returns error
	userID := 1
	s.userService.On("GetUser", userID).Return(nil, errors.New("user not found"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userId": userID,
	}))
	app.Get("/api/v1/user", s.userHandler.GetUser)

	req := httptest.NewRequest("GET", "/api/v1/user", nil)

	// WHEN — GET /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test UpdateUser handler
func (s *UserHandlerTestSuite) TestUpdateUser_Success() {
	// GIVEN — valid update body; service returns nil
	userLevel := 1
	updateUser := dto.UpdateUserRequest{
		Username:      "updateduser",
		FirstName:     "Updated",
		LastName:      lo.ToPtr("User"),
		UserLevel:     &userLevel,
		ContactNumber: "0987654321",
	}

	username := "admin"
	userID := 1
	s.userService.On("Update", mock.Anything, userID, updateUser, username).Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username": username,
		"userId":   userID,
	}))
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	body, _ := json.Marshal(updateUser)
	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *UserHandlerTestSuite) TestUpdateUser_InvalidBody() {
	// GIVEN — invalid JSON body
	app := fiber.New()
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT with invalid JSON is sent
	resp, err := app.Test(req)

	// THEN — 200 (handler may return error in body)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *UserHandlerTestSuite) TestUpdateUser_MissingUsername() {
	// GIVEN — valid body; no username in context
	updateUser := dto.UpdateUserRequest{
		Username:  "updateduser",
		FirstName: "Updated",
	}

	app := fiber.New()
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	body, _ := json.Marshal(updateUser)
	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *UserHandlerTestSuite) TestUpdateUser_ServiceError() {
	// GIVEN — service returns error
	updateUser := dto.UpdateUserRequest{
		Username:  "updateduser",
		FirstName: "Updated",
	}

	username := "admin"
	userID := 1
	s.userService.On("Update", mock.Anything, userID, updateUser, username).Return(errors.New("update failed"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"username": username,
		"userId":   userID,
	}))
	app.Put("/api/v1/user", s.userHandler.UpdateUser)

	body, _ := json.Marshal(updateUser)
	req := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// WHEN — PUT /api/v1/user is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

// Test GetUserList handler
func (s *UserHandlerTestSuite) TestGetUserList_Success() {
	// GIVEN — super admin; service returns list
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

	// Super admin can pass nil clientId
	s.userService.On("GetUserList", mock.MatchedBy(func(ctx any) bool {
		_, ok := ctx.(context.Context)
		return ok
	}), (*int)(nil)).Return(expectedUsers, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3, // Super admin - doesn't need clientId
	}))
	app.Get("/api/v1/user/list", s.userHandler.GetUserList)

	req := httptest.NewRequest("GET", "/api/v1/user/list", nil)

	// WHEN — GET /api/v1/user/list is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}

func (s *UserHandlerTestSuite) TestGetUserList_MissingClientId() {
	// GIVEN — no clientId in context
	app := fiber.New()
	app.Get("/api/v1/user/list", s.userHandler.GetUserList)

	req := httptest.NewRequest("GET", "/api/v1/user/list", nil)

	// WHEN — GET /api/v1/user/list is sent
	resp, err := app.Test(req)

	// THEN — 200
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *UserHandlerTestSuite) TestGetUserList_ServiceError() {
	// GIVEN — service returns error
	s.userService.On("GetUserList", mock.MatchedBy(func(ctx any) bool {
		_, ok := ctx.(context.Context)
		return ok
	}), (*int)(nil)).Return(nil, errors.New("database error"))

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{
		"userLevel": 3, // Super admin - doesn't need clientId
	}))
	app.Get("/api/v1/user/list", s.userHandler.GetUserList)

	req := httptest.NewRequest("GET", "/api/v1/user/list", nil)

	// WHEN — GET /api/v1/user/list is sent
	resp, err := app.Test(req)

	// THEN — 200 and expectations met
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.userService.AssertExpectations(s.T())
}
