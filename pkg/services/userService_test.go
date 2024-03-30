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

func TestCreate(t *testing.T) {
	t.Run("Create user", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		userService := services.NewUserService(mockUserRepo)

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
		user, err := userService.Create(request, userIdentity, clientId)

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
	t.Run("found user", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		userService := services.NewUserService(mockUserRepo)

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
		user, err := userService.Create(request, userIdentity, clientId)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("not validation", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		userService := services.NewUserService(mockUserRepo)

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
		user, err := userService.Create(request, userIdentity, clientId)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})
	t.Run("error on create", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		userService := services.NewUserService(mockUserRepo)

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
		user, err := userService.Create(request, userIdentity, clientId)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})
	t.Run("error on check", func(t *testing.T) {
		// Create a mock user repository
		mockUserRepo := new(mocks.IUserRepository)

		// Create a new instance of the user service with the mock repository
		userService := services.NewUserService(mockUserRepo)

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
		user, err := userService.Create(request, userIdentity, clientId)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
