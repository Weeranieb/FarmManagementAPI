package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"errors"
	"time"
)

type IDailyFeedService interface {
	Create(request models.AddDailyFeed, userIdentity string) (*models.DailyFeed, error)
	BulkCreate(dailyFeeds []*models.AddDailyFeed, userIdentity string) error
	Get(id int) (*models.DailyFeed, error)
	Update(request *models.DailyFeed, userIdentity string) error
	IsFeedOnDateAvailable(feedId, farmId, year int, month *int) (bool, error)
	GetDailyFeedList(farmId, feedId int, date string) ([]*models.DailyFeed, error)
}

type dailyFeedServiceImp struct {
	DailyFeedRepo repositories.IDailyFeedRepository
}

func NewDailyFeedService(dailyFeedRepo repositories.IDailyFeedRepository) IDailyFeedService {
	return &dailyFeedServiceImp{
		DailyFeedRepo: dailyFeedRepo,
	}
}

// FIXME: implement BulkCreate
func (sv dailyFeedServiceImp) Create(request models.AddDailyFeed, userIdentity string) (*models.DailyFeed, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check feed collection if exist
	checkDailyFeed, err := sv.DailyFeedRepo.FirstByQuery("\"ActivePondId\" = ? AND \"FeedCollectionId\" = ? AND \"FeedDate\" = ? AND \"DelFlag\" = ?", request.PondId, request.FeedCollectionId, request.FeedDate, false)
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

func (sv dailyFeedServiceImp) BulkCreate(dailyFeeds []*models.AddDailyFeed, userIdentity string) error {
	var dailyFeedList []*models.DailyFeed
	for _, dailyFeed := range dailyFeeds {
		var temp models.DailyFeed
		// validate request
		if err := dailyFeed.Validation(); err != nil {
			return err
		}

		dailyFeed.Transfer(&temp)
		temp.UpdatedBy = userIdentity
		temp.CreatedBy = userIdentity

		dailyFeedList = append(dailyFeedList, &temp)
	}

	// create feed collection
	if _, err := sv.DailyFeedRepo.BulkCreate(dailyFeedList); err != nil {
		return err
	}

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

func (sv dailyFeedServiceImp) IsFeedOnDateAvailable(feedId, farmId, year int, month *int) (bool, error) {
	// Construct date range based on year and optional month
	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)

	if month != nil {
		startDate = time.Date(year, time.Month(*month), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)
	}

	dailyFeed, err := sv.DailyFeedRepo.GetDailyFeedByFarm(feedId, farmId, startDate, endDate)
	if err != nil {
		return false, err
	}

	if dailyFeed != nil {
		return false, nil
	}

	return true, nil
}

func (sv dailyFeedServiceImp) GetDailyFeedList(farmId, feedId int, date string) ([]*models.DailyFeed, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	// take just month and year
	startDate := time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	return sv.DailyFeedRepo.TakeAllDailyFeed(feedId, farmId, startDate, endDate)
}
