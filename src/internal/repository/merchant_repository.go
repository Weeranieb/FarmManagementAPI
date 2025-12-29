package repository

import (
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=MerchantRepository --output=./mocks --outpkg=mocks --filename=merchant_repository.go --structname=MockMerchantRepository --with-expecter=false
type MerchantRepository interface {
	Create(merchant *model.Merchant) error
	GetByID(id int) (*model.Merchant, error)
	GetByContactNumberAndName(contactNumber, name string) (*model.Merchant, error)
	Update(merchant *model.Merchant) error
	List() ([]*model.Merchant, error)
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (r *merchantRepository) Create(merchant *model.Merchant) error {
	return r.db.Create(merchant).Error
}

func (r *merchantRepository) GetByID(id int) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &merchant, nil
}

func (r *merchantRepository) GetByContactNumberAndName(contactNumber, name string) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.Where("contact_number = ? AND name = ? AND deleted_at IS NULL", contactNumber, name).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &merchant, nil
}

func (r *merchantRepository) Update(merchant *model.Merchant) error {
	return r.db.Save(merchant).Error
}

func (r *merchantRepository) List() ([]*model.Merchant, error) {
	var merchants []*model.Merchant
	err := r.db.Where("deleted_at IS NULL").Find(&merchants).Error
	return merchants, err
}
