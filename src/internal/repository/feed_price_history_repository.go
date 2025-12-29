package repository

import (
	"errors"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

type FeedPriceHistoryRepository interface {
	Create(feedPriceHistory *model.FeedPriceHistory) error
	CreateBatch(feedPriceHistories []*model.FeedPriceHistory) error
	GetByID(id int) (*model.FeedPriceHistory, error)
	GetByFeedCollectionIdAndDate(feedCollectionId int, priceUpdatedDate time.Time) (*model.FeedPriceHistory, error)
	ListByFeedCollectionId(feedCollectionId int) ([]*model.FeedPriceHistory, error)
	Update(feedPriceHistory *model.FeedPriceHistory) error
}

type feedPriceHistoryRepository struct {
	db *gorm.DB
}

func NewFeedPriceHistoryRepository(db *gorm.DB) FeedPriceHistoryRepository {
	return &feedPriceHistoryRepository{db: db}
}

func (r *feedPriceHistoryRepository) Create(feedPriceHistory *model.FeedPriceHistory) error {
	return r.db.Create(feedPriceHistory).Error
}

func (r *feedPriceHistoryRepository) CreateBatch(feedPriceHistories []*model.FeedPriceHistory) error {
	if len(feedPriceHistories) == 0 {
		return nil
	}
	return r.db.Create(feedPriceHistories).Error
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

func (r *feedPriceHistoryRepository) Update(feedPriceHistory *model.FeedPriceHistory) error {
	return r.db.Save(feedPriceHistory).Error
}
