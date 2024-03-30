package repositories

import (
	"boonmafarm/api/pkg/models"
	"fmt"

	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) (*models.User, error)
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

func (rp userRepositoryImp) WithTrx(trxHandle *gorm.DB) IUserRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}
	rp.dbContext = trxHandle
	return rp
}
