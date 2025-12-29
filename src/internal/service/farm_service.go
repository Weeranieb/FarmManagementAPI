package service

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FarmService --output=./mocks --outpkg=service --filename=farm_service.go --structname=MockFarmService --with-expecter=false
type FarmService interface {
	Create(request dto.CreateFarmRequest, username string, clientId int) (*dto.FarmResponse, error)
	Get(id, clientId int) (*dto.FarmResponse, error)
	Update(request *model.Farm, username string) error
	GetList(clientId int) ([]*dto.FarmResponse, error)
	GetFarmIdByName(farmName string, clientId int) (int, error)
}

type farmService struct {
	farmRepo repository.FarmRepository
}

func NewFarmService(farmRepo repository.FarmRepository) FarmService {
	return &farmService{
		farmRepo: farmRepo,
	}
}

func (s *farmService) Create(request dto.CreateFarmRequest, username string, clientId int) (*dto.FarmResponse, error) {
	// Check if farm already exists
	checkFarm, err := s.farmRepo.GetByCodeAndClientId(request.Code, clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkFarm != nil {
		return nil, errors.ErrFarmAlreadyExists
	}

	newFarm := &model.Farm{
		ClientId: clientId,
		Code:     request.Code,
		Name:     request.Name,
		BaseModel: model.BaseModel{
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	// Create farm
	err = s.farmRepo.Create(newFarm)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toFarmResponse(newFarm), nil
}

func (s *farmService) Get(id, clientId int) (*dto.FarmResponse, error) {
	farm, err := s.farmRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if farm == nil {
		return nil, errors.ErrFarmNotFound
	}

	// Verify farm belongs to client
	if farm.ClientId != clientId {
		return nil, errors.ErrFarmNotFound
	}

	return s.toFarmResponse(farm), nil
}

func (s *farmService) Update(request *model.Farm, username string) error {
	// Update farm
	request.UpdatedBy = username
	if err := s.farmRepo.Update(request); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *farmService) GetList(clientId int) ([]*dto.FarmResponse, error) {
	farms, err := s.farmRepo.ListByClientId(clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.FarmResponse, 0, len(farms))
	for _, farm := range farms {
		responses = append(responses, s.toFarmResponse(farm))
	}

	return responses, nil
}

func (s *farmService) GetFarmIdByName(farmName string, clientId int) (int, error) {
	farm, err := s.farmRepo.GetByNameAndClientId(farmName, clientId)
	if err != nil {
		return 0, errors.ErrGeneric.Wrap(err)
	}

	if farm == nil {
		return 0, errors.ErrFarmNotFound
	}

	return farm.Id, nil
}

func (s *farmService) toFarmResponse(farm *model.Farm) *dto.FarmResponse {
	return &dto.FarmResponse{
		Id:        farm.Id,
		ClientId:  farm.ClientId,
		Code:      farm.Code,
		Name:      farm.Name,
		CreatedAt: farm.CreatedAt,
		CreatedBy: farm.CreatedBy,
		UpdatedAt: farm.UpdatedAt,
		UpdatedBy: farm.UpdatedBy,
	}
}

