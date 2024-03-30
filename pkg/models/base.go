package models

import "time"

type Base struct {
	DelFlag     bool      `json:"delFlag" gorm:"column:DelFlag"`
	CreatedDate time.Time `json:"createdDate" gorm:"column:CreatedDate;autoCreateTime:milli"`
	CreatedBy   string    `json:"createdBy" gorm:"column:CreatedBy"`
	UpdatedDate time.Time `json:"updatedDate" gorm:"column:UpdatedDate;autoUpdateTime:milli"`
	UpdatedBy   string    `json:"updatedBy" gorm:"column:UpdatedBy"`
}
