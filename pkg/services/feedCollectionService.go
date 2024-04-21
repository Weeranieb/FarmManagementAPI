package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IFeedCollectionService interface {
	Create(request models.AddFeedCollection, userIdentity string) (*models.FeedCollection, error)
	Get(id int) (*models.FeedCollection, error)
	Update(request *models.FeedCollection, userIdentity string) error
}

type feedCollectionServiceImp struct {
	FeedCollection repositories.IFeedCollectionRepository
}

func NewFeedCollectionService(feedCollectionRepo repositories.IFeedCollectionRepository) IFeedCollectionService {
	return &feedCollectionServiceImp{
		FeedCollection: feedCollectionRepo,
	}
}

func (sv feedCollectionServiceImp) Create(request models.AddFeedCollection, userIdentity string) (*models.FeedCollection, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check feed collection if exist
	checkFeedCollection, err := sv.FeedCollection.FirstByQuery("\"ClientId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", request.ClientId, request.Code, false)
	if err != nil {
		return nil, err
	}

	if checkFeedCollection != nil {
		return nil, errors.New("feed collection already exist")
	}

	newFeedCollection := &models.FeedCollection{}
	request.Transfer(newFeedCollection)
	newFeedCollection.UpdatedBy = userIdentity
	newFeedCollection.CreatedBy = userIdentity

	// create feed collection
	newFeedCollection, err = sv.FeedCollection.Create(newFeedCollection)
	if err != nil {
		return nil, err
	}

	return newFeedCollection, nil
}

func (sv feedCollectionServiceImp) Get(id int) (*models.FeedCollection, error) {
	return sv.FeedCollection.TakeById(id)
}

func (sv feedCollectionServiceImp) Update(request *models.FeedCollection, userIdentity string) error {
	// update feed collection
	request.UpdatedBy = userIdentity
	if err := sv.FeedCollection.Update(request); err != nil {
		return err
	}
	return nil
}
