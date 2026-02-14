package service

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=ClientService --output=./mocks --outpkg=service --filename=client_service.go --structname=MockClientService --with-expecter=false
type ClientService interface {
	Create(ctx context.Context, request dto.CreateClientRequest, username string) (*dto.ClientResponse, error)
	Get(id int) (*dto.ClientResponse, error)
	Update(ctx context.Context, request *model.Client, username string) error
	GetList() ([]*dto.ClientResponse, error)
	GetClientDropdown() ([]*dto.DropdownItem, error)
}

type clientService struct {
	clientRepo repository.ClientRepository
}

func NewClientService(clientRepo repository.ClientRepository) ClientService {
	return &clientService{
		clientRepo: clientRepo,
	}
}

func (s *clientService) Create(ctx context.Context, request dto.CreateClientRequest, username string) (*dto.ClientResponse, error) {
	// Check if client with same name already exists
	checkClient, err := s.clientRepo.GetByName(ctx, request.Name)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkClient != nil {
		return nil, errors.ErrClientAlreadyExists
	}

	newClient := &model.Client{
		Name:          request.Name,
		OwnerName:     request.OwnerName,
		ContactNumber: request.ContactNumber,
		IsActive:      true,
	}

	// Create client (CreatedBy/UpdatedBy set via BaseModel hook from ctx)
	err = s.clientRepo.Create(ctx, newClient)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toClientResponse(newClient), nil
}

func (s *clientService) Get(id int) (*dto.ClientResponse, error) {
	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if client == nil {
		return nil, errors.ErrClientNotFound
	}

	return s.toClientResponse(client), nil
}

func (s *clientService) Update(ctx context.Context, request *model.Client, username string) error {
	// Check if client exists
	existingClient, err := s.clientRepo.GetByID(request.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}

	if existingClient == nil {
		return errors.ErrClientNotFound
	}

	// Update client (UpdatedBy set via BaseModel hook from ctx)
	if err := s.clientRepo.Update(ctx, request); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}

	return nil
}

func (s *clientService) GetList() ([]*dto.ClientResponse, error) {
	clients, err := s.clientRepo.List()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.ClientResponse, 0, len(clients))
	for _, client := range clients {
		responses = append(responses, s.toClientResponse(client))
	}

	return responses, nil
}

func (s *clientService) GetClientDropdown() ([]*dto.DropdownItem, error) {
	clients, err := s.clientRepo.List()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	dropdown := make([]*dto.DropdownItem, 0, len(clients))
	for _, client := range clients {
		dropdown = append(dropdown, &dto.DropdownItem{
			Key:   client.Id,
			Value: client.Name,
		})
	}
	return dropdown, nil
}

func (s *clientService) toClientResponse(client *model.Client) *dto.ClientResponse {
	return &dto.ClientResponse{
		Id:            client.Id,
		Name:          client.Name,
		OwnerName:     client.OwnerName,
		ContactNumber: client.ContactNumber,
		IsActive:      client.IsActive,
		CreatedAt:     client.CreatedAt,
		CreatedBy:     client.CreatedBy,
		UpdatedAt:     client.UpdatedAt,
		UpdatedBy:     client.UpdatedBy,
	}
}
