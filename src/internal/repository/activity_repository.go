package repository

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=ActivityRepository --output=./mocks --outpkg=mocks --filename=activity_repository.go --structname=MockActivityRepository --with-expecter=false
type ActivityRepository interface {
	Create(ctx context.Context, activity *model.Activity) error
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(ctx context.Context, activity *model.Activity) error {
	return r.db.WithContext(ctx).Create(activity).Error
}
