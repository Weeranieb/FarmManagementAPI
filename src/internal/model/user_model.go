// src/internal/model/user.go
package model

import "time"

type User struct {
	Id            int     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId      *int    `json:"clientId" gorm:"column:client_id"`
	Username      string  `json:"username" gorm:"column:username;uniqueIndex"`
	Email         *string `json:"email" gorm:"column:email"`
	Password      string  `json:"-" gorm:"column:password"`
	FirstName     string  `json:"firstName" gorm:"column:first_name"`
	LastName      *string `json:"lastName" gorm:"column:last_name"`
	UserLevel     int     `json:"userLevel" gorm:"column:user_level"`
	ContactNumber string  `json:"contactNumber" gorm:"column:contact_number"`
	// PasswordUpdatedAt is when the current password was set — written only
	// alongside Password (create, self-change, admin reset), never by an
	// ordinary profile edit. NULL means "unknown" (pre-migration accounts).
	PasswordUpdatedAt *time.Time `json:"passwordUpdatedAt" gorm:"column:password_updated_at"`
	BaseModel
}

// UserCountPerClient holds the user total grouped by client_id.
type UserCountPerClient struct {
	ClientId int   `gorm:"column:client_id"`
	Total    int64 `gorm:"column:total"`
}
