package models

// AdditionalCost represents an additional cost in the system.
type AdditionalCost struct {
	Id         int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ActivityId int     `json:"activityId" gorm:"column:ActivityId"`
	Title      string  `json:"title" gorm:"column:Title"`
	Cost       float64 `json:"cost" gorm:"column:Cost"`
	Base
}
