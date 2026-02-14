package repository

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FeedCollectionRepository --output=./mocks --outpkg=mocks --filename=feed_collection_repository.go --structname=MockFeedCollectionRepository --with-expecter=false
type FeedCollectionRepository interface {
	Create(ctx context.Context, feedCollection *model.FeedCollection) error
	GetByID(id int) (*model.FeedCollection, error)
	GetByClientIdAndName(clientId int, name string) (*model.FeedCollection, error)
	Update(ctx context.Context, feedCollection *model.FeedCollection) error
	GetPage(clientId, page, pageSize int, orderBy, keyword string) ([]*model.FeedCollectionPage, int64, error)
}

type feedCollectionRepository struct {
	db *gorm.DB
}

func NewFeedCollectionRepository(db *gorm.DB) FeedCollectionRepository {
	return &feedCollectionRepository{db: db}
}

func (r *feedCollectionRepository) Create(ctx context.Context, feedCollection *model.FeedCollection) error {
	return r.db.WithContext(ctx).Create(feedCollection).Error
}

func (r *feedCollectionRepository) GetByID(id int) (*model.FeedCollection, error) {
	var feedCollection model.FeedCollection
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&feedCollection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feedCollection, nil
}

func (r *feedCollectionRepository) GetByClientIdAndName(clientId int, name string) (*model.FeedCollection, error) {
	var feedCollection model.FeedCollection
	err := r.db.Where("client_id = ? AND name = ? AND deleted_at IS NULL", clientId, name).First(&feedCollection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feedCollection, nil
}

func (r *feedCollectionRepository) Update(ctx context.Context, feedCollection *model.FeedCollection) error {
	return r.db.WithContext(ctx).Save(feedCollection).Error
}

func (r *feedCollectionRepository) GetPage(clientId, page, pageSize int, orderBy, keyword string) ([]*model.FeedCollectionPage, int64, error) {
	var feedCollections []*model.FeedCollectionPage
	var total int64

	// Subquery to find the latest price update date for each feed collection
	subQuery := r.db.Model(&model.FeedPriceHistory{}).
		Select("feed_price_histories.feed_collection_id, MAX(feed_price_histories.price_updated_date) as latest_price_updated_date").
		Group("feed_price_histories.feed_collection_id")

	// Main query with LEFT JOIN to get latest price
	query := r.db.Table("feed_collections").
		Select(`feed_collections.*, 
			feed_price_histories.price as latest_price,
			feed_price_histories.price_updated_date as latest_price_updated_date`).
		Joins("LEFT JOIN (?) as latest_price_history ON feed_collections.id = latest_price_history.feed_collection_id", subQuery).
		Joins("LEFT JOIN feed_price_histories ON feed_collections.id = feed_price_histories.feed_collection_id AND feed_price_histories.price_updated_date = latest_price_history.latest_price_updated_date").
		Where("feed_collections.client_id = ? AND feed_collections.deleted_at IS NULL", clientId)

	if keyword != "" {
		query = query.Where("(feed_collections.name LIKE ? OR feed_collections.unit LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	// Apply pagination
	offset := page * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&feedCollections).Error; err != nil {
		return nil, 0, err
	}

	return feedCollections, total, nil
}
