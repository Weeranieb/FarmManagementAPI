package processors

import (
	dbContext "boonmafarm/api/pkg/dbcontext"
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"errors"
)

type IActivityProcessor interface {
	CreateFill(request models.CreateFillActivityRequest, userIdentity string) (*models.Activity, *models.ActivePond, error)
	CreateMove(request models.CreateMoveActivityRequest, userIdentity string) (*models.Activity, *models.ActivePond, *models.ActivePond, error)
	CreateSell(request models.CreateSellActivityRequest, userIdentity string) (*models.Activity, *models.ActivePond, *[]models.SellDetail, error)
}

type activityProcessorImp struct {
	ActivePondService     services.IActivePondService
	ActivityService       services.IActivityService
	AdditionalCostService services.IAdditionalCostService
}

func NewActivityProcessor(activePondService services.IActivePondService, activityService services.IActivityService, additionalCostService services.IAdditionalCostService) IActivityProcessor {
	return &activityProcessorImp{
		ActivePondService:     activePondService,
		ActivityService:       activityService,
		AdditionalCostService: additionalCostService,
	}
}

func (sv activityProcessorImp) CreateFill(request models.CreateFillActivityRequest, userIdentity string) (*models.Activity, *models.ActivePond, error) {
	db := dbContext.Context.Postgresql
	tx := db.Begin()

	isNewPond := request.IsNewPond
	var currentActivePond *models.ActivePond
	if isNewPond {
		// check all activepond if not active left
		isAvailAble, err := sv.ActivePondService.CheckNewActivePondAvailable(request.PondId)
		if err != nil {
			return nil, nil, err
		}

		if !isAvailAble {
			tx.Rollback()
			return nil, nil, errors.New("the pond is already active")
		}

		// create active pond
		var addActivePond models.AddActivePond = models.AddActivePond{
			PondId:    request.PondId,
			StartDate: request.ActivityDate,
		}
		currentActivePond, err = sv.ActivePondService.WithTrx(tx).Create(addActivePond, userIdentity)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	} else {
		// get current active pond
		var err error
		currentActivePond, err = sv.ActivePondService.WithTrx(tx).GetActivePondByDate(request.PondId, request.ActivityDate)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}

		if currentActivePond == nil {
			tx.Rollback()
			return nil, nil, errors.New("the pond is not active")
		}
	}

	// update activity
	newActivity, err := sv.ActivityService.WithTrx(tx).CreateFill(request, userIdentity, currentActivePond.Id)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// create additional cost
	if request.AdditionalCost != nil {
		var addAdditionalCost models.AddAdditionalCostRequest = models.AddAdditionalCostRequest{
			ActivityId: newActivity.Id,
			Cost:       *request.AdditionalCost,
			Title:      "ค่าใช้จ่ายเพิ่มเติม (เติม)",
		}
		_, err := sv.AdditionalCostService.WithTrx(tx).Create(addAdditionalCost, userIdentity)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	// commit transaction
	tx.Commit()

	return newActivity, currentActivePond, nil
}

