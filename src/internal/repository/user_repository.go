package repository

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserRepository --output=./mocks --outpkg=mocks --filename=user_repository.go --structname=MockUserRepository --with-expecter=false
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
	ListByClientId(ctx context.Context, clientId *int) ([]*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
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

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) ListByClientId(ctx context.Context, clientId *int) ([]*model.User, error) {
	var users []*model.User
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if clientId != nil {
		query = query.Where("client_id = ?", *clientId)
	}
	err := query.Find(&users).Error
	return users, err
}
