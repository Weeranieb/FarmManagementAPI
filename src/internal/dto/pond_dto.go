package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreatePondsRequest is the body for POST /pond (create multiple ponds for a farm). New ponds are created with status maintenance.
type CreatePondsRequest struct {
	FarmId int      `json:"farmId" validate:"required"`
	Names  []string `json:"names" validate:"required,min=1,dive,required"`
}

// UpdatePondRequest is used by the service layer (id comes from path).
type UpdatePondRequest struct {
	Id     int    `json:"-"` // from path
	FarmId int    `json:"farmId"`
	Name   string `json:"name"`
	Status string `json:"status" validate:"omitempty,oneof=active maintenance"`
}

// UpdatePondBody is the request body for PUT /pond/:id (id in path).
type UpdatePondBody struct {
	FarmId int    `json:"farmId"`
	Name   string `json:"name"`
	Status string `json:"status" validate:"omitempty,oneof=active maintenance"`
}

type PondResponse struct {
	Id                 int        `json:"id"`
	FarmId             int        `json:"farmId"`
	Name               string     `json:"name"`
	TotalFish          *int       `json:"totalFish"`
	Status             string     `json:"status"`
	FishTypes          []string   `json:"fishTypes"`
	AgeDays            *int       `json:"ageDays"`
	StartDate          *time.Time `json:"startDate"`
	LatestActivityDate *time.Time `json:"latestActivityDate"`
	LatestActivityType *string    `json:"latestActivityType"`
	CreatedAt          time.Time  `json:"createdAt"`
	CreatedBy          string     `json:"createdBy"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	UpdatedBy          string     `json:"updatedBy"`
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

// PondSellDetailItem represents a single per-species line in a sell request.
type PondSellDetailItem struct {
	FishType     string          `json:"fishType" validate:"required"`
	Size         string          `json:"size" validate:"required"`
	Amount       decimal.Decimal `json:"amount" validate:"required,decimal_gt0" swaggertype:"number"`
	FishUnit     string          `json:"fishUnit" validate:"required"`
	PricePerUnit decimal.Decimal `json:"pricePerUnit" validate:"required,decimal_gt0" swaggertype:"number"`
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
