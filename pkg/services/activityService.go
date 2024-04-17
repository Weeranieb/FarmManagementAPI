package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/models/constants"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IActivityService interface {
	Create(request models.CreateActivityRequest, userIdentity string) (*models.CreateActivityResponse, error)
	Get(id int) (*models.Activity, error)
	Update(request *models.Activity, userIdentity string) error
}

type activityServiceImp struct {
	ActivityRepo   repositories.IActivityRepository
	SellDetailRepo repositories.ISellDetailRepository
}

func NewActivityService(activePondRepo repositories.IActivityRepository, sellDetailRepo repositories.ISellDetailRepository) IActivityService {
	return &activityServiceImp{
		ActivityRepo:   activePondRepo,
		SellDetailRepo: sellDetailRepo,
	}
}

// FIXME use transaction
func (sv activityServiceImp) Create(request models.CreateActivityRequest, userIdentity string) (*models.CreateActivityResponse, error) {
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check check activity if exist
	checkActivity, err := sv.ActivityRepo.FirstByQuery("\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", request.Mode, request.ActivityDate, false)
	if err != nil {
		return nil, err
	}

	if checkActivity != nil {
		return nil, errors.New("the activity already exist on the given date")
	}

	// declare variable
	newActivity := &models.Activity{}
	newSellDetail := []models.SellDetail{}
	var ret models.CreateActivityResponse

	request.Transfer(newActivity, &newSellDetail)
	newActivity.UpdatedBy = userIdentity
	newActivity.CreatedBy = userIdentity

	// create user
	newActivity, err = sv.ActivityRepo.Create(newActivity)
	if err != nil {
		return nil, err
	}

	ret.Activity = *newActivity
	sellId := newActivity.Id

	if newActivity.Mode == string(constants.SellType) {
		// add updated by and created by to sell detail
		for i := range newSellDetail {
			newSellDetail[i].SellId = sellId
			newSellDetail[i].UpdatedBy = userIdentity
			newSellDetail[i].CreatedBy = userIdentity
		}

		// create sell detail
		newSellDetail, err = sv.SellDetailRepo.BulkCreate(newSellDetail)
		if err != nil {
			return nil, err
		}

		ret.SellDetail = newSellDetail
	}

	return &ret, nil
}

func (sv activityServiceImp) Get(id int) (*models.Activity, error) {
	return sv.ActivityRepo.TakeById(id)
}

func (sv activityServiceImp) Update(request *models.Activity, userIdentity string) error {
	// update activePond
	request.UpdatedBy = userIdentity
	if err := sv.ActivityRepo.Update(request); err != nil {
		return err
	}
	return nil
}
