package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateDailyFeed(t *testing.T) {
	var (
		mockDailyFeedRepo *mocks.IDailyFeedRepository
		dailyFeedService  services.IDailyFeedService
	)

	beforeEach := func() {
		mockDailyFeedRepo = new(mocks.IDailyFeedRepository)
		dailyFeedService = services.NewDailyFeedService(mockDailyFeedRepo)
	}

	t.Run("Create daily feed", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddDailyFeed{
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		mockDailyFeedReturn := &models.DailyFeed{
			Id:               1,
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("FirstByQuery", "\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockDailyFeedRepo.On("Create", mock.Anything).Return(mockDailyFeedReturn, nil)

		// Call the Create method
		dailyFeed, err := dailyFeedService.Create(request, "testUser")

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, dailyFeed)
	})

	t.Run("Create daily feed with duplicate feed", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddDailyFeed{
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		mockDailyFeedReturn := &models.DailyFeed{
			Id:               1,
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("FirstByQuery", "\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockDailyFeedReturn, nil)

		// Call the Create method
		dailyFeed, err := dailyFeedService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, dailyFeed)
	})

	t.Run("Create daily feed with invalid data", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddDailyFeed{
			ActivePondId:     1,
			FeedCollectionId: 1,
			Amount:           1000,
		}

		// Call the Create method
		dailyFeed, err := dailyFeedService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Equal(t, "feedDate is empty", err.Error())
		assert.Nil(t, dailyFeed)
	})

	t.Run("Create daily feed with error on create", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddDailyFeed{
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("FirstByQuery", "\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockDailyFeedRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		// Call the Create method
		dailyFeed, err := dailyFeedService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, dailyFeed)
	})

	t.Run("Create daily feed with error on get", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddDailyFeed{
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("FirstByQuery", "\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

		// Call the Create method
		dailyFeed, err := dailyFeedService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, dailyFeed)
	})

	t.Run("Create daily feed with duplicate feed", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddDailyFeed{
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		mockDailyFeedReturn := &models.DailyFeed{
			Id:               1,
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("FirstByQuery", "\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockDailyFeedReturn, nil)

		// Call the Create method
		dailyFeed, err := dailyFeedService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, dailyFeed)
	})

	t.Run("Create daily feed with error on create", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddDailyFeed{
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("FirstByQuery", "\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockDailyFeedRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		// Call the Create method
		dailyFeed, err := dailyFeedService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, dailyFeed)
	})
}

func TestBulkCreateDailyFeed(t *testing.T) {
	var (
		mockDailyFeedRepo *mocks.IDailyFeedRepository
		dailyFeedService  services.IDailyFeedService
	)

	beforeEach := func() {
		mockDailyFeedRepo = new(mocks.IDailyFeedRepository)
		dailyFeedService = services.NewDailyFeedService(mockDailyFeedRepo)
	}

	t.Run("Bulk create daily feed", func(t *testing.T) {
		beforeEach()

		// Define test data
		// excelFile := "test.xlsx"

		// Call the BulkCreate method
		err := dailyFeedService.BulkCreate(nil, "testUser")

		// Assert the result
		assert.NoError(t, err)
	})
}

func TestGetDailyFeed(t *testing.T) {
	var (
		mockDailyFeedRepo *mocks.IDailyFeedRepository
		dailyFeedService  services.IDailyFeedService
	)

	beforeEach := func() {
		mockDailyFeedRepo = new(mocks.IDailyFeedRepository)
		dailyFeedService = services.NewDailyFeedService(mockDailyFeedRepo)
	}

	t.Run("Get daily feed", func(t *testing.T) {
		beforeEach()

		// Define test data
		mockDailyFeedReturn := &models.DailyFeed{
			Id:               1,
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("TakeById", mock.Anything).Return(mockDailyFeedReturn, nil)

		// Call the Get method
		dailyFeed, err := dailyFeedService.Get(1)

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, dailyFeed)
	})

	t.Run("Get daily feed with error on get", func(t *testing.T) {
		beforeEach()

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("TakeById", mock.Anything).Return(nil, assert.AnError)

		// Call the Get method
		dailyFeed, err := dailyFeedService.Get(1)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, dailyFeed)
	})
}

func TestUpdateDailyFeed(t *testing.T) {
	var (
		mockDailyFeedRepo *mocks.IDailyFeedRepository
		dailyFeedService  services.IDailyFeedService
	)

	beforeEach := func() {
		mockDailyFeedRepo = new(mocks.IDailyFeedRepository)
		dailyFeedService = services.NewDailyFeedService(mockDailyFeedRepo)
	}

	t.Run("Update daily feed", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := &models.DailyFeed{
			Id:               1,
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		// Call the Update method
		err := dailyFeedService.Update(request, "testUser")

		// Assert the result
		assert.NoError(t, err)
	})

	t.Run("Update daily feed with error on update", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := &models.DailyFeed{
			Id:               1,
			ActivePondId:     1,
			FeedCollectionId: 1,
			FeedDate:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount:           1000,
		}

		// Mock the necessary repository methods
		mockDailyFeedRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError)

		// Call the Update method
		err := dailyFeedService.Update(request, "testUser")

		// Assert the result
		assert.Error(t, err)
	})
}
