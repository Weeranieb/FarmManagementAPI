package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IActivePondService interface {
	Create(request models.AddActivePond, userIdentity string) (*models.ActivePond, error)
	Get(id int) (*models.ActivePond, error)
	GetList(farmId int) ([]*models.PondWithActive, error)
	Update(request *models.ActivePond, userIdentity string) error
}

type activePondServiceImp struct {
	ActivePondRepo repositories.IActivePondRepository
}

func NewActivePondService(activePondRepo repositories.IActivePondRepository) IActivePondService {
	return &activePondServiceImp{
		ActivePondRepo: activePondRepo,
	}
}

func (sv activePondServiceImp) Create(request models.AddActivePond, userIdentity string) (*models.ActivePond, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check pond if exist
	checkPond, err := sv.ActivePondRepo.FirstByQuery("\"PondId\" = ? AND \"IsActive\" = ? AND \"DelFlag\" = ?", request.PondId, true, false)
	if err != nil {
		return nil, err
	}

	if checkPond != nil {
		return nil, errors.New("pond already active")
	}

	newActivePond := &models.ActivePond{}
	request.Transfer(newActivePond)
	newActivePond.IsActive = true
	newActivePond.UpdatedBy = userIdentity
	newActivePond.CreatedBy = userIdentity

	// create user
	newActivePond, err = sv.ActivePondRepo.Create(newActivePond)
	if err != nil {
		return nil, err
	}

	return newActivePond, nil
}

func (sv activePondServiceImp) Get(id int) (*models.ActivePond, error) {
	return sv.ActivePondRepo.TakeById(id)
}

func (sv activePondServiceImp) GetList(farmId int) ([]*models.PondWithActive, error) {
	return sv.ActivePondRepo.GetListWithActive(farmId)
}

func (sv activePondServiceImp) Update(request *models.ActivePond, userIdentity string) error {
	// update activePond
	request.UpdatedBy = userIdentity
	if err := sv.ActivePondRepo.Update(request); err != nil {
		return err
	}
	return nil
}
