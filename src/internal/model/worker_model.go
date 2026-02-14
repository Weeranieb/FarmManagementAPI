package model

import "time"

type Worker struct {
	Id            int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId      int        `json:"clientId" gorm:"column:client_id"`
	FarmGroupId   int        `json:"farmGroupId" gorm:"column:farm_group_id"`
	FirstName     string     `json:"firstName" gorm:"column:first_name"`
	LastName      *string    `json:"lastName" gorm:"column:last_name"`
	ContactNumber *string    `json:"contactNumber" gorm:"column:contact_number"`
	Nationality   string     `json:"nationality" gorm:"column:nationality"`
	Salary        float64    `json:"salary" gorm:"column:salary"`
	HireDate      *time.Time `json:"hireDate" gorm:"column:hire_date"`
	IsActive      bool       `json:"isActive" gorm:"column:is_active"`
	BaseModel
}
