package service

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=PondService --output=./mocks --outpkg=service --filename=pond_service.go --structname=MockPondService --with-expecter=false
type PondService interface {
	Create(request dto.CreatePondRequest, username string) (*dto.PondResponse, error)
	CreateBatch(requests []dto.CreatePondRequest, username string) ([]*dto.PondResponse, error)
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

func (s *pondService) Create(request dto.CreatePondRequest, username string) (*dto.PondResponse, error) {
	// Check if pond already exists
	checkPond, err := s.pondRepo.GetByFarmIdAndCode(request.FarmId, request.Code)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkPond != nil {
		return nil, errors.ErrPondAlreadyExists
	}

	newPond := &model.Pond{
		FarmId: request.FarmId,
		Code:   request.Code,
		Name:   request.Name,
		BaseModel: model.BaseModel{
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	// Create pond
	err = s.pondRepo.Create(newPond)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toPondResponse(newPond), nil
}

func (s *pondService) CreateBatch(requests []dto.CreatePondRequest, username string) ([]*dto.PondResponse, error) {
	// Validate all requests
	for _, req := range requests {
		// Check if pond already exists
		checkPond, err := s.pondRepo.GetByFarmIdAndCode(req.FarmId, req.Code)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}

		if checkPond != nil {
			return nil, errors.ErrPondAlreadyExists
		}
	}

	// Create all ponds
	newPonds := make([]*model.Pond, 0, len(requests))
	for _, req := range requests {
		newPond := &model.Pond{
			FarmId: req.FarmId,
			Code:   req.Code,
			Name:   req.Name,
			BaseModel: model.BaseModel{
				CreatedBy: username,
				UpdatedBy: username,
			},
		}
		newPonds = append(newPonds, newPond)
	}

	// Create batch
	err := s.pondRepo.CreateBatch(newPonds)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	// Convert to responses
	responses := make([]*dto.PondResponse, 0, len(newPonds))
	for _, pond := range newPonds {
		responses = append(responses, s.toPondResponse(pond))
	}

	return responses, nil
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
	// Update pond
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
		Code:      pond.Code,
		Name:      pond.Name,
		CreatedAt: pond.CreatedAt,
		CreatedBy: pond.CreatedBy,
		UpdatedAt: pond.UpdatedAt,
		UpdatedBy: pond.UpdatedBy,
	}
}

