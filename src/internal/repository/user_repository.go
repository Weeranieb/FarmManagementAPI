package repository

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

// UserFilters are optional filters applied by List.
type UserFilters struct {
	Search    *string
	UserLevel *int
	ClientId  *int
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserRepository --output=./mocks --outpkg=mocks --filename=user_repository.go --structname=MockUserRepository --with-expecter=false
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(id int) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters UserFilters) ([]*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, filters UserFilters) ([]*model.User, error) {
	var users []*model.User
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if filters.ClientId != nil {
		query = query.Where("client_id = ?", *filters.ClientId)
	}
	if filters.UserLevel != nil {
		query = query.Where("user_level = ?", *filters.UserLevel)
	}
	if filters.Search != nil && *filters.Search != "" {
		like := "%" + *filters.Search + "%"
		query = query.Where(
			"username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			like, like, like, like,
		)
	}
	err := query.Order("id ASC").Find(&users).Error
	return users, err
}
