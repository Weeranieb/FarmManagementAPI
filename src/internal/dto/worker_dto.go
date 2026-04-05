package dto

import "time"

type CreateWorkerRequest struct {
	FarmGroupId   int        `json:"farmGroupId" validate:"required"`
	FirstName     string     `json:"firstName" validate:"required"`
	LastName      *string    `json:"lastName"`
	ContactNumber *string    `json:"contactNumber"`
	Nationality   string     `json:"nationality" validate:"required"`
	Salary        float64    `json:"salary" validate:"required"`
	HireDate      *time.Time `json:"hireDate"`
}

type UpdateWorkerRequest struct {
	Id            int        `json:"id" validate:"required"`
	FarmGroupId   int        `json:"farmGroupId"`
	FirstName     string     `json:"firstName"`
	LastName      *string    `json:"lastName"`
	ContactNumber *string    `json:"contactNumber"`
	Nationality   string     `json:"nationality"`
	Salary        float64    `json:"salary"`
	HireDate      *time.Time `json:"hireDate"`
	IsActive      *bool      `json:"isActive"`
}

type WorkerResponse struct {
	Id            int        `json:"id"`
	ClientId      int        `json:"clientId"`
	FarmGroupId   int        `json:"farmGroupId"`
	FirstName     string     `json:"firstName"`
	LastName      *string    `json:"lastName"`
	ContactNumber *string    `json:"contactNumber"`
	Nationality   string     `json:"nationality"`
	Salary        float64    `json:"salary"`
	HireDate      *time.Time `json:"hireDate"`
	IsActive      bool       `json:"isActive"`
	CreatedAt     time.Time  `json:"createdAt"`
	CreatedBy     string     `json:"createdBy"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	UpdatedBy     string     `json:"updatedBy"`
}

type PageResponse struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
}
