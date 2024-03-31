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
			Username:  "testUser",
			Password:  "testPassword",
			FirstName: "Test",
			IsAdmin:   false,
		}
		userIdentity := "testUser"
		clientId := 1
		mockUserReturn := &models.User{
			Id:            1,
			ClientId:      1,
			FirstName:     "Test",
			LastName:      nil,
			ContactNumber: nil,
			Username:      "testUser",
			Password:      "$2a$10$TTzfSpMsKJ4k4G/pAS99l.qc1ywYLiRgEcMf.mC.rf78qpAI4IN12",
			IsAdmin:       false,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(mockUserReturn, nil)
		mockUserRepo.On("FirstByQuery", "Username = ?", "testUser").Return(nil, nil)

		// Call the Create method
		user, err := authService.Create(request, userIdentity, clientId)

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
			Username:  "testUser",
			Password:  "testPassword",
			FirstName: "Test",
			IsAdmin:   false,
		}
		userIdentity := "testUser"
		clientId := 1
		mockUserReturn := &models.User{
			Id:            1,
			ClientId:      1,
			FirstName:     "Test",
			LastName:      nil,
			ContactNumber: nil,
			Username:      "testUser",
			Password:      "$2a$10$TTzfSpMsKJ4k4G/pAS99l.qc1ywYLiRgEcMf.mC.rf78qpAI4IN12",
			IsAdmin:       false,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(mockUserReturn, nil)
		mockUserRepo.On("FirstByQuery", "Username = ?", "testUser").Return(&models.User{}, nil)

		// Call the Create method
		user, err := authService.Create(request, userIdentity, clientId)

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
			Username: "testUser",
			Password: "testPassword",
			IsAdmin:  false,
		}
		userIdentity := "testUser"
		clientId := 1
		mockUserReturn := &models.User{
			Id:            1,
			ClientId:      1,
			FirstName:     "Test",
			LastName:      nil,
			ContactNumber: nil,
			Username:      "testUser",
			Password:      "$2a$10$TTzfSpMsKJ4k4G/pAS99l.qc1ywYLiRgEcMf.mC.rf78qpAI4IN12",
			IsAdmin:       false,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(mockUserReturn, nil)
		mockUserRepo.On("FirstByQuery", "Username = ?", "testUser").Return(&models.User{}, nil)

		// Call the Create method
		user, err := authService.Create(request, userIdentity, clientId)

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
			IsAdmin:   false,
		}
		userIdentity := "testUser"
		clientId := 1

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(nil, errors.New("error on create"))
		mockUserRepo.On("FirstByQuery", "Username = ?", "testUser").Return(nil, nil)

		// Call the Create method
		user, err := authService.Create(request, userIdentity, clientId)

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
			Username:  "testUser",
			Password:  "testPassword",
			FirstName: "Test",
			IsAdmin:   false,
		}
		userIdentity := "testUser"
		clientId := 1

		// Mock the necessary repository methods
		mockUserRepo.On("Create", mock.Anything).Return(nil, nil)
		mockUserRepo.On("FirstByQuery", "Username = ?", "testUser").Return(nil, errors.New("error on firstByQuery"))

		// Call the Create method
		user, err := authService.Create(request, userIdentity, clientId)

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
		mockRepo.On("FirstByQuery", "Username = ?", "testUser").Return(mockUser, nil)

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
		mockRepo.On("FirstByQuery", "Username = ?", "testUser").Return(mockUser, nil)

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
}
