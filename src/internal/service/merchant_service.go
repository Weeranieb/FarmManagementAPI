package service

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=MerchantService --output=./mocks --outpkg=service --filename=merchant_service.go --structname=MockMerchantService --with-expecter=false
type MerchantService interface {
	Create(request dto.CreateMerchantRequest, username string) (*dto.MerchantResponse, error)
	Get(id int) (*dto.MerchantResponse, error)
	Update(request *model.Merchant, username string) error
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

func (s *merchantService) Create(request dto.CreateMerchantRequest, username string) (*dto.MerchantResponse, error) {
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
		BaseModel: model.BaseModel{
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	// Create merchant
	err = s.merchantRepo.Create(newMerchant)
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

func (s *merchantService) Update(request *model.Merchant, username string) error {
	// Update merchant
	request.UpdatedBy = username
	if err := s.merchantRepo.Update(request); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
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

