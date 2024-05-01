package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Farm(t *testing.T) {
	var (
		mockFarmRepo *mocks.IFarmRepository
		farmService  services.IFarmService
	)

	beforeEach := func() {
		mockFarmRepo = new(mocks.IFarmRepository)
		farmService = services.NewFarmService(mockFarmRepo)
	}

	t.Run("Create farm success", func(t *testing.T) {
		beforeEach()

		request := models.AddFarm{
			Code: "F001",
			Name: "Farm 001",
		}

		mockFarmRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(nil, nil)
		mockFarmRepo.On("Create", mock.Anything).Return(&models.Farm{
			Id:   1,
			Code: "F001",
		}, nil)

		result, err := farmService.Create(request, "test", 1)

		mockFarmRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create farm failed by request validation", func(t *testing.T) {
		beforeEach()

		request := models.AddFarm{}

		result, err := farmService.Create(request, "test", 1)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm failed by farm already exist", func(t *testing.T) {
		beforeEach()

		request := models.AddFarm{
			Code: "F001",
			Name: "Farm 001",
		}

		mockFarmRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(&models.Farm{
			Id:   1,
			Code: "F001",
		}, nil)

		result, err := farmService.Create(request, "test", 1)

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm failed by farm already exist", func(t *testing.T) {
		beforeEach()

		request := models.AddFarm{
			Code: "F001",
		}

		mockFarmRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(&models.Farm{
			Id:   1,
			Code: "F001",
		}, nil)

		result, err := farmService.Create(request, "test", 1)

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm failed by repository error", func(t *testing.T) {
		beforeEach()

		request := models.AddFarm{
			Code: "F001",
			Name: "Farm 001",
		}

		mockFarmRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(nil, nil)
		mockFarmRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := farmService.Create(request, "test", 1)

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("FirstByQuery failed by repository error", func(t *testing.T) {
		beforeEach()

		request := models.AddFarm{
			Code: "F001",
			Name: "Farm 001",
		}

		mockFarmRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(nil, assert.AnError)

		result, err := farmService.Create(request, "test", 1)

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGet_Farm(t *testing.T) {
	var (
		mockFarmRepo *mocks.IFarmRepository
		farmService  services.IFarmService
	)

	beforeEach := func() {
		mockFarmRepo = new(mocks.IFarmRepository)
		farmService = services.NewFarmService(mockFarmRepo)
	}

	t.Run("Get farm success", func(t *testing.T) {
		beforeEach()

		mockFarmRepo.On("TakeById", 1).Return(&models.Farm{
			Id:       1,
			ClientId: 1,
		}, nil)

		result, err := farmService.Get(1, 1)

		mockFarmRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get farm failed by farm not found", func(t *testing.T) {
		beforeEach()

		mockFarmRepo.On("TakeById", 1).Return(&models.Farm{
			Id:       1,
			ClientId: 2,
		}, nil)

		result, err := farmService.Get(1, 1)

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Get farm failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFarmRepo.On("TakeById", 1).Return(nil, assert.AnError)

		result, err := farmService.Get(1, 1)

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate_Farm(t *testing.T) {
	var (
		mockFarmRepo *mocks.IFarmRepository
		farmService  services.IFarmService
	)

	beforeEach := func() {
		mockFarmRepo = new(mocks.IFarmRepository)
		farmService = services.NewFarmService(mockFarmRepo)
	}

	t.Run("Update farm success", func(t *testing.T) {
		beforeEach()

		mockFarmRepo.On("Update", mock.Anything).Return(nil)

		err := farmService.Update(&models.Farm{
			Id: 1,
		}, "test")

		mockFarmRepo.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Update farm failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFarmRepo.On("Update", mock.Anything).Return(assert.AnError)

		err := farmService.Update(&models.Farm{
			Id: 1,
		}, "test")

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
	})
}

func TestGetList_Farm(t *testing.T) {
	var (
		mockFarmRepo *mocks.IFarmRepository
		farmService  services.IFarmService
	)

	beforeEach := func() {
		mockFarmRepo = new(mocks.IFarmRepository)
		farmService = services.NewFarmService(mockFarmRepo)
	}

	t.Run("Get farm list success", func(t *testing.T) {
		beforeEach()

		mockFarmRepo.On("TakeAll", 1).Return([]*models.Farm{
			{
				Id:       1,
				ClientId: 1,
			},
		}, nil)

		result, err := farmService.GetList(1)

		mockFarmRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get farm list failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFarmRepo.On("TakeAll", 1).Return(nil, assert.AnError)

		result, err := farmService.GetList(1)

		mockFarmRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
