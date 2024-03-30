package models

// Farm represents a farm in the system.
type Farm struct {
	Id       int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:ClientId"`
	Code     string `json:"code" gorm:"column:Code"`
	Name     string `json:"name" gorm:"column:Name"`
	Base
}
