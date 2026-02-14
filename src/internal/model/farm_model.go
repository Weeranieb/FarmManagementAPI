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

// FarmWithPonds is a farm with its ponds, used as the result type for hierarchy.
type FarmWithPonds struct {
	Farm  Farm    `json:"farm"`
	Ponds []*Pond `json:"ponds"`
}

// FarmWithPondsLoad is a model used only for the hierarchy query with GORM Preload.
// It maps to the farms table and has a Ponds association so Preload("Ponds") can load ponds.
type FarmWithPondsLoad struct {
	Id       int    `gorm:"column:id;primaryKey;autoIncrement"`
	ClientId int    `gorm:"column:client_id"`
	Name     string `gorm:"column:name"`
	Status   string `gorm:"column:status"`
	BaseModel
	Ponds []*Pond `gorm:"foreignKey:FarmId"`
}

func (FarmWithPondsLoad) TableName() string { return "farms" }
