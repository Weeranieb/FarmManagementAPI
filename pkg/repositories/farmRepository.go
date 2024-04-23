package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"

	"gorm.io/gorm"
)

type IFarmRepository interface {
	Create(request *models.Farm) (*models.Farm, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Farm, error)
	Update(request *models.Farm) error
	TakeById(id int) (*models.Farm, error)
	TakeAll(clientId int) ([]*models.Farm, error)
}

type farmRepositoryImp struct {
	dbContext *gorm.DB
}

func NewFarmRepository(db *gorm.DB) IFarmRepository {
	return &farmRepositoryImp{
		dbContext: db,
	}
}

func (rp farmRepositoryImp) Create(request *models.Farm) (*models.Farm, error) {
	if err := rp.dbContext.Table(dbconst.TFarm).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp farmRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Farm, error) {
	var result *models.Farm
	if err := rp.dbContext.Table(dbconst.TFarm).Where(query, args...).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (rp farmRepositoryImp) TakeById(id int) (*models.Farm, error) {
	var result *models.Farm
	if err := rp.dbContext.Table(dbconst.TFarm).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (rp farmRepositoryImp) Update(request *models.Farm) error {
	if err := rp.dbContext.Table(dbconst.TFarm).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

func (rp farmRepositoryImp) TakeAll(clientId int) ([]*models.Farm, error) {
	var result []*models.Farm
	if err := rp.dbContext.Table(dbconst.TFarm).Where("\"ClientId\" = ? AND \"DelFlag\" = ?", clientId, false).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
