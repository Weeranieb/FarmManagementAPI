package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateClient(t *testing.T) {
	var (
		mockClientRepo *mocks.IClientRepository
		clientService  services.IClientService
	)

	beforeEach := func() {
		mockClientRepo = new(mocks.IClientRepository)
		clientService = services.NewClientService(mockClientRepo)
	}

	t.Run("Create client", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddClient{
			Name:          "testClient",
			OwnerName:     "testOwner",
			ContactNumber: "0123456789",
		}

		mockClientReturn := &models.Client{
			Id:   1,
			Name: "testClient",
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockClientRepo.On("FirstByQuery", "\"Name\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything).Return(nil, nil)
		mockClientRepo.On("Create", mock.Anything).Return(mockClientReturn, nil)

		// Call the Create method
		client, err := clientService.Create(request, "testUser")

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Create client with duplicate name", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddClient{
			Name:          "testClient",
			OwnerName:     "testOwner",
			ContactNumber: "0123456789",
		}

		mockClientReturn := &models.Client{
			Id:   1,
			Name: "testClient",
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockClientRepo.On("FirstByQuery", "\"Name\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything).Return(mockClientReturn, nil)

		// Call the Create method
		client, err := clientService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("FirstByQuery with error on get", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddClient{
			Name:          "testClient",
			OwnerName:     "testOwner",
			ContactNumber: "0123456789",
		}

		// Mock the necessary repository methods
		mockClientRepo.On("FirstByQuery", "\"Name\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		mockClientRepo.On("Create", mock.Anything).Return(nil, nil)

		// Call the Create method
		client, err := clientService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("Create client with invalid data", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddClient{
			Name:          "testClient",
			OwnerName:     "testOwner",
			ContactNumber: "",
		}

		// Call the Create method
		client, err := clientService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Equal(t, "contact number is empty", err.Error())
		assert.Nil(t, client)
	})

	t.Run("Create client with error on create", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddClient{
			Name:          "testClient",
			OwnerName:     "testOwner",
			ContactNumber: "0123456789",
		}

		// Mock the necessary repository methods
		mockClientRepo.On("FirstByQuery", "\"Name\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything).Return(nil, nil)
		mockClientRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		// Call the Create method
		client, err := clientService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestGetClient(t *testing.T) {
	var (
		mockClientRepo *mocks.IClientRepository
		clientService  services.IClientService
	)

	beforeEach := func() {
		mockClientRepo = new(mocks.IClientRepository)
		clientService = services.NewClientService(mockClientRepo)
	}

	t.Run("Get client", func(t *testing.T) {
		beforeEach()

		mockClientReturn := &models.Client{
			Id:   1,
			Name: "testClient",
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockClientRepo.On("TakeById", mock.Anything).Return(mockClientReturn, nil)

		// Call the Get method
		client, err := clientService.Get(1)

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Get client with error on get", func(t *testing.T) {
		beforeEach()

		// Mock the necessary repository methods
		mockClientRepo.On("TakeById", mock.Anything).Return(nil, assert.AnError)

		// Call the Get method
		client, err := clientService.Get(1)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestUpdateClient(t *testing.T) {
	var (
		mockClientRepo *mocks.IClientRepository
		clientService  services.IClientService
	)

	beforeEach := func() {
		mockClientRepo = new(mocks.IClientRepository)
		clientService = services.NewClientService(mockClientRepo)
	}

	t.Run("Update client", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := &models.Client{
			Id:   1,
			Name: "testClient",
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockClientRepo.On("Update", mock.Anything).Return(nil)

		// Call the Update method
		err := clientService.Update(request, "testUser")

		// Assert the result
		assert.NoError(t, err)
	})

	t.Run("Update client with error on update", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := &models.Client{
			Id:   1,
			Name: "testClient",
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockClientRepo.On("Update", mock.Anything).Return(assert.AnError)

		// Call the Update method
		err := clientService.Update(request, "testUser")

		// Assert the result
		assert.Error(t, err)
	})
}
