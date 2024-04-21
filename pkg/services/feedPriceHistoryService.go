package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IFeedPriceHistoryService interface {
	Create(request models.AddFeedPriceHistory, userIdentity string) (*models.FeedPriceHistory, error)
	Get(id int) (*models.FeedPriceHistory, error)
	Update(request *models.FeedPriceHistory, userIdentity string) error
}

type feedPriceHistoryServiceImp struct {
	FeedPriceHistory repositories.IFeedPriceHistoryRepository
}

func NewFeedPriceHistoryService(feedPriceHistoryRepo repositories.IFeedPriceHistoryRepository) IFeedPriceHistoryService {
	return &feedPriceHistoryServiceImp{
		FeedPriceHistory: feedPriceHistoryRepo,
	}
}

func (sv feedPriceHistoryServiceImp) Create(request models.AddFeedPriceHistory, userIdentity string) (*models.FeedPriceHistory, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check feed price history if exist
	checkFeedPriceHistory, err := sv.FeedPriceHistory.FirstByQuery("\"FeedCollectionId\" = ? AND \"PriceUpdatedDate\" = ? AND \"DelFlag\" = ?", request.FeedCollectionId, request.PriceUpdatedDate, false)
	if err != nil {
		return nil, err
	}

	if checkFeedPriceHistory != nil {
		return nil, errors.New("feed price history already exist")
	}

	newFeedPriceHistory := &models.FeedPriceHistory{}
	request.Transfer(newFeedPriceHistory)
	newFeedPriceHistory.UpdatedBy = userIdentity
	newFeedPriceHistory.CreatedBy = userIdentity

	// create user
	newFeedPriceHistory, err = sv.FeedPriceHistory.Create(newFeedPriceHistory)
	if err != nil {
		return nil, err
	}

	return newFeedPriceHistory, nil
}

func (sv feedPriceHistoryServiceImp) Get(id int) (*models.FeedPriceHistory, error) {
	return sv.FeedPriceHistory.TakeById(id)
}

func (sv feedPriceHistoryServiceImp) Update(request *models.FeedPriceHistory, userIdentity string) error {
	// update user
	request.UpdatedBy = userIdentity
	if err := sv.FeedPriceHistory.Update(request); err != nil {
		return err
	}
	return nil
}
