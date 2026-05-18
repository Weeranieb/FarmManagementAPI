package service

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
)

// newCalcOnlyService builds a pondService with nil dependencies. The Calc*
// methods are pure math and never touch repos or the DB, so leaving the
// dependencies nil is fine and lets these tests run without the full suite.
func newCalcOnlyService() PondService {
	return NewPondService(PondServiceParams{})
}

func TestCalcFillPond_FillFormulaIncludesWeight(t *testing.T) {
	// GIVEN — 2000 fish at ฿70/kg, 1.1 kg avg, no additional costs.
	// Fill price is per-kilogram, so weight is part of the formula.
	svc := newCalcOnlyService()
	resp := svc.CalcFillPond(context.Background(), dto.PondFillCalcRequest{
		Amount:       2000,
		FishWeight:   decimal.RequireFromString("1.1"),
		PricePerUnit: decimal.RequireFromString("70"),
	})

	// WHEN/THEN — totalCost = 2000 × 1.1 × 70 = 154000
	assert.Equal(t, 2000, resp.Quantity)
	assert.InDelta(t, 1.1, resp.AvgWeightKg, 1e-9)
	assert.InDelta(t, 2200.0, resp.TotalWeight, 1e-9)
	assert.InDelta(t, 70.0, resp.CostPerUnit, 1e-9)
	assert.InDelta(t, 154000.0, resp.BaseStockCost, 1e-9)
	assert.InDelta(t, 0.0, resp.AdditionalCostsTotal, 1e-9)
	assert.InDelta(t, 154000.0, resp.TotalCost, 1e-9)
}

func TestCalcFillPond_WithAdditionalCosts(t *testing.T) {
	svc := newCalcOnlyService()
	resp := svc.CalcFillPond(context.Background(), dto.PondFillCalcRequest{
		Amount:       100,
		FishWeight:   decimal.RequireFromString("1"),
		PricePerUnit: decimal.RequireFromString("10"),
		AdditionalCosts: []dto.AdditionalCostCalcItem{
			{Title: "ค่าขนส่ง", Cost: decimal.RequireFromString("500")},
			{Title: "ค่ายา", Cost: decimal.RequireFromString("200")},
		},
	})

	assert.InDelta(t, 1000.0, resp.BaseStockCost, 1e-9) // 100 × 1 × 10
	assert.InDelta(t, 700.0, resp.AdditionalCostsTotal, 1e-9)
	assert.InDelta(t, 1700.0, resp.TotalCost, 1e-9)
	assert.Len(t, resp.AdditionalCosts, 2)
}

func TestCalcFillPond_ZeroAmountReturnsZeros(t *testing.T) {
	svc := newCalcOnlyService()
	resp := svc.CalcFillPond(context.Background(), dto.PondFillCalcRequest{})
	assert.Equal(t, 0, resp.Quantity)
	assert.InDelta(t, 0.0, resp.TotalCost, 1e-9)
	assert.InDelta(t, 0.0, resp.TotalWeight, 1e-9)
}

func TestCalcMovePond_MoveFormulaIncludesWeight(t *testing.T) {
	// GIVEN — move uses amount × weight × pricePerUnit (price is per-kg)
	svc := newCalcOnlyService()
	resp := svc.CalcMovePond(context.Background(), dto.PondMoveCalcRequest{
		Amount:       2000,
		FishWeight:   decimal.RequireFromString("1.1"),
		PricePerUnit: decimal.RequireFromString("70"),
	})

	// WHEN/THEN — base transfer cost = 2000 × 1.1 × 70 = 154000
	assert.InDelta(t, 2200.0, resp.TotalWeight, 1e-9)
	assert.InDelta(t, 154000.0, resp.BaseTransferCost, 1e-9)
	assert.InDelta(t, 154000.0, resp.TotalCost, 1e-9)
}

func TestCalcMovePond_WithAdditionalCosts(t *testing.T) {
	svc := newCalcOnlyService()
	resp := svc.CalcMovePond(context.Background(), dto.PondMoveCalcRequest{
		Amount:       10,
		FishWeight:   decimal.RequireFromString("2"),
		PricePerUnit: decimal.RequireFromString("50"),
		AdditionalCosts: []dto.AdditionalCostCalcItem{
			{Title: "fuel", Cost: decimal.RequireFromString("300")},
		},
	})
	assert.InDelta(t, 1000.0, resp.BaseTransferCost, 1e-9)
	assert.InDelta(t, 300.0, resp.AdditionalCostsTotal, 1e-9)
	assert.InDelta(t, 1300.0, resp.TotalCost, 1e-9)
}

func TestCalcSellPond_SumsLineRevenue(t *testing.T) {
	svc := newCalcOnlyService()
	resp := svc.CalcSellPond(context.Background(), dto.PondSellCalcRequest{
		Details: []dto.PondSellCalcDetailItem{
			{FishSizeGradeId: 1, Weight: decimal.RequireFromString("100"), PricePerUnit: decimal.RequireFromString("80")},
			{FishSizeGradeId: 2, Weight: decimal.RequireFromString("50"), PricePerUnit: decimal.RequireFromString("100")},
		},
		AdditionalCosts: []dto.AdditionalCostCalcItem{
			{Title: "ice", Cost: decimal.RequireFromString("200")},
		},
	})

	assert.InDelta(t, 150.0, resp.TotalWeight, 1e-9)
	assert.InDelta(t, 13000.0, resp.TotalRevenue, 1e-9) // 100*80 + 50*100
	assert.InDelta(t, 200.0, resp.AdditionalCostsTotal, 1e-9)
	assert.InDelta(t, 12800.0, resp.NetTotal, 1e-9)
	assert.Len(t, resp.Items, 2)
	assert.InDelta(t, 8000.0, resp.Items[0].Subtotal, 1e-9)
	assert.InDelta(t, 5000.0, resp.Items[1].Subtotal, 1e-9)
}

func TestCalcSellPond_EmptyDetailsReturnsZeros(t *testing.T) {
	svc := newCalcOnlyService()
	resp := svc.CalcSellPond(context.Background(), dto.PondSellCalcRequest{})
	assert.Empty(t, resp.Items)
	assert.InDelta(t, 0.0, resp.TotalRevenue, 1e-9)
	assert.InDelta(t, 0.0, resp.NetTotal, 1e-9)
}
