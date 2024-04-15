package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"

	"gorm.io/gorm"
)

type IFarmOnFarmGroupRepository interface {
	Create(request *models.FarmOnFarmGroup) (*models.FarmOnFarmGroup, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FarmOnFarmGroup, error)
	Delete(id int) error
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
	if err := rp.dbContext.Table(dbconst.TFarmOnFarmGroup).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp FarmOnFarmGroupRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.FarmOnFarmGroup, error) {
	var result *models.FarmOnFarmGroup
	if err := rp.dbContext.Table(dbconst.TFarmOnFarmGroup).Where(query, args...).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (rp FarmOnFarmGroupRepository) Delete(id int) error {
	if err := rp.dbContext.Table(dbconst.TFarmOnFarmGroup).Where("\"Id\" = ?", id).Update("\"DelFlag\"", true).Error; err != nil {
		return err
	}
	return nil
}
