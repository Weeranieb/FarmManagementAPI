package repositories

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/dbconst"
	"boonmafarm/api/utils/dbutil"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IFeedCollectionRepository interface {
	Create(feedCollection *models.FeedCollection) (*models.FeedCollection, error)
	TakeById(id int) (*models.FeedCollection, error)
	TakePage(clientId, page, pageSize int, orderBy, keyword string) (*[]models.FeedCollectionPage, int64, error)
	FirstByQuery(query interface{}, args ...interface{}) (*models.FeedCollection, error)
	Update(feedCollection *models.FeedCollection) error
	WithTrx(trxHandle *gorm.DB) IFeedCollectionRepository
}

type feedCollectionRepositoryImp struct {
	dbContext *gorm.DB
}

func NewFeedCollectionRepository(db *gorm.DB) IFeedCollectionRepository {
	return &feedCollectionRepositoryImp{
		dbContext: db,
	}
}

func (rp feedCollectionRepositoryImp) WithTrx(trxHandle *gorm.DB) IFeedCollectionRepository {
	if trxHandle == nil {
		fmt.Println("Transaction Database not found")
		return rp
	}

	return &feedCollectionRepositoryImp{
		dbContext: trxHandle,
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

func (rp feedCollectionRepositoryImp) TakePage(clientId, page, pageSize int, orderBy, keyword string) (*[]models.FeedCollectionPage, int64, error) {
	var result []models.FeedCollectionPage
	var total int64

	// Subquery to find the latest price update date for each feed collection
	subQuery := rp.dbContext.Table(dbconst.TFeedPriceHistory).
		Select("\"FeedCollectionId\", MAX(\"PriceUpdatedDate\") as \"LatestPriceUpdatedDate\"").
		Group("\"FeedCollectionId\"")

	// Main query
	query := rp.dbContext.Table(dbconst.TFeedCollection).
		Select(fmt.Sprintf(`%s.*, 
            %s."Price" as "LatestPrice",
            %s."PriceUpdatedDate" as "LatestPriceUpdatedDate"`, dbconst.TFeedCollection, dbconst.TFeedPriceHistory, dbconst.TFeedPriceHistory)).
		Joins(fmt.Sprintf("LEFT JOIN (?) as LatestPriceHistory ON %s.\"Id\" = LatestPriceHistory.\"FeedCollectionId\"", dbconst.TFeedCollection), subQuery).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.\"Id\" = %s.\"FeedCollectionId\" AND %s.\"PriceUpdatedDate\" = LatestPriceHistory.\"LatestPriceUpdatedDate\"", dbconst.TFeedPriceHistory, dbconst.TFeedCollection, dbconst.TFeedPriceHistory, dbconst.TFeedPriceHistory)).
		Where(fmt.Sprintf("%s.\"ClientId\" = ? AND %s.\"DelFlag\" = ?", dbconst.TFeedCollection, dbconst.TFeedCollection), clientId, false).
		Order(orderBy)

	if keyword != "" {
		whereKeyword := fmt.Sprintf("(%s.\"Code\" LIKE ? OR %s.\"Name\" LIKE ? OR %s.\"Unit\" LIKE ?)", dbconst.TFeedCollection, dbconst.TFeedCollection, dbconst.TFeedCollection)
		query = query.Where(whereKeyword, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Limit(pageSize).Offset(page * pageSize).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	return &result, total, nil
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
	obj := dbutil.StructToMap(request)
	if err := rp.dbContext.Table(dbconst.TFeedCollection).Where("\"Id\" = ?", request.Id).Updates(obj).Error; err != nil {
		return err
	}
	return nil
}
