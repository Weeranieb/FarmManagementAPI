package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(request models.AddUser, userIdentity string, clientId int) (*models.User, error)
	GetUser(id int) (*models.User, error)
	Update(request *models.User, userIdentity string) error
}

type UserServiceImp struct {
	UserRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) IUserService {
	return &UserServiceImp{
		UserRepo: userRepo,
	}
}

func (sv UserServiceImp) Create(request models.AddUser, userIdentity string, clientId int) (*models.User, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check user if exist
	checkUser, err := sv.UserRepo.FirstByQuery("\"Username\" = ?", request.Username)
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

func (sv UserServiceImp) GetUser(id int) (*models.User, error) {
	return sv.UserRepo.TakeById(id)
}

func (sv UserServiceImp) Update(request *models.User, userIdentity string) error {
	// update user
	request.UpdatedBy = userIdentity
	if err := sv.UserRepo.Update(request); err != nil {
		return err
	}
	return nil
}
