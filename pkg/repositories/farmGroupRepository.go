package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"

	"gorm.io/gorm"
)

type IFarmGroupRepository interface {
	Create(request *models.FarmGroup) (*models.FarmGroup, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error)
	GetFarmList(farmGroupId int) (*[]models.Farm, error)
	Update(request *models.FarmGroup) error
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
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp FarmGroupRepository) FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where(query, args...).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (rp FarmGroupRepository) TakeById(id int) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where("\"Id\" = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (rp FarmGroupRepository) Update(request *models.FarmGroup) error {
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

func (rp FarmGroupRepository) GetFarmList(farmGroupId int) (*[]models.Farm, error) {
	var result []models.Farm

	// Execute the subquery to fetch distinct farm IDs
	var farmIDs []int
	if err := rp.dbContext.Table(dbconst.TFarmOnFarmGroup).Select("DISTINCT(\"FarmId\")").Where("\"FarmGroupId\" = (?)", farmGroupId).Where("\"DelFlag\" = ?", false).Pluck("FarmId", &farmIDs).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}

	// Fetch farms using the fetched farm IDs
	if err := rp.dbContext.Table(dbconst.TFarm).Where("\"Id\" IN (?)", farmIDs).Where("\"DelFlag\" = ?", false).Find(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}

	return &result, nil
}
