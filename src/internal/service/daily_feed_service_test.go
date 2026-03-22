package service

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
)

type DailyFeedServiceTestSuite struct {
	suite.Suite
	dailyFeedRepo      *mocks.MockDailyFeedRepository
	activePondRepo     *mocks.MockActivePondRepository
	feedCollectionRepo *mocks.MockFeedCollectionRepository
	priceHistoryRepo   *mocks.MockFeedPriceHistoryRepository
	svc                DailyFeedService
}

func (s *DailyFeedServiceTestSuite) SetupTest() {
	s.dailyFeedRepo = mocks.NewMockDailyFeedRepository(s.T())
	s.activePondRepo = mocks.NewMockActivePondRepository(s.T())
	s.feedCollectionRepo = mocks.NewMockFeedCollectionRepository(s.T())
	s.priceHistoryRepo = mocks.NewMockFeedPriceHistoryRepository(s.T())
	s.svc = NewDailyFeedService(
		s.dailyFeedRepo,
		s.activePondRepo,
		s.feedCollectionRepo,
		s.priceHistoryRepo,
	)
}

func (s *DailyFeedServiceTestSuite) TearDownTest() {
	s.dailyFeedRepo.ExpectedCalls = nil
	s.activePondRepo.ExpectedCalls = nil
	s.feedCollectionRepo.ExpectedCalls = nil
	s.priceHistoryRepo.ExpectedCalls = nil
}

func TestDailyFeedServiceSuite(t *testing.T) {
	suite.Run(t, new(DailyFeedServiceTestSuite))
}

func (s *DailyFeedServiceTestSuite) TestGetMonth_PondNotActive() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 5).Return(nil, nil)

	out, err := s.svc.GetMonth(context.Background(), 5, "2024-01")

	assert.Nil(s.T(), out)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotActive)
}

func (s *DailyFeedServiceTestSuite) TestGetMonth_InvalidMonth() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)

	out, err := s.svc.GetMonth(context.Background(), 1, "not-a-month")

	assert.Nil(s.T(), out)
	var appErr *errors.AppError
	assert.True(s.T(), stderrors.As(err, &appErr))
	assert.Equal(s.T(), errors.ErrValidationFailed.Code, appErr.Code)
}

func (s *DailyFeedServiceTestSuite) TestGetMonth_Success_EmptyTables() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	s.dailyFeedRepo.On("ListByActivePondAndMonth", 10, mock.Anything, mock.Anything).Return([]*model.DailyFeed{}, nil)
	s.dailyFeedRepo.On("ListFeedCollectionIdsByActivePond", 10).Return([]int{}, nil)

	out, err := s.svc.GetMonth(context.Background(), 1, "2024-01")

	require.NoError(s.T(), err)
	require.NotNil(s.T(), out)
	assert.Len(s.T(), out, 0)
}

func (s *DailyFeedServiceTestSuite) TestGetMonth_Success_WithEntriesAndPrices() {
	pondId := 1
	activePondId := 42
	fcId := 7
	month := "2024-01"
	feedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	s.activePondRepo.On("GetActiveByPondID", mock.Anything, pondId).Return(&model.ActivePond{Id: activePondId}, nil)
	s.dailyFeedRepo.On("ListByActivePondAndMonth", activePondId, mock.Anything, mock.Anything).Return([]*model.DailyFeed{
		{
			Id:               100,
			ActivePondId:     activePondId,
			FeedCollectionId: fcId,
			FeedDate:         feedDate,
			MorningAmount:    decimal.RequireFromString("1.5"),
			EveningAmount:    decimal.RequireFromString("2"),
		},
	}, nil)
	s.dailyFeedRepo.On("ListFeedCollectionIdsByActivePond", activePondId).Return([]int{fcId}, nil)
	s.feedCollectionRepo.On("GetByID", fcId).Return(&model.FeedCollection{
		Id:   fcId,
		Name: "Pellet A",
		Unit: "kg",
	}, nil)
	price10 := decimal.RequireFromString("42")
	s.priceHistoryRepo.On("ListByFeedCollectionId", fcId).Return([]*model.FeedPriceHistory{
		{FeedCollectionId: fcId, Price: decimal.RequireFromString("50"), PriceUpdatedDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)},
		{FeedCollectionId: fcId, Price: price10, PriceUpdatedDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)},
	}, nil)

	out, err := s.svc.GetMonth(context.Background(), pondId, month)

	require.NoError(s.T(), err)
	require.Len(s.T(), out, 1)
	assert.Equal(s.T(), fcId, out[0].FeedCollectionId)
	assert.Equal(s.T(), "Pellet A", out[0].FeedCollectionName)
	assert.Equal(s.T(), "kg", out[0].FeedUnit)
	require.Len(s.T(), out[0].Entries, 1)
	assert.Equal(s.T(), 15, out[0].Entries[0].Day)
	assert.True(s.T(), out[0].Entries[0].Morning.Equal(decimal.RequireFromString("1.5")))
	assert.True(s.T(), out[0].Entries[0].Evening.Equal(decimal.RequireFromString("2")))
	require.NotNil(s.T(), out[0].Entries[0].UnitPrice)
	assert.True(s.T(), out[0].Entries[0].UnitPrice.Equal(price10))
}

