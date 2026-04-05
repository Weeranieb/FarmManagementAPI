package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
)

type DailyLogHandlerTestSuite struct {
	suite.Suite
	dailyLogService *mocks.MockDailyLogService
	handler         DailyLogHandler
}

func (s *DailyLogHandlerTestSuite) SetupTest() {
	s.dailyLogService = mocks.NewMockDailyLogService(s.T())
	s.handler = NewDailyLogHandler(DailyLogHandlerParams{
		Config: &config.Config{
			App: config.AppConfig{DailyLogUploadPath: s.T().TempDir()},
		},
		DailyLogService: s.dailyLogService,
	})
}

func (s *DailyLogHandlerTestSuite) TearDownTest() {
	s.dailyLogService.ExpectedCalls = nil
}

func TestDailyLogHandlerSuite(t *testing.T) {
	suite.Run(t, new(DailyLogHandlerTestSuite))
}

func (s *DailyLogHandlerTestSuite) TestGetMonth_Success() {
	s.dailyLogService.On("GetMonth", mock.Anything, 7, "2024-03").Return(
		&dto.DailyLogMonthResponse{Entries: []dto.DailyLogEntryResponse{}}, nil)
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u", "userLevel": 1}))
	app.Get("/api/v1/pond/:pondId/daily-logs", s.handler.GetMonth)

	req := httptest.NewRequest("GET", "/api/v1/pond/7/daily-logs?month=2024-03", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), true, result["result"])
	s.dailyLogService.AssertExpectations(s.T())
}

func (s *DailyLogHandlerTestSuite) TestGetMonth_MissingMonthQuery() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u"}))
	app.Get("/api/v1/pond/:pondId/daily-logs", s.handler.GetMonth)

	req := httptest.NewRequest("GET", "/api/v1/pond/1/daily-logs", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyLogService.AssertNotCalled(s.T(), "GetMonth")
}

func (s *DailyLogHandlerTestSuite) TestGetMonth_InvalidPondId() {
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u"}))
	app.Get("/api/v1/pond/:pondId/daily-logs", s.handler.GetMonth)

	req := httptest.NewRequest("GET", "/api/v1/pond/x/daily-logs?month=2024-01", nil)
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyLogService.AssertNotCalled(s.T(), "GetMonth")
}

func (s *DailyLogHandlerTestSuite) TestBulkUpsert_Success() {
	body := dto.DailyLogBulkUpsertRequest{
		Month: "2024-04",
		Entries: []dto.DailyLogEntryInput{
			{
				Day: 1, FreshMorning: decimal.Zero, FreshEvening: decimal.Zero,
				PelletMorning: decimal.Zero, PelletEvening: decimal.Zero,
				DeathFishCount: 0, TouristCatchCount: 0,
			},
		},
	}
	s.dailyLogService.On("BulkUpsert", mock.Anything, 3, mock.MatchedBy(func(req dto.DailyLogBulkUpsertRequest) bool {
		return req.Month == body.Month && len(req.Entries) == 1 && req.Entries[0].Day == 1
	}), "alice").Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "alice", "userLevel": 1}))
	app.Put("/api/v1/pond/:pondId/daily-logs", s.handler.BulkUpsert)

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/pond/3/daily-logs", bytes.NewBuffer(raw))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), true, result["result"])
	s.dailyLogService.AssertExpectations(s.T())
}

func (s *DailyLogHandlerTestSuite) TestBulkUpsert_MissingUsername() {
	body := dto.DailyLogBulkUpsertRequest{
		Month: "2024-04",
		Entries: []dto.DailyLogEntryInput{
			{
				Day: 1, FreshMorning: decimal.Zero, FreshEvening: decimal.Zero,
				PelletMorning: decimal.Zero, PelletEvening: decimal.Zero,
			},
		},
	}
	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"userLevel": 1}))
	app.Put("/api/v1/pond/:pondId/daily-logs", s.handler.BulkUpsert)

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/pond/3/daily-logs", bytes.NewBuffer(raw))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyLogService.AssertNotCalled(s.T(), "BulkUpsert")
}

func (s *DailyLogHandlerTestSuite) TestUploadExcel_Success() {
	s.dailyLogService.On(
		"ImportFromExcelFile",
		mock.Anything,
		3,
		mock.MatchedBy(func(p *int) bool { return p == nil }),
		mock.MatchedBy(func(p *int) bool { return p != nil && *p == 9 }),
		"2024-01",
		mock.AnythingOfType("string"),
		"u",
	).Return(2, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(s.T(), writer.WriteField("month", "2024-01"))
	require.NoError(s.T(), writer.WriteField("pelletFeedCollectionId", "9"))
	part, err := writer.CreateFormFile("file", "test.xlsx")
	require.NoError(s.T(), err)
	_, err = io.WriteString(part, "dummy")
	require.NoError(s.T(), err)
	require.NoError(s.T(), writer.Close())

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u", "userLevel": 1}))
	app.Post("/api/v1/pond/:pondId/daily-logs/upload", s.handler.UploadExcel)

	req := httptest.NewRequest("POST", "/api/v1/pond/3/daily-logs/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), true, result["result"])
	s.dailyLogService.AssertExpectations(s.T())
}

func (s *DailyLogHandlerTestSuite) TestUploadExcel_ImportError_RemovesSavedFile() {
	s.dailyLogService.On(
		"ImportFromExcelFile",
		mock.Anything,
		3,
		mock.MatchedBy(func(p *int) bool { return p == nil }),
		mock.MatchedBy(func(p *int) bool { return p == nil }),
		"2024-01",
		mock.AnythingOfType("string"),
		"u",
	).Return(0, errors.ErrValidationFailed)

	uploadRoot := s.T().TempDir()
	s.handler = NewDailyLogHandler(DailyLogHandlerParams{
		Config: &config.Config{
			App: config.AppConfig{DailyLogUploadPath: uploadRoot},
		},
		DailyLogService: s.dailyLogService,
	})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(s.T(), writer.WriteField("month", "2024-01"))
	part, err := writer.CreateFormFile("file", "test.xlsx")
	require.NoError(s.T(), err)
	_, err = io.WriteString(part, "dummy")
	require.NoError(s.T(), err)
	require.NoError(s.T(), writer.Close())

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u", "userLevel": 1}))
	app.Post("/api/v1/pond/:pondId/daily-logs/upload", s.handler.UploadExcel)

	req := httptest.NewRequest("POST", "/api/v1/pond/3/daily-logs/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)

	pondDir := filepath.Join(uploadRoot, "pond_3")
	entries, err := os.ReadDir(pondDir)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), entries, "uploaded file should be removed when import fails")

	s.dailyLogService.AssertExpectations(s.T())
}
