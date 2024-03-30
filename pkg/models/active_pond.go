package models

// ActivePond represents a pond that is currently active.
type ActivePond struct {
	Id        int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	PondId    int     `json:"pondId" gorm:"column:PondId"`
	StartDate string  `json:"startDate" gorm:"column:StartDate"`
	EndDate   *string `json:"endDate" gorm:"column:EndDate"`
	IsActive  bool    `json:"isActive" gorm:"column:IsActive"`
	Base
}
