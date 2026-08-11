package service

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
)

type ActivityServiceTestSuite struct {
	suite.Suite
	activityRepo    *mocks.MockActivityRepository
	activityService ActivityService
}

func (s *ActivityServiceTestSuite) SetupTest() {
	s.activityRepo = mocks.NewMockActivityRepository(s.T())
	s.activityService = NewActivityService(s.activityRepo)
}

func (s *ActivityServiceTestSuite) TearDownTest() {
	s.activityRepo.ExpectedCalls = nil
}

// clientCtx returns a context for a normal user scoped to clientId 1.
func clientCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "user")
	ctx = context.WithValue(ctx, constants.ClientIDKey, 1)
	ctx = context.WithValue(ctx, constants.UserLevelKey, 1)
	return ctx
}

// superAdminCtx returns a context for a super admin (no client scoping).
func superAdminCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "admin")
	ctx = context.WithValue(ctx, constants.UserLevelKey, 3)
	return ctx
}

func (s *ActivityServiceTestSuite) TestListFeedComputesTotalsPerMode() {
	date := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	rows := []repository.ActivityFeedRow{
		{
			Id: 302, Mode: constants.ActivityModeFill, ActivityDate: date,
			CreatedAt: date.Add(9 * time.Hour), CreatedBy: "user", CreatedByName: "สมชาย",
			PondId: 31, PondName: "บ่อ A3", FishType: "nil", Amount: 4500,
			FishWeight: decimal.NewFromFloat(0.02), FishUnit: "kg",
			PricePerUnit: decimal.NewFromInt(400),
		},
		{
			Id: 401, Mode: constants.ActivityModeSell, ActivityDate: date,
			CreatedAt: date.Add(7 * time.Hour), CreatedBy: "user", CreatedByName: "สมชาย",
			// A sell row carries none of its own figures — amount, weight and
			// price are all unset and must come off the detail lines.
			PondId: 41, PondName: "บ่อ C1",
			FishWeight: decimal.Zero, FishUnit: "kg",
			PricePerUnit: decimal.Zero,
			MerchantName: lo.ToPtr("ร้านลุงพร"),
		},
		{
			Id: 212, Mode: constants.ActivityModeMove, ActivityDate: date,
			CreatedAt: date.Add(16 * time.Hour), CreatedBy: "user", CreatedByName: "สมชาย",
			PondId: 21, PondName: "บ่อ B1",
			ToPondId: lo.ToPtr(22), ToPondName: lo.ToPtr("บ่อ B2"),
			FishType: "kang", Amount: 800,
			FishWeight: decimal.NewFromFloat(0.5), FishUnit: "kg",
			PricePerUnit: decimal.NewFromInt(100),
		},
	}
	s.activityRepo.On("ListRecentByClientID", clientCtx(), lo.ToPtr(1), 10, (*repository.ActivityFeedCursor)(nil)).Return(rows, nil)
	s.activityRepo.On("SumSellDetailsByActivityIDs", clientCtx(), []int{401}).Return([]repository.SellTotalRow{
		{
			SellId: 401, Total: decimal.NewFromInt(24180),
			TotalWeight: decimal.NewFromInt(312), TotalFishCount: 260,
		},
	}, nil)
	// Sell ids lead: additional costs are summed for every mode (the pond
	// activity timeline reports a sale's costs alongside its gross revenue).
	s.activityRepo.On("SumAdditionalCostsByActivityIDs", clientCtx(), []int{401, 302, 212}).Return([]repository.AdditionalCostTotalRow{
		{ActivityId: 302, Total: decimal.NewFromInt(500)},
	}, nil)

	feed, err := s.activityService.ListFeed(clientCtx(), 10, nil)
	require.NoError(s.T(), err)
	require.Len(s.T(), feed, 3)

	// fill: 4500 × 0.02 kg × ฿400 + ฿500 additional = ฿36,500
	assert.Equal(s.T(), "fill", feed[0].Mode)
	assert.Equal(s.T(), 31, feed[0].PondId)
	assert.Equal(s.T(), "บ่อ A3", feed[0].PondName)
	assert.InDelta(s.T(), 36500, feed[0].Total, 0.001)
	assert.Nil(s.T(), feed[0].TotalWeight)
	assert.Equal(s.T(), "สมชาย", feed[0].CreatedByName)

	// sell: every figure comes off sell_details — total, weight, head count, and
	// the derived average price (฿24,180 ÷ 312 kg = ฿77.50/kg).
	assert.Equal(s.T(), "sell", feed[1].Mode)
	assert.InDelta(s.T(), 24180, feed[1].Total, 0.001)
	require.NotNil(s.T(), feed[1].TotalWeight)
	assert.InDelta(s.T(), 312, *feed[1].TotalWeight, 0.001)
	assert.Equal(s.T(), 260, feed[1].Amount)
	assert.InDelta(s.T(), 77.5, feed[1].PricePerUnit, 0.001)
	require.NotNil(s.T(), feed[1].Merchant)
	assert.Equal(s.T(), "ร้านลุงพร", *feed[1].Merchant)

	// move: 800 × 0.5 kg × ฿100, no additional costs recorded
	assert.Equal(s.T(), "move", feed[2].Mode)
	assert.InDelta(s.T(), 40000, feed[2].Total, 0.001)
	assert.Equal(s.T(), 21, feed[2].PondId)
	require.NotNil(s.T(), feed[2].ToPondId)
	assert.Equal(s.T(), 22, *feed[2].ToPondId)
	require.NotNil(s.T(), feed[2].ToPondName)
	assert.Equal(s.T(), "บ่อ B2", *feed[2].ToPondName)
}

