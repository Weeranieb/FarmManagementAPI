package dto

import "time"

type CreateUserRequest struct {
	Username      string  `json:"username" validate:"required"`
	Password      string  `json:"password" validate:"required"`
	Email         *string `json:"email" validate:"omitempty,email"`
	FirstName     string  `json:"firstName" validate:"required"`
	LastName      *string `json:"lastName"`
	UserLevel     int     `json:"userLevel"`
	ContactNumber string  `json:"contactNumber"`
	ClientId      *int    `json:"clientId"`
}

// UpdateUserRequest is used by the self-update endpoint (PUT /user).
// It deliberately excludes UserLevel and ClientId so callers cannot
// self-promote or reassign themselves to another client.
type UpdateUserRequest struct {
	Username      string  `json:"username"`
	Email         *string `json:"email" validate:"omitempty,email"`
	FirstName     string  `json:"firstName"`
	LastName      *string `json:"lastName"`
	ContactNumber string  `json:"contactNumber"`
}

// AdminUpdateUserRequest is used by super-admins via PUT /user/:id.
// It includes privileged fields (UserLevel, ClientId) that the service
// layer guards against promoting to SuperAdmin.
type AdminUpdateUserRequest struct {
	Username      string  `json:"username"`
	Email         *string `json:"email" validate:"omitempty,email"`
	FirstName     string  `json:"firstName"`
	LastName      *string `json:"lastName"`
	UserLevel     *int    `json:"userLevel"`
	ContactNumber string  `json:"contactNumber"`
	ClientId      *int    `json:"clientId"`
}

type UserListQuery struct {
	Search    *string `query:"search"`
	UserLevel *int    `query:"userLevel"`
	ClientId  *int    `query:"clientId"`
}

type UserResponse struct {
	Id            int       `json:"id"`
	ClientId      *int      `json:"clientId"`
	Username      string    `json:"username"`
	Email         *string   `json:"email"`
	FirstName     string    `json:"firstName"`
	LastName      *string   `json:"lastName"`
	UserLevel     int       `json:"userLevel"`
	ContactNumber string    `json:"contactNumber"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UpdatedBy     string    `json:"updatedBy"`
}
