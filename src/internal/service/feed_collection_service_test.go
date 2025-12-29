package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FeedCollectionServiceTestSuite struct {
	suite.Suite
	feedCollectionRepo   *mocks.MockFeedCollectionRepository
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository
	db                   *gorm.DB
	feedCollectionService FeedCollectionService
}

func (s *FeedCollectionServiceTestSuite) SetupTest() {
	s.feedCollectionRepo = mocks.NewMockFeedCollectionRepository(s.T())
	
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	// Use real repository for feedPriceHistory since it's simpler
	s.feedPriceHistoryRepo = repository.NewFeedPriceHistoryRepository(s.db)
	
	// Auto-migrate
	s.db.AutoMigrate(&model.FeedPriceHistory{})

	s.feedCollectionService = NewFeedCollectionService(s.feedCollectionRepo, s.feedPriceHistoryRepo, s.db)
}

func (s *FeedCollectionServiceTestSuite) TearDownTest() {
	s.feedCollectionRepo.ExpectedCalls = nil
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func TestFeedCollectionServiceSuite(t *testing.T) {
	suite.Run(t, new(FeedCollectionServiceTestSuite))
}

func (s *FeedCollectionServiceTestSuite) TestGet_Success() {
	feedCollectionId := 1
	expectedFarm := &model.FeedCollection{
		Id:       feedCollectionId,
		ClientId: 1,
		Code:     "FEED001",
		Name:     "Test Feed",
		Unit:     "kg",
	}

	s.feedCollectionRepo.On("GetByID", feedCollectionId).Return(expectedFarm, nil)

	result, err := s.feedCollectionService.Get(feedCollectionId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), feedCollectionId, result.Id)
	s.feedCollectionRepo.AssertExpectations(s.T())
}

func (s *FeedCollectionServiceTestSuite) TestGetPage_Success() {
	clientId := 1
	page := 0
	pageSize := 10
	feedCollections := []*model.FeedCollectionPage{
		{FeedCollection: model.FeedCollection{Id: 1, ClientId: clientId, Code: "FEED001", Name: "Feed 1", Unit: "kg"}},
		{FeedCollection: model.FeedCollection{Id: 2, ClientId: clientId, Code: "FEED002", Name: "Feed 2", Unit: "kg"}},
	}

	s.feedCollectionRepo.On("GetPage", clientId, page, pageSize, "", "").Return(feedCollections, int64(2), nil)

	result, err := s.feedCollectionService.GetPage(clientId, page, pageSize, "", "")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int64(2), result.Total)
	s.feedCollectionRepo.AssertExpectations(s.T())
}

