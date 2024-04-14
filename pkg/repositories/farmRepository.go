package repositories

import (
	"boonmafarm/api/pkg/models"

	"gorm.io/gorm"
)

type IFarmRepository interface {
	Create(request *models.Farm) (*models.Farm, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Farm, error)
	Update(request *models.Farm) error
	TakeById(id int) (*models.Farm, error)
}

type FarmRepositoryImp struct {
	dbContext *gorm.DB
}

func NewFarmRepository(db *gorm.DB) IFarmRepository {
	return &FarmRepositoryImp{
		dbContext: db,
	}
}

func (rp FarmRepositoryImp) Create(request *models.Farm) (*models.Farm, error) {
	if err := rp.dbContext.Table("Farms").Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp FarmRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Farm, error) {
	var result *models.Farm
	if err := rp.dbContext.Table("Farms").Where(query, args...).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (rp FarmRepositoryImp) TakeById(id int) (*models.Farm, error) {
	var result *models.Farm
	if err := rp.dbContext.Table("Farms").Where("\"Id\" = ?", id).Take(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (rp FarmRepositoryImp) Update(request *models.Farm) error {
	if err := rp.dbContext.Table("Farms").Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
