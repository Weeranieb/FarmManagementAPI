package utils

import (
	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
)

// CalculateFillCost returns cost for a single fill: amount × pricePerUnit + sum(additionalCosts).
func CalculateFillCost(amount int, pricePerUnit decimal.Decimal, additionalCosts []dto.AdditionalCostItem) decimal.Decimal {
	fishCost := decimal.NewFromInt(int64(amount)).Mul(pricePerUnit)
	fishCost = fishCost.Add(CalculateAdditionalCostsTotal(additionalCosts))
	return fishCost
}

func CalculateAdditionalCostsTotal(additionalCosts []dto.AdditionalCostItem) decimal.Decimal {
	total := decimal.Zero
	for _, c := range additionalCosts {
		total = total.Add(c.Cost)
	}
	return total
}

func CalculateMoveCost(amount int, pricePerUnit, fishWeight decimal.Decimal, additionalCosts []dto.AdditionalCostItem) (fishCost, additionalCost decimal.Decimal) {
	fishCost = decimal.NewFromInt(int64(amount)).Mul(fishWeight).Mul(pricePerUnit)
	additionalCost = CalculateAdditionalCostsTotal(additionalCosts)
	return fishCost, additionalCost
}
