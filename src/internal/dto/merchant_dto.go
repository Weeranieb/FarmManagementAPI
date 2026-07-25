package dto

import "time"

type CreateMerchantRequest struct {
	Name          string `json:"name" validate:"required"`
	ContactNumber string `json:"contactNumber"`
	Location      string `json:"location"`
	// ClientId is used only by super admins (no client in token) to target a
	// client. Normal users' client is taken from the JWT and this is ignored.
	ClientId *int `json:"clientId"`
}

type UpdateMerchantRequest struct {
	Id            int    `json:"id" validate:"required"`
	Name          string `json:"name"`
	ContactNumber string `json:"contactNumber"`
	Location      string `json:"location"`
}

type MerchantResponse struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	ContactNumber string    `json:"contactNumber"`
	Location      string    `json:"location"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UpdatedBy     string    `json:"updatedBy"`
}
