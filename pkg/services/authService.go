package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Create(request models.AddUser) (*models.User, error)
	Login(request models.Login) (string, *models.User, *time.Time, error)
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

	// create user
	res, err := sv.UserRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sv authServiceImp) Login(request models.Login) (string, *models.User, *time.Time, error) {
	// check user if exist
	checkUser, err := sv.UserRepo.FirstByQuery("\"Username\" = ? AND \"DelFlag\" = ?", request.Username, false)
	if err != nil {
		return "", nil, nil, err
	}

	if checkUser == nil {
		return "", nil, nil, errors.New("user or password is incorrect")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(request.Password))
	if err != nil {
		return "", nil, nil, errors.New("user or password is incorrect")
	}

	// create jwt token
	secretKey := viper.GetString("authentication.jwt_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// custom claims
	claims["username"] = checkUser.Username
	claims["userId"] = checkUser.Id
	claims["clientId"] = checkUser.ClientId
	claims["userLevel"] = checkUser.UserLevel

	// set expiration time
	expiredDate := jwt.TimeFunc().AddDate(0, 0, 1)
	claims["exp"] = expiredDate.Unix()

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, nil, err
	}

	return tokenString, checkUser, &expiredDate, nil
}
