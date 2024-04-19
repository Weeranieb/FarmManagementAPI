package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IWorkerRepository interface {
	Create(worker *models.Worker) (*models.Worker, error)
	TakeById(id int) (*models.Worker, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Worker, error)
	Update(worker *models.Worker) error
}

type workerRepositoryImp struct {
	dbContext *gorm.DB
}

func NewWorkerRepository(db *gorm.DB) IWorkerRepository {
	return &workerRepositoryImp{
		dbContext: db,
	}
}

func (rp workerRepositoryImp) Create(request *models.Worker) (*models.Worker, error) {
	if err := rp.dbContext.Table(dbconst.TWorker).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp workerRepositoryImp) TakeById(id int) (*models.Worker, error) {
	var result *models.Worker
	if err := rp.dbContext.Table(dbconst.TWorker).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Worker TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp workerRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Worker, error) {
	var result *models.Worker
	if err := rp.dbContext.Table(dbconst.TWorker).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Worker FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp workerRepositoryImp) Update(request *models.Worker) error {
	if err := rp.dbContext.Table(dbconst.TWorker).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
