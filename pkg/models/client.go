package models

import "errors"

// Client represents a client in the farm management system.
type Client struct {
	Id            int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"column:Name"`
	OwnerName     string `json:"ownerName" gorm:"column:OwnerName"`
	ContactNumber string `json:"contactNumber" gorm:"column:ContactNumber"`
	IsActive      bool   `json:"isActive" gorm:"column:IsActive"`
	Base
}

type AddClient struct {
	Name          string `json:"name" gorm:"column:Name"`
	OwnerName     string `json:"ownerName" gorm:"column:OwnerName"`
	ContactNumber string `json:"contactNumber" gorm:"column:ContactNumber"`
}

// Validation Add
func (a AddClient) Validation() error {
	if a.Name == "" {
		return errors.New(ErrNameEmpty)
	}
	if a.OwnerName == "" {
		return errors.New(ErrOwnerNameEmpty)
	}
	if a.ContactNumber == "" {
		return errors.New(ErrContactNumberEmpty)
	}

	return nil
}

// Transfer Add
func (a AddClient) Transfer(client *Client) error {
	client.Name = a.Name
	client.OwnerName = a.OwnerName
	client.ContactNumber = a.ContactNumber
	return nil
}

const (
	ErrNameEmpty          = "name is empty"
	ErrOwnerNameEmpty     = "owner name is empty"
	ErrContactNumberEmpty = "contact number is empty"
)
