package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IWorkerService interface {
	Create(request models.AddWorker, userIdentity string) (*models.Worker, error)
	Get(id int) (*models.Worker, error)
	Update(request *models.Worker, userIdentity string) error
}

type workerServiceImp struct {
	WorkerRepo repositories.IWorkerRepository
}

func NewWorkerService(worker repositories.IWorkerRepository) IWorkerService {
	return &workerServiceImp{
		WorkerRepo: worker,
	}
}

func (sv workerServiceImp) Create(request models.AddWorker, userIdentity string) (*models.Worker, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check pond if exist
	checkPond, err := sv.WorkerRepo.FirstByQuery("\"FarmGroupId\" = ? AND \"DelFlag\" = ?", request.FarmGroupId, false)
	if err != nil {
		return nil, err
	}

	if checkPond != nil {
		return nil, errors.New("worker already exist")
	}

	newWorker := &models.Worker{}
	request.Transfer(newWorker)
	// set is active
	newWorker.IsActive = true
	newWorker.UpdatedBy = userIdentity
	newWorker.CreatedBy = userIdentity

	// create user
	newWorker, err = sv.WorkerRepo.Create(newWorker)
	if err != nil {
		return nil, err
	}

	return newWorker, nil
}

func (sv workerServiceImp) Get(id int) (*models.Worker, error) {
	return sv.WorkerRepo.TakeById(id)
}

func (sv workerServiceImp) Update(request *models.Worker, userIdentity string) error {
	// update worker
	request.UpdatedBy = userIdentity
	if err := sv.WorkerRepo.Update(request); err != nil {
		return err
	}
	return nil
}