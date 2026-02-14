package repository

import (
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=PondRepository --output=./mocks --outpkg=mocks --filename=pond_repository.go --structname=MockPondRepository --with-expecter=false
type PondRepository interface {
	Create(pond *model.Pond) error
	CreateBatch(ponds []*model.Pond) error
	GetByID(id int) (*model.Pond, error)
	GetByFarmIdAndName(farmId int, name string) (*model.Pond, error)
	Update(pond *model.Pond) error
	ListByFarmId(farmId int) ([]*model.Pond, error)
	Delete(id int) error
}

type pondRepository struct {
	db *gorm.DB
}

func NewPondRepository(db *gorm.DB) PondRepository {
	return &pondRepository{db: db}
}

func (r *pondRepository) Create(pond *model.Pond) error {
	return r.db.Create(pond).Error
}

func (r *pondRepository) CreateBatch(ponds []*model.Pond) error {
	return r.db.Create(ponds).Error
}

func (r *pondRepository) GetByID(id int) (*model.Pond, error) {
	var pond model.Pond
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pond).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pond, nil
}

func (r *pondRepository) GetByFarmIdAndName(farmId int, name string) (*model.Pond, error) {
	var pond model.Pond
	err := r.db.Where("farm_id = ? AND name = ? AND deleted_at IS NULL", farmId, name).First(&pond).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pond, nil
}

func (r *pondRepository) Update(pond *model.Pond) error {
	return r.db.Save(pond).Error
}

func (r *pondRepository) ListByFarmId(farmId int) ([]*model.Pond, error) {
	var ponds []*model.Pond
	err := r.db.Where("farm_id = ? AND deleted_at IS NULL", farmId).Find(&ponds).Error
	return ponds, err
}

func (r *pondRepository) Delete(id int) error {
	return r.db.Delete(&model.Pond{}, id).Error
}

