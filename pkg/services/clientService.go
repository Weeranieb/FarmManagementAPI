package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IClientService interface {
	Create(request models.AddClient, userIdentity string) (*models.Client, error)
}

type ClientServiceImp struct {
	ClientRepo repositories.IClientRepository
}

func NewClientService(clientRepo repositories.IClientRepository) IClientService {
	return &ClientServiceImp{
		ClientRepo: clientRepo,
	}
}

func (sv ClientServiceImp) Create(request models.AddClient, userIdentity string) (*models.Client, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check name if exist
	checkUser, err := sv.ClientRepo.FirstByQuery("\"Name\" = ?", request.Name)
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
