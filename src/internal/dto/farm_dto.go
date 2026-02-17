package dto

type CreateFarmRequest struct {
	ClientId int    `json:"clientId" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

// UpdateFarmRequest is used by the service layer (id comes from path).
type UpdateFarmRequest struct {
	Id   int    `json:"-"` // from path
	Name string `json:"name"`
}

// UpdateFarmBody is the request body for PUT /farm/:id (id in path).
type UpdateFarmBody struct {
	Name string `json:"name"`
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

// FarmDetailSummary holds summary stats for the farm detail page cards
type FarmDetailSummary struct {
	TotalStock       int `json:"totalStock"`
	ActivePonds      int `json:"activePonds"`
	TotalPonds       int `json:"totalPonds"`
	MaintenancePonds int `json:"maintenancePonds"`
}

// FarmDetailPondItem is a pond entry in the farm detail response
type FarmDetailPondItem struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// FarmDetailResponse is the full payload for GET /farm/:id (farm detail page)
type FarmDetailResponse struct {
	Id        int                  `json:"id"`
	ClientId  int                  `json:"clientId"`
	Name      string               `json:"name"`
	Status    string               `json:"status"`
	CreatedAt string               `json:"createdAt"`
	Summary   FarmDetailSummary    `json:"summary"`
	Ponds     []FarmDetailPondItem `json:"ponds"`
}

// FarmHierarchyItem is a farm with its ponds for GET /farm/hierarchy
type FarmHierarchyItem struct {
	Id       int                  `json:"id"`
	ClientId int                  `json:"clientId"`
	Name     string               `json:"name"`
	Status   string               `json:"status"`
	Ponds    []FarmDetailPondItem `json:"ponds"`
}
