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
	Update(ctx context.Context, request dto.UpdateClientRequest, username string) error
	GetList() ([]*dto.ClientResponse, error)
	GetClientDropdown() ([]*dto.DropdownItem, error)
	GetSummaries() ([]*dto.ClientSummaryResponse, error)
}

type clientService struct {
	clientRepo repository.ClientRepository
	farmRepo   repository.FarmRepository
	pondRepo   repository.PondRepository
	userRepo   repository.UserRepository
}

func NewClientService(
	clientRepo repository.ClientRepository,
	farmRepo repository.FarmRepository,
	pondRepo repository.PondRepository,
	userRepo repository.UserRepository,
) ClientService {
	return &clientService{
		clientRepo: clientRepo,
		farmRepo:   farmRepo,
		pondRepo:   pondRepo,
		userRepo:   userRepo,
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
		Name:                    request.Name,
		OwnerName:               request.OwnerName,
		ContactNumber:           request.ContactNumber,
		IsActive:                true,
		IsTouristFishingEnabled: false,
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

func (s *clientService) Update(ctx context.Context, request dto.UpdateClientRequest, username string) error {
	// Check if client exists
	existingClient, err := s.clientRepo.GetByID(request.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}

	if existingClient == nil {
		return errors.ErrClientNotFound
	}

	if request.Name != "" {
		existingClient.Name = request.Name
	}
	if request.OwnerName != "" {
		existingClient.OwnerName = request.OwnerName
	}
	if request.ContactNumber != "" {
		existingClient.ContactNumber = request.ContactNumber
	}
	if request.IsActive != nil {
		existingClient.IsActive = *request.IsActive
	}
	if request.IsTouristFishingEnabled != nil {
		existingClient.IsTouristFishingEnabled = *request.IsTouristFishingEnabled
	}

	// Update client (UpdatedBy set via BaseModel hook from ctx)
	if err := s.clientRepo.Update(ctx, existingClient); err != nil {
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

// GetSummaries returns every client with aggregate counts of farms, ponds,
// and users. Counts are fetched in three GROUP BY queries (no N+1).
func (s *clientService) GetSummaries() ([]*dto.ClientSummaryResponse, error) {
	clients, err := s.clientRepo.List()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	farmCounts, err := s.farmRepo.CountAllByClient()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	pondCounts, err := s.pondRepo.CountAllByClient()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	userCounts, err := s.userRepo.CountAllByClient()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	farmByClient := make(map[int]int64, len(farmCounts))
	for _, row := range farmCounts {
		farmByClient[row.ClientId] = row.Total
	}
	pondByClient := make(map[int]int64, len(pondCounts))
	for _, row := range pondCounts {
		pondByClient[row.ClientId] = row.Total
	}
	userByClient := make(map[int]int64, len(userCounts))
	for _, row := range userCounts {
		userByClient[row.ClientId] = row.Total
	}

	summaries := make([]*dto.ClientSummaryResponse, 0, len(clients))
	for _, client := range clients {
		summaries = append(summaries, &dto.ClientSummaryResponse{
			Id:            client.Id,
			Name:          client.Name,
			OwnerName:     client.OwnerName,
			ContactNumber: client.ContactNumber,
			IsActive:      client.IsActive,
			FarmCount:     farmByClient[client.Id],
			PondCount:     pondByClient[client.Id],
			UserCount:     userByClient[client.Id],
		})
	}

	return summaries, nil
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
		Id:                      client.Id,
		Name:                    client.Name,
		OwnerName:               client.OwnerName,
		ContactNumber:           client.ContactNumber,
		IsActive:                client.IsActive,
		IsTouristFishingEnabled: client.IsTouristFishingEnabled,
		CreatedAt:               client.CreatedAt,
		CreatedBy:               client.CreatedBy,
		UpdatedAt:               client.UpdatedAt,
		UpdatedBy:               client.UpdatedBy,
	}
}
