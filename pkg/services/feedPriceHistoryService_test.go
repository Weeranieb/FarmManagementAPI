package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_FeedPriceHisotry(t *testing.T) {
	var (
		mockFeedPriceHistoryRepo *mocks.IFeedPriceHistoryRepository
		mockFeedPriceHistory     services.IFeedPriceHistoryService
	)

	beforeEach := func() {
		mockFeedPriceHistoryRepo = new(mocks.IFeedPriceHistoryRepository)
		mockFeedPriceHistory = services.NewFeedPriceHistoryService(mockFeedPriceHistoryRepo)
	}

	t.Run("Create feed price history success", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("FirstByQuery", "\"FeedCollectionId\" = ? AND \"PriceUpdatedDate\" = ? AND \"DelFlag\" = ?", 1, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, nil)
		mockFeedPriceHistoryRepo.On("Create", mock.Anything).Return(&models.FeedPriceHistory{
			Id:               1,
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

		result, err := mockFeedPriceHistory.Create(models.AddFeedPriceHistory{
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, "test")

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create feed price history failed by validation", func(t *testing.T) {
		beforeEach()

		result, err := mockFeedPriceHistory.Create(models.AddFeedPriceHistory{
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Time{},
		}, "test")

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create feed price history failed by feed price history already exist", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("FirstByQuery", "\"FeedCollectionId\" = ? AND \"PriceUpdatedDate\" = ? AND \"DelFlag\" = ?", 1, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(&models.FeedPriceHistory{
			Id:               1,
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

		result, err := mockFeedPriceHistory.Create(models.AddFeedPriceHistory{
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, "test")

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create feed price history failed by feed price history already exist", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("FirstByQuery", "\"FeedCollectionId\" = ? AND \"PriceUpdatedDate\" = ? AND \"DelFlag\" = ?", 1, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, errors.New("error"))

		result, err := mockFeedPriceHistory.Create(models.AddFeedPriceHistory{
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0,
				0, 0, 0, time.UTC),
		}, "test")

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create feed price history failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("FirstByQuery", "\"FeedCollectionId\" = ? AND \"PriceUpdatedDate\" = ? AND \"DelFlag\" = ?", 1, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, nil)
		mockFeedPriceHistoryRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := mockFeedPriceHistory.Create(models.AddFeedPriceHistory{
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, "test")

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGet_FeedPriceHisotry(t *testing.T) {
	var (
		mockFeedPriceHistoryRepo *mocks.IFeedPriceHistoryRepository
		mockFeedPriceHistory     services.IFeedPriceHistoryService
	)

	beforeEach := func() {
		mockFeedPriceHistoryRepo = new(mocks.IFeedPriceHistoryRepository)
		mockFeedPriceHistory = services.NewFeedPriceHistoryService(mockFeedPriceHistoryRepo)
	}

	t.Run("Get feed price history success", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("TakeById", 1).Return(&models.FeedPriceHistory{
			Id:               1,
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

		result, err := mockFeedPriceHistory.Get(1)

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get feed price history failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("TakeById", 1).Return(nil, assert.AnError)

		result, err := mockFeedPriceHistory.Get(1)

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate_FeedPriceHisotry(t *testing.T) {
	var (
		mockFeedPriceHistoryRepo *mocks.IFeedPriceHistoryRepository
		mockFeedPriceHistory     services.IFeedPriceHistoryService
	)

	beforeEach := func() {
		mockFeedPriceHistoryRepo = new(mocks.IFeedPriceHistoryRepository)
		mockFeedPriceHistory = services.NewFeedPriceHistoryService(mockFeedPriceHistoryRepo)
	}

	t.Run("Update feed price history success", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("Update", mock.Anything).Return(nil)

		err := mockFeedPriceHistory.Update(&models.FeedPriceHistory{
			Id:               1,
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, "test")

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Update feed price history failed by repository error", func(t *testing.T) {
		beforeEach()

		mockFeedPriceHistoryRepo.On("Update", mock.Anything).Return(assert.AnError)

		err := mockFeedPriceHistory.Update(&models.FeedPriceHistory{
			Id:               1,
			FeedCollectionId: 1,
			Price:            100,
			PriceUpdatedDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, "test")

		mockFeedPriceHistoryRepo.AssertExpectations(t)
		assert.NotNil(t, err)
	})
}
