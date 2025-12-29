package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FeedPriceHistoryHandlerTestSuite struct {
	suite.Suite
	feedPriceHistoryService service.FeedPriceHistoryService
	feedPriceHistoryHandler FeedPriceHistoryHandler
	db                      *gorm.DB
}

func (s *FeedPriceHistoryHandlerTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	s.db.AutoMigrate(&model.FeedPriceHistory{})
	feedPriceHistoryRepo := repository.NewFeedPriceHistoryRepository(s.db)
	s.feedPriceHistoryService = service.NewFeedPriceHistoryService(feedPriceHistoryRepo)
	s.feedPriceHistoryHandler = NewFeedPriceHistoryHandler(s.feedPriceHistoryService)
}

func (s *FeedPriceHistoryHandlerTestSuite) TearDownTest() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func TestFeedPriceHistoryHandlerSuite(t *testing.T) {
	suite.Run(t, new(FeedPriceHistoryHandlerTestSuite))
}

func (s *FeedPriceHistoryHandlerTestSuite) TestAddFeedPriceHistory_Success() {
	createReq := &dto.CreateFeedPriceHistoryRequest{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}

	username := "admin"
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
	}))
	app.Post("/api/v1/feedpricehistory", s.feedPriceHistoryHandler.AddFeedPriceHistory)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/feedpricehistory", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *FeedPriceHistoryHandlerTestSuite) TestGetFeedPriceHistory_Success() {
	// First create one
	createReq := dto.CreateFeedPriceHistoryRequest{
		FeedCollectionId: 1,
		Price:            100.50,
		PriceUpdatedDate: time.Now(),
	}
	_, _ = s.feedPriceHistoryService.Create(createReq, "admin")

	app := fiber.New()
	app.Get("/api/v1/feedpricehistory/:id", s.feedPriceHistoryHandler.GetFeedPriceHistory)

	req := httptest.NewRequest("GET", "/api/v1/feedpricehistory/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *FeedPriceHistoryHandlerTestSuite) TestGetAllFeedPriceHistory_Success() {
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

	app := fiber.New()
	app.Get("/api/v1/feedpricehistory", s.feedPriceHistoryHandler.GetAllFeedPriceHistory)

	req := httptest.NewRequest("GET", "/api/v1/feedpricehistory?feedCollectionId=1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}
