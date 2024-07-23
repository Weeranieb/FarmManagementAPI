package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"

	"gorm.io/gorm"
)

type IFeedPriceHistoryService interface {
	Create(request models.AddFeedPriceHistory, userIdentity string) (*models.FeedPriceHistory, error)
	Bulk(request []models.AddFeedPriceHistory, userIdentity string, feedCollectionId int) ([]models.FeedPriceHistory, error)
	Get(id int) (*models.FeedPriceHistory, error)
	Update(request *models.FeedPriceHistory, userIdentity string) error
	GetAll(feedCollectionId int) (*[]models.FeedPriceHistory, error)
	WithTrx(trxHandle *gorm.DB) IFeedPriceHistoryService
}

type feedPriceHistoryServiceImp struct {
	FeedPriceHistoryRepo repositories.IFeedPriceHistoryRepository
}

func NewFeedPriceHistoryService(feedPriceHistoryRepo repositories.IFeedPriceHistoryRepository) IFeedPriceHistoryService {
	return &feedPriceHistoryServiceImp{
		FeedPriceHistoryRepo: feedPriceHistoryRepo,
	}
}

func (sv feedPriceHistoryServiceImp) Create(request models.AddFeedPriceHistory, userIdentity string) (*models.FeedPriceHistory, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	if request.FeedCollectionId == 0 {
		return nil, errors.New("feed collection id is empty")
	}

	// check feed price history if exist
	checkFeedPriceHistory, err := sv.FeedPriceHistoryRepo.FirstByQuery("\"FeedCollectionId\" = ? AND \"PriceUpdatedDate\" = ? AND \"DelFlag\" = ?", request.FeedCollectionId, request.PriceUpdatedDate, false)
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
	newFeedPriceHistory, err = sv.FeedPriceHistoryRepo.Create(newFeedPriceHistory)
	if err != nil {
		return nil, err
	}

	return newFeedPriceHistory, nil
}

func (sv feedPriceHistoryServiceImp) Bulk(request []models.AddFeedPriceHistory, userIdentity string, feedCollectionId int) ([]models.FeedPriceHistory, error) {
	// validate request
	for _, req := range request {
		if err := req.Validation(); err != nil {
			return nil, err
		}
	}

	// check feed price history if exist
	for _, req := range request {
		checkFeedPriceHistory, err := sv.FeedPriceHistoryRepo.FirstByQuery("\"FeedCollectionId\" = ? AND \"PriceUpdatedDate\" = ? AND \"DelFlag\" = ?", feedCollectionId, req.PriceUpdatedDate, false)
		if err != nil {
			return nil, err
		}

		if checkFeedPriceHistory != nil {
			return nil, errors.New("feed price history already exist")
		}
	}

	newFeedPriceHistory := []models.FeedPriceHistory{}
	for _, req := range request {
		var temp models.FeedPriceHistory
		req.Transfer(&temp)
		temp.FeedCollectionId = feedCollectionId
		temp.UpdatedBy = userIdentity
		temp.CreatedBy = userIdentity
		newFeedPriceHistory = append(newFeedPriceHistory, temp)
	}

	// create user
	newFeedPriceHistory, err := sv.FeedPriceHistoryRepo.BulkCreate(newFeedPriceHistory)
	if err != nil {
		return nil, err
	}

	return newFeedPriceHistory, nil
}

func (sv feedPriceHistoryServiceImp) Get(id int) (*models.FeedPriceHistory, error) {
	return sv.FeedPriceHistoryRepo.TakeById(id)
}

func (sv feedPriceHistoryServiceImp) Update(request *models.FeedPriceHistory, userIdentity string) error {
	// update user
	request.UpdatedBy = userIdentity
	if err := sv.FeedPriceHistoryRepo.Update(request); err != nil {
		return err
	}
	return nil
}

func (sv feedPriceHistoryServiceImp) GetAll(feedCollectionId int) (*[]models.FeedPriceHistory, error) {
	return sv.FeedPriceHistoryRepo.TakeAll(feedCollectionId)
}

func (sv feedPriceHistoryServiceImp) WithTrx(trxHandle *gorm.DB) IFeedPriceHistoryService {
	sv.FeedPriceHistoryRepo = sv.FeedPriceHistoryRepo.WithTrx(trxHandle)
	return sv
}
