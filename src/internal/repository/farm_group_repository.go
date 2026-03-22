package repository

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

type FarmGroupRepository interface {
	Create(ctx context.Context, farmGroup *model.FarmGroup) error
	CreateJoins(ctx context.Context, joins []*model.FarmOnFarmGroup) error
	GetByID(id int) (*model.FarmGroup, error)
	GetFarmsByFarmGroupId(farmGroupId int) ([]*model.Farm, error)
	ListByClientId(clientId int) ([]*model.FarmGroup, error)
	Update(ctx context.Context, farmGroup *model.FarmGroup) error
	DeleteJoinsByFarmGroupId(ctx context.Context, farmGroupId int) error
}

type farmGroupRepository struct {
	db *gorm.DB
}

func NewFarmGroupRepository(db *gorm.DB) FarmGroupRepository {
	return &farmGroupRepository{db: db}
}

func (r *farmGroupRepository) Create(ctx context.Context, farmGroup *model.FarmGroup) error {
	return r.db.WithContext(ctx).Create(farmGroup).Error
}

func (r *farmGroupRepository) CreateJoins(ctx context.Context, joins []*model.FarmOnFarmGroup) error {
	if len(joins) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(joins).Error
}

func (r *farmGroupRepository) GetByID(id int) (*model.FarmGroup, error) {
	var fg model.FarmGroup
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&fg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &fg, nil
}

func (r *farmGroupRepository) GetFarmsByFarmGroupId(farmGroupId int) ([]*model.Farm, error) {
	var farms []*model.Farm
	err := r.db.
		Joins("JOIN farm_on_farm_group fofg ON fofg.farm_id = farms.id AND fofg.deleted_at IS NULL").
		Where("fofg.farm_group_id = ? AND farms.deleted_at IS NULL", farmGroupId).
		Find(&farms).Error
	if err != nil {
		return nil, err
	}
	return farms, nil
}

func (r *farmGroupRepository) ListByClientId(clientId int) ([]*model.FarmGroup, error) {
	var list []*model.FarmGroup
	err := r.db.Where("client_id = ? AND deleted_at IS NULL", clientId).
		Order("name ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *farmGroupRepository) Update(ctx context.Context, farmGroup *model.FarmGroup) error {
	return r.db.WithContext(ctx).Save(farmGroup).Error
}

func (r *farmGroupRepository) DeleteJoinsByFarmGroupId(ctx context.Context, farmGroupId int) error {
	return r.db.WithContext(ctx).
		Where("farm_group_id = ?", farmGroupId).
		Delete(&model.FarmOnFarmGroup{}).Error
}
