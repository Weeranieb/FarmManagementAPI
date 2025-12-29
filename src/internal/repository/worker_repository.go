package repository

import (
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=WorkerRepository --output=./mocks --outpkg=mocks --filename=worker_repository.go --structname=MockWorkerRepository --with-expecter=false
type WorkerRepository interface {
	Create(worker *model.Worker) error
	GetByID(id int) (*model.Worker, error)
	GetByFarmGroupId(farmGroupId int) (*model.Worker, error)
	Update(worker *model.Worker) error
	GetPage(clientId, page, pageSize int, orderBy, keyword string) ([]*model.Worker, int64, error)
}

type workerRepository struct {
	db *gorm.DB
}

func NewWorkerRepository(db *gorm.DB) WorkerRepository {
	return &workerRepository{db: db}
}

func (r *workerRepository) Create(worker *model.Worker) error {
	return r.db.Create(worker).Error
}

func (r *workerRepository) GetByID(id int) (*model.Worker, error) {
	var worker model.Worker
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&worker).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &worker, nil
}

func (r *workerRepository) GetByFarmGroupId(farmGroupId int) (*model.Worker, error) {
	var worker model.Worker
	err := r.db.Where("farm_group_id = ? AND deleted_at IS NULL", farmGroupId).First(&worker).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &worker, nil
}

func (r *workerRepository) Update(worker *model.Worker) error {
	return r.db.Save(worker).Error
}

func (r *workerRepository) GetPage(clientId, page, pageSize int, orderBy, keyword string) ([]*model.Worker, int64, error) {
	var workers []*model.Worker
	var total int64

	query := r.db.Model(&model.Worker{}).Where("client_id = ? AND deleted_at IS NULL", clientId)

	if keyword != "" {
		query = query.Where("(first_name LIKE ? OR last_name LIKE ? OR nationality LIKE ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	// Apply pagination
	offset := page * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&workers).Error; err != nil {
		return nil, 0, err
	}

	return workers, total, nil
}

