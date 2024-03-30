package models

// Pond represents a pond in the system.
type Pond struct {
	Id     int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	FarmId int    `json:"farmId" gorm:"column:FarmId"`
	Code   string `json:"code" gorm:"column:Code"`
	Name   string `json:"name" gorm:"column:Name"`
	Base
}
