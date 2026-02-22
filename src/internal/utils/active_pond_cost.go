package utils

import (
	"github.com/shopspring/decimal"
)

// FillCost returns cost for a single fill: amount × pricePerUnit + sum(additionalCosts).
func FillCost(amount int, pricePerUnit decimal.Decimal, additionalCosts []decimal.Decimal) decimal.Decimal {
	total := decimal.NewFromInt(int64(amount)).Mul(pricePerUnit)
	for _, c := range additionalCosts {
		total = total.Add(c)
	}
	return total
}
