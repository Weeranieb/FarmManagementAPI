package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreatePondItem is one pond entry in a bulk create request.
type CreatePondItem struct {
	Name string           `json:"name" validate:"required"`
	Area *decimal.Decimal `json:"area,omitempty" validate:"omitempty,decimal_gte0" swaggertype:"number"`
}

// BulkImportPondItem is one pond entry inside a farm in a bulk import.
type BulkImportPondItem struct {
	Name string           `json:"name" validate:"required,max=100"`
	Area *decimal.Decimal `json:"area,omitempty" validate:"omitempty,decimal_gte0" swaggertype:"number"`
}

// BulkImportFarmItem is one farm grouping in a bulk import.
type BulkImportFarmItem struct {
	Name  string               `json:"name" validate:"required,max=100"`
	Ponds []BulkImportPondItem `json:"ponds" validate:"required,min=1,dive"`
}

// BulkImportFarmPondRequest is the body for POST /pond/bulk-import/:clientId.
// Idempotent: missing farms/ponds are created, existing ponds get their area
// updated (when area is provided). Nothing is ever deleted.
type BulkImportFarmPondRequest struct {
	Farms []BulkImportFarmItem `json:"farms" validate:"required,min=1,max=1000,dive"`
}

// BulkImportFarmResult is the per-farm summary returned to the client.
// PondsUpdated counts existing ponds where an area was provided and applied.
// PondsUnchanged counts existing ponds that matched by name but had no area
// to apply (so no DB write happened for them).
type BulkImportFarmResult struct {
	Name           string `json:"name"`
	IsNew          bool   `json:"isNew"`
	PondsCreated   int    `json:"pondsCreated"`
	PondsUpdated   int    `json:"pondsUpdated"`
	PondsUnchanged int    `json:"pondsUnchanged"`
}

// BulkImportFarmPondResponse is the summary returned after a bulk import.
// PondsUpdated counts only ponds where an area was actually written; ponds
// that matched by name but came in with no area are reported separately in
// PondsUnchanged so the UI doesn't claim a change that didn't happen.
type BulkImportFarmPondResponse struct {
	FarmsCreated   int                    `json:"farmsCreated"`
	FarmsExisting  int                    `json:"farmsExisting"`
	PondsCreated   int                    `json:"pondsCreated"`
	PondsUpdated   int                    `json:"pondsUpdated"`
	PondsUnchanged int                    `json:"pondsUnchanged"`
	Farms          []BulkImportFarmResult `json:"farms"`
}

// CreatePondsRequest is the body for POST /pond (create multiple ponds for a farm). New ponds are created with status maintenance.
type CreatePondsRequest struct {
	FarmId int              `json:"farmId" validate:"required"`
	Ponds  []CreatePondItem `json:"ponds" validate:"required,min=1,dive"`
}

// UpdatePondRequest is used by the service layer (id comes from path).
type UpdatePondRequest struct {
	Id     int              `json:"-"` // from path
	FarmId int              `json:"farmId"`
	Name   string           `json:"name"`
	Status string           `json:"status" validate:"omitempty,oneof=active maintenance"`
	Area   *decimal.Decimal `json:"area,omitempty" validate:"omitempty,decimal_gte0" swaggertype:"number"`
}

// UpdatePondBody is the request body for PUT /pond/:id (id in path).
type UpdatePondBody struct {
	FarmId int              `json:"farmId"`
	Name   string           `json:"name"`
	Status string           `json:"status" validate:"omitempty,oneof=active maintenance"`
	Area   *decimal.Decimal `json:"area,omitempty" validate:"omitempty,decimal_gte0" swaggertype:"number"`
}

