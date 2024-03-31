package models

import "errors"

// User is the main user model.
type User struct {
	Id            int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId      int     `json:"clientId" gorm:"column:ClientId"`
	Username      string  `json:"username" gorm:"column:Username;unique"`
	Password      string  `json:"password" gorm:"column:Password"`
	FirstName     string  `json:"firstName" gorm:"column:FirstName"`
	LastName      *string `json:"lastName" gorm:"column:LastName"`
	ContactNumber *string `json:"contactNumber" gorm:"column:ContactNumber"`
	IsAdmin       bool    `json:"isAdmin" gorm:"column:IsAdmin"`
	Base
}

type AddUser struct {
	Username      string  `json:"username" gorm:"column:Username;unique"`
	Password      string  `json:"password" gorm:"column:Password"`
	FirstName     string  `json:"firstName" gorm:"column:FirstName"`
	LastName      *string `json:"lastName" gorm:"column:LastName"`
	ContactNumber *string `json:"contactNumber" gorm:"column:ContactNumber"`
	IsAdmin       bool    `json:"isAdmin" gorm:"column:IsAdmin"`
}

// Validation Add
func (a AddUser) Validation() error {
	if a.Username == "" {
		return errors.New(ErrUsernameEmpty)
	}
	if a.Password == "" {
		return errors.New(ErrPasswordEmpty)
	}
	if a.FirstName == "" {
		return errors.New(ErrFirstNameEmpty)
	}

	return nil
}

// Transfer Add
func (a AddUser) Transfer(user *User) error {
	user.Username = a.Username
	user.Password = a.Password
	user.FirstName = a.FirstName
	user.LastName = a.LastName
	user.ContactNumber = a.ContactNumber
	user.IsAdmin = a.IsAdmin
	return nil
}

const (
	ErrUsernameEmpty  = "username is empty"
	ErrPasswordEmpty  = "password is empty"
	ErrFirstNameEmpty = "first name is empty"
)

type Login struct {
	Username string `json:"username" gorm:"column:Username"`
	Password string `json:"password" gorm:"column:Password"`
}
