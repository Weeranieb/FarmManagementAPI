package utils

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
)

func TestFillCost(t *testing.T) {
	t.Run("amount and price only", func(t *testing.T) {
		// GIVEN — amount 10, price 5, no additional costs
		// WHEN — CalculateFillCost is called
		got := CalculateFillCost(10, decimal.RequireFromString("5"), nil)
		// THEN — total is 50
		assert.True(t, got.Equal(decimal.RequireFromString("50")))
	})
	t.Run("with additional costs", func(t *testing.T) {
		// GIVEN — amount 2, price 100, two additional costs 10 and 25
		// WHEN — CalculateFillCost is called
		got := CalculateFillCost(2, decimal.RequireFromString("100"), []dto.AdditionalCostItem{
			{Cost: decimal.RequireFromString("10")},
			{Cost: decimal.RequireFromString("25")},
		})
		// THEN — total is 235 (200 + 10 + 25)
		assert.True(t, got.Equal(decimal.RequireFromString("235")))
	})
}

func TestCalculateMoveCost(t *testing.T) {
	t.Run("amount, price and weight only", func(t *testing.T) {
		// GIVEN — amount 10, price 5, weight 2, no additional costs
		// WHEN — CalculateMoveCost is called
		fishCost, additionalCost := CalculateMoveCost(
			10,
			decimal.RequireFromString("5"),
			decimal.RequireFromString("2"),
			nil,
		)
		// THEN — fishCost 100, additionalCost 0
		assert.True(t, fishCost.Equal(decimal.RequireFromString("100")), "fishCost: 10 * 2 * 5 = 100")
		assert.True(t, additionalCost.Equal(decimal.Zero))
	})
	t.Run("with additional costs", func(t *testing.T) {
		// GIVEN — amount 3, price 20, weight 0.5, two additional costs
		// WHEN — CalculateMoveCost is called
		fishCost, additionalCost := CalculateMoveCost(
			3,
			decimal.RequireFromString("20"),
			decimal.RequireFromString("0.5"),
			[]dto.AdditionalCostItem{
				{Cost: decimal.RequireFromString("15")},
				{Cost: decimal.RequireFromString("5")},
			},
		)
		// THEN — fishCost 30, additionalCost 20
		assert.True(t, fishCost.Equal(decimal.RequireFromString("30")), "fishCost: 3 * 0.5 * 20 = 30")
		assert.True(t, additionalCost.Equal(decimal.RequireFromString("20")), "additionalCost: 15 + 5 = 20")
	})
}

func TestCalculateAdditionalCostsTotal(t *testing.T) {
	t.Run("nil returns zero", func(t *testing.T) {
		got := CalculateAdditionalCostsTotal(nil)
		assert.True(t, got.Equal(decimal.Zero))
	})
	t.Run("empty returns zero", func(t *testing.T) {
		got := CalculateAdditionalCostsTotal([]dto.AdditionalCostItem{})
		assert.True(t, got.Equal(decimal.Zero))
	})
	t.Run("sums costs with titles", func(t *testing.T) {
		got := CalculateAdditionalCostsTotal([]dto.AdditionalCostItem{
			{Title: "transport", Cost: decimal.RequireFromString("100.50")},
			{Title: "labour", Cost: decimal.RequireFromString("50.25")},
		})
		assert.True(t, got.Equal(decimal.RequireFromString("150.75")))
	})
}

func TestCalculateSellRevenue(t *testing.T) {
	t.Run("empty details returns zero", func(t *testing.T) {
		got := CalculateSellRevenue(nil)
		assert.True(t, got.Equal(decimal.Zero))
	})
	t.Run("single line weight * pricePerUnit", func(t *testing.T) {
		got := CalculateSellRevenue([]dto.PondSellDetailItem{
			{FishSizeGradeId: 1, Weight: decimal.RequireFromString("6.5"), PricePerUnit: decimal.RequireFromString("240")},
		})
		assert.True(t, got.Equal(decimal.RequireFromString("1560")), "6.5 * 240 = 1560")
	})
	t.Run("multiple lines sum subtotals", func(t *testing.T) {
		got := CalculateSellRevenue([]dto.PondSellDetailItem{
			{FishSizeGradeId: 1, Weight: decimal.RequireFromString("10"), PricePerUnit: decimal.RequireFromString("20")},
			{FishSizeGradeId: 2, Weight: decimal.RequireFromString("5"), PricePerUnit: decimal.RequireFromString("40")},
		})
		assert.True(t, got.Equal(decimal.RequireFromString("400")), "200 + 200 = 400")
	})
}

func TestCalculateSellTotals(t *testing.T) {
	t.Run("revenue and additional cost total", func(t *testing.T) {
		revenue, addCost := CalculateSellTotals(
			[]dto.PondSellDetailItem{
				{FishSizeGradeId: 1, Weight: decimal.RequireFromString("100"), PricePerUnit: decimal.RequireFromString("2")},
			},
			[]dto.AdditionalCostItem{{Title: "fee", Cost: decimal.RequireFromString("10")}},
		)
		assert.True(t, revenue.Equal(decimal.RequireFromString("200")), "100 * 2 = 200")
		assert.True(t, addCost.Equal(decimal.RequireFromString("10")))
	})
}

func TestCalculateSellDetailLines(t *testing.T) {
	t.Run("empty details returns empty slice", func(t *testing.T) {
		got := CalculateSellDetailLines(nil)
		assert.Empty(t, got)
	})
	t.Run("single detail line", func(t *testing.T) {
		got := CalculateSellDetailLines([]dto.PondSellDetailItem{
			{FishSizeGradeId: 1, Weight: decimal.RequireFromString("6.5"), PricePerUnit: decimal.RequireFromString("240")},
		})
		require.Len(t, got, 1)
		assert.Equal(t, 1, got[0].FishSizeGradeId)
		assert.Equal(t, 6.5, got[0].Weight)
		assert.Equal(t, 240.0, got[0].PricePerUnit)
		assert.Equal(t, 1560.0, got[0].Subtotal)
	})
	t.Run("multiple lines match CalculateSellRevenue sum", func(t *testing.T) {
		details := []dto.PondSellDetailItem{
			{FishSizeGradeId: 1, Weight: decimal.RequireFromString("10"), PricePerUnit: decimal.RequireFromString("2.5")},
			{FishSizeGradeId: 2, Weight: decimal.RequireFromString("4"), PricePerUnit: decimal.RequireFromString("10")},
		}
		lines := CalculateSellDetailLines(details)
		require.Len(t, lines, 2)
		assert.Equal(t, 25.0, lines[0].Subtotal, "10 * 2.5")
		assert.Equal(t, 40.0, lines[1].Subtotal, "4 * 10")
		revenue := CalculateSellRevenue(details)
		assert.True(t, revenue.Equal(decimal.RequireFromString("65")), "25 + 40 = 65")
	})
}