func (sv activityProcessorImp) CreateMove(request models.CreateMoveActivityRequest, userIdentity string) (*models.Activity, *models.ActivePond, *models.ActivePond, error) {
	db := dbContext.Context.Postgresql
	tx := db.Begin()

	isNewPond := request.IsNewPond
	isClose := request.IsClose

	var fromActivePond *models.ActivePond
	var toActivePond *models.ActivePond

	// get current active pond
	var err error
	fromActivePond, err = sv.ActivePondService.WithTrx(tx).GetActivePondByDate(request.PondId, request.ActivityDate)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if fromActivePond == nil {
		tx.Rollback()
		return nil, nil, nil, errors.New("the pond is not active")
	}

	if isClose {
		if !fromActivePond.IsActive {
			tx.Rollback()
			return nil, nil, nil, errors.New("the pond is already close")
		}

		fromActivePond.IsActive = false
		fromActivePond.EndDate = &request.ActivityDate

		// update active pond
		if err := sv.ActivePondService.WithTrx(tx).Update(fromActivePond, userIdentity); err != nil {
			tx.Rollback()
			return nil, nil, nil, err
		}
	}

	if isNewPond {
		// check all activepond if not active left
		isAvailAble, err := sv.ActivePondService.CheckNewActivePondAvailable(request.ToPondId)
		if err != nil {
			return nil, nil, nil, err
		}

		if !isAvailAble {
			tx.Rollback()
			return nil, nil, nil, errors.New("to the pond is already active")
		}

		// create active pond
		var addActivePond models.AddActivePond = models.AddActivePond{
			PondId:    request.ToPondId,
			StartDate: request.ActivityDate,
		}
		toActivePond, err = sv.ActivePondService.WithTrx(tx).Create(addActivePond, userIdentity)
		if err != nil {
			tx.Rollback()
			return nil, nil, nil, err
		}
	} else {
		// get current active pond
		var err error
		toActivePond, err = sv.ActivePondService.WithTrx(tx).GetActivePondByDate(request.ToPondId, request.ActivityDate)
		if err != nil {
			tx.Rollback()
			return nil, nil, nil, err
		}

		if toActivePond == nil {
			tx.Rollback()
			return nil, nil, nil, errors.New("to the pond is not active")
		}
	}

	// update activity
	newActivity, err := sv.ActivityService.WithTrx(tx).CreateMove(request, userIdentity, fromActivePond.Id, toActivePond.Id)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	// create additional cost
	if request.AdditionalCost != nil {
		var addAdditionalCost models.AddAdditionalCostRequest = models.AddAdditionalCostRequest{
			ActivityId: newActivity.Id,
			Cost:       *request.AdditionalCost,
			Title:      "ค่าใช้จ่ายเพิ่มเติม (ย้าย)",
		}
		_, err := sv.AdditionalCostService.WithTrx(tx).Create(addAdditionalCost, userIdentity)
		if err != nil {
			tx.Rollback()
			return nil, nil, nil, err
		}
	}

	// commit transaction
	tx.Commit()

	return newActivity, fromActivePond, toActivePond, nil
}

func (sv activityProcessorImp) CreateSell(request models.CreateSellActivityRequest, userIdentity string) (*models.Activity, *models.ActivePond, *[]models.SellDetail, error) {
	db := dbContext.Context.Postgresql
	tx := db.Begin()

	isClose := request.IsClose

	var currentActivePond *models.ActivePond

	// get current active pond
	var err error
	currentActivePond, err = sv.ActivePondService.WithTrx(tx).GetActivePondByDate(request.PondId, request.ActivityDate)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	if currentActivePond == nil {
		tx.Rollback()
		return nil, nil, nil, errors.New("the pond is not active")
	}

	// update active pond
	if isClose {
		if !currentActivePond.IsActive {
			tx.Rollback()
			return nil, nil, nil, errors.New("the pond is already close")
		}

		currentActivePond.IsActive = false
		currentActivePond.EndDate = &request.ActivityDate
		if err != sv.ActivePondService.WithTrx(tx).Update(currentActivePond, userIdentity) {
			tx.Rollback()
			return nil, nil, nil, err
		}
	}

	// update activity
	newActivity, sellDetails, err := sv.ActivityService.WithTrx(tx).CreateSell(request, userIdentity, currentActivePond.Id, tx)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	// create additional cost
	if request.AdditionalCost != nil {
		var addAdditionalCost models.AddAdditionalCostRequest = models.AddAdditionalCostRequest{
			ActivityId: newActivity.Id,
			Cost:       *request.AdditionalCost,
			Title:      "ค่าใช้จ่ายเพิ่มเติม (ขาย)",
		}
		_, err := sv.AdditionalCostService.WithTrx(tx).Create(addAdditionalCost, userIdentity)
		if err != nil {
			tx.Rollback()
			return nil, nil, nil, err
		}
	}

	// commit transaction
	tx.Commit()

	return newActivity, currentActivePond, sellDetails, nil
}
