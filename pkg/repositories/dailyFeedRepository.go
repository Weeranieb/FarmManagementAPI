package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type IDailyFeedRepository interface {
	Create(dailyFeed *models.DailyFeed) (*models.DailyFeed, error)
	BulkCreate(dailyFeeds []*models.DailyFeed) ([]*models.DailyFeed, error)
	TakeById(id int) (*models.DailyFeed, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.DailyFeed, error)
	Update(dailyFeed *models.DailyFeed) error
	GetDailyFeedByFarm(feedId, farmId int, startDate, endDate time.Time) (*models.DailyFeed, error)
	TakeAllDailyFeed(feedId, farmId int, startDate, endDate time.Time) ([]*models.DailyFeed, error)
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

func (rp dailyFeedRepositoryImp) GetDailyFeedByFarm(feedId, farmId int, startDate, endDate time.Time) (*models.DailyFeed, error) {
	var result *models.DailyFeed
	if err := rp.dbContext.Table(dbconst.TDailyFeed).
		Joins(fmt.Sprintf("JOIN %s ON %s.\"PondId\" = %s.\"Id\"", dbconst.TPond, dbconst.TDailyFeed, dbconst.TPond)).
		Where(fmt.Sprintf("%s.\"FeedCollectionId\" = ?", dbconst.TDailyFeed), feedId).
		Where(fmt.Sprintf("%s.\"FeedDate\" >= ? AND %s.\"FeedDate\" < ?", dbconst.TDailyFeed, dbconst.TDailyFeed), startDate, endDate).
		Where(fmt.Sprintf("%s.\"DelFlag\" = ?", dbconst.TDailyFeed), false).
		Where(fmt.Sprintf("%s.\"DelFlag\" = ?", dbconst.TPond), false).
		Where(fmt.Sprintf("%s.\"FarmId\" = ?", dbconst.TPond), farmId).
		First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		fmt.Println("Record not found DailyFeed GetDailyFeedByFarm", feedId, farmId, startDate, endDate)
		return nil, nil
	}
	return result, nil
}

func (rp dailyFeedRepositoryImp) TakeAllDailyFeed(feedId, farmId int, startDate, endDate time.Time) ([]*models.DailyFeed, error) {
	var result []*models.DailyFeed
	if err := rp.dbContext.Table(dbconst.TDailyFeed).
		Joins(fmt.Sprintf("JOIN %s ON %s.\"PondId\" = %s.\"Id\"", dbconst.TPond, dbconst.TDailyFeed, dbconst.TPond)).
		Where(fmt.Sprintf("%s.\"FeedCollectionId\" = ?", dbconst.TDailyFeed), feedId).
		Where(fmt.Sprintf("%s.\"FeedDate\" >= ? AND %s.\"FeedDate\" < ?", dbconst.TDailyFeed, dbconst.TDailyFeed), startDate, endDate).
		Where(fmt.Sprintf("%s.\"DelFlag\" = ?", dbconst.TDailyFeed), false).
		Where(fmt.Sprintf("%s.\"DelFlag\" = ?", dbconst.TPond), false).
		Where(fmt.Sprintf("%s.\"FarmId\" = ?", dbconst.TPond), farmId).
		Order(fmt.Sprintf("%s.\"Id\" ASC, %s.\"FeedDate\" ASC", dbconst.TPond, dbconst.TDailyFeed)).
		Find(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		fmt.Println("Record not found DailyFeed GetDailyFeedByFarm", feedId, farmId, startDate, endDate)
		return nil, nil
	}
	return result, nil
}
