package repository

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/gorm"
)

type ScanLogRepository interface {
	Create(ctx context.Context, scanLog *model.ScanLog) error
	Update(ctx context.Context, scanLog *model.ScanLog) error
}

type scanLogRepository struct {
	db *gorm.DB
}

func NewScanLogRepository(db *gorm.DB) ScanLogRepository {
	return &scanLogRepository{db: db}
}

func (r *scanLogRepository) Create(ctx context.Context, scanLog *model.ScanLog) error {
	return r.db.WithContext(ctx).Create(scanLog).Error
}

func (r *scanLogRepository) Update(ctx context.Context, scanLog *model.ScanLog) error {
	return r.db.WithContext(ctx).Save(scanLog).Error
}
