package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FeedPriceHistoryServiceTestSuite struct {
	suite.Suite
	feedPriceHistoryRepo   repository.FeedPriceHistoryRepository
	feedPriceHistoryService FeedPriceHistoryService
	db                     *gorm.DB
}

func (s *FeedPriceHistoryServiceTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	s.db.AutoMigrate(&model.FeedPriceHistory{})
	s.feedPriceHistoryRepo = repository.NewFeedPriceHistoryRepository(s.db)
	s.feedPriceHistoryService = NewFeedPriceHistoryService(s.feedPriceHistoryRepo)
}

func (s *FeedPriceHistoryServiceTestSuite) TearDownTest() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func TestFeedPriceHistoryServiceSuite(t *testing.T) {
	suite.Run(t, new(FeedPriceHistoryServiceTestSuite))
}

func (s *FeedPriceHistoryServiceTestSuite) TestCreate_Success() {
	req := dto.CreateFeedPriceHistoryRequest{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}
	username := "admin"

	result, err := s.feedPriceHistoryService.Create(req, username)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Price, result.Price)
	assert.Equal(s.T(), req.FeedCollectionId, result.FeedCollectionId)
}

func (s *FeedPriceHistoryServiceTestSuite) TestGet_Success() {
	// First create one
	req := dto.CreateFeedPriceHistoryRequest{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}
	created, _ := s.feedPriceHistoryService.Create(req, "admin")

	result, err := s.feedPriceHistoryService.Get(created.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), created.Id, result.Id)
}

func (s *FeedPriceHistoryServiceTestSuite) TestGetAll_Success() {
	feedCollectionId := 1
	req1 := dto.CreateFeedPriceHistoryRequest{
		FeedCollectionId: feedCollectionId,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}
	req2 := dto.CreateFeedPriceHistoryRequest{
		FeedCollectionId: feedCollectionId,
		Price:            110.00,
		PriceUpdatedDate: time.Now().AddDate(0, 0, 1),
	}
	s.feedPriceHistoryService.Create(req1, "admin")
	s.feedPriceHistoryService.Create(req2, "admin")

	result, err := s.feedPriceHistoryService.GetAll(feedCollectionId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	// Should be ordered by date DESC
	assert.Equal(s.T(), 110.00, result[0].Price)
}

