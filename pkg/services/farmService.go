package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IFarmService interface {
	Create(request models.AddFarm, userIdentity string) (*models.Farm, error)
	Get(id, clientId int) (*models.Farm, error)
	Update(request *models.Farm, userIdentity string) error
	GetList(clientId int) ([]*models.Farm, error)
	GetFarmIdByName(farmName string, clientId int) (int, error)
}

type farmServiceImp struct {
	FarmRepo repositories.IFarmRepository
}

func NewFarmService(farmRepo repositories.IFarmRepository) IFarmService {
	return &farmServiceImp{
		FarmRepo: farmRepo,
	}
}

func (sv farmServiceImp) Create(request models.AddFarm, userIdentity string) (*models.Farm, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check farm if exist
	checkFarm, err := sv.FarmRepo.FirstByQuery("\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, request.ClientId, false)
	if err != nil {
		return nil, err
	}

	if checkFarm != nil {
		return nil, errors.New("farm already exist")
	}

	newFarm := &models.Farm{}
	request.Transfer(newFarm)
	newFarm.UpdatedBy = userIdentity
	newFarm.CreatedBy = userIdentity

	// create farm
	newFarm, err = sv.FarmRepo.Create(newFarm)
	if err != nil {
		return nil, err
	}

	return newFarm, nil
}

func (sv farmServiceImp) Get(id, clientId int) (*models.Farm, error) {
	farm, err := sv.FarmRepo.TakeById(id)
	if err != nil {
		return nil, err
	}

	if farm.ClientId != clientId {
		return nil, errors.New("farm not found")
	}

	return farm, nil
}

func (sv farmServiceImp) Update(request *models.Farm, userIdentity string) error {
	// update farm
	request.UpdatedBy = userIdentity
	if err := sv.FarmRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv farmServiceImp) GetList(clientId int) ([]*models.Farm, error) {
	return sv.FarmRepo.TakeAll(clientId)
}

func (sv farmServiceImp) GetFarmIdByName(farmName string, clientId int) (int, error) {
	farm, err := sv.FarmRepo.FirstByQuery("\"Name\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", farmName, clientId, false)
	if err != nil {
		return 0, err
	}

	if farm == nil {
		return 0, errors.New("farm not found")
	}

	return farm.Id, nil
}
