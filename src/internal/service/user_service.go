package service

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
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
	Update(ctx context.Context, userId int, request dto.UpdateUserRequest, userIdentity string) (*dto.UserResponse, error)
	AdminUpdate(ctx context.Context, userId int, request dto.AdminUpdateUserRequest, userIdentity string) error
	AdminResetPassword(ctx context.Context, userId int, request dto.AdminResetPasswordRequest, userIdentity string) error
	ChangePassword(ctx context.Context, userId int, request dto.ChangePasswordRequest, userIdentity string) error
	Delete(ctx context.Context, userId int, userIdentity string) error
	GetUserList(ctx context.Context, filters dto.UserListQuery) ([]*dto.UserResponse, error)
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

	// Guardrail: super admins cannot be minted via the API.
	if request.UserLevel == constants.UserLevelSuperAdmin {
		return nil, errors.ErrUserCannotAssignSuperAdmin
	}

	// When a super admin creates a user they may specify the target clientId
	// in the request body; otherwise fall back to the caller's clientId.
	targetClientId := clientId
	if isSuperAdmin && request.ClientId != nil {
		targetClientId = request.ClientId
	}

	// check username uniqueness
	checkUser, err := s.userRepo.GetByUsername(request.Username)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if checkUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// check email uniqueness (if provided)
	if request.Email != nil && *request.Email != "" {
		existing, err := s.userRepo.GetByEmail(*request.Email)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		if existing != nil {
			return nil, errors.ErrUserEmailAlreadyExists
		}
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	newUser := &model.User{
		Username:          request.Username,
		Email:             request.Email,
		Password:          string(hashedPassword),
		PasswordUpdatedAt: lo.ToPtr(time.Now()),
		FirstName:         request.FirstName,
		LastName:          request.LastName,
		UserLevel:         request.UserLevel,
		ContactNumber:     request.ContactNumber,
		ClientId:          targetClientId,
	}

	// create user (CreatedBy/UpdatedBy set via BaseModel hook from ctx)
	err = s.userRepo.Create(ctx, newUser)
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

// Update is the self-update path. It intentionally does NOT allow changing
// UserLevel or ClientId; those are privileged fields handled by AdminUpdate.
// Returns the updated user so callers can refresh their local snapshot
// without an extra round-trip.
func (s *userService) Update(ctx context.Context, userId int, request dto.UpdateUserRequest, userIdentity string) (*dto.UserResponse, error) {
	existingUser, err := s.userRepo.GetByID(userId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if existingUser == nil {
		return nil, errors.ErrUserNotFound
	}

	if request.Username != "" && request.Username != existingUser.Username {
		clash, err := s.userRepo.GetByUsername(request.Username)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		if clash != nil && clash.Id != existingUser.Id {
			return nil, errors.ErrUserAlreadyExists
		}
		existingUser.Username = request.Username
	}
	if request.Email != nil && (existingUser.Email == nil || *request.Email != *existingUser.Email) {
		if *request.Email != "" {
			clash, err := s.userRepo.GetByEmail(*request.Email)
			if err != nil {
				return nil, errors.ErrGeneric.Wrap(err)
			}
			if clash != nil && clash.Id != existingUser.Id {
				return nil, errors.ErrUserEmailAlreadyExists
			}
		}
		existingUser.Email = request.Email
	}
	if request.FirstName != "" {
		existingUser.FirstName = request.FirstName
	}
	if request.LastName != nil {
		existingUser.LastName = request.LastName
	}
	if request.ContactNumber != "" {
		existingUser.ContactNumber = request.ContactNumber
	}

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	return s.toUserResponse(existingUser), nil
}

// assertCanManageUser gates the admin-facing user mutations. The handler only
// knows the caller's level; which client the *target* belongs to is only
// knowable here, so this is where cross-client access is stopped.
//
//	· super-admin records are never reachable through these endpoints
//	· a client admin may only act on users inside its own client
func (s *userService) assertCanManageUser(ctx context.Context, target *model.User) error {
	if target.UserLevel == constants.UserLevelSuperAdmin {
		return errors.ErrUserCannotModifySuperAdmin
	}
	if target.ClientId == nil {
		// Nothing to scope against, so only a super admin may act on it.
		if isSuper, _ := utils.IsSuperAdmin(ctx); !isSuper {
			return errors.ErrAuthPermissionDenied
		}
		return nil
	}
	ok, err := utils.CanAccessClient(ctx, *target.ClientId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return errors.ErrAuthPermissionDenied
	}
	return nil
}

// AdminUpdate is the admin path (client admin and above). It can change every
// field on a user except that it refuses to promote anyone to SuperAdmin or
// demote an existing SuperAdmin (those are out-of-band operations), and it
// keeps a client admin inside its own client.
func (s *userService) AdminUpdate(ctx context.Context, userId int, request dto.AdminUpdateUserRequest, userIdentity string) error {
	existingUser, err := s.userRepo.GetByID(userId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	// Never touch super-admin records, and never reach outside the caller's
	// client.
	if err := s.assertCanManageUser(ctx, existingUser); err != nil {
		return err
	}
	// Never promote to super-admin via this endpoint.
	if request.UserLevel != nil && *request.UserLevel == constants.UserLevelSuperAdmin {
		return errors.ErrUserCannotAssignSuperAdmin
	}
	// Moving a user between clients is a super-admin operation. Refusing beats
	// ignoring the field: a client admin sending it is trying to leave its own
	// scope, and a silent no-op would look like it worked.
	if request.ClientId != nil {
		isSuper, _ := utils.IsSuperAdmin(ctx)
		if !isSuper {
			return errors.ErrAuthPermissionDenied
		}
	}

	if request.Username != "" && request.Username != existingUser.Username {
		clash, err := s.userRepo.GetByUsername(request.Username)
		if err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		if clash != nil && clash.Id != existingUser.Id {
			return errors.ErrUserAlreadyExists
		}
		existingUser.Username = request.Username
	}
	if request.Email != nil && (existingUser.Email == nil || *request.Email != *existingUser.Email) {
		if *request.Email != "" {
			clash, err := s.userRepo.GetByEmail(*request.Email)
			if err != nil {
				return errors.ErrGeneric.Wrap(err)
			}
			if clash != nil && clash.Id != existingUser.Id {
				return errors.ErrUserEmailAlreadyExists
			}
		}
		existingUser.Email = request.Email
	}
	if request.FirstName != "" {
		existingUser.FirstName = request.FirstName
	}
	if request.LastName != nil {
		existingUser.LastName = request.LastName
	}
	if request.UserLevel != nil {
		existingUser.UserLevel = *request.UserLevel
	}
	if request.ContactNumber != "" {
		existingUser.ContactNumber = request.ContactNumber
	}
	if request.ClientId != nil {
		existingUser.ClientId = request.ClientId
	}

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

// AdminResetPassword overwrites a user's password. Callers must already be
// verified client-admin or above; this refuses a super admin's password and any
// target outside the caller's client.
func (s *userService) AdminResetPassword(ctx context.Context, userId int, request dto.AdminResetPasswordRequest, userIdentity string) error {
	existingUser, err := s.userRepo.GetByID(userId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	if err := s.assertCanManageUser(ctx, existingUser); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	existingUser.Password = string(hashedPassword)
	existingUser.PasswordUpdatedAt = lo.ToPtr(time.Now())

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

// ChangePassword lets an authenticated user change their own password by
// verifying the current password before hashing and storing the new one.
func (s *userService) ChangePassword(ctx context.Context, userId int, request dto.ChangePasswordRequest, userIdentity string) error {
	existingUser, err := s.userRepo.GetByID(userId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(request.CurrentPassword)); err != nil {
		return errors.ErrAuthInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	existingUser.Password = string(hashedPassword)
	existingUser.PasswordUpdatedAt = lo.ToPtr(time.Now())

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

// Delete soft-deletes a user. Callers must already be verified client-admin or
// above. Refuses the acting user, a super admin, and anyone outside the
// caller's client.
func (s *userService) Delete(ctx context.Context, userId int, userIdentity string) error {
	existingUser, err := s.userRepo.GetByID(userId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	actingUserId, err := utils.GetUserId(ctx)
	if err == nil && actingUserId == existingUser.Id {
		return errors.ErrUserCannotDeleteSelf
	}
	if err := s.assertCanManageUser(ctx, existingUser); err != nil {
		return err
	}

	if err := s.userRepo.Delete(ctx, existingUser.Id); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *userService) GetUserList(ctx context.Context, filters dto.UserListQuery) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.List(ctx, repository.UserFilters{
		Search:    filters.Search,
		UserLevel: filters.UserLevel,
		ClientId:  filters.ClientId,
	})
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
		Id:                user.Id,
		ClientId:          user.ClientId,
		Username:          user.Username,
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		UserLevel:         user.UserLevel,
		ContactNumber:     user.ContactNumber,
		PasswordUpdatedAt: user.PasswordUpdatedAt,
		CreatedAt:         user.CreatedAt,
		CreatedBy:         user.CreatedBy,
		UpdatedAt:         user.UpdatedAt,
		UpdatedBy:         user.UpdatedBy,
	}
}
