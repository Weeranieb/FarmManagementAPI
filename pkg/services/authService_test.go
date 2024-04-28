package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	t.Run("Create user", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		authService := services.NewAuthService(mockUserRepo)

		// Define test data
		request := models.AddUser{
			ClientId:      1,
			ContactNumber: "0123456789",
			Username:      "testUser",
			Password:      "testPassword",
			FirstName:     "Test",
			UserLevel:     1,
		}
		mockUserReturn := &models.User{
			Id:            1,
			ClientId:      1,
			FirstName:     "Test",
			LastName:      nil,
			ContactNumber: "0123456789",
			Username:      "testUser",
			Password:      "$2a$10$TTzfSpMsKJ4k4G/pAS99l.qc1ywYLiRgEcMf.mC.rf78qpAI4IN12",
			UserLevel:     1,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(mockUserReturn, nil)
		mockUserRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(nil, nil)

		// Call the Create method
		user, err := authService.Create(request)

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("found user", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		authService := services.NewAuthService(mockUserRepo)

		// Define test data
		request := models.AddUser{
			ClientId:      1,
			ContactNumber: "0123456789",
			Username:      "testUser",
			Password:      "testPassword",
			FirstName:     "Test",
			UserLevel:     1,
		}
		mockUserReturn := &models.User{
			Id:            1,
			ClientId:      1,
			FirstName:     "Test",
			LastName:      nil,
			ContactNumber: "0123456789",
			Username:      "testUser",
			Password:      "$2a$10$TTzfSpMsKJ4k4G/pAS99l.qc1ywYLiRgEcMf.mC.rf78qpAI4IN12",
			UserLevel:     1,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(mockUserReturn, nil)
		mockUserRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(&models.User{}, nil)

		// Call the Create method
		user, err := authService.Create(request)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("not validation", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		authService := services.NewAuthService(mockUserRepo)

		// Define test data
		request := models.AddUser{
			Username:  "testUser",
			Password:  "testPassword",
			UserLevel: 1,
		}
		mockUserReturn := &models.User{
			Id:            1,
			ClientId:      1,
			FirstName:     "Test",
			LastName:      nil,
			ContactNumber: "",
			Username:      "testUser",
			Password:      "$2a$10$TTzfSpMsKJ4k4G/pAS99l.qc1ywYLiRgEcMf.mC.rf78qpAI4IN12",
			UserLevel:     1,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(mockUserReturn, nil)
		mockUserRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(&models.User{}, nil)

		// Call the Create method
		user, err := authService.Create(request)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("error on create", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		authService := services.NewAuthService(mockUserRepo)

		// Define test data
		request := models.AddUser{
			Username:  "testUser",
			Password:  "testPassword",
			FirstName: "Test",
			UserLevel: 1,
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(nil, errors.New("error on create"))
		mockUserRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(nil, nil)

		// Call the Create method
		user, err := authService.Create(request)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("error on check", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		authService := services.NewAuthService(mockUserRepo)

		// Define test data
		request := models.AddUser{
			ClientId:      1,
			Username:      "testUser",
			Password:      "testPassword",
			FirstName:     "Test",
			ContactNumber: "1234567890",
			UserLevel:     1,
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(nil, nil)
		mockUserRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(nil, errors.New("error on firstByQuery"))

		// Call the Create method
		user, err := authService.Create(request)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("error on found user", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		authService := services.NewAuthService(mockUserRepo)

		// Define test data
		request := models.AddUser{
			ClientId:      1,
			Username:      "testUser",
			Password:      "testPassword",
			FirstName:     "Test",
			ContactNumber: "1234567890",
			UserLevel:     1,
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(nil, assert.AnError)
		mockUserRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(nil, nil)

		// Call the Create method
		user, err := authService.Create(request)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestLogin(t *testing.T) {
	// Your existing code here

	t.Run("TestLogin", func(t *testing.T) {

		mockRepo := new(mocks.IUserRepository)
		authService := services.NewAuthService(mockRepo)

		// Mock user
		mockUser := &models.User{
			Username: "testUser",
			Password: "$2a$10$uHzoME5/vYsXxVlhgIfbQOaWxqPAzaE1xVi0GGiAYU5nfiki5ZVAe", // bcrypt hash of "password"
			// Add other fields as needed
		}

		// Expectation for FirstByQuery
		mockRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(mockUser, nil)

		// Perform the login
		token, err := authService.Login(models.Login{
			Username: "testUser",
			Password: "testPassword", // Password to be compared with the hash
		})

		if err != nil {
			t.Errorf("Login failed: %v", err)
		}

		// Check if token is not empty
		if token == "" {
			t.Error("Token is empty")
		}
	})

	t.Run("TestLogin user error", func(t *testing.T) {

		mockRepo := new(mocks.IUserRepository)
		authService := services.NewAuthService(mockRepo)

		// Expectation for FirstByQuery
		mockRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(nil, assert.AnError)

		// Perform the login
		token, err := authService.Login(models.Login{
			Username: "testUser",
			Password: "testPassword", // Password to be compared with the hash
		})

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("TestLogin user not found", func(t *testing.T) {
		mockRepo := new(mocks.IUserRepository)
		authService := services.NewAuthService(mockRepo)

		// Expectation for FirstByQuery
		mockRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(nil, nil)

		// Perform the login
		token, err := authService.Login(models.Login{
			Username: "testUser",
			Password: "testPassword", // Password to be compared with the hash
		})

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("TestLogin password error", func(t *testing.T) {
		mockRepo := new(mocks.IUserRepository)
		authService := services.NewAuthService(mockRepo)

		// Mock user
		mockUser := &models.User{
			Username: "testUser",
			Password: "$2a$10$uHzoME5/vYsXxVlhgIfbQOaWxqPAzaE1xVi0GGiAYU5nfiki5ZVAe", // bcrypt hash of "password"
			// Add other fields as needed
		}

		// Expectation for FirstByQuery
		mockRepo.On("FirstByQuery", "\"Username\" = ? AND \"DelFlag\" = ?", "testUser", false).Return(mockUser, nil)

		// Perform the login
		token, err := authService.Login(models.Login{
			Username: "testUser",
			Password: "testPassord", // Password to be compared with the hash
		})

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
