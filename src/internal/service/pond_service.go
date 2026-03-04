package service

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/transaction"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=PondService --output=./mocks --outpkg=service --filename=pond_service.go --structname=MockPondService --with-expecter=false
type PondService interface {
	CreatePonds(ctx context.Context, request dto.CreatePondsRequest, username string) error
	Get(id int) (*dto.PondResponse, error)
	Update(ctx context.Context, request dto.UpdatePondRequest, username string) error
	GetList(farmId int) ([]*dto.PondResponse, error)
	Delete(id int, username string) error
	FillPond(ctx context.Context, pondId int, request dto.PondFillRequest, username string) (*dto.PondFillResponse, error)
	MovePond(ctx context.Context, sourcePondId int, request dto.PondMoveRequest, username string) (*dto.PondMoveResponse, error)
}

type PondServiceParams struct {
	dig.In

	PondRepo           repository.PondRepository
	FarmRepo           repository.FarmRepository
	ActivePondRepo     repository.ActivePondRepository
	ActivityRepo       repository.ActivityRepository
	AdditionalCostRepo repository.AdditionalCostRepository
	TxManager          transaction.Manager
}

type pondService struct {
	pondRepo           repository.PondRepository
	farmRepo           repository.FarmRepository
	activePondRepo     repository.ActivePondRepository
	activityRepo       repository.ActivityRepository
	additionalCostRepo repository.AdditionalCostRepository
	txManager          transaction.Manager
}

func NewPondService(params PondServiceParams) PondService {
	return &pondService{
		pondRepo:           params.PondRepo,
		farmRepo:           params.FarmRepo,
		activePondRepo:     params.ActivePondRepo,
		activityRepo:       params.ActivityRepo,
		additionalCostRepo: params.AdditionalCostRepo,
		txManager:          params.TxManager,
	}
}

func (s *pondService) CreatePonds(ctx context.Context, request dto.CreatePondsRequest, username string) error {
	normalizedNames := make([]string, 0, len(request.Names))
	for _, name := range request.Names {
		normalizedNames = append(normalizedNames, utils.NormalizePondNameForStore(name))
	}
	for _, name := range normalizedNames {
		checkPond, err := s.pondRepo.GetByFarmIdAndName(request.FarmId, name)
		if err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		if checkPond != nil {
			return errors.ErrPondAlreadyExists
		}
	}

	newPonds := make([]*model.Pond, 0, len(normalizedNames))
	for _, name := range normalizedNames {
		newPonds = append(newPonds, &model.Pond{
			FarmId: request.FarmId,
			Name:   name,
			Status: constants.FarmStatusMaintenance,
		})
	}

	// CreatedBy/UpdatedBy set via BaseModel hook from ctx
	return s.pondRepo.CreateBatch(ctx, newPonds)
}

