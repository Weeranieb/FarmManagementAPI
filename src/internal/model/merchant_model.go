package model

type Merchant struct {
	Id            int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId      int    `json:"clientId" gorm:"column:client_id"`
	Name          string `json:"name" gorm:"column:name"`
	ContactNumber string `json:"contactNumber" gorm:"column:contact_number"`
	Location      string `json:"location" gorm:"column:location"`
	BaseModel
}
