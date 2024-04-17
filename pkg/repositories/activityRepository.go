package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IActivityRepository interface {
	Create(request *models.Activity) (*models.Activity, error)
	TakeById(id int) (*models.Activity, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Activity, error)
	Update(request *models.Activity) error
	WithTrx(trxHandle *gorm.DB) IActivityRepository
}

type activityRepositoryImp struct {
	dbContext *gorm.DB
}

func NewActivityRepository(db *gorm.DB) IActivityRepository {
	return &activityRepositoryImp{
		dbContext: db,
	}
}

func (rp activityRepositoryImp) Create(request *models.Activity) (*models.Activity, error) {
	if err := rp.dbContext.Table(dbconst.TActivitiy).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp activityRepositoryImp) TakeById(id int) (*models.Activity, error) {
	var result *models.Activity
	if err := rp.dbContext.Table(dbconst.TActivitiy).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Activity TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp activityRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Activity, error) {
	var result *models.Activity
	if err := rp.dbContext.Table(dbconst.TActivitiy).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Activity FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp activityRepositoryImp) WithTrx(trxHandle *gorm.DB) IActivityRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}

func (rp activityRepositoryImp) Update(request *models.Activity) error {
	if err := rp.dbContext.Table(dbconst.TActivitiy).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
