package repository

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=SellDetailRepository --output=./mocks --outpkg=mocks --filename=sell_detail_repository.go --structname=MockSellDetailRepository --with-expecter=false
type SellDetailRepository interface {
	WithTx(tx *gorm.DB) SellDetailRepository
	CreateBatch(ctx context.Context, details []*model.SellDetail) error
}

type sellDetailRepository struct {
	db *gorm.DB
}

func NewSellDetailRepository(db *gorm.DB) SellDetailRepository {
	return &sellDetailRepository{db: db}
}

func (r *sellDetailRepository) WithTx(tx *gorm.DB) SellDetailRepository {
	return &sellDetailRepository{db: tx}
}

func (r *sellDetailRepository) CreateBatch(ctx context.Context, details []*model.SellDetail) error {
	return r.db.WithContext(ctx).Create(details).Error
}
