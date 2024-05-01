package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFarmOnFarmGroupService_Create(t *testing.T) {
	var (
		mockFarmOnFarmGroupRepo *mocks.IFarmOnFarmGroupRepository
		farmOnFarmGroupService  services.IFarmOnFarmGroupService
	)

	beforeEach := func() {
		mockFarmOnFarmGroupRepo = new(mocks.IFarmOnFarmGroupRepository)
		farmOnFarmGroupService = services.NewFarmOnFarmGroupService(mockFarmOnFarmGroupRepo)
	}

	t.Run("Create farm on farm group success", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmOnFarmGroup{
			FarmId:      1,
			FarmGroupId: 1,
		}

		mockFarmOnFarmGroupRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"FarmGroupId\" = ? AND \"DelFlag\" = ?", request.FarmId, request.FarmGroupId, false).Return(nil, nil)
		mockFarmOnFarmGroupRepo.On("Create", mock.Anything).Return(&models.FarmOnFarmGroup{
			Id:          1,
			FarmId:      1,
			FarmGroupId: 1,
		}, nil)

		result, err := farmOnFarmGroupService.Create(request, "test")

		mockFarmOnFarmGroupRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create farm on farm group failed by request validation", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmOnFarmGroup{}

		result, err := farmOnFarmGroupService.Create(request, "test")

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm on farm group failed by farm on farm group already exist", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmOnFarmGroup{
			FarmId:      1,
			FarmGroupId: 1,
		}

		mockFarmOnFarmGroupRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"FarmGroupId\" = ? AND \"DelFlag\" = ?", request.FarmId, request.FarmGroupId, false).Return(&models.FarmOnFarmGroup{}, nil)

		result, err := farmOnFarmGroupService.Create(request, "test")

		mockFarmOnFarmGroupRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("FirstByQuery failed by first by query", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmOnFarmGroup{
			FarmId:      1,
			FarmGroupId: 1,
		}

		mockFarmOnFarmGroupRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"FarmGroupId\" = ? AND \"DelFlag\" = ?", request.FarmId, request.FarmGroupId, false).Return(nil, assert.AnError)

		result, err := farmOnFarmGroupService.Create(request, "test")

		mockFarmOnFarmGroupRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm on farm group failed by create farm on farm group", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmOnFarmGroup{
			FarmId:      1,
			FarmGroupId: 1,
		}

		mockFarmOnFarmGroupRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"FarmGroupId\" = ? AND \"DelFlag\" = ?", request.FarmId, request.FarmGroupId, false).Return(nil, nil)
		mockFarmOnFarmGroupRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := farmOnFarmGroupService.Create(request, "test")

		mockFarmOnFarmGroupRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestFarmOnFarmGroupService_Delete(t *testing.T) {
	var (
		mockFarmOnFarmGroupRepo *mocks.IFarmOnFarmGroupRepository
		farmOnFarmGroupService  services.IFarmOnFarmGroupService
	)

	beforeEach := func() {
		mockFarmOnFarmGroupRepo = new(mocks.IFarmOnFarmGroupRepository)
		farmOnFarmGroupService = services.NewFarmOnFarmGroupService(mockFarmOnFarmGroupRepo)
	}

	t.Run("Delete farm on farm group success", func(t *testing.T) {
		beforeEach()

		mockFarmOnFarmGroupRepo.On("Delete", 1).Return(nil)

		err := farmOnFarmGroupService.Delete(1)

		mockFarmOnFarmGroupRepo.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Delete farm on farm group failed by delete farm on farm group", func(t *testing.T) {
		beforeEach()

		mockFarmOnFarmGroupRepo.On("Delete", 1).Return(assert.AnError)

		err := farmOnFarmGroupService.Delete(1)

		mockFarmOnFarmGroupRepo.AssertExpectations(t)
		assert.NotNil(t, err)
	})
}
