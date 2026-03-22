package dto

import "time"

type CreateFarmGroupRequest struct {
	ClientId int    `json:"clientId" validate:"required"`
	Name     string `json:"name" validate:"required"`
	FarmIds  []int  `json:"farmIds" validate:"required,min=1"`
}

type UpdateFarmGroupRequest struct {
	Id      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	FarmIds []int  `json:"farmIds" validate:"required,min=1"`
}

type FarmGroupFarmItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FarmGroupResponse struct {
	Id        int                 `json:"id"`
	ClientId  int                 `json:"clientId"`
	Name      string              `json:"name"`
	Farms     []FarmGroupFarmItem `json:"farms"`
	CreatedAt time.Time           `json:"createdAt"`
	CreatedBy string              `json:"createdBy"`
	UpdatedAt time.Time           `json:"updatedAt"`
	UpdatedBy string              `json:"updatedBy"`
}
