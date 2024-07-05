package services

import (
	dbContext "boonmafarm/api/pkg/dbcontext"
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/models/constants"
	"boonmafarm/api/pkg/repositories"
	"boonmafarm/api/utils/httputil"
	"errors"
)

type IActivityService interface {
	Create(request models.CreateActivityRequest, userIdentity string) (*models.ActivityWithSellDetail, error)
	Get(id int) (*models.ActivityWithSellDetail, error)
	Update(request *models.ActivityWithSellDetail, userIdentity string) ([]*models.SellDetail, error)
	TakePage(clientId, page, pageSize int, orderBy, keyword string, mode string, farmId int) (*httputil.PageModel, error)
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

		if len(sellDetails) == 0 {
			return nil, errors.New("sell detail not found")
		}

		payload.SellDetail = sellDetails
	}

	return &payload, nil
}

// update with sell detail check case
func (sv activityServiceImp) Update(request *models.ActivityWithSellDetail, userIdentity string) ([]*models.SellDetail, error) {
	// update activePond
	db := dbContext.Context.Postgresql
	tx := db.Begin()

	// update activity
	if err := sv.ActivityRepo.WithTrx(tx).Update(&request.Activity); err != nil {
		tx.Rollback()
		return nil, err
	}

	var newSellDetails []*models.SellDetail
	if request.Activity.Mode == string(constants.SellType) {
		for _, sellDetail := range request.SellDetail {
			if sellDetail.Id == 0 {
				sellDetail.SellId = request.Activity.Id
				sellDetail.CreatedBy = userIdentity
				sellDetail.UpdatedBy = userIdentity
				newSellDetail, err := sv.SellDetailRepo.WithTrx(tx).Create(&sellDetail)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				newSellDetails = append(newSellDetails, newSellDetail)
			} else {
				sellDetail.UpdatedBy = userIdentity
				if err := sv.SellDetailRepo.WithTrx(tx).Update(&sellDetail); err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	tx.Commit()

	return newSellDetails, nil
}

func (sv activityServiceImp) TakePage(clientId, page, pageSize int, orderBy, keyword string, mode string, farmId int) (*httputil.PageModel, error) {
	result := &httputil.PageModel{}
	var modePointer *string
	var farmIdPointer *int

	if mode != "" {
		modePointer = &mode
	}

	if farmId != 0 {
		farmIdPointer = &farmId
	}

	items, total, err := sv.ActivityRepo.TakePage(clientId, page, pageSize, orderBy, keyword, modePointer, farmIdPointer)
	if err != nil {
		return nil, err
	}

	result.Items = items
	result.Total = total

	return result, nil
}
