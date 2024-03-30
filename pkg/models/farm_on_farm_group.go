package models

// FarmOnFarmGroup represents a farm in a farm group.
type FarmOnFarmGroup struct {
	Id          int `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	FarmId      int `json:"farmId" gorm:"column:FarmId"`
	FarmGroupId int `json:"farmGroupId" gorm:"column:FarmGroupId"`
	Base
}
