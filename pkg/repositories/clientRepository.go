package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IClientRepository interface {
	Create(request *models.Client) (*models.Client, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Client, error)
	Update(request *models.Client) error
	TakeById(id int) (*models.Client, error)
	WithTrx(trxHandle *gorm.DB) IClientRepository
}

type clientRepositoryImp struct {
	dbContext *gorm.DB
}

func NewClientRepository(db *gorm.DB) IClientRepository {
	return &clientRepositoryImp{
		dbContext: db,
	}
}

func (rp clientRepositoryImp) Create(request *models.Client) (*models.Client, error) {
	if err := rp.dbContext.Table(dbconst.TClient).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp clientRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Client, error) {
	var result *models.Client
	if err := rp.dbContext.Table(dbconst.TClient).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Clien FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp clientRepositoryImp) TakeById(id int) (*models.Client, error) {
	var result *models.Client
	if err := rp.dbContext.Table(dbconst.TClient).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Client TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp clientRepositoryImp) WithTrx(trxHandle *gorm.DB) IClientRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}

func (rp clientRepositoryImp) Update(request *models.Client) error {
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TClient).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}
