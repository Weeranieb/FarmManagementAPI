package repositories

import (
	"boonmafarm/api/pkg/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IClientRepository interface {
	Create(request *models.Client) (*models.Client, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Client, error)
	WithTrx(trxHandle *gorm.DB) IClientRepository
}

type ClientRepositoryImp struct {
	dbContext *gorm.DB
}

func NewClientRepository(db *gorm.DB) IClientRepository {
	return &ClientRepositoryImp{
		dbContext: db,
	}
}

func (rp ClientRepositoryImp) Create(request *models.Client) (*models.Client, error) {
	if err := rp.dbContext.Table("Clients").Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp ClientRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Client, error) {
	var result *models.Client
	if err := rp.dbContext.Table("Clients").Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Clien FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp ClientRepositoryImp) WithTrx(trxHandle *gorm.DB) IClientRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}
