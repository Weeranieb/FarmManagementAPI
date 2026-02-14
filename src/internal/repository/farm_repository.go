package repository

import (
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FarmRepository --output=./mocks --outpkg=mocks --filename=farm_repository.go --structname=MockFarmRepository --with-expecter=false
type FarmRepository interface {
	Create(farm *model.Farm) error
	GetByID(id int) (*model.Farm, error)
	GetByNameAndClientId(name string, clientId int) (*model.Farm, error)
	Update(farm *model.Farm) error
	ListByClientId(clientId int) ([]*model.Farm, error)
	CountByClientId(clientId int) (*model.FarmCountByClientId, error)
}

type farmRepository struct {
	db *gorm.DB
}

func NewFarmRepository(db *gorm.DB) FarmRepository {
	return &farmRepository{db: db}
}

func (r *farmRepository) Create(farm *model.Farm) error {
	return r.db.Create(farm).Error
}

func (r *farmRepository) GetByID(id int) (*model.Farm, error) {
	var farm model.Farm
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&farm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &farm, nil
}

func (r *farmRepository) GetByNameAndClientId(name string, clientId int) (*model.Farm, error) {
	var farm model.Farm
	err := r.db.Where("name = ? AND client_id = ? AND deleted_at IS NULL", name, clientId).First(&farm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &farm, nil
}

func (r *farmRepository) Update(farm *model.Farm) error {
	return r.db.Save(farm).Error
}

func (r *farmRepository) ListByClientId(clientId int) ([]*model.Farm, error) {
	var farms []*model.Farm
	err := r.db.Where("client_id = ? AND deleted_at IS NULL", clientId).Find(&farms).Error
	return farms, err
}

func (r *farmRepository) CountByClientId(clientId int) (*model.FarmCountByClientId, error) {
	var result model.FarmCountByClientId
	err := r.db.Raw(
		"SELECT COUNT(*) as total, COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) as active_count FROM farms WHERE client_id = ? AND deleted_at IS NULL",
		constants.FarmStatusActive, clientId,
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
