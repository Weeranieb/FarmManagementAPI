package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IFarmOnFarmGroupService interface {
	Create(request models.AddFarmOnFarmGroup, userIdentity string) (*models.FarmOnFarmGroup, error)
	Delete(id int) error
}

type farmOnFarmGroupServiceImp struct {
	FarmOnFarmGroupRepo repositories.IFarmOnFarmGroupRepository
}

func NewFarmOnFarmGroupService(farmOnFarmGroupRepo repositories.IFarmOnFarmGroupRepository) IFarmOnFarmGroupService {
	return &farmOnFarmGroupServiceImp{
		FarmOnFarmGroupRepo: farmOnFarmGroupRepo,
	}
}

func (sv farmOnFarmGroupServiceImp) Create(request models.AddFarmOnFarmGroup, userIdentity string) (*models.FarmOnFarmGroup, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check if exist
	check, err := sv.FarmOnFarmGroupRepo.FirstByQuery("\"FarmId\" = ? AND \"FarmGroupId\" = ? AND \"DelFlag\" = ?", request.FarmId, request.FarmGroupId, false)
	if err != nil {
		return nil, err
	}

	if check != nil {
		return nil, errors.New("farm on farm group already exist")
	}

	newFarmOnFarmGroup := &models.FarmOnFarmGroup{}
	request.Transfer(newFarmOnFarmGroup)
	newFarmOnFarmGroup.UpdatedBy = userIdentity
	newFarmOnFarmGroup.CreatedBy = userIdentity

	// create farm
	newFarmOnFarmGroup, err = sv.FarmOnFarmGroupRepo.Create(newFarmOnFarmGroup)
	if err != nil {
		return nil, err
	}

	return newFarmOnFarmGroup, nil
}

func (sv farmOnFarmGroupServiceImp) Delete(id int) error {
	err := sv.FarmOnFarmGroupRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