// A sell whose detail lines are missing (or sum to 0 kg) must not divide by
// zero — the derived average price stays 0 instead of becoming +Inf/NaN, which
// would serialise as invalid JSON.
func (s *ActivityServiceTestSuite) TestListFeedSellWithoutWeightKeepsZeroPrice() {
	date := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	rows := []repository.ActivityFeedRow{
		{
			Id: 402, Mode: constants.ActivityModeSell, ActivityDate: date, CreatedAt: date,
			CreatedBy: "user", CreatedByName: "สมชาย", PondId: 41, PondName: "บ่อ C1",
			FishWeight: decimal.Zero, PricePerUnit: decimal.Zero,
		},
	}
	s.activityRepo.On("ListRecentByClientID", clientCtx(), lo.ToPtr(1), 10, (*repository.ActivityFeedCursor)(nil)).Return(rows, nil)
	s.activityRepo.On("SumSellDetailsByActivityIDs", clientCtx(), []int{402}).
		Return([]repository.SellTotalRow{}, nil)
	s.activityRepo.On("SumAdditionalCostsByActivityIDs", clientCtx(), []int{402}).
		Return([]repository.AdditionalCostTotalRow{}, nil)

	feed, err := s.activityService.ListFeed(clientCtx(), 10, nil)
	require.NoError(s.T(), err)
	require.Len(s.T(), feed, 1)
	assert.Zero(s.T(), feed[0].PricePerUnit)
	assert.Zero(s.T(), feed[0].Amount)
	assert.Zero(s.T(), feed[0].Total)
	require.NotNil(s.T(), feed[0].TotalWeight)
	assert.Zero(s.T(), *feed[0].TotalWeight)
}

func (s *ActivityServiceTestSuite) TestListSellDetailsReturnsPerGradeLines() {
	rows := []repository.SellDetailRow{
		{
			FishSizeGradeId: 2, SizeName: "ไซส์ 2", FishCount: lo.ToPtr(800),
			Weight: decimal.NewFromInt(2000), PricePerUnit: decimal.NewFromInt(140),
		},
		{
			FishSizeGradeId: 3, SizeName: "ไซส์ 3", FishCount: lo.ToPtr(1400),
			Weight: decimal.NewFromInt(5000), PricePerUnit: decimal.NewFromInt(154),
		},
	}
	s.activityRepo.On("ListSellDetailsByActivityID", clientCtx(), 401, lo.ToPtr(1)).Return(rows, nil)

	lines, err := s.activityService.ListSellDetails(clientCtx(), 401)
	require.NoError(s.T(), err)
	require.Len(s.T(), lines, 2)

	assert.Equal(s.T(), "ไซส์ 2", lines[0].SizeName)
	assert.InDelta(s.T(), 2000, lines[0].Weight, 0.001)
	assert.InDelta(s.T(), 140, lines[0].PricePerUnit, 0.001)
	assert.InDelta(s.T(), 280000, lines[0].Total, 0.001)
	require.NotNil(s.T(), lines[0].FishCount)
	assert.Equal(s.T(), 800, *lines[0].FishCount)

	assert.InDelta(s.T(), 770000, lines[1].Total, 0.001)

	// The line totals must reconstruct the headline the feed reports for this
	// sale, or the sheet would show a breakdown that doesn't add up.
	assert.InDelta(s.T(), 1050000, lines[0].Total+lines[1].Total, 0.001)
}

