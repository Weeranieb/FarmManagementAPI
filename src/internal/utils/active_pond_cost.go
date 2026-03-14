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

// CalculateSellRevenue sums amount * pricePerUnit for each sell detail line.
func CalculateSellRevenue(details []dto.PondSellDetailItem) decimal.Decimal {
	total := decimal.Zero
	for _, d := range details {
		total = total.Add(d.Amount.Mul(d.PricePerUnit))
	}
	return total
}

// CalculateSellTotals returns revenue from sell details and total of additional costs.
func CalculateSellTotals(details []dto.PondSellDetailItem, additionalCosts []dto.AdditionalCostItem) (revenue, additionalCostTotal decimal.Decimal) {
	revenue = CalculateSellRevenue(details)
	additionalCostTotal = CalculateAdditionalCostsTotal(additionalCosts)
	return revenue, additionalCostTotal
}

// SellDetailLine holds per-line sell calculation (same math as CalculateSellRevenue).
// Used by PreviewSellPond to build itemized breakdown.
type SellDetailLine struct {
	FishType     string
	Amount       float64
	PricePerUnit float64
	Subtotal     float64
}

// CalculateSellDetailLines returns per-line breakdown using the same logic as CalculateSellRevenue.
func CalculateSellDetailLines(details []dto.PondSellDetailItem) []SellDetailLine {
	lines := make([]SellDetailLine, 0, len(details))
	for _, d := range details {
		amt, _ := d.Amount.Float64()
		ppu, _ := d.PricePerUnit.Float64()
		lines = append(lines, SellDetailLine{
			FishType:     d.FishType,
			Amount:       amt,
			PricePerUnit: ppu,
			Subtotal:     amt * ppu,
		})
	}
	return lines
}
