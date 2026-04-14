package repository

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=ActivePondRepository --output=./mocks --outpkg=mocks --filename=active_pond_repository.go --structname=MockActivePondRepository --with-expecter=false
type ActivePondRepository interface {
	WithTx(tx *gorm.DB) ActivePondRepository
	GetActiveByPondID(ctx context.Context, pondId int) (*model.ActivePond, error)
	GetLatestByPondID(ctx context.Context, pondId int) (*model.ActivePond, error)
	Create(ctx context.Context, activePond *model.ActivePond) error
	Update(ctx context.Context, activePond *model.ActivePond) error
}

type activePondRepository struct {
	db *gorm.DB
}

func NewActivePondRepository(db *gorm.DB) ActivePondRepository {
	return &activePondRepository{db: db}
}

func (r *activePondRepository) WithTx(tx *gorm.DB) ActivePondRepository {
	return &activePondRepository{db: tx}
}

func (r *activePondRepository) GetActiveByPondID(ctx context.Context, pondId int) (*model.ActivePond, error) {
	var ap model.ActivePond
	err := r.db.WithContext(ctx).Where("pond_id = ? AND is_active = ? AND deleted_at IS NULL", pondId, true).First(&ap).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ap, nil
}

// GetLatestByPondID returns the most recent active_pond row for a pond, regardless of is_active status.
// Used for historical data entry on closed ponds.
func (r *activePondRepository) GetLatestByPondID(ctx context.Context, pondId int) (*model.ActivePond, error) {
	var ap model.ActivePond
	err := r.db.WithContext(ctx).Where("pond_id = ? AND deleted_at IS NULL", pondId).Order("id DESC").First(&ap).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ap, nil
}

func (r *activePondRepository) Create(ctx context.Context, activePond *model.ActivePond) error {
	return r.db.WithContext(ctx).Create(activePond).Error
}

func (r *activePondRepository) Update(ctx context.Context, activePond *model.ActivePond) error {
	return r.db.WithContext(ctx).Save(activePond).Error
}
