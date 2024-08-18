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
	GetClientWithFarms(userLevel int, clientId *int) ([]*models.ClientWithFarms, error)
	GetAllClient(userLevel int, keyword string) ([]*models.Client, error)
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
	return sv.ClientRepo.TakeById(id)
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
	if err := sv.ClientRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv clientServiceImp) GetClientWithFarms(userLevel int, clientId *int) ([]*models.ClientWithFarms, error) {
	// check user level
	if userLevel < 2 {
		return nil, errors.New("permission denied")
	}

	// get client by id
	response, err := sv.ClientRepo.GetClientWithFarms(clientId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (sv clientServiceImp) GetAllClient(userLevel int, keyword string) ([]*models.Client, error) {
	// check user level
	if userLevel < 2 {
		return nil, errors.New("permission denied")
	}

	// get all client
	response, err := sv.ClientRepo.TakePage(keyword)
	if err != nil {
		return nil, err
	}

	return response, nil
}
