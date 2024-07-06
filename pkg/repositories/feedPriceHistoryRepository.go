package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IFeedPriceHistoryRepository interface {
	Create(feedPriceHistory *models.FeedPriceHistory) (*models.FeedPriceHistory, error)
	TakeById(id int) (*models.FeedPriceHistory, error)
	TakeAll(feedCollectionId int) (*[]models.FeedPriceHistory, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FeedPriceHistory, error)
	Update(feedPriceHistory *models.FeedPriceHistory) error
}

type feedPriceHistoryRepositoryImp struct {
	dbContext *gorm.DB
}

func NewFeedPriceHistoryRepository(db *gorm.DB) IFeedPriceHistoryRepository {
	return &feedPriceHistoryRepositoryImp{
		dbContext: db,
	}
}

func (rp feedPriceHistoryRepositoryImp) Create(request *models.FeedPriceHistory) (*models.FeedPriceHistory, error) {
	if err := rp.dbContext.Table(dbconst.TFeedPriceHistory).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp feedPriceHistoryRepositoryImp) TakeById(id int) (*models.FeedPriceHistory, error) {
	var result *models.FeedPriceHistory
	if err := rp.dbContext.Table(dbconst.TFeedPriceHistory).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found FeedPriceHistory TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp feedPriceHistoryRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.FeedPriceHistory, error) {
	var result *models.FeedPriceHistory
	if err := rp.dbContext.Table(dbconst.TFeedPriceHistory).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found FeedPriceHistory FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp feedPriceHistoryRepositoryImp) Update(request *models.FeedPriceHistory) error {
	if err := rp.dbContext.Table(dbconst.TFeedPriceHistory).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}

func (rp feedPriceHistoryRepositoryImp) TakeAll(feedCollectionId int) (*[]models.FeedPriceHistory, error) {
	var result *[]models.FeedPriceHistory
	if err := rp.dbContext.Table(dbconst.TFeedPriceHistory).Where("\"FeedCollectionId\" = ? AND \"DelFlag\" = ?", feedCollectionId, false).Order("\"PriceUpdatedDate\" desc").Find(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found FeedPriceHistory TakeAll", feedCollectionId)
		return nil, nil
	}
	return result, nil
}
