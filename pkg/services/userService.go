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
	GetUserList(clientId int) ([]*models.User, error)
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
	checkUser, err := sv.UserRepo.FirstByQuery("\"Username\" = ? AND \"DelFlag\" = ?", request.Username, false)
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

func (sv userServiceImp) GetUser(id int) (*models.User, error) {
	return sv.UserRepo.TakeById(id)
}

func (sv userServiceImp) Update(request *models.User, userIdentity string) error {
	// update user
	request.UpdatedBy = userIdentity
	if err := sv.UserRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv userServiceImp) GetUserList(clientId int) ([]*models.User, error) {
	return sv.UserRepo.TakeAll(clientId)
}
