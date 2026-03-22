package model

type FarmGroup struct {
	Id       int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:client_id"`
	Name     string `json:"name" gorm:"column:name"`
	BaseModel
}

type FarmOnFarmGroup struct {
	Id          int `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FarmId      int `json:"farmId" gorm:"column:farm_id"`
	FarmGroupId int `json:"farmGroupId" gorm:"column:farm_group_id"`
	BaseModel
}

func (FarmOnFarmGroup) TableName() string {
	return "farm_on_farm_group"
}
