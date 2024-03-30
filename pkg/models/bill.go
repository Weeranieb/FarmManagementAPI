package models

import "time"

// Bill represents a bill in the system.
type Bill struct {
	Id          int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Type        string    `json:"type" gorm:"column:Type"`
	Other       *string   `json:"other" gorm:"column:Other"`
	FarmGroupId int       `json:"farmGroupId" gorm:"column:FarmGroupId"`
	PaidAmount  float64   `json:"paidAmount" gorm:"column:PaidAmount"`
	PaymentDate time.Time `json:"paymentDate" gorm:"column:PaymentDate"`
	Base
}
