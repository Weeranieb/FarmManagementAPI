package model

type FeedCollection struct {
	Id       int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:client_id"`
	Code     string `json:"code" gorm:"column:code"`
	Name     string `json:"name" gorm:"column:name"`
	Unit     string `json:"unit" gorm:"column:unit"`
	BaseModel
}

type FeedCollectionPage struct {
	FeedCollection
	LatestPrice           float64 `json:"latestPrice" gorm:"column:latest_price"`
	LatestPriceUpdatedDate *string `json:"latestPriceUpdatedDate" gorm:"column:latest_price_updated_date"`
}

