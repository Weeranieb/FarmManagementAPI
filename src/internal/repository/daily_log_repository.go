package repository

import (
	"context"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyLogRepository --output=./mocks --outpkg=mocks --filename=daily_log_repository.go --structname=MockDailyLogRepository --with-expecter=false
type DailyLogRepository interface {
	Upsert(ctx context.Context, logs []*model.DailyLog) error
	ListByActivePondAndMonth(activePondId int, start, end time.Time) ([]*model.DailyLog, error)
}

type dailyLogRepository struct {
	db *gorm.DB
}

func NewDailyLogRepository(db *gorm.DB) DailyLogRepository {
	return &dailyLogRepository{db: db}
}

func (r *dailyLogRepository) Upsert(ctx context.Context, logs []*model.DailyLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "active_pond_id"}, {Name: "feed_date"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"fresh_feed_collection_id", "pellet_feed_collection_id",
				"fresh_morning", "fresh_evening", "pellet_morning", "pellet_evening",
				"death_fish_count", "tourist_catch_count",
				"updated_by", "updated_at",
			}),
			Where: clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: "daily_logs.deleted_at IS NULL"}}},
		}).
		Create(logs).Error
}

func (r *dailyLogRepository) ListByActivePondAndMonth(activePondId int, start, end time.Time) ([]*model.DailyLog, error) {
	var logs []*model.DailyLog
	err := r.db.
		Where("active_pond_id = ? AND feed_date >= ? AND feed_date <= ? AND deleted_at IS NULL", activePondId, start, end).
		Order("feed_date").
		Find(&logs).Error
	return logs, err
}
