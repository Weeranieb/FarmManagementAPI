package models

import "errors"

// Farm represents a farm in the system.
type Farm struct {
	Id       int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:ClientId"`
	Code     string `json:"code" gorm:"column:Code"`
	Name     string `json:"name" gorm:"column:Name"`
	Base
}

type AddFarm struct {
	ClientId int    `json:"clientId" gorm:"column:ClientId"`
	Code     string `json:"code" gorm:"column:Code"`
	Name     string `json:"name" gorm:"column:Name"`
}

// Validation Add
func (a AddFarm) Validation() error {
	if a.ClientId == 0 {
		return errors.New(ErrClientIdEmpty)
	}
	if a.Code == "" {
		return errors.New(ErrCodeEmpty)
	}
	if a.Name == "" {
		return errors.New(ErrNameEmpty)
	}

	return nil
}

// Transfer Add
func (a AddFarm) Transfer(farm *Farm) {
	farm.ClientId = a.ClientId
	farm.Code = a.Code
	farm.Name = a.Name
}

const (
	ErrCodeEmpty = "code is empty"
)
