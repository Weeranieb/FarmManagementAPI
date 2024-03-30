package models

import "time"

// Worker represents a worker in the system.
type Worker struct {
	Id            int        `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId      int        `json:"clientId" gorm:"column:ClientId"`
	FarmGroupId   int        `json:"farmGroupId" gorm:"column:FarmGroupId"`
	FirstName     string     `json:"firstName" gorm:"column:FirstName"`
	LastName      *string    `json:"lastName" gorm:"column:LastName"`
	ContactNumber *string    `json:"contactNumber" gorm:"column:ContactNumber"`
	Salary        float64    `json:"salary" gorm:"column:Salary"`
	HireDate      *time.Time `json:"hireDate" gorm:"column:HireDate"`
	IsActive      bool       `json:"isActive" gorm:"column:IsActive"`
	Base
}
