package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IFarmGroupService interface {
	Create(request models.AddFarmGroup, userIdentity string, clientId int) (*models.FarmGroup, error)
	Get(id, clientId int) (*models.FarmGroup, error)
	GetFarmList(id int) (*[]models.Farm, error)
	Update(request *models.FarmGroup, userIdentity string) error
}

type farmGroupServiceImp struct {
	FarmGroupRepo repositories.IFarmGroupRepository
}

func NewFarmGroupService(farmGroupRepo repositories.IFarmGroupRepository) IFarmGroupService {
	return &farmGroupServiceImp{
		FarmGroupRepo: farmGroupRepo,
	}
}

func (sv farmGroupServiceImp) Create(request models.AddFarmGroup, userIdentity string, clientId int) (*models.FarmGroup, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check farm if exist
	checkFarmGroup, err := sv.FarmGroupRepo.FirstByQuery("\"Code\" = ? AND \"ClientId\" = ? AND \"DelFlag\" = ?", request.Code, clientId, false)
	if err != nil {
		return nil, err
	}

	if checkFarmGroup != nil {
		return nil, errors.New("farm group already exist")
	}

	newFarmGroup := &models.FarmGroup{}
	request.Transfer(newFarmGroup)
	newFarmGroup.ClientId = clientId
	newFarmGroup.UpdatedBy = userIdentity
	newFarmGroup.CreatedBy = userIdentity

	// create farm
	newFarmGroup, err = sv.FarmGroupRepo.Create(newFarmGroup)
	if err != nil {
		return nil, err
	}

	return newFarmGroup, nil
}

func (sv farmGroupServiceImp) Get(id, clientId int) (*models.FarmGroup, error) {
	farm, err := sv.FarmGroupRepo.TakeById(id)
	if err != nil {
		return nil, err
	}

	if farm.ClientId != clientId {
		return nil, errors.New("farm group not found")
	}
	return farm, nil
}

func (sv farmGroupServiceImp) Update(request *models.FarmGroup, userIdentity string) error {
	// update farm
	request.UpdatedBy = userIdentity
	if err := sv.FarmGroupRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv farmGroupServiceImp) GetFarmList(id int) (*[]models.Farm, error) {
	farm, err := sv.FarmGroupRepo.GetFarmList(id)
	if err != nil {
		return nil, err
	}
	return farm, nil
}
