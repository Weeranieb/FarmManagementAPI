package models

// FeedCollection represents a feed collection in the system.
type FeedCollection struct {
	Id       int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:ClientId"`
	Code     string `json:"code" gorm:"column:Code"`
	Name     string `json:"name" gorm:"column:Name"`
	Unit     string `json:"unit" gorm:"column:Unit"`
	Base
}
