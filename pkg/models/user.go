package models

// User is the main user model.
type User struct {
	Id            int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId      int     `json:"clientId" gorm:"column:ClientId"`
	Username      string  `json:"username" gorm:"column:Username"`
	Password      string  `json:"password" gorm:"column:Password"`
	FirstName     string  `json:"firstName" gorm:"column:FirstName"`
	LastName      *string `json:"lastName" gorm:"column:LastName"`
	ContactNumber *string `json:"contactNumber" gorm:"column:ContactNumber"`
	IsAdmin       bool    `json:"isAdmin" gorm:"column:IsAdmin"`
	Base
}
