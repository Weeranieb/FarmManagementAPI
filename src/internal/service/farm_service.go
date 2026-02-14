package service

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/mapper"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FarmService --output=./mocks --outpkg=service --filename=farm_service.go --structname=MockFarmService --with-expecter=false
type FarmService interface {
	Create(ctx context.Context, request dto.CreateFarmRequest, clientId int) (*dto.FarmResponse, error)
	Get(id int, clientId *int) (*dto.FarmDetailResponse, error)
	Update(ctx context.Context, request dto.UpdateFarmRequest) error
	GetList(clientId int) (*dto.FarmListResponse, error)
	GetHierarchy(clientId int) ([]*dto.FarmHierarchyItem, error)
	GetFarmIdByName(farmName string, clientId int) (int, error)
}

type farmService struct {
	farmRepo repository.FarmRepository
	pondRepo repository.PondRepository
}

func NewFarmService(farmRepo repository.FarmRepository, pondRepo repository.PondRepository) FarmService {
	return &farmService{
		farmRepo: farmRepo,
		pondRepo: pondRepo,
	}
}

func (s *farmService) Create(ctx context.Context, request dto.CreateFarmRequest, clientId int) (*dto.FarmResponse, error) {
	// Check if farm already exists
	checkFarm, err := s.farmRepo.GetByNameAndClientId(request.Name, clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkFarm != nil {
		return nil, errors.ErrFarmAlreadyExists
	}

	// Set default status if not provided
	status := constants.FarmStatusActive

	newFarm := &model.Farm{
		ClientId: clientId,
		Name:     request.Name,
		Status:   status,
	}

	// Create farm (CreatedBy/UpdatedBy set via BaseModel hook from ctx)
	if err = s.farmRepo.Create(ctx, newFarm); err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return mapper.ToFarmResponse(newFarm), nil
}

func (s *farmService) Get(id int, clientId *int) (*dto.FarmDetailResponse, error) {
	farm, err := s.farmRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if farm == nil {
		return nil, errors.ErrFarmNotFound
	}

	// Verify farm belongs to client
	if clientId != nil && farm.ClientId != *clientId {
		return nil, errors.ErrFarmNotFound
	}

	ponds, err := s.pondRepo.ListByFarmId(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if ponds == nil {
		ponds = []*model.Pond{}
	}

	return mapper.ToFarmDetailResponse(farm, ponds), nil
}

// Update loads the existing farm by id, merges request (id, name) with existing ClientId and Status, then updates.
// UpdatedBy is set from ctx in the repository (BaseModel hook).
func (s *farmService) Update(ctx context.Context, request dto.UpdateFarmRequest) error {
	existing, err := s.farmRepo.GetByID(request.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existing == nil {
		return errors.ErrFarmNotFound
	}
	updateFarm := &model.Farm{
		Id:       request.Id,
		ClientId: existing.ClientId,
		Name:     request.Name,
		Status:   existing.Status,
	}
	// Enforce unique farm name per client
	if updateFarm.Name != "" {
		byName, err := s.farmRepo.GetByNameAndClientId(updateFarm.Name, updateFarm.ClientId)
		if err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		if byName != nil && byName.Id != updateFarm.Id {
			return errors.ErrFarmAlreadyExists
		}
	}
	if err := s.farmRepo.Update(ctx, updateFarm); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *farmService) GetList(clientId int) (*dto.FarmListResponse, error) {
	list, err := s.farmRepo.ListByClientIdWithPonds(clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	countByClientId, err := s.farmRepo.CountByClientId(clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := mapper.ToFarmResponseListFromFarmWithPonds(list)

	return &dto.FarmListResponse{
		Farms:       responses,
		Total:       int(countByClientId.Total),
		TotalActive: int(countByClientId.ActiveCount),
	}, nil
}

func (s *farmService) GetHierarchy(clientId int) ([]*dto.FarmHierarchyItem, error) {
	list, err := s.farmRepo.ListByClientIdWithPonds(clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if len(list) == 0 {
		return []*dto.FarmHierarchyItem{}, nil
	}

	result := make([]*dto.FarmHierarchyItem, 0, len(list))
	for _, f := range list {
		result = append(result, mapper.ToFarmHierarchyItemFromFarmWithPonds(f))
	}

	return result, nil
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
