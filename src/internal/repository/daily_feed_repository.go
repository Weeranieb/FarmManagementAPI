package repository

import (
	"context"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyFeedRepository --output=./mocks --outpkg=mocks --filename=daily_feed_repository.go --structname=MockDailyFeedRepository --with-expecter=false
type DailyFeedRepository interface {
	Upsert(ctx context.Context, feeds []*model.DailyFeed) error
	ListByActivePondAndMonth(activePondId int, start, end time.Time) ([]*model.DailyFeed, error)
	ListFeedCollectionIdsByActivePond(activePondId int) ([]int, error)
	SoftDeleteByActivePondAndFeedCollection(ctx context.Context, activePondId, feedCollectionId int) error
}

type dailyFeedRepository struct {
	db *gorm.DB
}

func NewDailyFeedRepository(db *gorm.DB) DailyFeedRepository {
	return &dailyFeedRepository{db: db}
}

func (r *dailyFeedRepository) Upsert(ctx context.Context, feeds []*model.DailyFeed) error {
	if len(feeds) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "active_pond_id"}, {Name: "feed_collection_id"}, {Name: "feed_date"}},
			DoUpdates: clause.AssignmentColumns([]string{"morning_amount", "evening_amount", "updated_by", "updated_at"}),
			Where:     clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: "daily_feeds.deleted_at IS NULL"}}},
		}).
		Create(feeds).Error
}

func (r *dailyFeedRepository) ListByActivePondAndMonth(activePondId int, start, end time.Time) ([]*model.DailyFeed, error) {
	var feeds []*model.DailyFeed
	err := r.db.
		Where("active_pond_id = ? AND feed_date >= ? AND feed_date <= ? AND deleted_at IS NULL", activePondId, start, end).
		Order("feed_collection_id, feed_date").
		Find(&feeds).Error
	return feeds, err
}

func (r *dailyFeedRepository) ListFeedCollectionIdsByActivePond(activePondId int) ([]int, error) {
	var ids []int
	err := r.db.Model(&model.DailyFeed{}).
		Where("active_pond_id = ? AND deleted_at IS NULL", activePondId).
		Distinct("feed_collection_id").
		Pluck("feed_collection_id", &ids).Error
	return ids, err
}

func (r *dailyFeedRepository) SoftDeleteByActivePondAndFeedCollection(ctx context.Context, activePondId, feedCollectionId int) error {
	return r.db.WithContext(ctx).
		Where("active_pond_id = ? AND feed_collection_id = ? AND deleted_at IS NULL", activePondId, feedCollectionId).
		Delete(&model.DailyFeed{}).Error
}
