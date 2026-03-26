package dto

import "github.com/shopspring/decimal"

// --- Request DTOs ---

type DailyFeedEntryInput struct {
	Day     int             `json:"day" validate:"required,min=1,max=31"`
	Morning decimal.Decimal `json:"morning" validate:"decimal_gte0" swaggertype:"number"`
	Evening decimal.Decimal `json:"evening" validate:"decimal_gte0" swaggertype:"number"`
}

type DailyFeedBulkUpsertRequest struct {
	FeedCollectionId int                   `json:"feedCollectionId" validate:"required"`
	Month            string                `json:"month" validate:"required"` // "YYYY-MM"
	Entries          []DailyFeedEntryInput `json:"entries" validate:"required,min=1,dive"`
}

// --- Response DTOs ---

type DailyFeedEntryResponse struct {
	Id        int              `json:"id"`
	Day       int              `json:"day"`
	Morning   decimal.Decimal  `json:"morning"`
	Evening   decimal.Decimal  `json:"evening"`
	UnitPrice *decimal.Decimal `json:"unitPrice"` // resolved from feed_price_history; nil if no history
}

type DailyFeedTableResponse struct {
	FeedCollectionId   int                      `json:"feedCollectionId"`
	FeedCollectionName string                   `json:"feedCollectionName"`
	FeedUnit           string                   `json:"feedUnit"`
	Entries            []DailyFeedEntryResponse `json:"entries"`
}

type DailyFeedExcelUploadResponse struct {
	RowsImported int    `json:"rowsImported"`
	SavedPath    string `json:"savedPath"`
}
