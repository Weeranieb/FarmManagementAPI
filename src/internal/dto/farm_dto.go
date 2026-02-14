package dto

type CreateFarmRequest struct {
	ClientId int    `json:"clientId" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type UpdateFarmRequest struct {
	Id       int    `json:"id" validate:"required"`
	ClientId int    `json:"clientId"`
	Name     string `json:"name"`
}

type FarmResponse struct {
	Id        int    `json:"id"`
	ClientId  int    `json:"clientId"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	PondCount int    `json:"pondCount"`
}

type FarmListResponse struct {
	Farms       []*FarmResponse `json:"farms"`
	Total       int             `json:"total"`
	TotalActive int             `json:"totalActive"`
}
