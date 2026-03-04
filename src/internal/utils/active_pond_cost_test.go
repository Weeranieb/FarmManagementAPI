package utils

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
)

func TestFillCost(t *testing.T) {
	t.Run("amount and price only", func(t *testing.T) {
		got := CalculateFillCost(10, decimal.RequireFromString("5"), nil)
		assert.True(t, got.Equal(decimal.RequireFromString("50")))
	})
	t.Run("with additional costs", func(t *testing.T) {
		got := CalculateFillCost(2, decimal.RequireFromString("100"), []dto.AdditionalCostItem{
			{Cost: decimal.RequireFromString("10")},
			{Cost: decimal.RequireFromString("25")},
		})
		assert.True(t, got.Equal(decimal.RequireFromString("235"))) // 200 + 10 + 25
	})
}

func TestCalculateMoveCost(t *testing.T) {
	t.Run("amount, price and weight only", func(t *testing.T) {
		fishCost, additionalCost := CalculateMoveCost(
			10,
			decimal.RequireFromString("5"),
			decimal.RequireFromString("2"),
			nil,
		)
		assert.True(t, fishCost.Equal(decimal.RequireFromString("100")), "fishCost: 10 * 2 * 5 = 100")
		assert.True(t, additionalCost.Equal(decimal.Zero))
	})
	t.Run("with additional costs", func(t *testing.T) {
		fishCost, additionalCost := CalculateMoveCost(
			3,
			decimal.RequireFromString("20"),
			decimal.RequireFromString("0.5"),
			[]dto.AdditionalCostItem{
				{Cost: decimal.RequireFromString("15")},
				{Cost: decimal.RequireFromString("5")},
			},
		)
		assert.True(t, fishCost.Equal(decimal.RequireFromString("30")), "fishCost: 3 * 0.5 * 20 = 30")
		assert.True(t, additionalCost.Equal(decimal.RequireFromString("20")), "additionalCost: 15 + 5 = 20")
	})
}
