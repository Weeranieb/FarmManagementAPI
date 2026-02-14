package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
)

type UserServiceTestSuite struct {
	suite.Suite
	userRepo    *mocks.MockUserRepository
	userService UserService
}

func (s *UserServiceTestSuite) SetupTest() {
	s.userRepo = mocks.NewMockUserRepository(s.T())
	s.userService = NewUserService(s.userRepo)
}

func (s *UserServiceTestSuite) TearDownTest() {
	s.userRepo.ExpectedCalls = nil
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) TestCreate_Success() {
	ctx := context.Background()
	req := dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		LastName:      lo.ToPtr("User"),
		UserLevel:     1,
		ContactNumber: "1234567890",
	}
	userIdentity := "admin"
	clientId := lo.ToPtr(1)

	// Mock repository calls - GetByUsername returns (nil, nil) when not found
	s.userRepo.On("GetByUsername", req.Username).Return(nil, nil)

	expectedTime := time.Now()
	expectedUser := &model.User{
		Id:            1,
		Username:      req.Username,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		UserLevel:     req.UserLevel,
		ContactNumber: req.ContactNumber,
		ClientId:      clientId,
		Password:      "hashed_password",
		BaseModel: model.BaseModel{
			CreatedAt: expectedTime,
			UpdatedAt: expectedTime,
			CreatedBy: userIdentity,
			UpdatedBy: userIdentity,
		},
	}

	s.userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*model.User)
		user.Id = expectedUser.Id
		user.CreatedAt = expectedUser.CreatedAt
		user.UpdatedAt = expectedUser.UpdatedAt
	})

	// Execute
	result, err := s.userService.Create(ctx, req, userIdentity, clientId)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Username, result.Username)
	assert.Equal(s.T(), req.FirstName, result.FirstName)
	assert.Equal(s.T(), expectedUser.Id, result.Id)
	assert.Equal(s.T(), clientId, result.ClientId)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestCreate_UsernameExists() {
	ctx := context.Background()
	req := dto.CreateUserRequest{
		Username:      "existinguser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}
	userIdentity := "admin"
	clientId := lo.ToPtr(1)

	existingUser := &model.User{
		Id:        1,
		Username:  req.Username,
		FirstName: "Existing",
	}

	// GetByUsername returns the existing user when found
	s.userRepo.On("GetByUsername", req.Username).Return(existingUser, nil)

	// Execute
	result, err := s.userService.Create(ctx, req, userIdentity, clientId)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	// Error message format changed to include error code
	assert.Contains(s.T(), err.Error(), "User already exists")
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestCreate_GetByUsernameError() {
	ctx := context.Background()
	req := dto.CreateUserRequest{
		Username:      "testuser",
		Password:      "password123",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
	}
	userIdentity := "admin"
	clientId := lo.ToPtr(1)

	// GetByUsername returns an error (not a not-found case)
	s.userRepo.On("GetByUsername", req.Username).Return(nil, errors.New("database error"))

	// Execute
	result, err := s.userService.Create(ctx, req, userIdentity, clientId)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUser_Success() {
	userID := 1
	expectedTime := time.Now()
	expectedUser := &model.User{
		Id:            userID,
		Username:      "testuser",
		FirstName:     "Test",
		LastName:      lo.ToPtr("User"),
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
		BaseModel: model.BaseModel{
			CreatedAt: expectedTime,
			UpdatedAt: expectedTime,
			CreatedBy: "admin",
			UpdatedBy: "admin",
		},
	}

	s.userRepo.On("GetByID", userID).Return(expectedUser, nil)

	// Execute
	result, err := s.userService.GetUser(userID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.Id, result.Id)
	assert.Equal(s.T(), expectedUser.Username, result.Username)
	assert.Equal(s.T(), expectedUser.FirstName, result.FirstName)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUser_NotFound() {
	userID := 999

	s.userRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	// Execute
	result, err := s.userService.GetUser(userID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestUpdate_Success() {
	userIdentity := "admin"
	updateUser := &model.User{
		Id:            1,
		ClientId:      lo.ToPtr(1),
		Username:      "updateduser",
		FirstName:     "Updated",
		LastName:      lo.ToPtr("User"),
		UserLevel:     1,
		ContactNumber: "0987654321",
		BaseModel: model.BaseModel{
			UpdatedBy: userIdentity,
		},
	}

	s.userRepo.On("Update", updateUser).Return(nil)

	// Execute
	err := s.userService.Update(updateUser, userIdentity)

	// Assert
	assert.NoError(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestUpdate_Error() {
	userIdentity := "admin"
	updateUser := &model.User{
		Id:        1,
		Username:  "updateduser",
		FirstName: "Updated",
	}

	s.userRepo.On("Update", updateUser).Return(errors.New("update failed"))

	// Execute
	err := s.userService.Update(updateUser, userIdentity)

	// Assert
	assert.Error(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUserList_Success() {
	ctx := context.Background()
	clientId := 1
	clientIdPtr := &clientId

	expectedTime := time.Now()
	expectedUsers := []*model.User{
		{
			Id:            1,
			ClientId:      &clientId,
			Username:      "user1",
			FirstName:     "User",
			LastName:      lo.ToPtr("One"),
			UserLevel:     1,
			ContactNumber: "1111111111",
			BaseModel: model.BaseModel{
				CreatedAt: expectedTime,
				UpdatedAt: expectedTime,
				CreatedBy: "admin",
				UpdatedBy: "admin",
			},
		},
		{
			Id:            2,
			ClientId:      &clientId,
			Username:      "user2",
			FirstName:     "User",
			LastName:      lo.ToPtr("Two"),
			UserLevel:     1,
			ContactNumber: "2222222222",
			BaseModel: model.BaseModel{
				CreatedAt: expectedTime,
				UpdatedAt: expectedTime,
				CreatedBy: "admin",
				UpdatedBy: "admin",
			},
		},
	}

	s.userRepo.On("ListByClientId", ctx, clientIdPtr).Return(expectedUsers, nil)

	// Execute
	result, err := s.userService.GetUserList(ctx, clientIdPtr)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedUsers[0].Id, result[0].Id)
	assert.Equal(s.T(), expectedUsers[1].Id, result[1].Id)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUserList_Empty() {
	ctx := context.Background()
	clientId := 999
	clientIdPtr := &clientId

	s.userRepo.On("ListByClientId", ctx, clientIdPtr).Return([]*model.User{}, nil)

	// Execute
	result, err := s.userService.GetUserList(ctx, clientIdPtr)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 0)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUserList_Error() {
	ctx := context.Background()
	clientId := 1
	clientIdPtr := &clientId

	s.userRepo.On("ListByClientId", ctx, clientIdPtr).Return(nil, errors.New("database error"))

	// Execute
	result, err := s.userService.GetUserList(ctx, clientIdPtr)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.userRepo.AssertExpectations(s.T())
}
