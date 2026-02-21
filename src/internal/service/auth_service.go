package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=AuthService --output=./mocks --outpkg=service --filename=auth_service.go --structname=MockAuthService --with-expecter=false
type AuthService interface {
	Register(request dto.RegisterRequest) (*dto.UserResponse, error)
	Login(request dto.LoginRequest) (string, *dto.UserResponse, *time.Time, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(request dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	checkUser, err := s.userRepo.GetByUsername(request.Username)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkUser != nil {
		return nil, errors.ErrAuthUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	newUser := &model.User{
		ClientId:      &request.ClientId,
		Username:      request.Username,
		Password:      string(hashedPassword),
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		UserLevel:     request.UserLevel,
		ContactNumber: request.ContactNumber,
		BaseModel: model.BaseModel{
			CreatedBy: request.Username,
			UpdatedBy: request.Username,
		},
	}

	// Create user (no request user context for registration; CreatedBy set on model)
	err = s.userRepo.Create(context.Background(), newUser)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toUserResponse(newUser), nil
}

func (s *authService) Login(request dto.LoginRequest) (string, *dto.UserResponse, *time.Time, error) {
	// Check if user exists
	checkUser, err := s.userRepo.GetByUsername(request.Username)
	if err != nil {
		return "", nil, nil, errors.ErrGeneric.Wrap(err)
	}

	if checkUser == nil {
		return "", nil, nil, errors.ErrAuthInvalidCredentials
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(request.Password))
	if err != nil {
		return "", nil, nil, errors.ErrAuthInvalidCredentials
	}

	// Create JWT token
	secretKey := viper.GetString("authentication.jwt_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Custom claims
	claims["username"] = checkUser.Username
	claims["userId"] = checkUser.Id
	claims["clientId"] = checkUser.ClientId
	claims["userLevel"] = checkUser.UserLevel

	var expiredDate time.Time
	if request.RememberMe {
		expiredDate = time.Now().AddDate(0, 1, 0) // 30 days
	} else {
		expiredDate = time.Now().AddDate(0, 0, 1) // 1 day
	}
	claims["exp"] = expiredDate.Unix()

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, nil, errors.ErrGeneric.Wrap(err)
	}

	return tokenString, s.toUserResponse(checkUser), &expiredDate, nil
}

func (s *authService) toUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:            user.Id,
		ClientId:      user.ClientId,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		UserLevel:     user.UserLevel,
		ContactNumber: user.ContactNumber,
		CreatedAt:     user.CreatedAt,
		CreatedBy:     user.CreatedBy,
		UpdatedAt:     user.UpdatedAt,
		UpdatedBy:     user.UpdatedBy,
	}
}
