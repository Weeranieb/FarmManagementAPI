package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IFarmService interface {
	Create(request models.AddFarm, userIdentity string, clientId int) (*models.Farm, error)
	Get(id, clientId int) (*models.Farm, error)
	Update(request *models.Farm, userIdentity string) error
}

type FarmServiceImp struct {
	FarmRepo repositories.IFarmRepository
}

func NewFarmService(farmRepo repositories.IFarmRepository) IFarmService {
	return &FarmServiceImp{
		FarmRepo: farmRepo,
	}
}

func (sv FarmServiceImp) Create(request models.AddFarm, userIdentity string, clientId int) (*models.Farm, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check farm if exist
	checkFarm, err := sv.FarmRepo.FirstByQuery("\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, clientId, false)
	if err != nil {
		return nil, err
	}

	if checkFarm != nil {
		return nil, errors.New("farm already exist")
	}

	newFarm := &models.Farm{}
	request.Transfer(newFarm)
	newFarm.ClientId = clientId
	newFarm.UpdatedBy = userIdentity
	newFarm.CreatedBy = userIdentity

	// create farm
	newFarm, err = sv.FarmRepo.Create(newFarm)
	if err != nil {
		return nil, err
	}

	return newFarm, nil
}

func (sv FarmServiceImp) Get(id, clientId int) (*models.Farm, error) {
	farm, err := sv.FarmRepo.TakeById(id)
	if err != nil {
		return nil, err
	}

	if farm.ClientId != clientId {
		return nil, errors.New("farm not found")
	}

	return farm, nil
}

func (sv FarmServiceImp) Update(request *models.Farm, userIdentity string) error {
	// update farm
	request.UpdatedBy = userIdentity
	if err := sv.FarmRepo.Update(request); err != nil {
		return err
	}
	return nil
}
