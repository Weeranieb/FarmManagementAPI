package service

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
	"github.com/weeranieb/boonmafarm-backend/src/internal/transaction"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DailyLogServiceTestSuite struct {
	suite.Suite
	db                 *gorm.DB
	dailyLogRepo       *mocks.MockDailyLogRepository
	activePondRepo     *mocks.MockActivePondRepository
	feedCollectionRepo *mocks.MockFeedCollectionRepository
	priceHistoryRepo   *mocks.MockFeedPriceHistoryRepository
	pondRepo           *mocks.MockPondRepository
	svc                DailyLogService
}

func (s *DailyLogServiceTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	s.dailyLogRepo = mocks.NewMockDailyLogRepository(s.T())
	s.activePondRepo = mocks.NewMockActivePondRepository(s.T())
	s.feedCollectionRepo = mocks.NewMockFeedCollectionRepository(s.T())
	s.priceHistoryRepo = mocks.NewMockFeedPriceHistoryRepository(s.T())
	s.pondRepo = mocks.NewMockPondRepository(s.T())
	s.svc = NewDailyLogService(
		s.dailyLogRepo,
		s.activePondRepo,
		s.feedCollectionRepo,
		s.priceHistoryRepo,
		s.pondRepo,
		transaction.NewManager(s.db),
	)
	s.dailyLogRepo.On("WithTx", mock.Anything).Maybe().Return(s.dailyLogRepo)
}

func (s *DailyLogServiceTestSuite) TearDownTest() {
	s.dailyLogRepo.ExpectedCalls = nil
	s.activePondRepo.ExpectedCalls = nil
	s.feedCollectionRepo.ExpectedCalls = nil
	s.priceHistoryRepo.ExpectedCalls = nil
	s.pondRepo.ExpectedCalls = nil
}

func TestDailyLogServiceSuite(t *testing.T) {
	suite.Run(t, new(DailyLogServiceTestSuite))
}

func (s *DailyLogServiceTestSuite) TestGetMonth_PondNotActive() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 3).Return(nil, nil)
	_, err := s.svc.GetMonth(context.Background(), 3, "2024-01")
	assert.ErrorIs(s.T(), err, errors.ErrPondNotActive)
}

func (s *DailyLogServiceTestSuite) TestGetMonth_InvalidMonth() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	_, err := s.svc.GetMonth(context.Background(), 1, "bad")
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), errors.ErrValidationFailed.Message)
}

func (s *DailyLogServiceTestSuite) TestGetMonth_Success_Empty() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	s.dailyLogRepo.On("ListByActivePondAndMonth", 10, mock.Anything, mock.Anything).Return([]*model.DailyLog{}, nil)

	out, err := s.svc.GetMonth(context.Background(), 1, "2024-03")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), out)
	assert.Len(s.T(), out.Entries, 0)
}

func (s *DailyLogServiceTestSuite) TestGetMonth_Success_WithPrices() {
	activePondId := 10
	freshID, pelletID := 11, 12
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: activePondId}, nil)
	fd := time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)
	s.dailyLogRepo.On("ListByActivePondAndMonth", activePondId, mock.Anything, mock.Anything).Return([]*model.DailyLog{
		{
			Id:                     1,
			ActivePondId:           activePondId,
			FeedDate:               fd,
			FreshFeedCollectionId:  intPtr(freshID),
			PelletFeedCollectionId: intPtr(pelletID),
			FreshMorning:           decimal.RequireFromString("1"),
			PelletMorning:          decimal.RequireFromString("2"),
		},
	}, nil)
	s.feedCollectionRepo.On("GetByID", freshID).Return(&model.FeedCollection{Id: freshID, Name: "F", Unit: "kg", FeedType: constants.FeedTypeFresh}, nil)
	s.feedCollectionRepo.On("GetByID", pelletID).Return(&model.FeedCollection{Id: pelletID, Name: "P", Unit: "kg", FeedType: constants.FeedTypePellet}, nil)
	price10 := decimal.RequireFromString("10")
	s.priceHistoryRepo.On("ListByFeedCollectionId", freshID).Return([]*model.FeedPriceHistory{
		{FeedCollectionId: freshID, Price: price10, PriceUpdatedDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
	}, nil)
	s.priceHistoryRepo.On("ListByFeedCollectionId", pelletID).Return([]*model.FeedPriceHistory{}, nil)

	out, err := s.svc.GetMonth(context.Background(), 1, "2024-03")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "F", out.FreshFeedCollectionName)
	assert.Equal(s.T(), "P", out.PelletFeedCollectionName)
	require.Len(s.T(), out.Entries, 1)
	assert.True(s.T(), out.Entries[0].FreshUnitPrice.Equal(price10))
}

