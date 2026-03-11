package service

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=MerchantService --output=./mocks --outpkg=service --filename=merchant_service.go --structname=MockMerchantService --with-expecter=false
type MerchantService interface {
	Create(ctx context.Context, request dto.CreateMerchantRequest, username string) (*dto.MerchantResponse, error)
	Get(id int) (*dto.MerchantResponse, error)
	Update(ctx context.Context, request dto.UpdateMerchantRequest, username string) error
	Delete(ctx context.Context, id int) error
	GetList() ([]*dto.MerchantResponse, error)
}

type merchantService struct {
	merchantRepo repository.MerchantRepository
}

func NewMerchantService(merchantRepo repository.MerchantRepository) MerchantService {
	return &merchantService{
		merchantRepo: merchantRepo,
	}
}

func (s *merchantService) Create(ctx context.Context, request dto.CreateMerchantRequest, username string) (*dto.MerchantResponse, error) {
	// Check if merchant already exists
	checkMerchant, err := s.merchantRepo.GetByContactNumberAndName(request.ContactNumber, request.Name)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkMerchant != nil {
		return nil, errors.ErrMerchantAlreadyExists
	}

	newMerchant := &model.Merchant{
		Name:          request.Name,
		ContactNumber: request.ContactNumber,
		Location:      request.Location,
	}

	// Create merchant (CreatedBy/UpdatedBy set via BaseModel hook from ctx)
	err = s.merchantRepo.Create(ctx, newMerchant)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toMerchantResponse(newMerchant), nil
}

func (s *merchantService) Get(id int) (*dto.MerchantResponse, error) {
	merchant, err := s.merchantRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if merchant == nil {
		return nil, errors.ErrMerchantNotFound
	}

	return s.toMerchantResponse(merchant), nil
}

func (s *merchantService) Update(ctx context.Context, request dto.UpdateMerchantRequest, username string) error {
	existing, err := s.merchantRepo.GetByID(request.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existing == nil {
		return errors.ErrMerchantNotFound
	}
	existing.Name = request.Name
	existing.ContactNumber = request.ContactNumber
	existing.Location = request.Location
	if err := s.merchantRepo.Update(ctx, existing); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *merchantService) Delete(ctx context.Context, id int) error {
	existing, err := s.merchantRepo.GetByID(id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existing == nil {
		return errors.ErrMerchantNotFound
	}
	return s.merchantRepo.Delete(ctx, id)
}

func (s *merchantService) GetList() ([]*dto.MerchantResponse, error) {
	merchants, err := s.merchantRepo.List()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.MerchantResponse, 0, len(merchants))
	for _, merchant := range merchants {
		responses = append(responses, s.toMerchantResponse(merchant))
	}

	return responses, nil
}

func (s *merchantService) toMerchantResponse(merchant *model.Merchant) *dto.MerchantResponse {
	return &dto.MerchantResponse{
		Id:            merchant.Id,
		Name:          merchant.Name,
		ContactNumber: merchant.ContactNumber,
		Location:      merchant.Location,
		CreatedAt:     merchant.CreatedAt,
		CreatedBy:     merchant.CreatedBy,
		UpdatedAt:     merchant.UpdatedAt,
		UpdatedBy:     merchant.UpdatedBy,
	}
}
