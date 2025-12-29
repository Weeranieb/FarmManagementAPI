package model

type Merchant struct {
	Id            int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"column:name"`
	ContactNumber string `json:"contactNumber" gorm:"column:contact_number"`
	Location      string `json:"location" gorm:"column:location"`
	BaseModel
}
