package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IActivePondRepository interface {
	Create(activePond *models.ActivePond) (*models.ActivePond, error)
	TakeById(id int) (*models.ActivePond, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.ActivePond, error)
	Update(activePond *models.ActivePond) error
}

type activePondRepositoryImp struct {
	dbContext *gorm.DB
}

func NewActivePondRepository(db *gorm.DB) IActivePondRepository {
	return &activePondRepositoryImp{
		dbContext: db,
	}
}

func (rp activePondRepositoryImp) Create(request *models.ActivePond) (*models.ActivePond, error) {
	if err := rp.dbContext.Table(dbconst.TActivePond).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp activePondRepositoryImp) TakeById(id int) (*models.ActivePond, error) {
	var result *models.ActivePond
	if err := rp.dbContext.Table(dbconst.TActivePond).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Active Pond TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp activePondRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.ActivePond, error) {
	var result *models.ActivePond
	if err := rp.dbContext.Table(dbconst.TActivePond).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found active Pond FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp activePondRepositoryImp) Update(request *models.ActivePond) error {
	if err := rp.dbContext.Table(dbconst.TActivePond).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