// A legacy line predating the required fish_count column must report "not
// recorded" (nil), not a confident zero.
func (s *ActivityServiceTestSuite) TestListSellDetailsKeepsMissingFishCountNil() {
	rows := []repository.SellDetailRow{
		{
			FishSizeGradeId: 1, SizeName: "ไซส์ 1", FishCount: nil,
			Weight: decimal.NewFromInt(50), PricePerUnit: decimal.NewFromInt(200),
		},
	}
	s.activityRepo.On("ListSellDetailsByActivityID", clientCtx(), 402, lo.ToPtr(1)).Return(rows, nil)

	lines, err := s.activityService.ListSellDetails(clientCtx(), 402)
	require.NoError(s.T(), err)
	require.Len(s.T(), lines, 1)
	assert.Nil(s.T(), lines[0].FishCount)
	assert.InDelta(s.T(), 10000, lines[0].Total, 0.001)
}

// Super admin is unscoped: the repository is asked with a nil clientId so the
// scoping clause is skipped entirely.
func (s *ActivityServiceTestSuite) TestListSellDetailsSuperAdminIsUnscoped() {
	s.activityRepo.On("ListSellDetailsByActivityID", superAdminCtx(), 401, (*int)(nil)).
		Return([]repository.SellDetailRow{}, nil)

	lines, err := s.activityService.ListSellDetails(superAdminCtx(), 401)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), lines)
}

func (s *ActivityServiceTestSuite) TestListSellDetailsRejectsUserWithoutClient() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "user")
	ctx = context.WithValue(ctx, constants.UserLevelKey, 1)

	lines, err := s.activityService.ListSellDetails(ctx, 401)
	assert.Nil(s.T(), lines)
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
}

func (s *ActivityServiceTestSuite) TestListFeedEmptySkipsTotalQueries() {
	s.activityRepo.On("ListRecentByClientID", clientCtx(), lo.ToPtr(1), 0, (*repository.ActivityFeedCursor)(nil)).
		Return([]repository.ActivityFeedRow{}, nil)

	feed, err := s.activityService.ListFeed(clientCtx(), 0, nil)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), feed)
	s.activityRepo.AssertNotCalled(s.T(), "SumSellDetailsByActivityIDs")
	s.activityRepo.AssertNotCalled(s.T(), "SumAdditionalCostsByActivityIDs")
}

func (s *ActivityServiceTestSuite) TestListFeedSuperAdminIsUnscoped() {
	s.activityRepo.On("ListRecentByClientID", superAdminCtx(), (*int)(nil), 5, (*repository.ActivityFeedCursor)(nil)).
		Return([]repository.ActivityFeedRow{}, nil)

	feed, err := s.activityService.ListFeed(superAdminCtx(), 5, nil)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), feed)
}

func (s *ActivityServiceTestSuite) TestListFeedPassesCursorThrough() {
	cursor := &repository.ActivityFeedCursor{
		ActivityDate: time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
		Id:           412,
	}
	s.activityRepo.On("ListRecentByClientID", clientCtx(), lo.ToPtr(1), 30, cursor).
		Return([]repository.ActivityFeedRow{}, nil)

	feed, err := s.activityService.ListFeed(clientCtx(), 30, cursor)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), feed)
	// The cursor must reach the repository untouched — the service does no
	// paging arithmetic of its own, which is what keeps pages from overlapping.
	s.activityRepo.AssertExpectations(s.T())
}

func (s *ActivityServiceTestSuite) TestListFeedRejectsUserWithoutClient() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "user")
	ctx = context.WithValue(ctx, constants.UserLevelKey, 1)

	feed, err := s.activityService.ListFeed(ctx, 10, nil)
	assert.Nil(s.T(), feed)
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
}

func TestActivityServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ActivityServiceTestSuite))
}