type PondResponse struct {
	Id                 int              `json:"id"`
	FarmId             int              `json:"farmId"`
	Name               string           `json:"name"`
	TotalFish          *int             `json:"totalFish"`
	Status             string           `json:"status"`
	Area               *decimal.Decimal `json:"area,omitempty" swaggertype:"number"`
	FishTypes          []string         `json:"fishTypes"`
	AgeDays            *int             `json:"ageDays"`
	StartDate          *time.Time       `json:"startDate"`
	LatestActivityDate *time.Time       `json:"latestActivityDate"`
	LatestActivityType *string          `json:"latestActivityType"`
	CreatedAt          time.Time        `json:"createdAt"`
	CreatedBy          string           `json:"createdBy"`
	UpdatedAt          time.Time        `json:"updatedAt"`
	UpdatedBy          string           `json:"updatedBy"`
}

// AdditionalCostItem represents a single additional cost with a title and amount.
type AdditionalCostItem struct {
	Title string          `json:"title" validate:"required"`
	Cost  decimal.Decimal `json:"cost" validate:"required,decimal_gte0" swaggertype:"number"`
}

// PondFillRequest is the body for POST /pond/:pondId/fill (add fish to pond).
type PondFillRequest struct {
	FishType        string               `json:"fishType" validate:"required"`
	Amount          int                  `json:"amount" validate:"required,min=1"`
	FishWeight      decimal.Decimal      `json:"fishWeight,omitempty" validate:"omitempty,decimal_gt0" swaggertype:"number"`
	PricePerUnit    decimal.Decimal      `json:"pricePerUnit" validate:"required,decimal_gt0" swaggertype:"number"`
	AdditionalCosts []AdditionalCostItem `json:"additionalCosts,omitempty" validate:"dive"`
	ActivityDate    string               `json:"activityDate" validate:"required"`
	Remark          *string              `json:"remark,omitempty"`
}

// PondFillResponse is the response for POST /pond/:pondId/fill.
type PondFillResponse struct {
	ActivityId   int64 `json:"activityId"`
	ActivePondId int64 `json:"activePondId"`
}

// PondMoveRequest is the body for POST /pond/:pondId/move (transfer fish to another pond).
type PondMoveRequest struct {
	ToPondId        int                  `json:"toPondId" validate:"required"`
	FishType        string               `json:"fishType" validate:"required"`
	Amount          int                  `json:"amount" validate:"required,min=1"`
	FishWeight      decimal.Decimal      `json:"fishWeight,omitempty" validate:"omitempty,decimal_gte0" swaggertype:"number"`
	PricePerUnit    decimal.Decimal      `json:"pricePerUnit" validate:"required,decimal_gt0" swaggertype:"number"`
	AdditionalCosts []AdditionalCostItem `json:"additionalCosts,omitempty" validate:"dive"`
	ActivityDate    string               `json:"activityDate" validate:"required"`
	Remark          *string              `json:"remark,omitempty"`
	MarkToClose     bool                 `json:"markToClose"`
}

// PondMoveResponse is the response for POST /pond/:pondId/move.
type PondMoveResponse struct {
	ActivityId     int64 `json:"activityId"`
	ActivePondId   int64 `json:"activePondId"`
	ToActivePondId int64 `json:"toActivePondId"`
}

// PondSellDetailItem represents a single fish-size-grade line in a sell request.
type PondSellDetailItem struct {
	FishSizeGradeId int             `json:"fishSizeGradeId" validate:"required"`
	Weight          decimal.Decimal `json:"weight" validate:"required,decimal_gt0" swaggertype:"number"`
	PricePerUnit    decimal.Decimal `json:"pricePerUnit" validate:"required,decimal_gt0" swaggertype:"number"`
	FishCount       *int            `json:"fishCount,omitempty"`
}

// PondSellRequest is the body for POST /pond/:pondId/sell.
type PondSellRequest struct {
	ActivityDate    string               `json:"activityDate" validate:"required"`
	Details         []PondSellDetailItem `json:"details" validate:"required,min=1,dive"`
	MerchantId      *int                 `json:"merchantId,omitempty"`
	MarkToClose     bool                 `json:"markToClose"`
	AdditionalCosts []AdditionalCostItem `json:"additionalCosts,omitempty" validate:"dive"`
}

