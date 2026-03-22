package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
)

type DailyFeedHandlerTestSuite struct {
	suite.Suite
	dailyFeedService *mocks.MockDailyFeedService
	handler          DailyFeedHandler
}

func (s *DailyFeedHandlerTestSuite) SetupTest() {
	s.dailyFeedService = mocks.NewMockDailyFeedService(s.T())
	s.handler = NewDailyFeedHandler(s.dailyFeedService)
}

func (s *DailyFeedHandlerTestSuite) TearDownTest() {
	s.dailyFeedService.ExpectedCalls = nil
}

func TestDailyFeedHandlerSuite(t *testing.T) {
	suite.Run(t, new(DailyFeedHandlerTestSuite))
}

func (s *DailyFeedHandlerTestSuite) TestGetMonth_Success() {
	s.dailyFeedService.On("GetMonth", mock.Anything, 7, "2024-03").Return([]*dto.DailyFeedTableResponse{}, nil)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u", "userLevel": 1}))
	app.Get("/api/v1/pond/:pondId/daily-feed", s.handler.GetMonth)

	req := httptest.NewRequest("GET", "/api/v1/pond/7/daily-feed?month=2024-03", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), true, result["result"])
	s.dailyFeedService.AssertExpectations(s.T())
}

func (s *DailyFeedHandlerTestSuite) TestGetMonth_MissingMonthQuery() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u"}))
	app.Get("/api/v1/pond/:pondId/daily-feed", s.handler.GetMonth)

	req := httptest.NewRequest("GET", "/api/v1/pond/1/daily-feed", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyFeedService.AssertNotCalled(s.T(), "GetMonth")
}

func (s *DailyFeedHandlerTestSuite) TestGetMonth_InvalidPondId() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u"}))
	app.Get("/api/v1/pond/:pondId/daily-feed", s.handler.GetMonth)

	req := httptest.NewRequest("GET", "/api/v1/pond/x/daily-feed?month=2024-01", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyFeedService.AssertNotCalled(s.T(), "GetMonth")
}

func (s *DailyFeedHandlerTestSuite) TestBulkUpsert_Success() {
	body := dto.DailyFeedBulkUpsertRequest{
		FeedCollectionId: 2,
		Month:            "2024-04",
		Entries: []dto.DailyFeedEntryInput{
			{Day: 1, Morning: decimal.RequireFromString("0"), Evening: decimal.RequireFromString("0")},
		},
	}
	s.dailyFeedService.On("BulkUpsert", mock.Anything, 3, body, "alice").Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "alice", "userLevel": 1}))
	app.Put("/api/v1/pond/:pondId/daily-feed", s.handler.BulkUpsert)

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/pond/3/daily-feed", bytes.NewBuffer(raw))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), true, result["result"])
	s.dailyFeedService.AssertExpectations(s.T())
}

func (s *DailyFeedHandlerTestSuite) TestBulkUpsert_MissingUsername() {
	body := dto.DailyFeedBulkUpsertRequest{
		FeedCollectionId: 2,
		Month:            "2024-04",
		Entries: []dto.DailyFeedEntryInput{
			{Day: 1, Morning: decimal.Zero, Evening: decimal.Zero},
		},
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 1}))
	app.Put("/api/v1/pond/:pondId/daily-feed", s.handler.BulkUpsert)

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/pond/3/daily-feed", bytes.NewBuffer(raw))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyFeedService.AssertNotCalled(s.T(), "BulkUpsert")
}

func (s *DailyFeedHandlerTestSuite) TestDeleteTable_Success() {
	s.dailyFeedService.On("DeleteTable", mock.Anything, 4, 9).Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u"}))
	app.Delete("/api/v1/pond/:pondId/daily-feed/:feedCollectionId", s.handler.DeleteTable)

	req := httptest.NewRequest("DELETE", "/api/v1/pond/4/daily-feed/9", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), true, result["result"])
	s.dailyFeedService.AssertExpectations(s.T())
}

func (s *DailyFeedHandlerTestSuite) TestDeleteTable_InvalidPondId() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u"}))
	app.Delete("/api/v1/pond/:pondId/daily-feed/:feedCollectionId", s.handler.DeleteTable)

	req := httptest.NewRequest("DELETE", "/api/v1/pond/bad/daily-feed/9", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyFeedService.AssertNotCalled(s.T(), "DeleteTable")
}
