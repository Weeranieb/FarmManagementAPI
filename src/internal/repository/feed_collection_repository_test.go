package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FeedCollectionRepositoryTestSuite struct {
	suite.Suite
	db                *gorm.DB
	feedCollectionRepo FeedCollectionRepository
}

func (s *FeedCollectionRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	err = s.db.AutoMigrate(&model.FeedCollection{}, &model.FeedPriceHistory{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.feedCollectionRepo = NewFeedCollectionRepository(s.db)
}

func (s *FeedCollectionRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *FeedCollectionRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM feed_price_histories")
	s.db.Exec("DELETE FROM feed_collections")
}

func TestFeedCollectionRepositorySuite(t *testing.T) {
	suite.Run(t, new(FeedCollectionRepositoryTestSuite))
}

func (s *FeedCollectionRepositoryTestSuite) TestCreate_Success() {
	feedCollection := &model.FeedCollection{
		ClientId: 1,
		Code:     "FEED001",
		Name:     "Test Feed",
		Unit:     "kg",
	}

	err := s.feedCollectionRepo.Create(feedCollection)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), feedCollection.Id)
	assert.Equal(s.T(), "FEED001", feedCollection.Code)
}

func (s *FeedCollectionRepositoryTestSuite) TestGetByID_Success() {
	feedCollection := &model.FeedCollection{
		ClientId: 1,
		Code:     "FEED001",
		Name:     "Test Feed",
		Unit:     "kg",
	}
	s.feedCollectionRepo.Create(feedCollection)

	result, err := s.feedCollectionRepo.GetByID(feedCollection.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), feedCollection.Id, result.Id)
}

func (s *FeedCollectionRepositoryTestSuite) TestGetByClientIdAndCode_Success() {
	feedCollection := &model.FeedCollection{
		ClientId: 1,
		Code:     "FEED001",
		Name:     "Test Feed",
		Unit:     "kg",
	}
	s.feedCollectionRepo.Create(feedCollection)

	result, err := s.feedCollectionRepo.GetByClientIdAndCode(1, "FEED001")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "FEED001", result.Code)
}

func (s *FeedCollectionRepositoryTestSuite) TestGetPage_Success() {
	feed1 := &model.FeedCollection{ClientId: 1, Code: "FEED001", Name: "Feed 1", Unit: "kg"}
	feed2 := &model.FeedCollection{ClientId: 1, Code: "FEED002", Name: "Feed 2", Unit: "kg"}
	feed3 := &model.FeedCollection{ClientId: 2, Code: "FEED003", Name: "Feed 3", Unit: "kg"}
	s.feedCollectionRepo.Create(feed1)
	s.feedCollectionRepo.Create(feed2)
	s.feedCollectionRepo.Create(feed3)

	results, total, err := s.feedCollectionRepo.GetPage(1, 0, 10, "", "")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), total)
	assert.Len(s.T(), results, 2)
}

func (s *FeedCollectionRepositoryTestSuite) TestUpdate_Success() {
	feedCollection := &model.FeedCollection{
		ClientId: 1,
		Code:     "FEED001",
		Name:     "Test Feed",
		Unit:     "kg",
	}
	s.feedCollectionRepo.Create(feedCollection)

	feedCollection.Name = "Updated Feed"
	err := s.feedCollectionRepo.Update(feedCollection)

	assert.NoError(s.T(), err)
	
	updated, _ := s.feedCollectionRepo.GetByID(feedCollection.Id)
	assert.Equal(s.T(), "Updated Feed", updated.Name)
}

