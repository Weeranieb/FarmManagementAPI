package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(request models.AddUser, userIdentity string, clientId int) (*models.User, error)
}

type userServiceImp struct {
	UserRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) IUserService {
	return &userServiceImp{
		UserRepo: userRepo,
	}
}

func (sv userServiceImp) Create(request models.AddUser, userIdentity string, clientId int) (*models.User, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check user if exist
	checkUser, err := sv.UserRepo.FirstByQuery("Username = ?", request.Username)
	if err != nil {
		return nil, err
	}

	if checkUser != nil {
		return nil, errors.New("user already exist")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{}
	request.Transfer(newUser)
	newUser.Password = string(hashedPassword)
	newUser.ClientId = clientId
	newUser.UpdatedBy = userIdentity
	newUser.CreatedBy = userIdentity

	// create user
	newUser, err = sv.UserRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
