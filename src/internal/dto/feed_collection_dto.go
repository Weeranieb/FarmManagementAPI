package dto

import "time"

type CreateFeedCollectionRequest struct {
	Name               string                              `json:"name" validate:"required"`
	Unit               string                              `json:"unit" validate:"required"`
	FeedType           string                              `json:"feedType" validate:"omitempty,oneof=fresh pellet"`
	Fcr                *float64                            `json:"fcr,omitempty"`
	ClientId           *int                                `json:"clientId,omitempty"` // when JWT has no clientId (e.g. super admin), required for create
	FeedPriceHistories []CreateFeedPriceHistoryItemRequest `json:"feedPriceHistories"`
}

type CreateFeedPriceHistoryItemRequest struct {
	Price            float64   `json:"price" validate:"required"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate" validate:"required"`
}

type UpdateFeedCollectionRequest struct {
	Id       int      `json:"id" validate:"required"`
	Name     string   `json:"name"`
	Unit     string   `json:"unit"`
	FeedType string   `json:"feedType" validate:"omitempty,oneof=fresh pellet"`
	Fcr      *float64 `json:"fcr,omitempty"`
}

type FeedCollectionResponse struct {
	Id        int       `json:"id"`
	ClientId  int       `json:"clientId"`
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	FeedType  string    `json:"feedType"`
	Fcr       *float64  `json:"fcr,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

type FeedCollectionPageResponse struct {
	FeedCollectionResponse
	LatestPrice            *float64 `json:"latestPrice"`
	LatestPriceUpdatedDate *string  `json:"latestPriceUpdatedDate"`
}

type CreateFeedCollectionResponse struct {
	FeedCollection   *FeedCollectionResponse `json:"feedCollection"`
	FeedPriceHistory []any                   `json:"feedPriceHistory"` // Will be FeedPriceHistoryResponse when that model is created
}