func (s *DailyFeedServiceTestSuite) TestBulkUpsert_PondNotActive() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 3).Return(nil, nil)

	err := s.svc.BulkUpsert(context.Background(), 3, dto.DailyFeedBulkUpsertRequest{
		FeedCollectionId: 1,
		Month:            "2024-02",
		Entries:          []dto.DailyFeedEntryInput{{Day: 1, Morning: decimal.Zero, Evening: decimal.Zero}},
	}, "user")

	assert.ErrorIs(s.T(), err, errors.ErrPondNotActive)
}

func (s *DailyFeedServiceTestSuite) TestBulkUpsert_InvalidMonth() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)

	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyFeedBulkUpsertRequest{
		FeedCollectionId: 1,
		Month:            "02/2024",
		Entries:          []dto.DailyFeedEntryInput{{Day: 1, Morning: decimal.Zero, Evening: decimal.Zero}},
	}, "user")

	var appErr *errors.AppError
	assert.True(s.T(), stderrors.As(err, &appErr))
	assert.Equal(s.T(), errors.ErrValidationFailed.Code, appErr.Code)
}

func (s *DailyFeedServiceTestSuite) TestBulkUpsert_FeedCollectionNotFound() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	s.feedCollectionRepo.On("GetByID", 99).Return(nil, nil)

	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyFeedBulkUpsertRequest{
		FeedCollectionId: 99,
		Month:            "2024-02",
		Entries:          []dto.DailyFeedEntryInput{{Day: 1, Morning: decimal.Zero, Evening: decimal.Zero}},
	}, "user")

	assert.ErrorIs(s.T(), err, errors.ErrFeedCollectionNotFound)
}

func (s *DailyFeedServiceTestSuite) TestBulkUpsert_Success_UpsertsExpectedRows() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	s.feedCollectionRepo.On("GetByID", 5).Return(&model.FeedCollection{Id: 5}, nil)
	s.dailyFeedRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(feeds []*model.DailyFeed) bool {
		if len(feeds) != 2 {
			return false
		}
		d1 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
		d28 := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)
		return feeds[0].ActivePondId == 10 && feeds[0].FeedCollectionId == 5 &&
			feeds[0].FeedDate.Equal(d1) && feeds[1].FeedDate.Equal(d28)
	})).Return(nil)

	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyFeedBulkUpsertRequest{
		FeedCollectionId: 5,
		Month:            "2024-02",
		Entries: []dto.DailyFeedEntryInput{
			{Day: 1, Morning: decimal.RequireFromString("1"), Evening: decimal.RequireFromString("2")},
			{Day: 28, Morning: decimal.RequireFromString("3"), Evening: decimal.RequireFromString("4")},
		},
	}, "editor")

	assert.NoError(s.T(), err)
}

func (s *DailyFeedServiceTestSuite) TestBulkUpsert_SkipsDayNotInMonth() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 1).Return(&model.ActivePond{Id: 10}, nil)
	s.feedCollectionRepo.On("GetByID", 5).Return(&model.FeedCollection{Id: 5}, nil)
	s.dailyFeedRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(feeds []*model.DailyFeed) bool {
		// Feb 2024 has no day 30 — time.Date rolls to March; entry skipped
		return len(feeds) == 1 && feeds[0].FeedDate.Day() == 1
	})).Return(nil)

	err := s.svc.BulkUpsert(context.Background(), 1, dto.DailyFeedBulkUpsertRequest{
		FeedCollectionId: 5,
		Month:            "2024-02",
		Entries: []dto.DailyFeedEntryInput{
			{Day: 1, Morning: decimal.Zero, Evening: decimal.Zero},
			{Day: 30, Morning: decimal.Zero, Evening: decimal.Zero},
		},
	}, "editor")

	assert.NoError(s.T(), err)
}

func (s *DailyFeedServiceTestSuite) TestDeleteTable_Success() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 2).Return(&model.ActivePond{Id: 88}, nil)
	s.dailyFeedRepo.On("SoftDeleteByActivePondAndFeedCollection", mock.Anything, 88, 3).Return(nil)

	err := s.svc.DeleteTable(context.Background(), 2, 3)

	assert.NoError(s.T(), err)
}

func (s *DailyFeedServiceTestSuite) TestDeleteTable_PondNotActive() {
	s.activePondRepo.On("GetActiveByPondID", mock.Anything, 2).Return(nil, nil)

	err := s.svc.DeleteTable(context.Background(), 2, 3)

	assert.ErrorIs(s.T(), err, errors.ErrPondNotActive)
}
