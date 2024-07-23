package processors

import (
	dbContext "boonmafarm/api/pkg/dbcontext"
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
)

type IFeedProcessor interface {
	CreateFeedCollection(request models.CreateFeedRequest, userIdentity string, clientId int) (*models.FeedCollection, []models.FeedPriceHistory, error)
}

type feedProcessorImp struct {
	FeedCollectionService   services.IFeedCollectionService
	FeedPriceHistoryService services.IFeedPriceHistoryService
}

func NewFeedProcessor(feedCollectionService services.IFeedCollectionService, feedPriceHistoryService services.IFeedPriceHistoryService) IFeedProcessor {
	return &feedProcessorImp{
		FeedCollectionService:   feedCollectionService,
		FeedPriceHistoryService: feedPriceHistoryService,
	}
}

func (p feedProcessorImp) CreateFeedCollection(request models.CreateFeedRequest, userIdentity string, clientId int) (*models.FeedCollection, []models.FeedPriceHistory, error) {
	db := dbContext.Context.Postgresql
	tx := db.Begin()

	// create feed collection
	var feedCollection *models.FeedCollection
	var addFeedCollection models.AddFeedCollection = models.AddFeedCollection{
		Code: request.Code,
		Name: request.Name,
		Unit: request.Unit,
	}

	feedCollection, err := p.FeedCollectionService.WithTrx(tx).Create(addFeedCollection, userIdentity, clientId)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// create feed price history
	var feedPriceHistories []models.FeedPriceHistory
	feedPriceHistories, err = p.FeedPriceHistoryService.WithTrx(tx).Bulk(request.FeedPriceHistories, userIdentity, feedCollection.Id)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	tx.Commit()

	return feedCollection, feedPriceHistories, nil
}
