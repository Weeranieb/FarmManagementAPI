package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"

	"gorm.io/gorm"
)

type IAdditionalCostService interface {
	Create(request models.AddAdditionalCostRequest, userIdentity string) (*models.AdditionalCost, error)
	Get(id int) (*models.AdditionalCost, error)
	Update(request *models.AdditionalCost, userIdentity string) error
	WithTrx(trxHandle *gorm.DB) IAdditionalCostService
	// TakePage(clientId, page, pageSize int, orderBy, keyword string, billType string, farmGroupId int) (*httputil.PageModel, error)
}

type additionalCostServiceImp struct {
	AdditionalCostRepo repositories.IAdditionalCostRepository
}

func NewAdditionalCostService(additionalCostRepo repositories.IAdditionalCostRepository) IAdditionalCostService {
	return &additionalCostServiceImp{
		AdditionalCostRepo: additionalCostRepo,
	}
}

func (sv additionalCostServiceImp) Create(request models.AddAdditionalCostRequest, userIdentity string) (*models.AdditionalCost, error) {
	var err error
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check bill if exist
	newAdditionalCost := &models.AdditionalCost{}
	request.Transfer(newAdditionalCost)
	newAdditionalCost.UpdatedBy = userIdentity
	newAdditionalCost.CreatedBy = userIdentity

	// create bill
	newAdditionalCost, err = sv.AdditionalCostRepo.Create(newAdditionalCost)
	if err != nil {
		return nil, err
	}

	return newAdditionalCost, nil
}

func (sv additionalCostServiceImp) Get(id int) (*models.AdditionalCost, error) {
	return sv.AdditionalCostRepo.TakeById(id)
}

func (sv additionalCostServiceImp) Update(request *models.AdditionalCost, userIdentity string) error {
	// update bill
	request.UpdatedBy = userIdentity
	if err := sv.AdditionalCostRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv additionalCostServiceImp) WithTrx(trxHandle *gorm.DB) IAdditionalCostService {
	sv.AdditionalCostRepo = sv.AdditionalCostRepo.WithTrx(trxHandle)
	return sv
}
