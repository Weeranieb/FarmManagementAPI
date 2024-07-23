package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"boonmafarm/api/utils/httputil"
	"errors"

	"gorm.io/gorm"
)

type IFeedCollectionService interface {
	Create(request models.AddFeedCollection, userIdentity string, clientId int) (*models.FeedCollection, error)
	Get(id int) (*models.FeedCollection, error)
	Update(request *models.FeedCollection, userIdentity string) error
	TakePage(clientId, page, pageSize int, orderBy, keyword string) (*httputil.PageModel, error)
	WithTrx(trxHandle *gorm.DB) IFeedCollectionService
}

type feedCollectionServiceImp struct {
	FeedCollectionRepo repositories.IFeedCollectionRepository
}

func NewFeedCollectionService(feedCollectionRepo repositories.IFeedCollectionRepository) IFeedCollectionService {
	return &feedCollectionServiceImp{
		FeedCollectionRepo: feedCollectionRepo,
	}
}

func (sv feedCollectionServiceImp) Create(request models.AddFeedCollection, userIdentity string, clientId int) (*models.FeedCollection, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check feed collection if exist
	checkFeedCollection, err := sv.FeedCollectionRepo.FirstByQuery("\"ClientId\" = ? AND \"Code\" = ? AND \"DelFlag\" = ?", clientId, request.Code, false)
	if err != nil {
		return nil, err
	}

	if checkFeedCollection != nil {
		return nil, errors.New("feed collection already exist")
	}

	newFeedCollection := &models.FeedCollection{}
	request.Transfer(newFeedCollection)
	newFeedCollection.ClientId = clientId
	newFeedCollection.UpdatedBy = userIdentity
	newFeedCollection.CreatedBy = userIdentity

	// create feed collection
	newFeedCollection, err = sv.FeedCollectionRepo.Create(newFeedCollection)
	if err != nil {
		return nil, err
	}

	return newFeedCollection, nil
}

func (sv feedCollectionServiceImp) Get(id int) (*models.FeedCollection, error) {
	return sv.FeedCollectionRepo.TakeById(id)
}

func (sv feedCollectionServiceImp) Update(request *models.FeedCollection, userIdentity string) error {
	// update feed collection
	request.UpdatedBy = userIdentity
	if err := sv.FeedCollectionRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv feedCollectionServiceImp) TakePage(clientId, page, pageSize int, orderBy, keyword string) (*httputil.PageModel, error) {
	result := &httputil.PageModel{}
	items, total, err := sv.FeedCollectionRepo.TakePage(clientId, page, pageSize, orderBy, keyword)
	if err != nil {
		return nil, err
	}

	result.Items = items
	result.Total = total

	return result, nil
}

func (sv feedCollectionServiceImp) WithTrx(trxHandle *gorm.DB) IFeedCollectionService {
	sv.FeedCollectionRepo = sv.FeedCollectionRepo.WithTrx(trxHandle)
	return sv
}