// PondSellResponse is the response for POST /pond/:pondId/sell.
type PondSellResponse struct {
	ActivityId   int64 `json:"activityId"`
	ActivePondId int64 `json:"activePondId"`
}

// --- Preview (Review & Confirm) DTOs ---

// AdditionalCostLine is a single row in the additional-costs summary.
type AdditionalCostLine struct {
	Title string  `json:"title"`
	Cost  float64 `json:"cost"`
}

// PondFillPreviewResponse is returned by POST /pond/:pondId/fill/preview.
type PondFillPreviewResponse struct {
	Valid           bool                 `json:"valid"`
	Species         string               `json:"species"`
	Quantity        int                  `json:"quantity"`
	AvgWeightKg     float64              `json:"avgWeightKg"`
	TotalWeight     float64              `json:"totalWeight"`
	CostPerUnit     float64              `json:"costPerUnit"`
	BaseStockCost   float64              `json:"baseStockCost"`
	AdditionalCosts []AdditionalCostLine `json:"additionalCosts"`
	TotalCost       float64              `json:"totalCost"`
	StockBefore     int                  `json:"stockBefore"`
	StockAfter      int                  `json:"stockAfter"`
	StockDelta      int                  `json:"stockDelta"`
	ValidationError string               `json:"validationError,omitempty"`
}

// PondMovePreviewResponse is returned by POST /pond/:pondId/move/preview.
type PondMovePreviewResponse struct {
	Valid            bool                 `json:"valid"`
	Species          string               `json:"species"`
	Quantity         int                  `json:"quantity"`
	AvgWeightKg      float64              `json:"avgWeightKg"`
	TotalWeight      float64              `json:"totalWeight"`
	CostPerUnit      float64              `json:"costPerUnit"`
	BaseTransferCost float64              `json:"baseTransferCost"`
	AdditionalCosts  []AdditionalCostLine `json:"additionalCosts"`
	TotalCost        float64              `json:"totalCost"`
	StockBefore      int                  `json:"stockBefore"`
	StockAfter       int                  `json:"stockAfter"`
	StockDelta       int                  `json:"stockDelta"`
	ValidationError  string               `json:"validationError,omitempty"`
}

// PondSellPreviewItem is one row in the sale details summary.
type PondSellPreviewItem struct {
	FishSizeGradeId   int     `json:"fishSizeGradeId"`
	FishSizeGradeName string  `json:"fishSizeGradeName"`
	Weight            float64 `json:"weight"`
	PricePerKg        float64 `json:"pricePerKg"`
	Subtotal          float64 `json:"subtotal"`
	FishCount         *int    `json:"fishCount,omitempty"`
}

// PondSellPreviewResponse is returned by POST /pond/:pondId/sell/preview.
type PondSellPreviewResponse struct {
	Valid           bool                  `json:"valid"`
	Items           []PondSellPreviewItem `json:"items"`
	TotalRevenue    float64               `json:"totalRevenue"`
	TotalWeight     float64               `json:"totalWeight"`
	ValidationError string                `json:"validationError,omitempty"`
}

// ActivityResponse is one row of the pond activity history timeline returned
// by GET /pond/:pondId/activities. Total is computed server-side so the
// client doesn't need to know the fill/move (amount*pricePerUnit) vs sell
// (sum of sell_details.weight * price_per_unit) shape.
type ActivityResponse struct {
	Id           int       `json:"id"`
	Mode         string    `json:"mode"` // fill | move | sell
	ActivityDate time.Time `json:"activityDate"`
	FishType     string    `json:"fishType"`
	Amount       int       `json:"amount"`
	PricePerUnit float64   `json:"pricePerUnit"`
	Total        float64   `json:"total"`
	Merchant     *string   `json:"merchant,omitempty"`
	ToPondName   *string   `json:"toPondName,omitempty"`
}
