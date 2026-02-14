package dto

import "time"

// CreatePondsRequest is the body for POST /pond (create multiple ponds for a farm). New ponds are created with status maintenance.
type CreatePondsRequest struct {
	FarmId int      `json:"farmId" validate:"required"`
	Names  []string `json:"names" validate:"required,min=1,dive,required"`
}

type UpdatePondRequest struct {
	Id     int    `json:"id" validate:"required"`
	FarmId int    `json:"farmId"`
	Name   string `json:"name"`
	Status string `json:"status" validate:"omitempty,oneof=active maintenance"`
}

type PondResponse struct {
	Id        int       `json:"id"`
	FarmId    int       `json:"farmId"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

