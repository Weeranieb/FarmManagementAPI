package dto

import "time"

type CreateFarmRequest struct {
	ClientId int    `json:"clientId" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type UpdateFarmRequest struct {
	Id       int    `json:"id" validate:"required"`
	ClientId int    `json:"clientId"`
	Code     string `json:"code"`
	Name     string `json:"name"`
}

type FarmResponse struct {
	Id        int       `json:"id"`
	ClientId  int       `json:"clientId"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

