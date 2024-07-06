package models

import (
	"errors"
	"time"
)

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

// AddBill represents a bill in the system.
type AddBill struct {
	Type        string    `json:"type" gorm:"column:Type"`
	Other       *string   `json:"other" gorm:"column:Other"`
	FarmGroupId int       `json:"farmGroupId" gorm:"column:FarmGroupId"`
	PaidAmount  float64   `json:"paidAmount" gorm:"column:PaidAmount"`
	PaymentDate time.Time `json:"paymentDate" gorm:"column:PaymentDate"`
}

type BillWithFarmGroupName struct {
	Bill
	Name string `json:"name" gorm:"column:Name"`
}

// Validation Add
func (a AddBill) Validation() error {
	if a.Type == "" {
		return errors.New(ErrTypeEmpty)
	}
	if a.FarmGroupId == 0 {
		return errors.New(ErrFarmGroupIdEmpty)
	}
	if a.PaidAmount == 0 {
		return errors.New(ErrPaidAmountEmpty)
	}
	if a.PaymentDate.IsZero() {
		return errors.New(ErrPaymentDateEmpty)
	}
	return nil
}

// Transfer Add
func (a AddBill) Transfer(bill *Bill) {
	bill.Type = a.Type
	bill.Other = a.Other
	bill.FarmGroupId = a.FarmGroupId
	bill.PaidAmount = a.PaidAmount
	bill.PaymentDate = a.PaymentDate
}

const (
	ErrTypeEmpty        = "type is empty"
	ErrPaidAmountEmpty  = "paid amount is empty"
	ErrPaymentDateEmpty = "payment date is empty"
)
