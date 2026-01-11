package dto

import "time"

type CreateUserRequest struct {
	Username      string  `json:"username" validate:"required"`
	Password      string  `json:"password" validate:"required"`
	FirstName     string  `json:"firstName" validate:"required"`
	LastName      *string `json:"lastName"`
	UserLevel     int     `json:"userLevel"`
	ContactNumber string  `json:"contactNumber"`
	ClientId      *int    `json:"clientId"`
}

type UpdateUserRequest struct {
	Username      string  `json:"username"`
	FirstName     string  `json:"firstName"`
	LastName      *string `json:"lastName"`
	UserLevel     int     `json:"userLevel"`
	ContactNumber string  `json:"contactNumber"`
}

type UserResponse struct {
	Id            int       `json:"id"`
	ClientId      *int      `json:"clientId"`
	Username      string    `json:"username"`
	FirstName     string    `json:"firstName"`
	LastName      *string   `json:"lastName"`
	UserLevel     int       `json:"userLevel"`
	ContactNumber string    `json:"contactNumber"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UpdatedBy     string    `json:"updatedBy"`
}
