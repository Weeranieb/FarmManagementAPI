package models

import "errors"

// Pond represents a pond in the system.
type Pond struct {
	Id     int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	FarmId int    `json:"farmId" gorm:"column:FarmId"`
	Code   string `json:"code" gorm:"column:Code"`
	Name   string `json:"name" gorm:"column:Name"`
	Base
}

type AddPond struct {
	FarmId int    `json:"farmId" gorm:"column:FarmId"`
	Code   string `json:"code" gorm:"column:Code"`
	Name   string `json:"name" gorm:"column:Name"`
}

// Validation Add
func (a AddPond) Validation() error {
	if a.FarmId == 0 {
		return errors.New(ErrFarmIdEmpty)
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
func (a AddPond) Transfer(pond *Pond) {
	pond.FarmId = a.FarmId
	pond.Code = a.Code
	pond.Name = a.Name
}
