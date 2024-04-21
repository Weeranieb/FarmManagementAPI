package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IMerchantService interface {
	Create(request models.AddMerchant, userIdentity string) (*models.Merchant, error)
	Get(id int) (*models.Merchant, error)
	Update(request *models.Merchant, userIdentity string) error
}

type merchantServiceImp struct {
	MerchantRepo repositories.IMerchantRepository
}

func NewMerchantService(merchantRepo repositories.IMerchantRepository) IMerchantService {
	return &merchantServiceImp{
		MerchantRepo: merchantRepo,
	}
}

func (sv merchantServiceImp) Create(request models.AddMerchant, userIdentity string) (*models.Merchant, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check pond if exist
	checkMerchant, err := sv.MerchantRepo.FirstByQuery("\"ContactNumber\" = ? AND \"Name\" = ? AND \"DelFlag\" = ?", request.ContactNumber, request.Name, false)
	if err != nil {
		return nil, err
	}

	if checkMerchant != nil {
		return nil, errors.New("merchant already exist")
	}

	newMerchant := &models.Merchant{}
	request.Transfer(newMerchant)
	newMerchant.UpdatedBy = userIdentity
	newMerchant.CreatedBy = userIdentity

	// create user
	newMerchant, err = sv.MerchantRepo.Create(newMerchant)
	if err != nil {
		return nil, err
	}

	return newMerchant, nil
}

func (sv merchantServiceImp) Get(id int) (*models.Merchant, error) {
	return sv.MerchantRepo.TakeById(id)
}

func (sv merchantServiceImp) Update(request *models.Merchant, userIdentity string) error {
	// update user
	request.UpdatedBy = userIdentity
	if err := sv.MerchantRepo.Update(request); err != nil {
		return err
	}
	return nil
}
