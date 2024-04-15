package models

import "errors"

// FarmOnFarmGroup represents a farm in a farm group.
type FarmOnFarmGroup struct {
	Id          int `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	FarmId      int `json:"farmId" gorm:"column:FarmId"`
	FarmGroupId int `json:"farmGroupId" gorm:"column:FarmGroupId"`
	Base
}

type AddFarmOnFarmGroup struct {
	FarmId      int `json:"farmId" gorm:"column:FarmId"`
	FarmGroupId int `json:"farmGroupId" gorm:"column:FarmGroupId"`
}

// Validation Add
func (a AddFarmOnFarmGroup) Validation() error {
	if a.FarmId == 0 {
		return errors.New(ErrFarmIdEmpty)
	}
	if a.FarmGroupId == 0 {
		return errors.New(ErrFarmGroupIdEmpty)
	}

	return nil
}

// Transfer Add
func (a AddFarmOnFarmGroup) Transfer(farmOnFarmGroup *FarmOnFarmGroup) {
	farmOnFarmGroup.FarmId = a.FarmId
	farmOnFarmGroup.FarmGroupId = a.FarmGroupId
}

const (
	ErrFarmIdEmpty      = "farm id is empty"
	ErrFarmGroupIdEmpty = "farm group id is empty"
)
