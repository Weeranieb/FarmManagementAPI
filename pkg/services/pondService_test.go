package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Pond(t *testing.T) {
	var (
		mockPondRepo *mocks.IPondRepository
		mockPond     services.IPondService
	)

	beforeEach := func() {
		mockPondRepo = new(mocks.IPondRepository)
		mockPond = services.NewPondService(mockPondRepo)
	}

	t.Run("Create pond success", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(nil, nil)
		mockPondRepo.On("Create", mock.Anything).Return(&models.Pond{
			Id:     1,
			FarmId: 1,
			Code:   "code",
		}, nil)

		result, err := mockPond.Create(models.AddPond{
			FarmId: 1,
			Code:   "code",
			Name:   "name",
		}, "test")

		mockPondRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create pond failed by validation", func(t *testing.T) {
		beforeEach()

		result, err := mockPond.Create(models.AddPond{
			FarmId: 1,
			Code:   "code",
			Name:   "",
		}, "test")

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create pond failed by pond already exist", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(&models.Pond{
			Id:     1,
			FarmId: 1,
			Code:   "code",
		}, nil)

		result, err := mockPond.Create(models.AddPond{
			FarmId: 1,
			Code:   "code",
			Name:   "name",
		}, "test")

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create pond failed by repository error by FirstByQuery", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(nil, assert.AnError)

		result, err := mockPond.Create(models.AddPond{
			FarmId: 1,
			Code:   "code",
			Name:   "name",
		}, "test")

		mockPondRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create pond failed by repository error by Create", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("FirstByQuery", "\"FarmId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(nil, nil)
		mockPondRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := mockPond.Create(models.AddPond{
			FarmId: 1,
			Code:   "code",
			Name:   "name",
		}, "test")

		mockPondRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGet_Pond(t *testing.T) {
	var (
		mockPondRepo *mocks.IPondRepository
		mockPond     services.IPondService
	)

	beforeEach := func() {
		mockPondRepo = new(mocks.IPondRepository)
		mockPond = services.NewPondService(mockPondRepo)
	}

	t.Run("Get pond success", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("TakeById", 1).Return(&models.Pond{
			Id:     1,
			FarmId: 1,
			Code:   "code",
		}, nil)

		result, err := mockPond.Get(1)

		mockPondRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get pond failed by repository error by TakeById", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("TakeById", 1).Return(
			nil, assert.AnError,
		)

		result, err := mockPond.Get(1)

		mockPondRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate_Pond(t *testing.T) {
	var (
		mockPondRepo *mocks.IPondRepository
		mockPond     services.IPondService
	)

	beforeEach := func() {
		mockPondRepo = new(mocks.IPondRepository)
		mockPond = services.NewPondService(mockPondRepo)
	}

	t.Run("Update pond success", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("Update", mock.Anything).Return(nil)

		err := mockPond.Update(&models.Pond{
			Id:     1,
			FarmId: 1,
			Code:   "code",
		}, "test")

		mockPondRepo.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Update pond failed by repository error by Update", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("Update", mock.Anything).Return(assert.AnError)

		err := mockPond.Update(&models.Pond{
			Id:     1,
			FarmId: 1,
			Code:   "code",
		}, "test")

		mockPondRepo.AssertExpectations(t)
		assert.NotNil(t, err)
	})
}

func TestGetList_Pond(t *testing.T) {
	var (
		mockPondRepo *mocks.IPondRepository
		mockPond     services.IPondService
	)

	beforeEach := func() {
		mockPondRepo = new(mocks.IPondRepository)
		mockPond = services.NewPondService(mockPondRepo)
	}

	t.Run("Get list pond success", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("TakeAll", 1).Return([]*models.Pond{
			{
				Id:     1,
				FarmId: 1,
				Code:   "code",
			},
		}, nil)

		result, err := mockPond.GetList(1)

		mockPondRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get list pond failed by repository error by GetList", func(t *testing.T) {
		beforeEach()

		mockPondRepo.On("TakeAll", 1).Return(nil, assert.AnError)

		result, err := mockPond.GetList(1)

		mockPondRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
