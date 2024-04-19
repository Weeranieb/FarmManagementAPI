package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IBillRepository interface {
	Create(pond *models.Bill) (*models.Bill, error)
	TakeById(id int) (*models.Bill, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Bill, error)
	Update(pond *models.Bill) error
}

type billRepositoryImp struct {
	dbContext *gorm.DB
}

func NewBillRepository(db *gorm.DB) IBillRepository {
	return &billRepositoryImp{
		dbContext: db,
	}
}

func (rp billRepositoryImp) Create(request *models.Bill) (*models.Bill, error) {
	if err := rp.dbContext.Table(dbconst.TBill).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp billRepositoryImp) TakeById(id int) (*models.Bill, error) {
	var result *models.Bill
	if err := rp.dbContext.Table(dbconst.TBill).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Bill TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp billRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Bill, error) {
	var result *models.Bill
	if err := rp.dbContext.Table(dbconst.TBill).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Bill FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp billRepositoryImp) Update(request *models.Bill) error {
	if err := rp.dbContext.Table(dbconst.TBill).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
