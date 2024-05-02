package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"testing"

	"boonmafarm/api/pkg/repositories/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_FeedCollection(t *testing.T) {
	var (
		mockFeedCollectionRepo *mocks.IFeedCollectionRepository
		feedCollectionService  services.IFeedCollectionService
	)

	beforeEach := func() {
		mockFeedCollectionRepo = new(mocks.IFeedCollectionRepository)
		feedCollectionService = services.NewFeedCollectionService(mockFeedCollectionRepo)
	}

	t.Run("Create feed collection success", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("FirstByQuery", "\"ClientId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(nil, nil)
		mockFeedCollectionRepo.On("Create", mock.Anything).Return(&models.FeedCollection{
			Id:       1,
			ClientId: 1,
		}, nil)

		result, err := feedCollectionService.Create(models.AddFeedCollection{
			ClientId: 1,
			Code:     "code",
			Name:     "name",
			Unit:     "Kilogram",
		}, "test")

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create feed collection failed by validation", func(t *testing.T) {
		beforeEach()

		result, err := feedCollectionService.Create(models.AddFeedCollection{
			ClientId: 1,
			Code:     "code",
			Name:     "",
			Unit:     "Kilogram",
		}, "test")

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create feed collection failed by feed collection already exist", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("FirstByQuery", "\"ClientId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(&models.FeedCollection{
			Id:       1,
			ClientId: 1,
		}, nil)

		result, err := feedCollectionService.Create(models.AddFeedCollection{
			ClientId: 1,
			Code:     "code",
			Name:     "name",
			Unit:     "Kilogram",
		}, "test")

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create feed collection failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("FirstByQuery", "\"ClientId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(nil, assert.AnError)

		result, err := feedCollectionService.Create(models.AddFeedCollection{
			ClientId: 1,
			Code:     "code",
			Name:     "name",
			Unit:     "Kilogram",
		}, "test")

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create feed collection failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("FirstByQuery", "\"ClientId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", 1, "code", false).Return(nil, nil)
		mockFeedCollectionRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := feedCollectionService.Create(models.AddFeedCollection{
			ClientId: 1,
			Code:     "code",
			Name:     "name",
			Unit:     "Kilogram",
		}, "test")

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGet_FeedCollection(t *testing.T) {
	var (
		mockFeedCollectionRepo *mocks.IFeedCollectionRepository
		feedCollectionService  services.IFeedCollectionService
	)

	beforeEach := func() {
		mockFeedCollectionRepo = new(mocks.IFeedCollectionRepository)
		feedCollectionService = services.NewFeedCollectionService(mockFeedCollectionRepo)
	}

	t.Run("Get feed collection success", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("TakeById", 1).Return(
			&models.FeedCollection{
				Id:       1,
				ClientId: 1,
			}, nil)

		result, err := feedCollectionService.Get(1)

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get feed collection failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("TakeById", 1).Return(nil, assert.AnError)

		result, err := feedCollectionService.Get(1)

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Get feed collection failed by feed collection not found", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("TakeById", 1).Return(nil, nil)

		result, err := feedCollectionService.Get(1)

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate_FeedCollection(t *testing.T) {
	var (
		mockFeedCollectionRepo *mocks.IFeedCollectionRepository
		feedCollectionService  services.IFeedCollectionService
	)

	beforeEach := func() {
		mockFeedCollectionRepo = new(mocks.IFeedCollectionRepository)
		feedCollectionService = services.NewFeedCollectionService(mockFeedCollectionRepo)
	}

	t.Run("Update feed collection success", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("Update", mock.Anything).Return(nil)

		err := feedCollectionService.Update(&models.FeedCollection{
			Id:       1,
			ClientId: 1,
		}, "test")

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Update feed collection failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFeedCollectionRepo.On("Update", mock.Anything).Return(assert.AnError)

		err := feedCollectionService.Update(&models.FeedCollection{
			Id:       1,
			ClientId: 1,
		}, "test")

		mockFeedCollectionRepo.AssertExpectations(t)
		assert.NotNil(t, err)
	})
}
