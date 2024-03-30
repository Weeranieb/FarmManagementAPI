package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"

	"gorm.io/gorm"
)

type IUserService interface {
	Create(request models.AddUsers, userIdentity string) (*models.User, error)
	WithTrx(trxHandle *gorm.DB) IUserService
}

type userServiceImp struct {
	UserRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) IUserService {
	return &userServiceImp{
		UserRepo: userRepo,
	}
}

func (sv userServiceImp) WithTrx(trxHandle *gorm.DB) IUserService {
	sv.UserRepo = sv.UserRepo.WithTrx(trxHandle)
	return sv
}

func (sv userServiceImp) Create(request models.AddUsers, userIdentity string) (*models.User, error) {
	// TODO implement business logic
	return sv.UserRepo.Create(&models.User{})
}
