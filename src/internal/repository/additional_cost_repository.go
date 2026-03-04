package repository

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=AdditionalCostRepository --output=./mocks --outpkg=mocks --filename=additional_cost_repository.go --structname=MockAdditionalCostRepository --with-expecter=false
type AdditionalCostRepository interface {
	WithTx(tx *gorm.DB) AdditionalCostRepository
	Create(ctx context.Context, ac *model.AdditionalCost) error
	CreateBatch(ctx context.Context, items []*model.AdditionalCost) error
}

type additionalCostRepository struct {
	db *gorm.DB
}

func NewAdditionalCostRepository(db *gorm.DB) AdditionalCostRepository {
	return &additionalCostRepository{db: db}
}

func (r *additionalCostRepository) WithTx(tx *gorm.DB) AdditionalCostRepository {
	return &additionalCostRepository{db: tx}
}

func (r *additionalCostRepository) Create(ctx context.Context, ac *model.AdditionalCost) error {
	return r.db.WithContext(ctx).Create(ac).Error
}

func (r *additionalCostRepository) CreateBatch(ctx context.Context, items []*model.AdditionalCost) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(items).Error
}
