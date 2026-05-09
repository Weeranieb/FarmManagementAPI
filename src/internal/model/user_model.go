// src/internal/model/user.go
package model

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
	BaseModel
}

// UserCountPerClient holds the user total grouped by client_id.
type UserCountPerClient struct {
	ClientId int   `gorm:"column:client_id"`
	Total    int64 `gorm:"column:total"`
}
