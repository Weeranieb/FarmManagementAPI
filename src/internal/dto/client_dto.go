package dto

import "time"

type CreateClientRequest struct {
	Name          string `json:"name" validate:"required"`
	OwnerName     string `json:"ownerName" validate:"required"`
	ContactNumber string `json:"contactNumber" validate:"required"`
}

type UpdateClientRequest struct {
	Id            int    `json:"id" validate:"required"`
	Name          string `json:"name"`
	OwnerName     string `json:"ownerName"`
	ContactNumber string `json:"contactNumber"`
	IsActive      *bool  `json:"isActive"`
}

type ClientResponse struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	OwnerName     string    `json:"ownerName"`
	ContactNumber string    `json:"contactNumber"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UpdatedBy     string    `json:"updatedBy"`
}
