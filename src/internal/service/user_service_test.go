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
	apperrors "github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
	"golang.org/x/crypto/bcrypt"
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

	updated, err := s.userService.Update(context.Background(), userID, req, "self")
	assert.NoError(s.T(), err)
	if assert.NotNil(s.T(), updated) {
		assert.Equal(s.T(), "Updated", updated.FirstName)
	}
}

func (s *UserServiceTestSuite) TestAdminUpdate_RejectsPromotionToSuperAdmin() {
	userID := 1
	s.userRepo.On("GetByID", userID).Return(&model.User{Id: userID, UserLevel: constants.UserLevelNormal}, nil)

	level := constants.UserLevelSuperAdmin
	err := s.userService.AdminUpdate(superAdminCtx(), userID, dto.AdminUpdateUserRequest{UserLevel: &level}, "admin")
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "super admin")
}

func (s *UserServiceTestSuite) TestAdminUpdate_RejectsModifyingSuperAdmin() {
	userID := 1
	s.userRepo.On("GetByID", userID).Return(&model.User{Id: userID, UserLevel: constants.UserLevelSuperAdmin}, nil)

	err := s.userService.AdminUpdate(superAdminCtx(), userID, dto.AdminUpdateUserRequest{FirstName: "X"}, "admin")
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
	// Super-admin context: this reassigns clientId, which a client admin may
	// not do (see TestAdminUpdate_ClientAdminCannotReassignClient).
	err := s.userService.AdminUpdate(superAdminCtx(), userID, dto.AdminUpdateUserRequest{
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

	err := s.userService.Delete(superAdminCtx(), userID, "admin")
	assert.Error(s.T(), err)
}

func (s *UserServiceTestSuite) TestDelete_HappyPath() {
	userID := 2
	s.userRepo.On("GetByID", userID).Return(&model.User{Id: userID, UserLevel: constants.UserLevelNormal}, nil)
	s.userRepo.On("Delete", mock.Anything, userID).Return(nil)

	err := s.userService.Delete(superAdminCtx(), userID, "admin")
	assert.NoError(s.T(), err)
}

// --- password_updated_at -----------------------------------------------------
//
// The column exists so clients can say "you last changed your password on X".
// It must move on exactly the paths that write users.password and on no others,
// otherwise the date silently becomes a lie.

// userWithPassword builds a user whose stored hash matches `plain`.
func userWithPassword(id int, plain string) *model.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	return &model.User{
		Id:        id,
		Username:  "u",
		UserLevel: constants.UserLevelNormal,
		Password:  string(hash),
	}
}

func (s *UserServiceTestSuite) TestCreate_StampsPasswordUpdatedAt() {
	req := dto.CreateUserRequest{
		Username:  "fresh",
		Password:  "password123",
		FirstName: "Fresh",
		UserLevel: constants.UserLevelNormal,
	}
	s.userRepo.On("GetByUsername", req.Username).Return(nil, nil)
	s.userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	before := time.Now()
	result, err := s.userService.Create(context.Background(), req, "admin", lo.ToPtr(1))

	assert.NoError(s.T(), err)
	if assert.NotNil(s.T(), result) && assert.NotNil(s.T(), result.PasswordUpdatedAt) {
		assert.False(s.T(), result.PasswordUpdatedAt.Before(before),
			"a new account's password was set now, so the stamp must be now")
	}
}

func (s *UserServiceTestSuite) TestChangePassword_StampsPasswordUpdatedAt() {
	userID := 1
	existing := userWithPassword(userID, "oldpassword")
	oldHash := existing.Password
	s.userRepo.On("GetByID", userID).Return(existing, nil)
	s.userRepo.On("Update", mock.Anything, existing).Return(nil)

	before := time.Now()
	err := s.userService.ChangePassword(context.Background(), userID, dto.ChangePasswordRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword123",
	}, "self")

	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), oldHash, existing.Password, "the hash should have been replaced")
	if assert.NotNil(s.T(), existing.PasswordUpdatedAt) {
		assert.False(s.T(), existing.PasswordUpdatedAt.Before(before))
	}
}

