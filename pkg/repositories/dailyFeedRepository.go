package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IDailyFeedRepository interface {
	Create(dailyFeed *models.DailyFeed) (*models.DailyFeed, error)
	BulkCreate(dailyFeeds []*models.DailyFeed) ([]*models.DailyFeed, error)
	TakeById(id int) (*models.DailyFeed, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.DailyFeed, error)
	Update(dailyFeed *models.DailyFeed) error
}

type dailyFeedRepositoryImp struct {
	dbContext *gorm.DB
}

func NewDailyFeedRepository(db *gorm.DB) IDailyFeedRepository {
	return &dailyFeedRepositoryImp{
		dbContext: db,
	}
}

func (rp dailyFeedRepositoryImp) Create(request *models.DailyFeed) (*models.DailyFeed, error) {
	if err := rp.dbContext.Table(dbconst.TDailyFeed).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rp dailyFeedRepositoryImp) BulkCreate(requests []*models.DailyFeed) ([]*models.DailyFeed, error) {
	if err := rp.dbContext.Table(dbconst.TDailyFeed).Create(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (rp dailyFeedRepositoryImp) TakeById(id int) (*models.DailyFeed, error) {
	var result *models.DailyFeed
	if err := rp.dbContext.Table(dbconst.TDailyFeed).Where("\"Id\" = ? AND \"DelFlag\" = ?", id, false).Take(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found DailyFeed TakeById", id)
		return nil, nil
	}
	return result, nil
}

func (rp dailyFeedRepositoryImp) FirstByQuery(query interface{}, args ...interface{}) (*models.DailyFeed, error) {
	var result *models.DailyFeed
	if err := rp.dbContext.Table(dbconst.TDailyFeed).Where(query, args...).First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		fmt.Println("Record not found DailyFeed FirstByQuery", query)
		return nil, nil
	}
	return result, nil
}

func (rp dailyFeedRepositoryImp) Update(request *models.DailyFeed) error {
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TDailyFeed).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}
