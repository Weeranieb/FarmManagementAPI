package service

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=PondService --output=./mocks --outpkg=service --filename=pond_service.go --structname=MockPondService --with-expecter=false
type PondService interface {
	CreatePonds(request dto.CreatePondsRequest, username string) error
	Get(id int) (*dto.PondResponse, error)
	Update(request *model.Pond, username string) error
	GetList(farmId int) ([]*dto.PondResponse, error)
	Delete(id int, username string) error
}

type pondService struct {
	pondRepo repository.PondRepository
}

func NewPondService(pondRepo repository.PondRepository) PondService {
	return &pondService{
		pondRepo: pondRepo,
	}
}

func (s *pondService) CreatePonds(request dto.CreatePondsRequest, username string) error {
	for _, name := range request.Names {
		checkPond, err := s.pondRepo.GetByFarmIdAndName(request.FarmId, name)
		if err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		if checkPond != nil {
			return errors.ErrPondAlreadyExists
		}
	}

	newPonds := make([]*model.Pond, 0, len(request.Names))
	for _, name := range request.Names {
		newPonds = append(newPonds, &model.Pond{
			FarmId: request.FarmId,
			Name:   name,
			Status: constants.FarmStatusMaintenance,
			BaseModel: model.BaseModel{
				CreatedBy: username,
				UpdatedBy: username,
			},
		})
	}

	return s.pondRepo.CreateBatch(newPonds)
}

func (s *pondService) Get(id int) (*dto.PondResponse, error) {
	pond, err := s.pondRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if pond == nil {
		return nil, errors.ErrPondNotFound
	}

	return s.toPondResponse(pond), nil
}

func (s *pondService) Update(request *model.Pond, username string) error {
	// Enforce unique pond name per farm
	if request.Name != "" {
		existing, err := s.pondRepo.GetByFarmIdAndName(request.FarmId, request.Name)
		if err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		if existing != nil && existing.Id != request.Id {
			return errors.ErrPondAlreadyExists
		}
	}

	request.UpdatedBy = username
	if err := s.pondRepo.Update(request); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *pondService) GetList(farmId int) ([]*dto.PondResponse, error) {
	ponds, err := s.pondRepo.ListByFarmId(farmId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.PondResponse, 0, len(ponds))
	for _, pond := range ponds {
		responses = append(responses, s.toPondResponse(pond))
	}

	return responses, nil
}

func (s *pondService) Delete(id int, username string) error {
	// Delete pond (soft delete)
	if err := s.pondRepo.Delete(id); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *pondService) toPondResponse(pond *model.Pond) *dto.PondResponse {
	return &dto.PondResponse{
		Id:        pond.Id,
		FarmId:    pond.FarmId,
		Name:      pond.Name,
		Status:    pond.Status,
		CreatedAt: pond.CreatedAt,
		CreatedBy: pond.CreatedBy,
		UpdatedAt: pond.UpdatedAt,
		UpdatedBy: pond.UpdatedBy,
	}
}
