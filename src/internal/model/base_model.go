package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}
