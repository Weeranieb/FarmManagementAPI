//go:build cgo

package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
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

type DailyLogHandlerTestSuite struct {
	suite.Suite
	dailyLogService *mocks.MockDailyLogService
	handler         DailyLogHandler
}

func (s *DailyLogHandlerTestSuite) SetupTest() {
	s.dailyLogService = mocks.NewMockDailyLogService(s.T())
	s.handler = NewDailyLogHandler(DailyLogHandlerParams{
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
	tc := 0
	body := dto.DailyLogBulkUpsertRequest{
		Month: "2024-04",
		Entries: []dto.DailyLogEntryInput{
			{
				Day: 1, FreshMorning: decimal.Zero, FreshEvening: decimal.Zero,
				PelletMorning: decimal.Zero, PelletEvening: decimal.Zero,
				DeathFishCount: 0, TouristCatchCount: &tc,
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

func (s *DailyLogHandlerTestSuite) TestUploadTemplate_Success() {
	expected := &dto.DailyLogTemplateImportResponse{
		Results: []dto.DailyLogTemplateImportResult{
			{PondId: 1, PondName: "Pond A", RowsImported: 5},
		},
		Skipped: []string{"Sheet2"},
	}
	s.dailyLogService.On("ImportFromTemplate",
		mock.Anything,
		10,
		[]int{1, 3},
		mock.AnythingOfType("[]uint8"),
		"alice",
	).Return(expected, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(s.T(), writer.WriteField("selectedPondIds", "1"))
	require.NoError(s.T(), writer.WriteField("selectedPondIds", "3"))
	part, err := writer.CreateFormFile("file", "template.xlsx")
	require.NoError(s.T(), err)
	_, err = io.WriteString(part, "dummy")
	require.NoError(s.T(), err)
	require.NoError(s.T(), writer.Close())

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "alice", "userLevel": 1}))
	app.Post("/api/v1/farm/:farmId/daily-logs/import-template", s.handler.UploadTemplate)

	req := httptest.NewRequest("POST", "/api/v1/farm/10/daily-logs/import-template", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(s.T(), true, result["result"])
	s.dailyLogService.AssertExpectations(s.T())
}

func (s *DailyLogHandlerTestSuite) TestUploadTemplate_MissingPondIds() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "template.xlsx")
	require.NoError(s.T(), err)
	_, err = io.WriteString(part, "dummy")
	require.NoError(s.T(), err)
	require.NoError(s.T(), writer.Close())

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u", "userLevel": 1}))
	app.Post("/api/v1/farm/:farmId/daily-logs/import-template", s.handler.UploadTemplate)

	req := httptest.NewRequest("POST", "/api/v1/farm/10/daily-logs/import-template", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyLogService.AssertNotCalled(s.T(), "ImportFromTemplate")
}

func (s *DailyLogHandlerTestSuite) TestUploadTemplate_InvalidFileExtension() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(s.T(), writer.WriteField("selectedPondIds", "1"))
	part, err := writer.CreateFormFile("file", "data.csv")
	require.NoError(s.T(), err)
	_, err = io.WriteString(part, "dummy")
	require.NoError(s.T(), err)
	require.NoError(s.T(), writer.Close())

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{"username": "u", "userLevel": 1}))
	app.Post("/api/v1/farm/:farmId/daily-logs/import-template", s.handler.UploadTemplate)

	req := httptest.NewRequest("POST", "/api/v1/farm/10/daily-logs/import-template", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(s.T(), result["error"])
	s.dailyLogService.AssertNotCalled(s.T(), "ImportFromTemplate")
}
