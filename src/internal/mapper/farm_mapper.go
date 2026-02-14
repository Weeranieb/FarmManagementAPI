package mapper

import (
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
