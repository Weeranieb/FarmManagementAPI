package dto

import "time"

type CreateFeedPriceHistoryRequest struct {
	FeedCollectionId int       `json:"feedCollectionId" validate:"required"`
	Price            float64   `json:"price" validate:"required"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate" validate:"required"`
}

type UpdateFeedPriceHistoryRequest struct {
	Id               int       `json:"id" validate:"required"`
	FeedCollectionId int       `json:"feedCollectionId"`
	Price            float64   `json:"price"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate"`
}

type FeedPriceHistoryResponse struct {
	Id               int       `json:"id"`
	FeedCollectionId int       `json:"feedCollectionId"`
	Price            float64   `json:"price"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate"`
	CreatedAt        time.Time `json:"createdAt"`
	CreatedBy        string    `json:"createdBy"`
	UpdatedAt        time.Time `json:"updatedAt"`
	UpdatedBy        string    `json:"updatedBy"`
}