func (s *UserServiceTestSuite) TestChangePassword_WrongCurrent_LeavesStampAlone() {
	userID := 1
	existing := userWithPassword(userID, "oldpassword")
	s.userRepo.On("GetByID", userID).Return(existing, nil)
	// No Update expectation: a rejected change must not write anything.

	err := s.userService.ChangePassword(context.Background(), userID, dto.ChangePasswordRequest{
		CurrentPassword: "wrong",
		NewPassword:     "newpassword123",
	}, "self")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), existing.PasswordUpdatedAt, "a failed attempt is not a password change")
	s.userRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *UserServiceTestSuite) TestAdminResetPassword_StampsPasswordUpdatedAt() {
	userID := 2
	existing := userWithPassword(userID, "whatever")
	s.userRepo.On("GetByID", userID).Return(existing, nil)
	s.userRepo.On("Update", mock.Anything, existing).Return(nil)

	before := time.Now()
	err := s.userService.AdminResetPassword(context.Background(), userID, dto.AdminResetPasswordRequest{
		Password: "resetpassword123",
	}, "admin")

	assert.NoError(s.T(), err)
	if assert.NotNil(s.T(), existing.PasswordUpdatedAt) {
		assert.False(s.T(), existing.PasswordUpdatedAt.Before(before))
	}
}

func (s *UserServiceTestSuite) TestUpdate_ProfileEdit_DoesNotStampPasswordUpdatedAt() {
	userID := 1
	stamped := time.Now().Add(-72 * time.Hour)
	existing := &model.User{
		Id:                userID,
		Username:          "u",
		UserLevel:         constants.UserLevelNormal,
		PasswordUpdatedAt: lo.ToPtr(stamped),
	}
	s.userRepo.On("GetByID", userID).Return(existing, nil)
	s.userRepo.On("Update", mock.Anything, existing).Return(nil)

	updated, err := s.userService.Update(context.Background(), userID, dto.UpdateUserRequest{
		FirstName:     "Renamed",
		ContactNumber: "0899999999",
	}, "self")

	assert.NoError(s.T(), err)
	// This is the whole reason the column exists: updated_at would have moved
	// here, and reusing it would have claimed a password change that never
	// happened.
	assert.Equal(s.T(), stamped, *existing.PasswordUpdatedAt)
	if assert.NotNil(s.T(), updated) && assert.NotNil(s.T(), updated.PasswordUpdatedAt) {
		assert.Equal(s.T(), stamped, *updated.PasswordUpdatedAt)
	}
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

// --- client-admin scoping ----------------------------------------------------
//
// These mutations used to be super-admin only, so the service never had to ask
// which client the *target* belonged to. Now that a farm owner can manage its
// own staff, that question is the whole security boundary.

// clientAdminCtx is a client admin (level 2) of client 1.
func clientAdminCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "owner")
	ctx = context.WithValue(ctx, constants.ClientIDKey, 1)
	ctx = context.WithValue(ctx, constants.UserLevelKey, constants.UserLevelClientAdmin)
	return ctx
}

func (s *UserServiceTestSuite) TestAdminUpdate_ClientAdminOwnClientAllowed() {
	existing := &model.User{Id: 5, ClientId: lo.ToPtr(1), Username: "worker", UserLevel: constants.UserLevelNormal}
	s.userRepo.On("GetByID", 5).Return(existing, nil)
	s.userRepo.On("Update", mock.Anything, existing).Return(nil)

	err := s.userService.AdminUpdate(clientAdminCtx(), 5, dto.AdminUpdateUserRequest{FirstName: "Renamed"}, "owner")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Renamed", existing.FirstName)
}

