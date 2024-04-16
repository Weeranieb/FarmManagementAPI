package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IPondRepository interface {
	Create(pond *models.Pond) (*models.Pond, error)
	TakeById(id int) (*models.Pond, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.Pond, error)
	Update(pond *models.Pond) error
}

type pondRepositoryImp struct {
	dbContext *gorm.DB
}

func NewPondRepository(db *gorm.DB) IPondRepository {
	return &pondRepositoryImp{
		dbContext: db,
	}
}

func (rp pondRepositoryImp) Create(request *models.Pond) (*models.Pond, error) {
	if err := rp.dbContext.Table(dbconst.TPond).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp pondRepositoryImp) TakeById(id int) (*models.Pond, error) {
	var result *models.Pond
	if err := rp.dbContext.Table(dbconst.TPond).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Pond TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp pondRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.Pond, error) {
	var result *models.Pond
	if err := rp.dbContext.Table(dbconst.TPond).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found Pond FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp pondRepositoryImp) Update(request *models.Pond) error {
	if err := rp.dbContext.Table(dbconst.TPond).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
