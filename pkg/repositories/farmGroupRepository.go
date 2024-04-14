package repositories

import (
	"boonmafarm/api/pkg/models"

	"gorm.io/gorm"
)

type IFarmGroupRepository interface {
	Create(request *models.FarmGroup) (*models.FarmGroup, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error)
	TakeById(id int) (*models.FarmGroup, error)
}

type FarmGroupRepository struct {
	dbContext *gorm.DB
}

func NewFarmGroupRepository(db *gorm.DB) IFarmGroupRepository {
	return &FarmGroupRepository{
		dbContext: db,
	}
}

func (rp FarmGroupRepository) Create(request *models.FarmGroup) (*models.FarmGroup, error) {
	if err := rp.dbContext.Table("FarmGroups").Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp FarmGroupRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table("FarmGroups").Where(query, args...).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (rp FarmGroupRepository) TakeById(id int) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table("FarmGroups").Where("\"Id\" = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
