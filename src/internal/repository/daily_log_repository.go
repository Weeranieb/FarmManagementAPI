package repository

import (
	"context"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DailyLogIDFeedDate is a minimal row projection for template-import reconcile.
type DailyLogIDFeedDate struct {
	Id       int       `gorm:"column:id"`
	FeedDate time.Time `gorm:"column:feed_date"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyLogRepository --output=./mocks --outpkg=mocks --filename=daily_log_repository.go --structname=MockDailyLogRepository --with-expecter=false
type DailyLogRepository interface {
	WithTx(tx *gorm.DB) DailyLogRepository
	Upsert(ctx context.Context, logs []*model.DailyLog) error
	ListIDAndFeedDateByActivePondRange(ctx context.Context, activePondId int, min, max time.Time) ([]DailyLogIDFeedDate, error)
	HardDeleteByIDs(ctx context.Context, ids []int) error
	ListByActivePondAndMonth(ctx context.Context, activePondId int, start, end time.Time) ([]*model.DailyLog, error)
	HardDeleteByActivePondAndDates(ctx context.Context, activePondId int, dates []time.Time) error
}

type dailyLogRepository struct {
	db *gorm.DB
}

func NewDailyLogRepository(db *gorm.DB) DailyLogRepository {
	return &dailyLogRepository{db: db}
}

func (r *dailyLogRepository) WithTx(tx *gorm.DB) DailyLogRepository {
	return &dailyLogRepository{db: tx}
}

func (r *dailyLogRepository) ListIDAndFeedDateByActivePondRange(ctx context.Context, activePondId int, min, max time.Time) ([]DailyLogIDFeedDate, error) {
	var rows []DailyLogIDFeedDate
	err := r.db.WithContext(ctx).Model(&model.DailyLog{}).
		Select("id", "feed_date").
		Where("active_pond_id = ? AND feed_date >= ? AND feed_date <= ? AND deleted_at IS NULL", activePondId, min, max).
		Order("feed_date").
		Find(&rows).Error
	return rows, err
}

func (r *dailyLogRepository) HardDeleteByIDs(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Delete(&model.DailyLog{}).Error
}

func (r *dailyLogRepository) Upsert(ctx context.Context, logs []*model.DailyLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:     []clause.Column{{Name: "active_pond_id"}, {Name: "feed_date"}},
			TargetWhere: clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: "deleted_at IS NULL"}}},
			DoUpdates: clause.AssignmentColumns([]string{
				"fresh_feed_collection_id", "pellet_feed_collection_id",
				"fresh_morning", "fresh_evening", "pellet_morning", "pellet_evening",
				"death_fish_count", "tourist_catch_count",
				"updated_by", "updated_at",
			}),
		}).
		Create(logs).Error
}

func (r *dailyLogRepository) HardDeleteByActivePondAndDates(ctx context.Context, activePondId int, dates []time.Time) error {
	if len(dates) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().
		Where("active_pond_id = ? AND feed_date IN ?", activePondId, dates).
		Delete(&model.DailyLog{}).Error
}

func (r *dailyLogRepository) ListByActivePondAndMonth(ctx context.Context, activePondId int, start, end time.Time) ([]*model.DailyLog, error) {
	var logs []*model.DailyLog
	err := r.db.WithContext(ctx).
		Where("active_pond_id = ? AND feed_date >= ? AND feed_date <= ? AND deleted_at IS NULL", activePondId, start, end).
		Order("feed_date").
		Find(&logs).Error
	return logs, err
}
