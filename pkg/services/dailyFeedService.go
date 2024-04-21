package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
	"mime/multipart"
)

type IDailyFeedService interface {
	Create(request models.AddDailyFeed, userIdentity string) (*models.DailyFeed, error)
	BulkCreate(excelFile multipart.File, userIdentity string) error
	Get(id int) (*models.DailyFeed, error)
	Update(request *models.DailyFeed, userIdentity string) error
}

type dailyFeedServiceImp struct {
	DailyFeedRepo repositories.IFDailyFeedRepository
}

func NewDailyFeedService(dailyFeedRepo repositories.IFDailyFeedRepository) IDailyFeedService {
	return &dailyFeedServiceImp{
		DailyFeedRepo: dailyFeedRepo,
	}
}

func (sv dailyFeedServiceImp) Create(request models.AddDailyFeed, userIdentity string) (*models.DailyFeed, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check feed collection if exist
	checkDailyFeed, err := sv.DailyFeedRepo.FirstByQuery("\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", request.ActivePondId, request.FeedCollectionId, request.FeedDate, false)
	if err != nil {
		return nil, err
	}

	if checkDailyFeed != nil {
		return nil, errors.New("daily feed already exist")
	}

	newDailyFeed := &models.DailyFeed{}
	request.Transfer(newDailyFeed)
	newDailyFeed.UpdatedBy = userIdentity
	newDailyFeed.CreatedBy = userIdentity

	// create feed collection
	newDailyFeed, err = sv.DailyFeedRepo.Create(newDailyFeed)
	if err != nil {
		return nil, err
	}

	return newDailyFeed, nil
}

func (sv dailyFeedServiceImp) BulkCreate(excelFile multipart.File, userIdentity string) error {
	// read excel file
	// FIXME: implement read excel file
	return nil
}

func (sv dailyFeedServiceImp) Get(id int) (*models.DailyFeed, error) {
	return sv.DailyFeedRepo.TakeById(id)
}

func (sv dailyFeedServiceImp) Update(request *models.DailyFeed, userIdentity string) error {
	// update feed collection
	request.UpdatedBy = userIdentity
	if err := sv.DailyFeedRepo.Update(request); err != nil {
		return err
	}
	return nil
}
