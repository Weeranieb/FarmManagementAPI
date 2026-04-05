package dto

import "github.com/shopspring/decimal"

// --- Request DTOs ---

type DailyLogEntryInput struct {
	Day               int             `json:"day" validate:"required,min=1,max=31"`
	FreshMorning      decimal.Decimal `json:"freshMorning" validate:"decimal_gte0" swaggertype:"number"`
	FreshEvening      decimal.Decimal `json:"freshEvening" validate:"decimal_gte0" swaggertype:"number"`
	PelletMorning     decimal.Decimal `json:"pelletMorning" validate:"decimal_gte0" swaggertype:"number"`
	PelletEvening     decimal.Decimal `json:"pelletEvening" validate:"decimal_gte0" swaggertype:"number"`
	DeathFishCount    int             `json:"deathFishCount" validate:"gte=0"`
	TouristCatchCount *int            `json:"touristCatchCount,omitempty" validate:"omitempty,gte=0"`
}

type DailyLogBulkUpsertRequest struct {
	Month                  string               `json:"month" validate:"required"` // YYYY-MM
	FreshFeedCollectionId  *int                 `json:"freshFeedCollectionId,omitempty"`
	PelletFeedCollectionId *int                 `json:"pelletFeedCollectionId,omitempty"`
	Entries                []DailyLogEntryInput `json:"entries" validate:"required,min=1,dive"`
}

// --- Response DTOs ---

type DailyLogEntryResponse struct {
	Id                int              `json:"id"`
	Day               int              `json:"day"`
	FreshMorning      decimal.Decimal  `json:"freshMorning"`
	FreshEvening      decimal.Decimal  `json:"freshEvening"`
	PelletMorning     decimal.Decimal  `json:"pelletMorning"`
	PelletEvening     decimal.Decimal  `json:"pelletEvening"`
	DeathFishCount    int              `json:"deathFishCount"`
	TouristCatchCount *int             `json:"touristCatchCount"`
	FreshUnitPrice    *decimal.Decimal `json:"freshUnitPrice,omitempty"`
	PelletUnitPrice   *decimal.Decimal `json:"pelletUnitPrice,omitempty"`
}

type DailyLogMonthResponse struct {
	FreshFeedCollectionId    *int                    `json:"freshFeedCollectionId,omitempty"`
	PelletFeedCollectionId   *int                    `json:"pelletFeedCollectionId,omitempty"`
	FreshFeedCollectionName  string                  `json:"freshFeedCollectionName"`
	PelletFeedCollectionName string                  `json:"pelletFeedCollectionName"`
	FreshUnit                string                  `json:"freshUnit"`
	PelletUnit               string                  `json:"pelletUnit"`
	Entries                  []DailyLogEntryResponse `json:"entries"`
}

type DailyLogExcelUploadResponse struct {
	RowsImported int    `json:"rowsImported"`
	SavedPath    string `json:"savedPath"`
}
