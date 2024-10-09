package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"

	"gorm.io/gorm"
)

type IFarmGroupRepository interface {
	Create(request *models.FarmGroup) (*models.FarmGroup, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error)
	GetFarmList(farmGroupId int) (*[]models.Farm, error)
	TakeAllByClientId(clientId int) (*[]models.FarmGroup, error)
	Update(request *models.FarmGroup) error
	Delete(id int, username string) error
	TakeById(id int) (*models.FarmGroup, error)
}

type farmGroupRepositoryImp struct {
	dbContext *gorm.DB
}

func NewFarmGroupRepository(db *gorm.DB) IFarmGroupRepository {
	return &farmGroupRepositoryImp{
		dbContext: db,
	}
}

func (rp farmGroupRepositoryImp) Create(request *models.FarmGroup) (*models.FarmGroup, error) {
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp farmGroupRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where(query, args...).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (rp farmGroupRepositoryImp) TakeById(id int) (*models.FarmGroup, error) {
	var result *models.FarmGroup
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (rp farmGroupRepositoryImp) TakeAllByClientId(clientId int) (*[]models.FarmGroup, error) {
	var result []models.FarmGroup
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where("\"ClientId\" = ? AND \"DelFlag\" = ?", clientId, false).Find(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}
	return &result, nil
}

func (rp farmGroupRepositoryImp) Update(request *models.FarmGroup) error {
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}

func (rp farmGroupRepositoryImp) GetFarmList(farmGroupId int) (*[]models.Farm, error) {
	var result []models.Farm

	// Execute a single query with join to fetch farms directly
	if err := rp.dbContext.Table(dbconst.TFarm).
		Joins("JOIN "+dbconst.TFarmOnFarmGroup+" ON "+dbconst.TFarmOnFarmGroup+".\"FarmId\" = "+dbconst.TFarm+".\"Id\"").
		Where(dbconst.TFarmOnFarmGroup+".\"FarmGroupId\" = ?", farmGroupId).
		Where(dbconst.TFarmOnFarmGroup+".\"DelFlag\" = ?", false).
		Where(dbconst.TFarm+".\"DelFlag\" = ?", false).
		Find(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}

	return &result, nil
}

func (rp farmGroupRepositoryImp) Delete(id int, username string) error {
	if err := rp.dbContext.Table(dbconst.TFarmGroup).Where("\"Id\" = ?", id).Updates(map[string]interface{}{"\"DelFlag\"": true, "\"UpdatedBy\"": username}).Error; err != nil {
		return err
	}
	return nil
}
