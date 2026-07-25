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
			PondName: "บ่อ A3", FishType: "nil", Amount: 4500,
			FishWeight: decimal.NewFromFloat(0.02), FishUnit: "kg",
			PricePerUnit: decimal.NewFromInt(400),
		},
		{
			Id: 401, Mode: constants.ActivityModeSell, ActivityDate: date,
			CreatedAt: date.Add(7 * time.Hour), CreatedBy: "user", CreatedByName: "สมชาย",
			PondName: "บ่อ C1", FishType: "nil", Amount: 312,
			FishWeight: decimal.Zero, FishUnit: "kg",
			PricePerUnit: decimal.Zero,
			MerchantName: lo.ToPtr("ร้านลุงพร"),
		},
		{
			Id: 212, Mode: constants.ActivityModeMove, ActivityDate: date,
			CreatedAt: date.Add(16 * time.Hour), CreatedBy: "user", CreatedByName: "สมชาย",
			PondName: "บ่อ B1", ToPondName: lo.ToPtr("บ่อ B2"), FishType: "kang", Amount: 800,
			FishWeight: decimal.NewFromFloat(0.5), FishUnit: "kg",
			PricePerUnit: decimal.NewFromInt(100),
		},
	}
	s.activityRepo.On("ListRecentByClientID", clientCtx(), lo.ToPtr(1), 10).Return(rows, nil)
	s.activityRepo.On("SumSellDetailsByActivityIDs", clientCtx(), []int{401}).Return([]repository.SellTotalRow{
		{SellId: 401, Total: decimal.NewFromInt(24180), TotalWeight: decimal.NewFromInt(312)},
	}, nil)
	// Sell ids lead: additional costs are summed for every mode (the pond
	// activity timeline reports a sale's costs alongside its gross revenue).
	s.activityRepo.On("SumAdditionalCostsByActivityIDs", clientCtx(), []int{401, 302, 212}).Return([]repository.AdditionalCostTotalRow{
		{ActivityId: 302, Total: decimal.NewFromInt(500)},
	}, nil)

	feed, err := s.activityService.ListFeed(clientCtx(), 10)
	require.NoError(s.T(), err)
	require.Len(s.T(), feed, 3)

	// fill: 4500 × 0.02 kg × ฿400 + ฿500 additional = ฿36,500
	assert.Equal(s.T(), "fill", feed[0].Mode)
	assert.Equal(s.T(), "บ่อ A3", feed[0].PondName)
	assert.InDelta(s.T(), 36500, feed[0].Total, 0.001)
	assert.Nil(s.T(), feed[0].TotalWeight)
	assert.Equal(s.T(), "สมชาย", feed[0].CreatedByName)

	// sell: totals come straight from sell_details
	assert.Equal(s.T(), "sell", feed[1].Mode)
	assert.InDelta(s.T(), 24180, feed[1].Total, 0.001)
	require.NotNil(s.T(), feed[1].TotalWeight)
	assert.InDelta(s.T(), 312, *feed[1].TotalWeight, 0.001)
	require.NotNil(s.T(), feed[1].Merchant)
	assert.Equal(s.T(), "ร้านลุงพร", *feed[1].Merchant)

	// move: 800 × 0.5 kg × ฿100, no additional costs recorded
	assert.Equal(s.T(), "move", feed[2].Mode)
	assert.InDelta(s.T(), 40000, feed[2].Total, 0.001)
	require.NotNil(s.T(), feed[2].ToPondName)
	assert.Equal(s.T(), "บ่อ B2", *feed[2].ToPondName)
}

func (s *ActivityServiceTestSuite) TestListFeedEmptySkipsTotalQueries() {
	s.activityRepo.On("ListRecentByClientID", clientCtx(), lo.ToPtr(1), 0).
		Return([]repository.ActivityFeedRow{}, nil)

	feed, err := s.activityService.ListFeed(clientCtx(), 0)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), feed)
	s.activityRepo.AssertNotCalled(s.T(), "SumSellDetailsByActivityIDs")
	s.activityRepo.AssertNotCalled(s.T(), "SumAdditionalCostsByActivityIDs")
}

func (s *ActivityServiceTestSuite) TestListFeedSuperAdminIsUnscoped() {
	s.activityRepo.On("ListRecentByClientID", superAdminCtx(), (*int)(nil), 5).
		Return([]repository.ActivityFeedRow{}, nil)

	feed, err := s.activityService.ListFeed(superAdminCtx(), 5)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), feed)
}

func (s *ActivityServiceTestSuite) TestListFeedRejectsUserWithoutClient() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "user")
	ctx = context.WithValue(ctx, constants.UserLevelKey, 1)

	feed, err := s.activityService.ListFeed(ctx, 10)
	assert.Nil(s.T(), feed)
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
}

func TestActivityServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ActivityServiceTestSuite))
}
