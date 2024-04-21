package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IMerchantRepository interface {
	Create(merchant *models.Merchant) (*models.Merchant, error)
	TakeById(id int) (*models.Merchant, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Merchant, error)
	Update(merchant *models.Merchant) error
}

type merchantRepositoryImp struct {
	dbContext *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) IMerchantRepository {
	return &merchantRepositoryImp{
		dbContext: db,
	}
}

func (rp merchantRepositoryImp) Create(request *models.Merchant) (*models.Merchant, error) {
	if err := rp.dbContext.Table(dbconst.TMerchant).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp merchantRepositoryImp) TakeById(id int) (*models.Merchant, error) {
	var result *models.Merchant
	if err := rp.dbContext.Table(dbconst.TMerchant).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Merchant TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp merchantRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Merchant, error) {
	var result *models.Merchant
	if err := rp.dbContext.Table(dbconst.TMerchant).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Merchant FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp merchantRepositoryImp) Update(request *models.Merchant) error {
	if err := rp.dbContext.Table(dbconst.TMerchant).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
