package models

// Client represents a client in the farm management system.
type Client struct {
	Id            int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Name          string  `json:"name" gorm:"column:Name"`
	OwnerName     string  `json:"ownerName" gorm:"column:OwnerName"`
	ContactNumber *string `json:"contactNumber" gorm:"column:ContactNumber"`
	IsActive      bool    `json:"isActive" gorm:"column:IsActive"`
	Base
}
