package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Worker(t *testing.T) {
	var (
		mockWorkerRepo *mocks.IWorkerRepository
		mockWorker     services.IWorkerService
	)

	beforeEach := func() {
		mockWorkerRepo = new(mocks.IWorkerRepository)
		mockWorker = services.NewWorkerService(mockWorkerRepo)
	}

	t.Run("Create worker success", func(t *testing.T) {
		beforeEach()

		mockWorkerRepo.On("FirstByQuery", "\"FarmGroupId\" = ? AND \"DelFlag\" = ?", 1, false).Return(nil, nil)
		mockWorkerRepo.On("Create", mock.Anything).Return(&models.Worker{
			Id:          1,
			ClientId:    1,
			FarmGroupId: 1,
			FirstName:   "name",
			Nationality: "ไทย",
			Salary:      10000,
		}, nil)

		result, err := mockWorker.Create(models.AddWorker{
			FarmGroupId: 1,
			FirstName:   "name",
			Salary:      10000,
			Nationality: "ไทย",
		}, "test", 1)

		mockWorkerRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create worker failed by validation", func(t *testing.T) {
		beforeEach()

		result, err := mockWorker.Create(models.AddWorker{
			Nationality: "ไทย",
			FarmGroupId: 1,
			FirstName:   "",
			Salary:      10000,
		}, "test", 1)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create worker failed by worker already exist", func(t *testing.T) {
		beforeEach()

		mockWorkerRepo.On("FirstByQuery", "\"FarmGroupId\" = ? AND \"DelFlag\" = ?", 1, false).Return(&models.Worker{
			Id:          1,
			ClientId:    1,
			FarmGroupId: 1,
			FirstName:   "name",
			Salary:      10000,
		}, nil)

		result, err := mockWorker.Create(models.AddWorker{
			Nationality: "ไทย",
			FarmGroupId: 1,
			FirstName:   "name",
			Salary:      10000,
		}, "test", 1)

		mockWorkerRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create worker failed by repository error by FirstByQuery", func(t *testing.T) {
		beforeEach()

		mockWorkerRepo.On("FirstByQuery", "\"FarmGroupId\" = ? AND \"DelFlag\" = ?", 1, false).Return(nil, assert.AnError)

		result, err := mockWorker.Create(models.AddWorker{
			Nationality: "ไทย",
			FarmGroupId: 1,
			FirstName:   "name",
			Salary:      10000,
		}, "test", 1)

		mockWorkerRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create worker failed by repository error by Create", func(t *testing.T) {
		beforeEach()

		mockWorkerRepo.On("FirstByQuery", "\"FarmGroupId\" = ? AND \"DelFlag\" = ?", 1, false).Return(nil, nil)
		mockWorkerRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := mockWorker.Create(models.AddWorker{
			Nationality: "ไทย",
			FarmGroupId: 1,
			FirstName:   "name",
			Salary:      10000,
		}, "test", 1)

		mockWorkerRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGet_Worker(t *testing.T) {
	var (
		mockWorkerRepo *mocks.IWorkerRepository
		mockWorker     services.IWorkerService
	)

	beforeEach := func() {
		mockWorkerRepo = new(mocks.IWorkerRepository)
		mockWorker = services.NewWorkerService(mockWorkerRepo)
	}

	t.Run("Get worker success", func(t *testing.T) {
		beforeEach()

		mockWorkerRepo.On("TakeById", 1).Return(&models.Worker{
			Id:          1,
			ClientId:    1,
			FarmGroupId: 1,
			FirstName:   "name",
			Nationality: "ไทย",
			Salary:      10000,
		}, nil)

		result, err := mockWorker.Get(1)

		mockWorkerRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get worker failed by repository error", func(t *testing.T) {
		beforeEach()

		mockWorkerRepo.On("TakeById", 1).Return(nil, assert.AnError)

		result, err := mockWorker.Get(1)

		mockWorkerRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate_Worker(t *testing.T) {
	var (
		mockWorkerRepo *mocks.IWorkerRepository
		mockWorker     services.IWorkerService
	)

	beforeEach := func() {
		mockWorkerRepo = new(mocks.IWorkerRepository)
		mockWorker = services.NewWorkerService(mockWorkerRepo)
	}

	t.Run("Update worker success", func(t *testing.T) {
		beforeEach()

		mockWorkerRepo.On("Update", mock.Anything).Return(nil)

		err := mockWorker.Update(&models.Worker{
			ClientId:    1,
			FarmGroupId: 1,
			FirstName:   "name",
			Salary:      10000,
			Nationality: "ไทย",
		}, "test")

		mockWorkerRepo.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Update worker failed by repository", func(t *testing.T) {
		beforeEach()

		// Update
		mockWorkerRepo.On("Update", mock.Anything).Return(assert.AnError)

		err := mockWorker.Update(&models.Worker{
			ClientId:    1,
			FarmGroupId: 1,
			FirstName:   "name",
			Salary:      10000,
		}, "test")

		mockWorkerRepo.AssertExpectations(t)
		assert.NotNil(t, err)
	})

}
