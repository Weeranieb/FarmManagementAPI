package services_test

// import (
// 	"boonmafarm/api/pkg/models"
// 	"boonmafarm/api/pkg/repositories/mocks"
// 	"boonmafarm/api/pkg/services"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestCreate_ActivePond(t *testing.T) {
// 	var (
// 		mockActivePondRepo *mocks.IActivePondRepository
// 		mockActivePond     services.IActivePondService
// 	)

// 	beforeEach := func() {
// 		mockActivePondRepo = new(mocks.IActivePondRepository)
// 		mockActivePond = services.NewActivePondService(mockActivePondRepo)
// 	}

// 	t.Run("Create active pond success", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("FirstByQuery", "\"PondId\" = ? AND \"IsActive\" = ? AND \"DelFlag\" = ?", 1, true, false).Return(nil, nil)
// 		mockActivePondRepo.On("Create", mock.Anything).Return(&models.ActivePond{
// 			Id:        1,
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, nil)

// 		result, err := mockActivePond.Create(models.AddActivePond{
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, "test")

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.Nil(t, err)
// 		assert.NotNil(t, result)
// 	})

// 	t.Run("Create active pond failed by validation", func(t *testing.T) {
// 		beforeEach()

// 		result, err := mockActivePond.Create(models.AddActivePond{
// 			PondId:    1,
// 			StartDate: time.Time{},
// 		}, "test")

// 		assert.NotNil(t, err)
// 		assert.Nil(t, result)
// 	})

// 	t.Run("Create active pond failed by repository error by FirstByQuery", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("FirstByQuery", "\"PondId\" = ? AND \"IsActive\" = ? AND \"DelFlag\" = ?", 1, true, false).Return(
// 			nil, assert.AnError,
// 		)

// 		result, err := mockActivePond.Create(models.AddActivePond{
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, "test")

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.NotNil(t, err)
// 		assert.Nil(t, result)
// 	})

// 	t.Run("Create active pond failed by repository error by Create", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("FirstByQuery", "\"PondId\" = ? AND \"IsActive\" = ? AND \"DelFlag\" = ?", 1, true, false).Return(nil, nil)
// 		mockActivePondRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

// 		result, err := mockActivePond.Create(models.AddActivePond{
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, "test")

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.NotNil(t, err)
// 		assert.Nil(t, result)
// 	})

// 	t.Run("Create active pond failed by active pond already exist", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("FirstByQuery", "\"PondId\" = ? AND \"IsActive\" = ? AND \"DelFlag\" = ?", 1, true, false).Return(&models.ActivePond{
// 			Id:        1,
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, nil)

// 		result, err := mockActivePond.Create(models.AddActivePond{
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, "test")

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.NotNil(t, err)
// 		assert.Nil(t, result)
// 	})
// }

// func TestGet_ActivePond(t *testing.T) {
// 	var (
// 		mockActivePondRepo *mocks.IActivePondRepository
// 		mockActivePond     services.IActivePondService
// 	)

// 	beforeEach := func() {
// 		mockActivePondRepo = new(mocks.IActivePondRepository)
// 		mockActivePond = services.NewActivePondService(mockActivePondRepo)
// 	}

// 	t.Run("Get active pond success", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("TakeById", 1).Return(&models.ActivePond{
// 			Id:        1,
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, nil)
// 		result, err := mockActivePond.Get(1)

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.Nil(t, err)
// 		assert.NotNil(t, result)
// 	})

// 	t.Run("Get active pond failed by repository error", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("TakeById", 1).Return(
// 			nil, assert.AnError,
// 		)

// 		result, err := mockActivePond.Get(1)

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.NotNil(t, err)
// 		assert.Nil(t, result)
// 	})
// }

// func TestUpdate_ActivePond(t *testing.T) {
// 	var (
// 		mockActivePondRepo *mocks.IActivePondRepository
// 		mockActivePond     services.IActivePondService
// 	)

// 	beforeEach := func() {
// 		mockActivePondRepo = new(mocks.IActivePondRepository)
// 		mockActivePond = services.NewActivePondService(mockActivePondRepo)
// 	}

// 	t.Run("Update active pond success", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("Update", mock.Anything).Return(nil)
// 		err := mockActivePond.Update(&models.ActivePond{
// 			Id:        1,
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, "test")

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("Update active pond failed by repository error", func(t *testing.T) {
// 		beforeEach()

// 		mockActivePondRepo.On("Update", mock.Anything).Return(assert.AnError)
// 		err := mockActivePond.Update(&models.ActivePond{
// 			Id:        1,
// 			PondId:    1,
// 			StartDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
// 		}, "test")

// 		mockActivePondRepo.AssertExpectations(t)
// 		assert.NotNil(t, err)
// 	})
// }
