//go:build cgo

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
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
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
		UserLevel:     constants.UserLevelNormal,
		ContactNumber: "1234567890",
	}
	userIdentity := "admin"
	clientId := lo.ToPtr(1)
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
	s.userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*model.User)
		user.Id = expectedUser.Id
		user.CreatedAt = expectedUser.CreatedAt
		user.UpdatedAt = expectedUser.UpdatedAt
	})

	result, err := s.userService.Create(ctx, req, userIdentity, clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Username, result.Username)
	assert.Equal(s.T(), req.FirstName, result.FirstName)
	assert.Equal(s.T(), expectedUser.Id, result.Id)
	assert.Equal(s.T(), clientId, result.ClientId)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestCreate_WithEmail_ChecksUniqueness() {
	ctx := context.Background()
	email := "new@example.com"
	req := dto.CreateUserRequest{
		Username:  "u1",
		Password:  "pw",
		Email:     &email,
		FirstName: "A",
		UserLevel: constants.UserLevelNormal,
	}
	clientId := lo.ToPtr(1)
	s.userRepo.On("GetByUsername", req.Username).Return(nil, nil)
	s.userRepo.On("GetByEmail", email).Return(nil, nil)
	s.userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	_, err := s.userService.Create(ctx, req, "admin", clientId)
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestCreate_EmailAlreadyExists() {
	ctx := context.Background()
	email := "dup@example.com"
	req := dto.CreateUserRequest{
		Username:  "u1",
		Password:  "pw",
		Email:     &email,
		FirstName: "A",
		UserLevel: constants.UserLevelNormal,
	}
	s.userRepo.On("GetByUsername", req.Username).Return(nil, nil)
	s.userRepo.On("GetByEmail", email).Return(&model.User{Id: 9, Email: &email}, nil)

	_, err := s.userService.Create(ctx, req, "admin", lo.ToPtr(1))
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "Email already exists")
}

func (s *UserServiceTestSuite) TestCreate_RejectsSuperAdminLevel() {
	req := dto.CreateUserRequest{
		Username:  "root",
		Password:  "pw",
		FirstName: "R",
		UserLevel: constants.UserLevelSuperAdmin,
	}
	_, err := s.userService.Create(context.Background(), req, "admin", lo.ToPtr(1))
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "super admin")
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
	existingUser := &model.User{Id: 1, Username: req.Username}
	s.userRepo.On("GetByUsername", req.Username).Return(existingUser, nil)

	result, err := s.userService.Create(ctx, req, "admin", lo.ToPtr(1))
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "User already exists")
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

	result, err := s.userService.GetUser(userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.Id, result.Id)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUser_NotFound() {
	userID := 999
	s.userRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	result, err := s.userService.GetUser(userID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *UserServiceTestSuite) TestUpdate_SelfUpdate_StripsPrivilegedFields() {
	userID := 1
	existingUser := &model.User{
		Id:        userID,
		ClientId:  lo.ToPtr(1),
		Username:  "olduser",
		FirstName: "Old",
		LastName:  lo.ToPtr("Name"),
		UserLevel: constants.UserLevelNormal,
	}
	req := dto.UpdateUserRequest{
		FirstName: "Updated",
	}
	s.userRepo.On("GetByID", userID).Return(existingUser, nil)
	s.userRepo.On("Update", mock.Anything, existingUser).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*model.User)
		assert.Equal(s.T(), "Updated", u.FirstName)
		assert.Equal(s.T(), constants.UserLevelNormal, u.UserLevel, "self-update must not change user level")
	})

	err := s.userService.Update(context.Background(), userID, req, "self")
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestAdminUpdate_RejectsPromotionToSuperAdmin() {
	userID := 1
	s.userRepo.On("GetByID", userID).Return(&model.User{Id: userID, UserLevel: constants.UserLevelNormal}, nil)

	level := constants.UserLevelSuperAdmin
	err := s.userService.AdminUpdate(context.Background(), userID, dto.AdminUpdateUserRequest{UserLevel: &level}, "admin")
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "super admin")
}

