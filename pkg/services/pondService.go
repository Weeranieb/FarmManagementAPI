package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IPondService interface {
	Create(request models.AddPond, userIdentity string) (*models.Pond, error)
	Get(id int) (*models.Pond, error)
	Update(request *models.Pond, userIdentity string) error
	GetList(clientId int) ([]*models.Pond, error)
}

type pondServiceImp struct {
	PondRepo repositories.IPondRepository
}

func NewPondService(pondRepo repositories.IPondRepository) IPondService {
	return &pondServiceImp{
		PondRepo: pondRepo,
	}
}

func (sv pondServiceImp) Create(request models.AddPond, userIdentity string) (*models.Pond, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check pond if exist
	checkPond, err := sv.PondRepo.FirstByQuery("\"FarmId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", request.FarmId, request.Code, false)
	if err != nil {
		return nil, err
	}

	if checkPond != nil {
		return nil, errors.New("pond already exist")
	}

	newPond := &models.Pond{}
	request.Transfer(newPond)
	newPond.UpdatedBy = userIdentity
	newPond.CreatedBy = userIdentity

	// create user
	newPond, err = sv.PondRepo.Create(newPond)
	if err != nil {
		return nil, err
	}

	return newPond, nil
}

func (sv pondServiceImp) Get(id int) (*models.Pond, error) {
	return sv.PondRepo.TakeById(id)
}

func (sv pondServiceImp) Update(request *models.Pond, userIdentity string) error {
	// update user
	request.UpdatedBy = userIdentity
	if err := sv.PondRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv pondServiceImp) GetList(clientId int) ([]*models.Pond, error) {
	return sv.PondRepo.TakeAll(clientId)
}
