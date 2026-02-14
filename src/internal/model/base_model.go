package model

import (
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

// BeforeCreate sets CreatedBy and UpdatedBy from context when username is provided (constants.UsernameKey).
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if tx.Statement.Context == nil {
		return nil
	}
	v := tx.Statement.Context.Value(constants.UsernameKey)
	if v == nil {
		return nil
	}
	username, ok := v.(string)
	if !ok || username == "" {
		return nil
	}
	b.CreatedBy = username
	b.UpdatedBy = username
	return nil
}

// BeforeUpdate sets UpdatedBy from context when username is provided (constants.UsernameKey).
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Context == nil {
		return nil
	}
	v := tx.Statement.Context.Value(constants.UsernameKey)
	if v == nil {
		return nil
	}
	username, ok := v.(string)
	if !ok || username == "" {
		return nil
	}
	b.UpdatedBy = username
	return nil
}