func (s *UserServiceTestSuite) TestAdminUpdate_RejectsModifyingSuperAdmin() {
	userID := 1
	s.userRepo.On("GetByID", userID).Return(&model.User{Id: userID, UserLevel: constants.UserLevelSuperAdmin}, nil)

	err := s.userService.AdminUpdate(context.Background(), userID, dto.AdminUpdateUserRequest{FirstName: "X"}, "admin")
	assert.Error(s.T(), err)
}

func (s *UserServiceTestSuite) TestAdminUpdate_HappyPath() {
	userID := 1
	existing := &model.User{
		Id:        userID,
		Username:  "u",
		UserLevel: constants.UserLevelNormal,
	}
	s.userRepo.On("GetByID", userID).Return(existing, nil)
	s.userRepo.On("Update", mock.Anything, existing).Return(nil)

	newLevel := constants.UserLevelClientAdmin
	newCid := 7
	err := s.userService.AdminUpdate(context.Background(), userID, dto.AdminUpdateUserRequest{
		FirstName: "New",
		UserLevel: &newLevel,
		ClientId:  &newCid,
	}, "admin")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "New", existing.FirstName)
	assert.Equal(s.T(), constants.UserLevelClientAdmin, existing.UserLevel)
	assert.Equal(s.T(), &newCid, existing.ClientId)
}

func (s *UserServiceTestSuite) TestDelete_RejectsSuperAdmin() {
	userID := 1
	s.userRepo.On("GetByID", userID).Return(&model.User{Id: userID, UserLevel: constants.UserLevelSuperAdmin}, nil)

	err := s.userService.Delete(context.Background(), userID, "admin")
	assert.Error(s.T(), err)
}

func (s *UserServiceTestSuite) TestDelete_HappyPath() {
	userID := 2
	s.userRepo.On("GetByID", userID).Return(&model.User{Id: userID, UserLevel: constants.UserLevelNormal}, nil)
	s.userRepo.On("Delete", mock.Anything, userID).Return(nil)

	err := s.userService.Delete(context.Background(), userID, "admin")
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestGetUserList_Success() {
	ctx := context.Background()
	clientId := 1
	filters := dto.UserListQuery{ClientId: &clientId}
	repoFilters := repository.UserFilters{ClientId: &clientId}
	expectedTime := time.Now()
	expectedUsers := []*model.User{
		{
			Id:       1,
			ClientId: &clientId,
			Username: "user1",
			BaseModel: model.BaseModel{
				CreatedAt: expectedTime,
				UpdatedAt: expectedTime,
				CreatedBy: "admin",
				UpdatedBy: "admin",
			},
		},
		{
			Id:       2,
			ClientId: &clientId,
			Username: "user2",
			BaseModel: model.BaseModel{
				CreatedAt: expectedTime,
				UpdatedAt: expectedTime,
				CreatedBy: "admin",
				UpdatedBy: "admin",
			},
		},
	}
	s.userRepo.On("List", ctx, repoFilters).Return(expectedUsers, nil)

	result, err := s.userService.GetUserList(ctx, filters)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
}

func (s *UserServiceTestSuite) TestGetUserList_Empty() {
	ctx := context.Background()
	filters := dto.UserListQuery{}
	repoFilters := repository.UserFilters{}
	s.userRepo.On("List", ctx, repoFilters).Return([]*model.User{}, nil)

	result, err := s.userService.GetUserList(ctx, filters)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0)
}

func (s *UserServiceTestSuite) TestGetUserList_Error() {
	ctx := context.Background()
	filters := dto.UserListQuery{}
	repoFilters := repository.UserFilters{}
	s.userRepo.On("List", ctx, repoFilters).Return(nil, errors.New("database error"))

	result, err := s.userService.GetUserList(ctx, filters)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}
