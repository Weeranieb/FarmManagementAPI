package dto

import "time"

type RegisterRequest struct {
	ClientId      int     `json:"clientId" validate:"required"`
	Username      string  `json:"username" validate:"required"`
	Password      string  `json:"password" validate:"required"`
	FirstName     string  `json:"firstName" validate:"required"`
	LastName      *string `json:"lastName"`
	UserLevel     int     `json:"userLevel"`
	ContactNumber string  `json:"contactNumber"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string        `json:"accessToken"`
	ExpiredAt   *time.Time    `json:"expiredAt"`
	User        *UserResponse `json:"user"`
}
