package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
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
	TakeAll(clientId int) ([]*models.User, error)
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
	if err := rp.dbContext.Table(dbconst.TUser).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp userRepositoryImp) TakeById(id int) (*models.User, error) {
	var result *models.User
	if err := rp.dbContext.Table(dbconst.TUser).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
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
	if err := rp.dbContext.Table(dbconst.TUser).Where(query, args...).First(&result).Error; err != nil {
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
	if err := rp.dbContext.Table(dbconst.TUser).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

func (rp userRepositoryImp) TakeAll(clientId int) ([]*models.User, error) {
	var result []*models.User
	if err := rp.dbContext.Table(dbconst.TUser).Where("\"ClientId\" = ? AND \"DelFlag\" = ?", clientId, false).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
