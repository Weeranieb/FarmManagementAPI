package models

import "errors"

// Merchant represents a merchant in the system.
type Merchant struct {
	Id            int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"column:Name"`
	ContactNumber string `json:"contactNumber" gorm:"column:ContactNumber"`
	Location      string `json:"location" gorm:"column:Location"`
	Base
}

// AddMerchant represents a merchant to be added in the system.
type AddMerchant struct {
	Name          string `json:"name" gorm:"column:Name"`
	ContactNumber string `json:"contactNumber" gorm:"column:ContactNumber"`
	Location      string `json:"location" gorm:"column:Location"`
}

// Validation Add
func (a AddMerchant) Validation() error {
	if a.Name == "" {
		return errors.New(ErrNameEmpty)
	}
	return nil
}

// Transfer Add
func (a AddMerchant) Transfer(merchant *Merchant) {
	merchant.Name = a.Name
	merchant.ContactNumber = a.ContactNumber
	merchant.Location = a.Location
}
