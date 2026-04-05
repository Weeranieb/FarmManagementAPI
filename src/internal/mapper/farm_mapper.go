package mapper

import (
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

// ToFarmResponse maps model.Farm to dto.FarmResponse (PondCount remains 0).
func ToFarmResponse(farm *model.Farm) *dto.FarmResponse {
	if farm == nil {
		return nil
	}
	return &dto.FarmResponse{
		Id:        farm.Id,
		ClientId:  farm.ClientId,
		Name:      farm.Name,
		Status:    farm.Status,
		PondCount: 0,
	}
}

// ToFarmResponseFromFarmWithPonds maps model.FarmWithPonds to dto.FarmResponse with PondCount set from ponds length.
func ToFarmResponseFromFarmWithPonds(f *model.FarmWithPonds) *dto.FarmResponse {
	if f == nil {
		return nil
	}
	ponds := f.Ponds
	if ponds == nil {
		ponds = []*model.Pond{}
	}
	return &dto.FarmResponse{
		Id:        f.Farm.Id,
		ClientId:  f.Farm.ClientId,
		Name:      f.Farm.Name,
		Status:    utils.DeriveFarmStatusFromPonds(ponds),
		PondCount: len(ponds),
	}
}

// ToFarmResponseListFromFarmWithPonds maps a slice of model.FarmWithPonds to a slice of dto.FarmResponse with PondCount.
func ToFarmResponseListFromFarmWithPonds(list []*model.FarmWithPonds) []*dto.FarmResponse {
	if list == nil {
		return nil
	}
	responses := make([]*dto.FarmResponse, 0, len(list))
	for _, f := range list {
		responses = append(responses, ToFarmResponseFromFarmWithPonds(f))
	}
	return responses
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
	pondList := ponds
	if pondList == nil {
		pondList = []*model.Pond{}
	}
	pondItems := make([]dto.FarmDetailPondItem, 0, len(pondList))
	var activePonds, maintenancePonds int
	for _, p := range pondList {
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
		Status:    utils.DeriveFarmStatusFromPonds(pondList),
		CreatedAt: createdAt,
		Summary: dto.FarmDetailSummary{
			TotalStock:       0, // FIXME: no stock source yet
			ActivePonds:      activePonds,
			TotalPonds:       len(pondList),
			MaintenancePonds: maintenancePonds,
		},
		Ponds: pondItems,
	}
}

// ToFarmHierarchyItem maps model.Farm and its ponds to dto.FarmHierarchyItem
func ToFarmHierarchyItem(farm *model.Farm, ponds []*model.Pond) *dto.FarmHierarchyItem {
	if farm == nil {
		return nil
	}
	pondList := ponds
	if pondList == nil {
		pondList = []*model.Pond{}
	}
	pondItems := make([]dto.FarmDetailPondItem, 0, len(pondList))
	for _, p := range pondList {
		pondItems = append(pondItems, dto.FarmDetailPondItem{Id: p.Id, Name: p.Name, Status: p.Status})
	}
	return &dto.FarmHierarchyItem{
		Id:       farm.Id,
		ClientId: farm.ClientId,
		Name:     farm.Name,
		Status:   utils.DeriveFarmStatusFromPonds(pondList),
		Ponds:    pondItems,
	}
}

// ToFarmHierarchyItemFromFarmWithPonds maps model.FarmWithPonds to dto.FarmHierarchyItem
func ToFarmHierarchyItemFromFarmWithPonds(f *model.FarmWithPonds) *dto.FarmHierarchyItem {
	if f == nil {
		return nil
	}
	return ToFarmHierarchyItem(&f.Farm, f.Ponds)
}
