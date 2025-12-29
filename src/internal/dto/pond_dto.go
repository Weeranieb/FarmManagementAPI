package dto

import "time"

type CreatePondRequest struct {
	FarmId int    `json:"farmId" validate:"required"`
	Code   string `json:"code" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type CreatePondsRequest []CreatePondRequest

type UpdatePondRequest struct {
	Id     int    `json:"id" validate:"required"`
	FarmId int    `json:"farmId"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

type PondResponse struct {
	Id        int       `json:"id"`
	FarmId    int       `json:"farmId"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

