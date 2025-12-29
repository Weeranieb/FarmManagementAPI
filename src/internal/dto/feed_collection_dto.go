package dto

import "time"

type CreateFeedCollectionRequest struct {
	Code               string                           `json:"code" validate:"required"`
	Name               string                           `json:"name" validate:"required"`
	Unit               string                           `json:"unit" validate:"required"`
	FeedPriceHistories []CreateFeedPriceHistoryItemRequest `json:"feedPriceHistories"`
}

type CreateFeedPriceHistoryItemRequest struct {
	Price            float64   `json:"price" validate:"required"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate" validate:"required"`
}

type UpdateFeedCollectionRequest struct {
	Id       int    `json:"id" validate:"required"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Unit     string `json:"unit"`
}

type FeedCollectionResponse struct {
	Id       int       `json:"id"`
	ClientId int       `json:"clientId"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Unit     string    `json:"unit"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

type FeedCollectionPageResponse struct {
	FeedCollectionResponse
	LatestPrice           *float64 `json:"latestPrice"`
	LatestPriceUpdatedDate *string  `json:"latestPriceUpdatedDate"`
}

type CreateFeedCollectionResponse struct {
	FeedCollection   *FeedCollectionResponse   `json:"feedCollection"`
	FeedPriceHistory []interface{}             `json:"feedPriceHistory"` // Will be FeedPriceHistoryResponse when that model is created
}

