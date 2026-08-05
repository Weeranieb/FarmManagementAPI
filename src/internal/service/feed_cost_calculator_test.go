package service

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
)

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// TestCalcCycleFeedCost_WithFallback verifies whole-cycle feed cost aggregation:
// fresh + pellet priced per day, and a log dated before any recorded price
// falls back to the nearest (earliest) available price instead of costing 0.
func TestCalcCycleFeedCost_WithFallback(t *testing.T) {
	dailyLogRepo := mocks.NewMockDailyLogRepository(t)
	priceRepo := mocks.NewMockFeedPriceHistoryRepository(t)
	calc := NewFeedCostCalculator(dailyLogRepo, priceRepo)

	freshFc, pelletFc := 10, 20
	ap := &model.ActivePond{
		Id:                     1,
		FreshFeedCollectionId:  &freshFc,
		PelletFeedCollectionId: &pelletFc,
	}

	logs := []*model.DailyLog{
		// In-effect day: fresh 100×5 + pellet (50+50)×2 = 500 + 200 = 700
		{ActivePondId: 1, FeedDate: day(2026, 3, 1),
			Fresh: decimal.NewFromInt(100), PelletMorning: decimal.NewFromInt(50), PelletEvening: decimal.NewFromInt(50)},
		// Before any price exists: fresh 10 × fallback(5) = 50, pellet 0
		{ActivePondId: 1, FeedDate: day(2026, 2, 1),
			Fresh: decimal.NewFromInt(10)},
	}
	histories := []*model.FeedPriceHistory{
		// Fresh: priced per ลัง (Price drives cost).
		{FeedCollectionId: freshFc, Price: decimal.NewFromInt(5), PriceUpdatedDate: day(2026, 3, 1)},
		// Pellet: fed and priced per ถุง — Price drives cost, same as fresh.
		{FeedCollectionId: pelletFc, Price: decimal.NewFromInt(2), PriceUpdatedDate: day(2026, 3, 1)},
	}

	dailyLogRepo.On("ListByActivePondIds", mock.Anything, mock.Anything).Return(logs, nil)
	priceRepo.On("ListByFeedCollectionIds", mock.Anything).Return(histories, nil)

	got, err := calc.CalcCycleFeedCost(context.Background(), ap)
	require.NoError(t, err)
	assert.True(t, decimal.NewFromInt(750).Equal(got), "expected 750, got %s", got)
}

// TestCalcCycleFeedCost_MissingPriceHistory verifies that feed logged against a
// collection with no price history at all is counted as zero (and does not error).
func TestCalcCycleFeedCost_MissingPriceHistory(t *testing.T) {
	dailyLogRepo := mocks.NewMockDailyLogRepository(t)
	priceRepo := mocks.NewMockFeedPriceHistoryRepository(t)
	calc := NewFeedCostCalculator(dailyLogRepo, priceRepo)

	freshFc := 10
	ap := &model.ActivePond{Id: 1, FreshFeedCollectionId: &freshFc}
	logs := []*model.DailyLog{
		{ActivePondId: 1, FeedDate: day(2026, 3, 1), Fresh: decimal.NewFromInt(100)},
		{ActivePondId: 1, FeedDate: day(2026, 3, 2), Fresh: decimal.NewFromInt(80)},
	}
	dailyLogRepo.On("ListByActivePondIds", mock.Anything, mock.Anything).Return(logs, nil)
	priceRepo.On("ListByFeedCollectionIds", mock.Anything).Return([]*model.FeedPriceHistory{}, nil)

	got, err := calc.CalcCycleFeedCost(context.Background(), ap)
	require.NoError(t, err)
	assert.True(t, got.IsZero(), "no price history → feed cost 0, got %s", got)
}

// TestCalcCycleFeedCostBatch_NoLogs verifies every input pond is present in the
// result (zero) even when there are no daily logs, and that a nil pond is skipped.
func TestCalcCycleFeedCostBatch_NoLogs(t *testing.T) {
	dailyLogRepo := mocks.NewMockDailyLogRepository(t)
	priceRepo := mocks.NewMockFeedPriceHistoryRepository(t)
	calc := NewFeedCostCalculator(dailyLogRepo, priceRepo)

	fresh := 10
	aps := []*model.ActivePond{
		{Id: 1, FreshFeedCollectionId: &fresh},
		{Id: 2},
		nil,
	}
	dailyLogRepo.On("ListByActivePondIds", mock.Anything, mock.Anything).Return([]*model.DailyLog{}, nil)

	got, err := calc.CalcCycleFeedCostBatch(context.Background(), aps)
	require.NoError(t, err)
	assert.True(t, got[1].IsZero())
	assert.True(t, got[2].IsZero())
	_, hasNil := got[0]
	assert.False(t, hasNil, "nil pond should not create a result entry")
}
