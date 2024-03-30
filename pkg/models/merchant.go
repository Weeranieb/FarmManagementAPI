package models

// Merchant represents a merchant in the system.
type Merchant struct {
	Id            int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"column:Name"`
	ContactNumber string `json:"contactNumber" gorm:"column:ContactNumber"`
	Location      string `json:"location" gorm:"column:Location"`
	Base
}
