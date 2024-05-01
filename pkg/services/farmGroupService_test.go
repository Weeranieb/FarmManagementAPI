package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_FarmGroup(t *testing.T) {
	var (
		mockFarmGroupRepo *mocks.IFarmGroupRepository
		farmGroupService  services.IFarmGroupService
	)

	beforeEach := func() {
		mockFarmGroupRepo = new(mocks.IFarmGroupRepository)
		farmGroupService = services.NewFarmGroupService(mockFarmGroupRepo)
	}

	t.Run("Create farm group success", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmGroup{
			Code: "FG001",
			Name: "Farm Group 001",
		}

		mockFarmGroupRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(nil, nil)
		mockFarmGroupRepo.On("Create", mock.Anything).Return(&models.FarmGroup{
			Id:   1,
			Code: "FG001",
		}, nil)

		result, err := farmGroupService.Create(request, "test", 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create farm group failed by request validation", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmGroup{}

		result, err := farmGroupService.Create(request, "test", 1)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm group failed by farm group already exist", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmGroup{
			Code: "FG001",
			Name: "Farm Group 001",
		}

		mockFarmGroupRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(&models.FarmGroup{
			Id:   1,
			Code: "FG001",
		}, nil)

		result, err := farmGroupService.Create(request, "test", 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm group failed by repository error", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmGroup{
			Code: "FG001",
			Name: "Farm Group 001",
		}

		mockFarmGroupRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(nil, assert.AnError)

		result, err := farmGroupService.Create(request, "test", 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("FirstByQuery failed by first by query", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmGroup{
			Code: "FG001",
			Name: "Farm Group 001",
		}

		mockFarmGroupRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(nil, assert.AnError)

		result, err := farmGroupService.Create(request, "test", 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create farm group failed by create farm group", func(t *testing.T) {
		beforeEach()

		request := models.AddFarmGroup{
			Code: "FG001",
			Name: "Farm Group 001",
		}

		mockFarmGroupRepo.On("FirstByQuery", "\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, 1, false).Return(nil, nil)
		mockFarmGroupRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := farmGroupService.Create(request, "test", 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGet_FarmGroup(t *testing.T) {
	var (
		mockFarmGroupRepo *mocks.IFarmGroupRepository
		farmGroupService  services.IFarmGroupService
	)

	beforeEach := func() {
		mockFarmGroupRepo = new(mocks.IFarmGroupRepository)
		farmGroupService = services.NewFarmGroupService(mockFarmGroupRepo)
	}

	t.Run("Get farm group success", func(t *testing.T) {
		beforeEach()

		mockFarmGroupRepo.On("TakeById", 1).Return(&models.FarmGroup{
			Id:       1,
			Code:     "FG001",
			ClientId: 1,
		}, nil)

		result, err := farmGroupService.Get(1, 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get farm group failed by farm group not found", func(t *testing.T) {
		beforeEach()

		mockFarmGroupRepo.On("TakeById", 1).Return(&models.FarmGroup{
			Id:       1,
			Code:     "FG001",
			ClientId: 2,
		}, nil)

		result, err := farmGroupService.Get(1, 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Get farm group failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFarmGroupRepo.On("TakeById", 1).Return(nil, assert.AnError)

		result, err := farmGroupService.Get(1, 1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate_FarmGroup(t *testing.T) {
	var (
		mockFarmGroupRepo *mocks.IFarmGroupRepository
		farmGroupService  services.IFarmGroupService
	)

	beforeEach := func() {
		mockFarmGroupRepo = new(mocks.IFarmGroupRepository)
		farmGroupService = services.NewFarmGroupService(mockFarmGroupRepo)
	}

	t.Run("Update farm group success", func(t *testing.T) {
		beforeEach()

		request := &models.FarmGroup{
			Id:       1,
			Code:     "FG001",
			ClientId: 1,
			Name:     "Farm Group 001",
		}

		mockFarmGroupRepo.On("Update", request).Return(nil)

		err := farmGroupService.Update(request, "test")

		mockFarmGroupRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Update farm group failed by repository error", func(t *testing.T) {
		beforeEach()

		request := &models.FarmGroup{
			Id:       1,
			Code:     "FG001",
			ClientId: 1,
			Name:     "Farm Group 001",
		}

		mockFarmGroupRepo.On("Update", request).Return(assert.AnError)

		err := farmGroupService.Update(request, "test")

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Error(t, err)
	})
}

func TestGetFarmList_FarmGroup(t *testing.T) {
	var (
		mockFarmGroupRepo *mocks.IFarmGroupRepository
		farmGroupService  services.IFarmGroupService
	)

	beforeEach := func() {
		mockFarmGroupRepo = new(mocks.IFarmGroupRepository)
		farmGroupService = services.NewFarmGroupService(mockFarmGroupRepo)
	}

	t.Run("Get farm list success", func(t *testing.T) {
		beforeEach()

		mockFarmGroupRepo.On("GetFarmList", 1).Return(&[]models.Farm{
			{
				Id:       1,
				Code:     "F001",
				ClientId: 1,
			},
		}, nil)

		result, err := farmGroupService.GetFarmList(1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get farm list failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFarmGroupRepo.On("GetFarmList", 1).Return(nil, assert.AnError)

		result, err := farmGroupService.GetFarmList(1)

		mockFarmGroupRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