func (s *UserServiceTestSuite) TestAdminUpdate_ClientAdminForeignClientDenied() {
	s.userRepo.On("GetByID", 6).Return(
		&model.User{Id: 6, ClientId: lo.ToPtr(2), Username: "theirs", UserLevel: constants.UserLevelNormal}, nil)

	err := s.userService.AdminUpdate(clientAdminCtx(), 6, dto.AdminUpdateUserRequest{FirstName: "Hacked"}, "owner")
	assert.ErrorIs(s.T(), err, apperrors.ErrAuthPermissionDenied)
	s.userRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *UserServiceTestSuite) TestAdminUpdate_ClientAdminCannotReassignClient() {
	s.userRepo.On("GetByID", 5).Return(
		&model.User{Id: 5, ClientId: lo.ToPtr(1), Username: "worker", UserLevel: constants.UserLevelNormal}, nil)

	// Moving a user to another client would be a way out of the caller's scope.
	err := s.userService.AdminUpdate(clientAdminCtx(), 5, dto.AdminUpdateUserRequest{ClientId: lo.ToPtr(2)}, "owner")
	assert.ErrorIs(s.T(), err, apperrors.ErrAuthPermissionDenied)
	s.userRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *UserServiceTestSuite) TestAdminUpdate_ClientAdminCannotTouchSuperAdmin() {
	s.userRepo.On("GetByID", 1).Return(
		&model.User{Id: 1, Username: "sys_admin", UserLevel: constants.UserLevelSuperAdmin}, nil)

	err := s.userService.AdminUpdate(clientAdminCtx(), 1, dto.AdminUpdateUserRequest{FirstName: "X"}, "owner")
	assert.ErrorIs(s.T(), err, apperrors.ErrUserCannotModifySuperAdmin)
}

func (s *UserServiceTestSuite) TestAdminResetPassword_ClientAdminForeignClientDenied() {
	s.userRepo.On("GetByID", 6).Return(
		&model.User{Id: 6, ClientId: lo.ToPtr(2), Username: "theirs", UserLevel: constants.UserLevelNormal}, nil)

	err := s.userService.AdminResetPassword(clientAdminCtx(), 6,
		dto.AdminResetPasswordRequest{Password: "Passw0rd123"}, "owner")
	assert.ErrorIs(s.T(), err, apperrors.ErrAuthPermissionDenied)
	s.userRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *UserServiceTestSuite) TestAdminResetPassword_ClientAdminOwnClientAllowed() {
	existing := &model.User{Id: 5, ClientId: lo.ToPtr(1), Username: "worker", UserLevel: constants.UserLevelNormal}
	s.userRepo.On("GetByID", 5).Return(existing, nil)
	s.userRepo.On("Update", mock.Anything, existing).Return(nil)

	before := existing.Password
	err := s.userService.AdminResetPassword(clientAdminCtx(), 5,
		dto.AdminResetPasswordRequest{Password: "Passw0rd123"}, "owner")
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), before, existing.Password)
}

func (s *UserServiceTestSuite) TestDelete_ClientAdminForeignClientDenied() {
	s.userRepo.On("GetByID", 6).Return(
		&model.User{Id: 6, ClientId: lo.ToPtr(2), Username: "theirs", UserLevel: constants.UserLevelNormal}, nil)

	err := s.userService.Delete(clientAdminCtx(), 6, "owner")
	assert.ErrorIs(s.T(), err, apperrors.ErrAuthPermissionDenied)
	s.userRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *UserServiceTestSuite) TestDelete_ClientAdminOwnClientAllowed() {
	s.userRepo.On("GetByID", 5).Return(
		&model.User{Id: 5, ClientId: lo.ToPtr(1), Username: "worker", UserLevel: constants.UserLevelNormal}, nil)
	s.userRepo.On("Delete", mock.Anything, 5).Return(nil)

	err := s.userService.Delete(clientAdminCtx(), 5, "owner")
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestCreate_ClientAdminUsesOwnClientAndIgnoresRequestedOne() {
	req := dto.CreateUserRequest{
		Username:  "newworker",
		Password:  "Passw0rd123",
		FirstName: "New",
		UserLevel: constants.UserLevelNormal,
		// A client admin naming another client must not be honoured.
		ClientId: lo.ToPtr(2),
	}
	s.userRepo.On("GetByUsername", req.Username).Return(nil, nil)
	s.userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).
		Run(func(args mock.Arguments) {
			u := args.Get(1).(*model.User)
			assert.Equal(s.T(), lo.ToPtr(1), u.ClientId, "must land in the caller's own client")
		})

	_, err := s.userService.Create(clientAdminCtx(), req, "owner", lo.ToPtr(1))
	assert.NoError(s.T(), err)
}
