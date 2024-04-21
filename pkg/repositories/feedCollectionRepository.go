package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IFeedCollectionRepository interface {
	Create(feedCollection *models.FeedCollection) (*models.FeedCollection, error)
	TakeById(id int) (*models.FeedCollection, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FeedCollection, error)
	Update(feedCollection *models.FeedCollection) error
}

type feedCollectionRepositoryImp struct {
	dbContext *gorm.DB
}

func NewFeedCollectionRepository(db *gorm.DB) IFeedCollectionRepository {
	return &feedCollectionRepositoryImp{
		dbContext: db,
	}
}

func (rp feedCollectionRepositoryImp) Create(request *models.FeedCollection) (*models.FeedCollection, error) {
	if err := rp.dbContext.Table(dbconst.TFeedCollection).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp feedCollectionRepositoryImp) TakeById(id int) (*models.FeedCollection, error) {
	var result *models.FeedCollection
	if err := rp.dbContext.Table(dbconst.TFeedCollection).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found FeedCollection TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp feedCollectionRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.FeedCollection, error) {
	var result *models.FeedCollection
	if err := rp.dbContext.Table(dbconst.TFeedCollection).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found FeedCollection FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp feedCollectionRepositoryImp) Update(request *models.FeedCollection) error {
	if err := rp.dbContext.Table(dbconst.TFeedCollection).Where("\"Id\" = ?", request.Id).Updates(&request).Error; err != nil {
		return err
	}
	return nil
}
