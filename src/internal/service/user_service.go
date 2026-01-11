package service

import (
	"context"
	"fmt"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserService --output=./mocks --outpkg=service --filename=user_service.go --structname=MockUserService --with-expecter=false
type UserService interface {
	Create(ctx context.Context, request dto.CreateUserRequest, userIdentity string, clientId *int) (*dto.UserResponse, error)
	GetUser(id int) (*dto.UserResponse, error)
	Update(request *model.User, userIdentity string) error
	GetUserList(ctx context.Context, clientId *int) ([]*dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, request dto.CreateUserRequest, userIdentity string, clientId *int) (*dto.UserResponse, error) {
	// validate request
	isSuperAdmin, _ := utils.IsSuperAdmin(ctx)
	if !isSuperAdmin {
		if clientId == nil {
			return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("client id is required"))
		}
	}

	// check user if exist
	checkUser, err := s.userRepo.GetByUsername(request.Username)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	newUser := &model.User{
		Username:      request.Username,
		Password:      string(hashedPassword),
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		UserLevel:     request.UserLevel,
		ContactNumber: request.ContactNumber,
		ClientId:      clientId,
		BaseModel: model.BaseModel{
			CreatedBy: userIdentity,
			UpdatedBy: userIdentity,
		},
	}

	// create user
	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toUserResponse(newUser), nil
}

func (s *userService) GetUser(id int) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return s.toUserResponse(user), nil
}

func (s *userService) Update(request *model.User, userIdentity string) error {
	// update user
	request.UpdatedBy = userIdentity
	if err := s.userRepo.Update(request); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *userService) GetUserList(ctx context.Context, clientId *int) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.ListByClientId(ctx, clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, s.toUserResponse(user))
	}

	return responses, nil
}

func (s *userService) toUserResponse(user *model.User) *dto.UserResponse {
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
