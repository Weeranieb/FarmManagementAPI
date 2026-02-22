package utils

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestFillCost(t *testing.T) {
	t.Run("amount and price only", func(t *testing.T) {
		got := FillCost(10, decimal.RequireFromString("5"), nil)
		assert.True(t, got.Equal(decimal.RequireFromString("50")))
	})
	t.Run("with additional costs", func(t *testing.T) {
		got := FillCost(2, decimal.RequireFromString("100"), []decimal.Decimal{
			decimal.RequireFromString("10"),
			decimal.RequireFromString("25"),
		})
		assert.True(t, got.Equal(decimal.RequireFromString("235"))) // 200 + 10 + 25
	})
}