func (s *pondService) Get(id int) (*dto.PondResponse, error) {
	pa, err := s.pondRepo.GetByIDWithFarmAndActivePond(context.Background(), id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if pa == nil {
		return nil, errors.ErrPondNotFound
	}
	return s.toPondResponseFromPondWithActive(pa), nil
}

func (s *pondService) Update(ctx context.Context, req dto.UpdatePondRequest, username string) error {
	existing, err := s.pondRepo.GetByID(req.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existing == nil {
		return errors.ErrPondNotFound
	}

	// Apply only provided fields (non-zero / non-empty so partial update is safe)
	if req.FarmId != 0 {
		existing.FarmId = req.FarmId
	}
	if req.Name != "" {
		existing.Name = utils.NormalizePondNameForStore(req.Name)
	}
	if req.Status != "" {
		existing.Status = req.Status
	}

	// Enforce unique pond name per farm when name was updated
	if req.Name != "" {
		dup, err := s.pondRepo.GetByFarmIdAndName(existing.FarmId, existing.Name)
		if err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		if dup != nil && dup.Id != existing.Id {
			return errors.ErrPondAlreadyExists
		}
	}

	// UpdatedBy set via BaseModel hook from ctx
	if err := s.pondRepo.Update(ctx, existing); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *pondService) GetList(farmId int) ([]*dto.PondResponse, error) {
	list, err := s.pondRepo.ListByFarmIdWithActivePond(context.Background(), farmId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	responses := make([]*dto.PondResponse, 0, len(list))
	for _, pa := range list {
		responses = append(responses, s.toPondResponseFromPondWithActive(pa))
	}
	return responses, nil
}

func (s *pondService) Delete(id int, username string) error {
	// Delete pond (soft delete)
	if err := s.pondRepo.Delete(id); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *pondService) FillPond(ctx context.Context, pondId int, request dto.PondFillRequest, username string) (*dto.PondFillResponse, error) {
	data, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if data == nil || data.Pond == nil {
		return nil, errors.ErrPondNotFound
	}
	pond := data.Pond
	if data.ClientId == 0 {
		return nil, errors.ErrFarmNotFound
	}
	ok, err := utils.CanAccessClient(ctx, data.ClientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, errors.ErrAuthPermissionDenied
	}

	if !constants.IsValidFishType(request.FishType) {
		return nil, errors.ErrInvalidFishType
	}

	activityDate, err := time.Parse("2006-01-02", request.ActivityDate)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}

	activePond := data.ActivePond
	fillCost := utils.CalculateFillCost(request.Amount, request.PricePerUnit, request.AdditionalCosts)

	var resp *dto.PondFillResponse
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		pondRepo := s.pondRepo.WithTx(tx)
		activePondRepo := s.activePondRepo.WithTx(tx)

		if activePond == nil {
			activePond = &model.ActivePond{
				PondId:    pondId,
				StartDate: activityDate,
				IsActive:  true,
				TotalCost: fillCost,
				NetResult: decimal.Zero.Sub(fillCost),
				TotalFish: request.Amount,
				FishTypes: []string{request.FishType},
			}
			if err := activePondRepo.Create(ctx, activePond); err != nil {
				return err
			}
			if pond.Status == constants.FarmStatusMaintenance {
				pond.Status = constants.FarmStatusActive
				if err := pondRepo.Update(ctx, pond); err != nil {
					return err
				}
			}
		} else {
			activePond.TotalCost = activePond.TotalCost.Add(fillCost)
			activePond.NetResult = activePond.TotalProfit.Sub(activePond.TotalCost)
			activePond.TotalFish += request.Amount
			activePond.FishTypes = utils.AppendStringIfMissing(activePond.FishTypes, request.FishType)
			if err := activePondRepo.Update(ctx, activePond); err != nil {
				return err
			}
		}

		activity := &model.Activity{
			ActivePondId: activePond.Id,
			Mode:         constants.ActivityModeFill,
			Amount:       request.Amount,
			FishType:     request.FishType,
			FishWeight:   request.FishWeight,
			FishUnit:     constants.FishUnitKg,
			PricePerUnit: request.PricePerUnit,
			ActivityDate: activityDate,
		}
		if err := s.createActivityWithAdditionalCosts(ctx, tx, activity, request.AdditionalCosts); err != nil {
			return err
		}

		resp = &dto.PondFillResponse{
			ActivityId:   int64(activity.Id),
			ActivePondId: int64(activePond.Id),
		}
		return nil
	})
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	return resp, nil
}

// validatePondWithFarmAndActivePondSource validates that data from GetByIDWithFarmAndActivePond
// represents a valid source pond (has pond, active pond, and farm/client).
func (s *pondService) validatePondWithFarmAndActivePondSource(data *repository.PondWithFarmAndActivePond) error {
	if data == nil || data.Pond == nil {
		return errors.ErrPondNotFound
	}
	if data.ActivePond == nil {
		return errors.ErrPondSourceNotActive
	}
	if data.ClientId == 0 {
		return errors.ErrFarmNotFound
	}
	return nil
}

// validatePondWithFarmAndActivePondDest validates that data from GetByIDWithFarmAndActivePond
// represents a valid destination pond (has pond) and belongs to the same client as the source.
func (s *pondService) validatePondWithFarmAndActivePondDest(data *repository.PondWithFarmAndActivePond, expectedClientId int) error {
	if data == nil || data.Pond == nil {
		return errors.ErrPondNotFound
	}
	if data.ClientId != expectedClientId {
		return errors.ErrAuthPermissionDenied
	}
	return nil
}

