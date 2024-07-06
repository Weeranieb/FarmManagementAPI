package models

import (
	"errors"
	"time"
)

// Worker represents a worker in the system.
type Worker struct {
	Id            int        `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId      int        `json:"clientId" gorm:"column:ClientId"`
	FarmGroupId   int        `json:"farmGroupId" gorm:"column:FarmGroupId"`
	FirstName     string     `json:"firstName" gorm:"column:FirstName"`
	LastName      *string    `json:"lastName" gorm:"column:LastName"`
	ContactNumber *string    `json:"contactNumber" gorm:"column:ContactNumber"`
	Nationality   string     `json:"nationality" gorm:"column:Nationality"`
	Salary        float64    `json:"salary" gorm:"column:Salary"`
	HireDate      *time.Time `json:"hireDate" gorm:"column:HireDate"`
	IsActive      bool       `json:"isActive" gorm:"column:IsActive"`
	Base
}

type AddWorker struct {
	FarmGroupId   int        `json:"farmGroupId" gorm:"column:FarmGroupId"`
	FirstName     string     `json:"firstName" gorm:"column:FirstName"`
	LastName      *string    `json:"lastName" gorm:"column:LastName"`
	ContactNumber *string    `json:"contactNumber" gorm:"column:ContactNumber"`
	Nationality   string     `json:"nationality" gorm:"column:Nationality"`
	Salary        float64    `json:"salary" gorm:"column:Salary"`
	HireDate      *time.Time `json:"hireDate" gorm:"column:HireDate"`
}

// Validation Add
func (a AddWorker) Validation() error {
	if a.FarmGroupId == 0 {
		return errors.New(ErrFarmGroupIdEmpty)
	}
	if a.FirstName == "" {
		return errors.New(ErrFirstNameEmpty)
	}
	if a.Nationality == "" {
		return errors.New(ErrNationalityEmpty)
	}
	if a.Salary == 0 {
		return errors.New(ErrSalaryEmpty)
	}
	return nil
}

// Transfer Add
func (a AddWorker) Transfer(worker *Worker) {
	worker.FarmGroupId = a.FarmGroupId
	worker.FirstName = a.FirstName
	worker.LastName = a.LastName
	worker.ContactNumber = a.ContactNumber
	worker.Nationality = a.Nationality
	worker.Salary = a.Salary
	worker.HireDate = a.HireDate
}

const (
	ErrSalaryEmpty      = "salary is empty"
	ErrNationalityEmpty = "nationality is empty"
)