func intPtr(v int) *int { return &v }

func readTestXlsx(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("../excel/excel_dailylog/test_no_fishing.xlsx")
	require.NoError(t, err)
	return data
}

func firstSheetName(t *testing.T, xlsxBytes []byte) string {
	t.Helper()
	f, err := excelize.OpenReader(bytes.NewReader(xlsxBytes))
	require.NoError(t, err)
	defer f.Close()
	sheets := f.GetSheetList()
	require.NotEmpty(t, sheets)
	return sheets[0]
}

func (s *DailyLogServiceTestSuite) TestBulkUpsert_PondNotActive() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 3).Return(nil, nil)
	err := s.svc.BulkUpsert(context.Background(), 3, dto.DailyLogBulkUpsertRequest{
		Month: "2024-01",
		Entries: []dto.DailyLogEntryInput{
			{Day: 1, FreshMorning: decimal.Zero, FreshEvening: decimal.Zero, PelletMorning: decimal.Zero, PelletEvening: decimal.Zero},
		},
	}, "u")
	assert.ErrorIs(s.T(), err, errors.ErrPondNotActive)
}

func (s *DailyLogServiceTestSuite) TestBulkUpsert_InvalidMonth() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyLogBulkUpsertRequest{
		Month: "xx",
		Entries: []dto.DailyLogEntryInput{
			{Day: 1, FreshMorning: decimal.Zero, FreshEvening: decimal.Zero, PelletMorning: decimal.Zero, PelletEvening: decimal.Zero},
		},
	}, "u")
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), errors.ErrValidationFailed.Message)
}

func (s *DailyLogServiceTestSuite) TestBulkUpsert_FeedCollectionWrongType() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	pid := 5
	s.feedCollectionRepo.On("GetByID", 5).Return(&model.FeedCollection{Id: 5, FeedType: constants.FeedTypePellet}, nil)
	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyLogBulkUpsertRequest{
		Month:                 "2024-01",
		FreshFeedCollectionId: &pid,
		Entries: []dto.DailyLogEntryInput{
			{Day: 1, FreshMorning: decimal.RequireFromString("1"), FreshEvening: decimal.Zero, PelletMorning: decimal.Zero, PelletEvening: decimal.Zero},
		},
	}, "u")
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), errors.ErrValidationFailed.Message)
}

func (s *DailyLogServiceTestSuite) TestBulkUpsert_Success() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	fid, pid := 4, 5
	s.feedCollectionRepo.On("GetByID", 4).Return(&model.FeedCollection{Id: 4, FeedType: constants.FeedTypeFresh}, nil)
	s.feedCollectionRepo.On("GetByID", 5).Return(&model.FeedCollection{Id: 5, FeedType: constants.FeedTypePellet}, nil)
	s.dailyLogRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(logs []*model.DailyLog) bool {
		return len(logs) == 1 && logs[0].ActivePondId == 10 && *logs[0].FreshFeedCollectionId == 4 &&
			logs[0].FreshMorning.Equal(decimal.RequireFromString("1"))
	})).Return(nil)

	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyLogBulkUpsertRequest{
		Month:                  "2024-01",
		FreshFeedCollectionId:  &fid,
		PelletFeedCollectionId: &pid,
		Entries: []dto.DailyLogEntryInput{
			{Day: 1, FreshMorning: decimal.RequireFromString("1"), FreshEvening: decimal.Zero, PelletMorning: decimal.Zero, PelletEvening: decimal.Zero},
		},
	}, "u")
	assert.NoError(s.T(), err)
}

func (s *DailyLogServiceTestSuite) TestImportFromTemplate_PondNotFound_Skipped() {
	xlsxBytes := readTestXlsx(s.T())

	s.pondRepo.On("ListByFarmId", 1).Return([]*model.Pond{
		{Id: 99, FarmId: 1, Name: "NoMatch"},
	}, nil)

	resp, err := s.svc.ImportFromTemplate(context.Background(), 1, []int{99}, xlsxBytes, "tester")
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), resp.Results)
	assert.NotEmpty(s.T(), resp.Skipped)
}

