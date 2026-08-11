package repository

import (
	"context"
	"errors"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FeedPriceHistoryRepository --output=./mocks --outpkg=mocks --filename=feed_price_history_repository.go --structname=MockFeedPriceHistoryRepository --with-expecter=false
type FeedPriceHistoryRepository interface {
	Create(ctx context.Context, feedPriceHistory *model.FeedPriceHistory) error
	CreateBatch(ctx context.Context, feedPriceHistories []*model.FeedPriceHistory) error
	GetByID(id int) (*model.FeedPriceHistory, error)
	GetByFeedCollectionIdAndDate(feedCollectionId int, priceUpdatedDate time.Time) (*model.FeedPriceHistory, error)
	ListByFeedCollectionId(feedCollectionId int) ([]*model.FeedPriceHistory, error)
	ListByFeedCollectionIds(feedCollectionIds []int) ([]*model.FeedPriceHistory, error)
	Update(ctx context.Context, feedPriceHistory *model.FeedPriceHistory) error
	Delete(ctx context.Context, id int) error
}

type feedPriceHistoryRepository struct {
	db *gorm.DB
}

func NewFeedPriceHistoryRepository(db *gorm.DB) FeedPriceHistoryRepository {
	return &feedPriceHistoryRepository{db: db}
}

func (r *feedPriceHistoryRepository) Create(ctx context.Context, feedPriceHistory *model.FeedPriceHistory) error {
	return r.db.WithContext(ctx).Create(feedPriceHistory).Error
}

func (r *feedPriceHistoryRepository) CreateBatch(ctx context.Context, feedPriceHistories []*model.FeedPriceHistory) error {
	if len(feedPriceHistories) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(feedPriceHistories).Error
}

func (r *feedPriceHistoryRepository) GetByID(id int) (*model.FeedPriceHistory, error) {
	var feedPriceHistory model.FeedPriceHistory
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&feedPriceHistory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feedPriceHistory, nil
}

func (r *feedPriceHistoryRepository) GetByFeedCollectionIdAndDate(feedCollectionId int, priceUpdatedDate time.Time) (*model.FeedPriceHistory, error) {
	var feedPriceHistory model.FeedPriceHistory
	err := r.db.Where("feed_collection_id = ? AND price_updated_date = ? AND deleted_at IS NULL", feedCollectionId, priceUpdatedDate).First(&feedPriceHistory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feedPriceHistory, nil
}

func (r *feedPriceHistoryRepository) ListByFeedCollectionId(feedCollectionId int) ([]*model.FeedPriceHistory, error) {
	var feedPriceHistories []*model.FeedPriceHistory
	err := r.db.Where("feed_collection_id = ? AND deleted_at IS NULL", feedCollectionId).
		Order("price_updated_date DESC").
		Find(&feedPriceHistories).Error
	return feedPriceHistories, err
}

// ListByFeedCollectionIds returns price history for several feed collections in
// one query (used by whole-cycle feed-cost aggregation). Returns an empty slice
// when no ids are given. Callers group and order by FeedCollectionId themselves
// (see groupPriceHistoryAscending), so no ORDER BY is imposed here.
func (r *feedPriceHistoryRepository) ListByFeedCollectionIds(feedCollectionIds []int) ([]*model.FeedPriceHistory, error) {
	if len(feedCollectionIds) == 0 {
		return []*model.FeedPriceHistory{}, nil
	}
	var feedPriceHistories []*model.FeedPriceHistory
	err := r.db.Where("feed_collection_id IN ? AND deleted_at IS NULL", feedCollectionIds).
		Find(&feedPriceHistories).Error
	return feedPriceHistories, err
}

func (r *feedPriceHistoryRepository) Update(ctx context.Context, feedPriceHistory *model.FeedPriceHistory) error {
	return r.db.WithContext(ctx).Save(feedPriceHistory).Error
}

// Delete soft-deletes a price-history row (BaseModel.DeletedAt), so it drops out
// of every `deleted_at IS NULL` read but stays recoverable.
func (r *feedPriceHistoryRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.FeedPriceHistory{}, id).Error
}
