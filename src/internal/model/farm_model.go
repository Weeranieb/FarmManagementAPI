package model

type Farm struct {
	Id       int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:client_id"`
	Name     string `json:"name" gorm:"column:name"`
	Status   string `json:"status" gorm:"column:status;default:'active'"`
	BaseModel
}

// FarmCountByClientId holds total and active farm counts for a client
type FarmCountByClientId struct {
	Total       int64 `json:"total" gorm:"column:total"`
	ActiveCount int64 `json:"activeCount" gorm:"column:active_count"`
}
