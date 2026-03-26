package service

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=WorkerService --output=./mocks --outpkg=service --filename=worker_service.go --structname=MockWorkerService --with-expecter=false
type WorkerService interface {
	Create(ctx context.Context, request dto.CreateWorkerRequest, username string, clientId int) (*dto.WorkerResponse, error)
	Get(id int) (*dto.WorkerResponse, error)
	Update(ctx context.Context, request dto.UpdateWorkerRequest, username string) error
	GetPage(clientId, page, pageSize int, orderBy, keyword string) (*dto.PageResponse, error)
}

type workerService struct {
	workerRepo repository.WorkerRepository
}

func NewWorkerService(workerRepo repository.WorkerRepository) WorkerService {
	return &workerService{
		workerRepo: workerRepo,
	}
}

func (s *workerService) Create(ctx context.Context, request dto.CreateWorkerRequest, username string, clientId int) (*dto.WorkerResponse, error) {
	// Check if worker already exists (by FarmGroupId - this seems odd but matches old logic)
	checkWorker, err := s.workerRepo.GetByFarmGroupId(request.FarmGroupId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkWorker != nil {
		return nil, errors.ErrWorkerAlreadyExists
	}

	newWorker := &model.Worker{
		ClientId:      clientId,
		FarmGroupId:   request.FarmGroupId,
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		ContactNumber: request.ContactNumber,
		Nationality:   request.Nationality,
		Salary:        decimal.NewFromFloat(request.Salary),
		HireDate:      request.HireDate,
		IsActive:      true,
	}

	// Create worker (CreatedBy/UpdatedBy set via BaseModel hook from ctx)
	err = s.workerRepo.Create(ctx, newWorker)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toWorkerResponse(newWorker), nil
}

func (s *workerService) Get(id int) (*dto.WorkerResponse, error) {
	worker, err := s.workerRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if worker == nil {
		return nil, errors.ErrWorkerNotFound
	}

	return s.toWorkerResponse(worker), nil
}

func (s *workerService) Update(ctx context.Context, request dto.UpdateWorkerRequest, username string) error {
	existingWorker, err := s.workerRepo.GetByID(request.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existingWorker == nil {
		return errors.ErrWorkerNotFound
	}

	if request.FarmGroupId != 0 {
		existingWorker.FarmGroupId = request.FarmGroupId
	}
	if request.FirstName != "" {
		existingWorker.FirstName = request.FirstName
	}
	if request.LastName != nil {
		existingWorker.LastName = request.LastName
	}
	if request.ContactNumber != nil {
		existingWorker.ContactNumber = request.ContactNumber
	}
	if request.Nationality != "" {
		existingWorker.Nationality = request.Nationality
	}
	if request.Salary != 0 {
		existingWorker.Salary = decimal.NewFromFloat(request.Salary)
	}
	if request.HireDate != nil {
		existingWorker.HireDate = request.HireDate
	}
	if request.IsActive != nil {
		existingWorker.IsActive = *request.IsActive
	}

	// Update worker (UpdatedBy set via BaseModel hook from ctx)
	if err := s.workerRepo.Update(ctx, existingWorker); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *workerService) GetPage(clientId, page, pageSize int, orderBy, keyword string) (*dto.PageResponse, error) {
	workers, total, err := s.workerRepo.GetPage(clientId, page, pageSize, orderBy, keyword)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.WorkerResponse, 0, len(workers))
	for _, worker := range workers {
		responses = append(responses, s.toWorkerResponse(worker))
	}

	return &dto.PageResponse{
		Items: responses,
		Total: total,
	}, nil
}

func (s *workerService) toWorkerResponse(worker *model.Worker) *dto.WorkerResponse {
	return &dto.WorkerResponse{
		Id:            worker.Id,
		ClientId:      worker.ClientId,
		FarmGroupId:   worker.FarmGroupId,
		FirstName:     worker.FirstName,
		LastName:      worker.LastName,
		ContactNumber: worker.ContactNumber,
		Nationality:   worker.Nationality,
		Salary:        worker.Salary.InexactFloat64(),
		HireDate:      worker.HireDate,
		IsActive:      worker.IsActive,
		CreatedAt:     worker.CreatedAt,
		CreatedBy:     worker.CreatedBy,
		UpdatedAt:     worker.UpdatedAt,
		UpdatedBy:     worker.UpdatedBy,
	}
}
