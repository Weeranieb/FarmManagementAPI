package repositories

import (
	"boonmafarm/api/pkg/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) (*models.User, error)
	TakeById(id int) (*models.User, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.User, error)
	Update(user *models.User) error
	WithTrx(trxHandle *gorm.DB) IUserRepository
}

type userRepositoryImp struct {
	dbContext *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepositoryImp{
		dbContext: db,
	}
}

func (rp userRepositoryImp) Create(request *models.User) (*models.User, error) {
	if err := rp.dbContext.Table("Users").Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp userRepositoryImp) TakeById(id int) (*models.User, error) {
	var result *models.User
	if err := rp.dbContext.Table("Users").Where("\"Id\" = ?", id).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found User TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp userRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.User, error) {
	var result *models.User
	if err := rp.dbContext.Table("Users").Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found User FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp userRepositoryImp) WithTrx(trxHandle *gorm.DB) IUserRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}

func (rp userRepositoryImp) Update(request *models.User) error {
	if err := rp.dbContext.Table("Users").Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
