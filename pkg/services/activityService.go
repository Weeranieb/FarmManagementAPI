package services

import (
	dbContext "boonmafarm/api/pkg/dbcontext"
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/models/constants"
	"boonmafarm/api/pkg/repositories"
	"errors"
)

type IActivityService interface {
	Create(request models.CreateActivityRequest, userIdentity string) (*models.ActivityWithSellDetail, error)
	Get(id int) (*models.ActivityWithSellDetail, error)
	Update(request *models.ActivityWithSellDetail, userIdentity string) error
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

func (sv activityServiceImp) Create(request models.CreateActivityRequest, userIdentity string) (*models.ActivityWithSellDetail, error) {
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
	var ret models.ActivityWithSellDetail

	request.Transfer(newActivity, &newSellDetail)
	newActivity.UpdatedBy = userIdentity
	newActivity.CreatedBy = userIdentity

	db := dbContext.Context.Postgresql
	tx := db.Begin()
	// create user
	newActivity, err = sv.ActivityRepo.WithTrx(tx).Create(newActivity)
	if err != nil {
		tx.Rollback()
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
		newSellDetail, err = sv.SellDetailRepo.WithTrx(tx).BulkCreate(newSellDetail)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		ret.SellDetail = newSellDetail
	}

	// commit transaction
	tx.Commit()

	return &ret, nil
}

// get with sell detail check case
func (sv activityServiceImp) Get(id int) (*models.ActivityWithSellDetail, error) {
	// get activity
	var payload models.ActivityWithSellDetail
	activity, err := sv.ActivityRepo.TakeById(id)
	if err != nil {
		return nil, err
	}

	payload.Activity = *activity

	if activity.Mode == string(constants.SellType) {
		// get sell details
		sellDetails, err := sv.SellDetailRepo.ListByQuery("\"SellId\" = ? AND \"DelFlag\" = ?", id, false)
		if err != nil {
			return nil, err
		}
		payload.SellDetail = sellDetails
	}

	return &payload, nil
}

// update with sell detail check case
func (sv activityServiceImp) Update(request *models.ActivityWithSellDetail, userIdentity string) error {
	// update activePond
	// request.UpdatedBy = userIdentity
	// if err := sv.ActivityRepo.Update(request); err != nil {
	// 	return err
	// }
	return nil
}
