package repository

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=ActivePondRepository --output=./mocks --outpkg=mocks --filename=active_pond_repository.go --structname=MockActivePondRepository --with-expecter=false
type ActivePondRepository interface {
	WithTx(tx *gorm.DB) ActivePondRepository
	GetActiveByPondID(ctx context.Context, pondId int) (*model.ActivePond, error)
	GetByIDForUpdate(ctx context.Context, id int) (*model.ActivePond, error)
	ListByPondID(ctx context.Context, pondId int) ([]*model.ActivePond, error)
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

// GetByIDForUpdate loads one cycle by its active_pond id and locks the row FOR
// UPDATE. Must be called inside a transaction (via WithTx): a concurrent
// sell/move on the same cycle then serializes behind this lock instead of
// racing on a stale head count. Returns (nil, nil) when the row is absent.
func (r *activePondRepository) GetByIDForUpdate(ctx context.Context, id int) (*model.ActivePond, error) {
	var ap model.ActivePond
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&ap).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ap, nil
}

// ListByPondID returns every cycle (active and closed) for a pond, newest first,
// so callers can present the pond's cycle history with each cycle's P&L.
func (r *activePondRepository) ListByPondID(ctx context.Context, pondId int) ([]*model.ActivePond, error) {
	var cycles []*model.ActivePond
	err := r.db.WithContext(ctx).
		Where("pond_id = ? AND deleted_at IS NULL", pondId).
		Order("start_date DESC, id DESC").
		Find(&cycles).Error
	return cycles, err
}

func (r *activePondRepository) Create(ctx context.Context, activePond *model.ActivePond) error {
	return r.db.WithContext(ctx).Create(activePond).Error
}

func (r *activePondRepository) Update(ctx context.Context, activePond *model.ActivePond) error {
	return r.db.WithContext(ctx).Save(activePond).Error
}
