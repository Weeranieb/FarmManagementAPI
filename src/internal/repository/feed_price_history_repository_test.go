package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FeedPriceHistoryRepositoryTestSuite struct {
	suite.Suite
	db                  *gorm.DB
	feedPriceHistoryRepo FeedPriceHistoryRepository
}

func (s *FeedPriceHistoryRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	err = s.db.AutoMigrate(&model.FeedPriceHistory{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.feedPriceHistoryRepo = NewFeedPriceHistoryRepository(s.db)
}

func (s *FeedPriceHistoryRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *FeedPriceHistoryRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM feed_price_histories")
}

func TestFeedPriceHistoryRepositorySuite(t *testing.T) {
	suite.Run(t, new(FeedPriceHistoryRepositoryTestSuite))
}

func (s *FeedPriceHistoryRepositoryTestSuite) TestCreate_Success() {
	feedPriceHistory := &model.FeedPriceHistory{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}

	err := s.feedPriceHistoryRepo.Create(feedPriceHistory)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), feedPriceHistory.Id)
	assert.Equal(s.T(), 100.50, feedPriceHistory.Price)
}

func (s *FeedPriceHistoryRepositoryTestSuite) TestCreateBatch_Success() {
	feedPriceHistories := []*model.FeedPriceHistory{
		{FeedCollectionId: 1, Price: 100.50, PriceUpdatedDate: time.Now()},
		{FeedCollectionId: 1, Price: 110.00, PriceUpdatedDate: time.Now().AddDate(0, 0, 1)},
	}

	err := s.feedPriceHistoryRepo.CreateBatch(feedPriceHistories)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), feedPriceHistories[0].Id)
	assert.NotZero(s.T(), feedPriceHistories[1].Id)
}

func (s *FeedPriceHistoryRepositoryTestSuite) TestGetByID_Success() {
	feedPriceHistory := &model.FeedPriceHistory{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}
	s.feedPriceHistoryRepo.Create(feedPriceHistory)

	result, err := s.feedPriceHistoryRepo.GetByID(feedPriceHistory.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), feedPriceHistory.Id, result.Id)
}

func (s *FeedPriceHistoryRepositoryTestSuite) TestGetByFeedCollectionIdAndDate_Success() {
	date := time.Now()
	feedPriceHistory := &model.FeedPriceHistory{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: date,
	}
	s.feedPriceHistoryRepo.Create(feedPriceHistory)

	result, err := s.feedPriceHistoryRepo.GetByFeedCollectionIdAndDate(1, date)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), 100.50, result.Price)
}

func (s *FeedPriceHistoryRepositoryTestSuite) TestListByFeedCollectionId_Success() {
	date1 := time.Now()
	date2 := time.Now().AddDate(0, 0, 1)
	feedPriceHistory1 := &model.FeedPriceHistory{FeedCollectionId: 1, Price: 100.50, PriceUpdatedDate: date1}
	feedPriceHistory2 := &model.FeedPriceHistory{FeedCollectionId: 1, Price: 110.00, PriceUpdatedDate: date2}
	feedPriceHistory3 := &model.FeedPriceHistory{FeedCollectionId: 2, Price: 120.00, PriceUpdatedDate: date1}
	s.feedPriceHistoryRepo.Create(feedPriceHistory1)
	s.feedPriceHistoryRepo.Create(feedPriceHistory2)
	s.feedPriceHistoryRepo.Create(feedPriceHistory3)

	results, err := s.feedPriceHistoryRepo.ListByFeedCollectionId(1)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 2)
	// Should be ordered by date DESC
	assert.Equal(s.T(), 110.00, results[0].Price)
}

func (s *FeedPriceHistoryRepositoryTestSuite) TestUpdate_Success() {
	feedPriceHistory := &model.FeedPriceHistory{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}
	s.feedPriceHistoryRepo.Create(feedPriceHistory)

	feedPriceHistory.Price = 150.00
	err := s.feedPriceHistoryRepo.Update(feedPriceHistory)

	assert.NoError(s.T(), err)
	
	updated, _ := s.feedPriceHistoryRepo.GetByID(feedPriceHistory.Id)
	assert.Equal(s.T(), 150.00, updated.Price)
}

