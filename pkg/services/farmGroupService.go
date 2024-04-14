package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IFarmGroupService interface {
	Create(request models.AddFarmGroup, userIdentity string, clientId int) (*models.FarmGroup, error)
	Get(id, clientId int) (*models.GetFarmGroupResponse, error)
	Update(request *models.FarmGroup, userIdentity string) error
}

type FarmGroupServiceImp struct {
	FarmGroupRepo repositories.IFarmGroupRepository
}

func NewFarmGroupService(farmGroupRepo repositories.IFarmGroupRepository) IFarmGroupService {
	return &FarmGroupServiceImp{
		FarmGroupRepo: farmGroupRepo,
	}
}

func (sv FarmGroupServiceImp) Create(request models.AddFarmGroup, userIdentity string, clientId int) (*models.FarmGroup, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check farm if exist
	checkFarmGroup, err := sv.FarmGroupRepo.FirstByQuery("\"Code\" = ? AND \"ClientId\" = ?", request.Code, clientId)
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

func (sv FarmGroupServiceImp) Get(id, clientId int) (*models.GetFarmGroupResponse, error) {
	farm, err := sv.FarmGroupRepo.GetFarmGroupWithFarms(id, clientId)
	if err != nil {
		return nil, err
	}

	return farm, nil
}

func (sv FarmGroupServiceImp) Update(request *models.FarmGroup, userIdentity string) error {
	// update farm
	request.UpdatedBy = userIdentity
	if err := sv.FarmGroupRepo.Update(request); err != nil {
		return err
	}
	return nil
}
