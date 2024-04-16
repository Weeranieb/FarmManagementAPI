package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IClientService interface {
	Create(request models.AddClient, userIdentity string) (*models.Client, error)
	Get(id int) (*models.Client, error)
	Update(request *models.Client, userIdentity string) error
}

type clientServiceImp struct {
	ClientRepo repositories.IClientRepository
}

func NewClientService(clientRepo repositories.IClientRepository) IClientService {
	return &clientServiceImp{
		ClientRepo: clientRepo,
	}
}

func (sv clientServiceImp) Get(id int) (*models.Client, error) {
	// get client by id
	res, err := sv.ClientRepo.TakeById(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (sv clientServiceImp) Create(request models.AddClient, userIdentity string) (*models.Client, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check name if exist
	checkUser, err := sv.ClientRepo.FirstByQuery("\"Name\" = ? AND \"DelFlag\" = ?", request.Name, false)
	if err != nil {
		return nil, err
	}

	if checkUser != nil {
		return nil, errors.New("client name duplicate")
	}

	newClient := &models.Client{}
	request.Transfer(newClient)
	newClient.UpdatedBy = userIdentity
	newClient.CreatedBy = userIdentity
	newClient.IsActive = true

	// create client
	res, err := sv.ClientRepo.Create(newClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sv clientServiceImp) Update(request *models.Client, userIdentity string) error {
	// update client
	request.UpdatedBy = userIdentity
	// request.UpdatedDate = time.Now()
	if err := sv.ClientRepo.Update(request); err != nil {
		return err
	}
	return nil
}
