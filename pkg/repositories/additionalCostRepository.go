package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IAdditionalCostRepository interface {
	Create(pond *models.AdditionalCost) (*models.AdditionalCost, error)
	TakeById(id int) (*models.AdditionalCost, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.AdditionalCost, error)
	Update(pond *models.AdditionalCost) error
	WithTrx(trxHandle *gorm.DB) IAdditionalCostRepository
}

type additionalCostRepositoryImp struct {
	dbContext *gorm.DB
}

func NewAdditionalCostRepository(db *gorm.DB) IAdditionalCostRepository {
	return &additionalCostRepositoryImp{
		dbContext: db,
	}
}

func (rp additionalCostRepositoryImp) Create(request *models.AdditionalCost) (*models.AdditionalCost, error) {
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp additionalCostRepositoryImp) TakeById(id int) (*models.AdditionalCost, error) {
	var result *models.AdditionalCost
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found AdditionalCost TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp additionalCostRepositoryImp) WithTrx(trxHandle *gorm.DB) IAdditionalCostRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}

	return &additionalCostRepositoryImp{
		dbContext: trxHandle,
	}
}

func (rp additionalCostRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.AdditionalCost, error) {
	var result *models.AdditionalCost
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found AdditionalCost FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp additionalCostRepositoryImp) Update(request *models.AdditionalCost) error {
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TAdditionalCost).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}
