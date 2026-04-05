package model

type Client struct {
	Id                      int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name                    string `json:"name" gorm:"column:name"`
	OwnerName               string `json:"ownerName" gorm:"column:owner_name"`
	ContactNumber           string `json:"contactNumber" gorm:"column:contact_number"`
	IsActive                bool   `json:"isActive" gorm:"column:is_active"`
	IsTouristFishingEnabled bool   `json:"isTouristFishingEnabled" gorm:"column:is_tourist_fishing_enabled"`
	BaseModel
}
