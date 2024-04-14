package models

import "errors"

// FarmGroup represents a group of farms.
type FarmGroup struct {
	Id       int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:ClientId"`
	Code     string `json:"code" gorm:"column:Code"`
	Name     string `json:"name" gorm:"column:Name"`
	Base
}

type AddFarmGroup struct {
	Code string `json:"code" gorm:"column:Code"`
	Name string `json:"name" gorm:"column:Name"`
}

// Validation Add
func (a AddFarmGroup) Validation() error {
	if a.Code == "" {
		return errors.New(ErrCodeEmpty)
	}
	if a.Name == "" {
		return errors.New(ErrNameEmpty)
	}

	return nil
}

// Transfer Add
func (a AddFarmGroup) Transfer(farmGroup *FarmGroup) {
	farmGroup.Code = a.Code
	farmGroup.Name = a.Name
}

type GetFarmGroupResponse struct {
	FarmGroup
	Farms []Farm `json:"farms"`
}
