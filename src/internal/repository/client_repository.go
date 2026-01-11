package repository

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=ClientRepository --output=./mocks --outpkg=mocks --filename=client_repository.go --structname=MockClientRepository --with-expecter=false
type ClientRepository interface {
	Create(client *model.Client) error
	GetByID(id int) (*model.Client, error)
	GetByName(ctx context.Context, name string) (*model.Client, error)
	Update(client *model.Client) error
	List() ([]*model.Client, error)
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(client *model.Client) error {
	return r.db.Create(client).Error
}

func (r *clientRepository) GetByID(id int) (*model.Client, error) {
	var client model.Client
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) GetByName(ctx context.Context, name string) (*model.Client, error) {
	var client model.Client
	err := r.db.WithContext(ctx).Where("name = ? AND deleted_at IS NULL", name).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) Update(client *model.Client) error {
	return r.db.Save(client).Error
}

func (r *clientRepository) List() ([]*model.Client, error) {
	var clients []*model.Client
	err := r.db.Where("deleted_at IS NULL").Find(&clients).Error
	return clients, err
}