func (s *pondService) MovePond(ctx context.Context, sourcePondId int, request dto.PondMoveRequest, username string) (*dto.PondMoveResponse, error) {
	sourceData, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, sourcePondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if err := s.validatePondWithFarmAndActivePondSource(sourceData); err != nil {
		return nil, err
	}
	ok, err := utils.CanAccessClient(ctx, sourceData.ClientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, errors.ErrAuthPermissionDenied
	}

	if sourcePondId == request.ToPondId {
		return nil, errors.ErrPondInvalidInput
	}

	destData, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, request.ToPondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if err := s.validatePondWithFarmAndActivePondDest(destData, sourceData.ClientId); err != nil {
		return nil, err
	}

	if !constants.IsValidFishType(request.FishType) {
		return nil, errors.ErrInvalidFishType
	}
	activityDate, err := time.Parse("2006-01-02", request.ActivityDate)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}

	// Calculate price part = total fish weight * price per kg
	fishCost, additionalCost := utils.CalculateMoveCost(request.Amount, request.PricePerUnit, request.FishWeight, request.AdditionalCosts)
	halfAdditional := additionalCost.Div(decimal.NewFromInt(2))
	destMoveCost := fishCost.Add(halfAdditional)

	sourceActive := sourceData.ActivePond
	destPond := destData.Pond
	destActive := destData.ActivePond

	var resp *dto.PondMoveResponse
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		pondRepo := s.pondRepo.WithTx(tx)
		activePondRepo := s.activePondRepo.WithTx(tx)

		if destActive == nil {
			newDestActive := &model.ActivePond{
				PondId:      request.ToPondId,
				StartDate:   activityDate,
				IsActive:    true,
				TotalCost:   destMoveCost,
				TotalProfit: decimal.Zero,
				NetResult:   decimal.Zero.Sub(destMoveCost),
				TotalFish:   request.Amount,
				FishTypes:   []string{request.FishType},
			}
			if err := activePondRepo.Create(ctx, newDestActive); err != nil {
				return err
			}
			destActive = newDestActive
			if destPond.Status == constants.FarmStatusMaintenance {
				destPond.Status = constants.FarmStatusActive
				if err := pondRepo.Update(ctx, destPond); err != nil {
					return err
				}
			}
		} else {
			destActive.TotalCost = destActive.TotalCost.Add(destMoveCost)
			destActive.NetResult = destActive.TotalProfit.Sub(destActive.TotalCost)
			destActive.TotalFish += request.Amount
			destActive.FishTypes = utils.AppendStringIfMissing(destActive.FishTypes, request.FishType)
			if err := activePondRepo.Update(ctx, destActive); err != nil {
				return err
			}
		}

		sourceActive.TotalCost = sourceActive.TotalCost.Add(halfAdditional)
		sourceActive.TotalProfit = sourceActive.TotalProfit.Add(fishCost)
		sourceActive.NetResult = sourceActive.TotalProfit.Sub(sourceActive.TotalCost)
		sourceActive.TotalFish -= request.Amount
		sourceActive.IsActive = !request.IsClose // if true then the pond is close
		if sourceActive.TotalFish < 0 {
			sourceActive.TotalFish = 0
		}
		if err := activePondRepo.Update(ctx, sourceActive); err != nil {
			return err
		}

		toActivePondId := destActive.Id
		activity := &model.Activity{
			ActivePondId:   sourceActive.Id,
			ToActivePondId: &toActivePondId,
			Mode:           constants.ActivityModeMove,
			Amount:         request.Amount,
			FishType:       request.FishType,
			FishWeight:     request.FishWeight,
			FishUnit:       constants.FishUnitKg,
			PricePerUnit:   request.PricePerUnit,
			ActivityDate:   activityDate,
		}
		if err := s.createActivityWithAdditionalCosts(ctx, tx, activity, request.AdditionalCosts); err != nil {
			return err
		}

		resp = &dto.PondMoveResponse{
			ActivityId:     int64(activity.Id),
			ActivePondId:   int64(sourceActive.Id),
			ToActivePondId: int64(destActive.Id),
		}
		return nil
	})
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	return resp, nil
}

// createActivityWithAdditionalCosts creates the activity and then each additional cost linked to it.
func (s *pondService) createActivityWithAdditionalCosts(
	ctx context.Context,
	tx *gorm.DB,
	activity *model.Activity,
	additionalCosts []dto.AdditionalCostItem,
) error {
	if err := s.activityRepo.WithTx(tx).Create(ctx, activity); err != nil {
		return err
	}
	if len(additionalCosts) > 0 {
		items := make([]*model.AdditionalCost, 0, len(additionalCosts))
		for _, item := range additionalCosts {
			items = append(items, &model.AdditionalCost{
				ActivityId: activity.Id,
				Title:      item.Title,
				Cost:       item.Cost,
			})
		}
		if err := s.additionalCostRepo.WithTx(tx).CreateBatch(ctx, items); err != nil {
			return err
		}
	}
	return nil
}

func (s *pondService) toPondResponseFromPondWithActive(pa *repository.PondWithFarmAndActivePond) *dto.PondResponse {
	if pa == nil || pa.Pond == nil {
		return nil
	}
	pond := pa.Pond
	resp := &dto.PondResponse{
		Id:        pond.Id,
		FarmId:    pond.FarmId,
		Name:      pond.Name,
		Status:    pond.Status,
		CreatedAt: pond.CreatedAt,
		CreatedBy: pond.CreatedBy,
		UpdatedAt: pond.UpdatedAt,
		UpdatedBy: pond.UpdatedBy,
	}
	if pa.ActivePond != nil {
		ap := pa.ActivePond
		totalFish := ap.TotalFish
		resp.TotalFish = &totalFish
		resp.FishTypes = ap.FishTypes
		if !ap.StartDate.IsZero() {
			// Start date = day 1; each full day after adds 1
			daysSince := int(time.Since(ap.StartDate).Hours() / 24)
			ageDays := daysSince + 1
			if ageDays < 1 {
				ageDays = 0
			}
			resp.AgeDays = &ageDays
		}
	}
	resp.LatestActivityDate = pa.LatestActivityDate
	resp.LatestActivityType = pa.LatestActivityType
	return resp
}