func (s *DailyLogServiceTestSuite) TestImportFromTemplate_PondNotInSelectedIds_Skipped() {
	xlsxBytes := readTestXlsx(s.T())
	sheetName := firstSheetName(s.T(), xlsxBytes)

	s.pondRepo.On("ListByFarmId", 1).Return([]*model.Pond{
		{Id: 5, FarmId: 1, Name: sheetName},
	}, nil)

	resp, err := s.svc.ImportFromTemplate(context.Background(), 1, []int{999}, xlsxBytes, "tester")
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), resp.Results)
	assert.NotEmpty(s.T(), resp.Skipped)
}

func (s *DailyLogServiceTestSuite) TestImportFromTemplate_Success() {
	xlsxBytes := readTestXlsx(s.T())
	sheetName := firstSheetName(s.T(), xlsxBytes)

	s.pondRepo.On("ListByFarmId", 1).Return([]*model.Pond{
		{Id: 5, FarmId: 1, Name: sheetName},
	}, nil)
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 5).Return(&model.ActivePond{
		Id:        50,
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}, nil)
	s.dailyLogRepo.On("ListIDAndFeedDateByActivePondRange", mock.Anything, 50, mock.Anything, mock.Anything).Return([]repository.DailyLogIDFeedDate{}, nil).Once()
	s.dailyLogRepo.On("HardDeleteByIDs", mock.Anything, mock.MatchedBy(func(ids []int) bool { return len(ids) == 0 })).Return(nil).Once()
	s.dailyLogRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(logs []*model.DailyLog) bool {
		return len(logs) > 0 && logs[0].ActivePondId == 50
	})).Return(nil)

	resp, err := s.svc.ImportFromTemplate(context.Background(), 1, []int{5}, xlsxBytes, "tester")
	assert.NoError(s.T(), err)
	require.Len(s.T(), resp.Results, 1)
	assert.Equal(s.T(), 5, resp.Results[0].PondId)
	assert.Equal(s.T(), sheetName, resp.Results[0].PondName)
	assert.Greater(s.T(), resp.Results[0].RowsImported, 0)
}

func (s *DailyLogServiceTestSuite) TestImportFromTemplate_HardDeletesStaleRowsByID() {
	xlsxBytes := readTestXlsx(s.T())
	sheetName := firstSheetName(s.T(), xlsxBytes)

	s.pondRepo.On("ListByFarmId", 1).Return([]*model.Pond{
		{Id: 5, FarmId: 1, Name: sheetName},
	}, nil)
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 5).Return(&model.ActivePond{
		Id:        50,
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}, nil)
	// A day inside the fixture's Feb–Mar 2026 span that the sparse template does not emit
	stale := repository.DailyLogIDFeedDate{
		Id:       99,
		FeedDate: time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC),
	}
	s.dailyLogRepo.On("ListIDAndFeedDateByActivePondRange", mock.Anything, 50, mock.Anything, mock.Anything).Return([]repository.DailyLogIDFeedDate{stale}, nil).Once()
	s.dailyLogRepo.On("HardDeleteByIDs", mock.Anything, []int{99}).Return(nil).Once()
	s.dailyLogRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(logs []*model.DailyLog) bool {
		return len(logs) > 0 && logs[0].ActivePondId == 50
	})).Return(nil)

	resp, err := s.svc.ImportFromTemplate(context.Background(), 1, []int{5}, xlsxBytes, "tester")
	assert.NoError(s.T(), err)
	require.Len(s.T(), resp.Results, 1)
	assert.Equal(s.T(), 5, resp.Results[0].PondId)
	s.dailyLogRepo.AssertExpectations(s.T())
}

func (s *DailyLogServiceTestSuite) TestBulkUpsert_SkipsInvalidDayForMonth() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	s.dailyLogRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(logs []*model.DailyLog) bool {
		return len(logs) == 0
	})).Return(nil)

	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyLogBulkUpsertRequest{
		Month: "2024-02",
		Entries: []dto.DailyLogEntryInput{
			{Day: 31, FreshMorning: decimal.Zero, FreshEvening: decimal.Zero, PelletMorning: decimal.Zero, PelletEvening: decimal.Zero},
		},
	}, "u")
	assert.NoError(s.T(), err)
}
