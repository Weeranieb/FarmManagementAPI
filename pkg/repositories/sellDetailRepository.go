package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ISellDetailRepository interface {
	Create(request *models.SellDetail) (*models.SellDetail, error)
	BulkCreate(request []models.SellDetail) ([]models.SellDetail, error)
	TakeById(id int) (*models.SellDetail, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.SellDetail, error)
	ListByQuery(query interface{}, args ...interface{}) ([]models.SellDetail, error)
	Update(request *models.SellDetail) error
	WithTrx(trxHandle *gorm.DB) ISellDetailRepository
}

type sellDetailRepositoryImp struct {
	dbContext *gorm.DB
}

func NewSellDetailRepository(db *gorm.DB) ISellDetailRepository {
	return &sellDetailRepositoryImp{
		dbContext: db,
	}
}

func (rp sellDetailRepositoryImp) Create(request *models.SellDetail) (*models.SellDetail, error) {
	if err := rp.dbContext.Table(dbconst.TSellDetail).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp sellDetailRepositoryImp) BulkCreate(request []models.SellDetail) ([]models.SellDetail, error) {
	if err := rp.dbContext.Table(dbconst.TSellDetail).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp sellDetailRepositoryImp) TakeById(id int) (*models.SellDetail, error) {
	var result *models.SellDetail
	if err := rp.dbContext.Table(dbconst.TSellDetail).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Activity TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp sellDetailRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.SellDetail, error) {
	var result *models.SellDetail
	if err := rp.dbContext.Table(dbconst.TSellDetail).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Activity FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp sellDetailRepositoryImp) WithTrx(trxHandle *gorm.DB) ISellDetailRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}

func (rp sellDetailRepositoryImp) Update(request *models.SellDetail) error {
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TSellDetail).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}

func (rp sellDetailRepositoryImp) ListByQuery(query interface{}, args ...interface{}) ([]models.SellDetail, error) {
	var result []models.SellDetail
	if err := rp.dbContext.Table(dbconst.TSellDetail).Where(query, args...).Find(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Activity ListByQuery", query)
		return nil, nil
	}
	return result, nil
}
