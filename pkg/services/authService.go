package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Create(request models.AddUser) (*models.User, error)
	Login(request models.Login) (string, error)
}

type authServiceImp struct {
	UserRepo repositories.IUserRepository
}

func NewAuthService(userRepo repositories.IUserRepository) IAuthService {
	return &authServiceImp{
		UserRepo: userRepo,
	}
}

func (sv authServiceImp) Create(request models.AddUser) (*models.User, error) {
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

	// create user
	newUser, err = sv.UserRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (sv authServiceImp) Login(request models.Login) (string, error) {
	// check user if exist
	checkUser, err := sv.UserRepo.FirstByQuery("Username = ?", request.Username)
	if err != nil {
		return "", err
	}

	if checkUser == nil {
		return "", errors.New("user or password is incorrect")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(request.Password))
	if err != nil {
		return "", errors.New("user or password is incorrect")
	}

	// create jwt token
	secretKey := viper.GetString("jwt.secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = checkUser.Username
	claims["clientId"] = checkUser.ClientId
	claims["isAdmin"] = checkUser.IsAdmin
	claims["exp"] = jwt.TimeFunc().AddDate(0, 0, 1).Unix()

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
