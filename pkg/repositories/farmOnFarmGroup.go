package repositories

import (
	"boonmafarm/api/pkg/models"

	"gorm.io/gorm"
)

type IFarmOnFarmGroupRepository interface {
	Create(request *models.FarmOnFarmGroup) (*models.FarmOnFarmGroup, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FarmOnFarmGroup, error)
	TakeById(id int) (*models.FarmOnFarmGroup, error)
}

type FarmOnFarmGroupRepository struct {
	dbContext *gorm.DB
}

func NewFarmOnFarmGroupRepository(db *gorm.DB) IFarmOnFarmGroupRepository {
	return &FarmOnFarmGroupRepository{
		dbContext: db,
	}
}

func (rp FarmOnFarmGroupRepository) Create(request *models.FarmOnFarmGroup) (*models.FarmOnFarmGroup, error) {
	if err := rp.dbContext.Table("FarmOnFarmGroups").Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp FarmOnFarmGroupRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.FarmOnFarmGroup, error) {
	var result *models.FarmOnFarmGroup
	if err := rp.dbContext.Table("FarmOnFarmGroups").Where(query, args...).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (rp FarmOnFarmGroupRepository) TakeById(id int) (*models.FarmOnFarmGroup, error) {
	var result *models.FarmOnFarmGroup
	if err := rp.dbContext.Table("FarmOnFarmGroups").Where("\"Id\" = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
