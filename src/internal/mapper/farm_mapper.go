package mapper

import (
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
)

// ToFarmResponse maps model.Farm to dto.FarmResponse
func ToFarmResponse(farm *model.Farm) *dto.FarmResponse {
	if farm == nil {
		return nil
	}
	return &dto.FarmResponse{
		Id:       farm.Id,
		ClientId: farm.ClientId,
		Name:     farm.Name,
		Status:   farm.Status,
	}
}

// ToFarmResponseList maps a slice of model.Farm to a slice of dto.FarmResponse
func ToFarmResponseList(farms []*model.Farm) []*dto.FarmResponse {
	if farms == nil {
		return nil
	}
	responses := make([]*dto.FarmResponse, 0, len(farms))
	for _, farm := range farms {
		responses = append(responses, ToFarmResponse(farm))
	}
	return responses
}

// ToFarmDetailResponse maps model.Farm and its ponds to dto.FarmDetailResponse
func ToFarmDetailResponse(farm *model.Farm, ponds []*model.Pond) *dto.FarmDetailResponse {
	if farm == nil {
		return nil
	}
	pondItems := make([]dto.FarmDetailPondItem, 0, len(ponds))
	var activePonds, maintenancePonds int
	for _, p := range ponds {
		pondItems = append(pondItems, dto.FarmDetailPondItem{Id: p.Id, Name: p.Name, Status: p.Status})
		switch p.Status {
		case constants.FarmStatusActive:
			activePonds++
		case constants.FarmStatusMaintenance:
			maintenancePonds++
		}
	}
	createdAt := ""
	if !farm.CreatedAt.IsZero() {
		createdAt = farm.CreatedAt.UTC().Format(time.RFC3339)
	}
	return &dto.FarmDetailResponse{
		Id:        farm.Id,
		ClientId:  farm.ClientId,
		Name:      farm.Name,
		Status:    farm.Status,
		CreatedAt: createdAt,
		Summary: dto.FarmDetailSummary{
			TotalStock:       0, // FIXME: no stock source yet
			ActivePonds:      activePonds,
			TotalPonds:       len(ponds),
			MaintenancePonds: maintenancePonds,
		},
		Ponds: pondItems,
	}
}
