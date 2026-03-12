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
	SellPond(ctx context.Context, pondId int, request dto.PondSellRequest, username string) (*dto.PondSellResponse, error)
}

type PondServiceParams struct {
	dig.In

	PondRepo           repository.PondRepository
	FarmRepo           repository.FarmRepository
	ActivePondRepo     repository.ActivePondRepository
	ActivityRepo       repository.ActivityRepository
	AdditionalCostRepo repository.AdditionalCostRepository
	SellDetailRepo     repository.SellDetailRepository
	MerchantRepo       repository.MerchantRepository
	TxManager          transaction.Manager
}

type pondService struct {
	pondRepo           repository.PondRepository
	farmRepo           repository.FarmRepository
	activePondRepo     repository.ActivePondRepository
	activityRepo       repository.ActivityRepository
	additionalCostRepo repository.AdditionalCostRepository
	sellDetailRepo     repository.SellDetailRepository
	merchantRepo       repository.MerchantRepository
	txManager          transaction.Manager
}

func NewPondService(params PondServiceParams) PondService {
	return &pondService{
		pondRepo:           params.PondRepo,
		farmRepo:           params.FarmRepo,
		activePondRepo:     params.ActivePondRepo,
		activityRepo:       params.ActivityRepo,
		additionalCostRepo: params.AdditionalCostRepo,
		sellDetailRepo:     params.SellDetailRepo,
		merchantRepo:       params.MerchantRepo,
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
	// Calculate
	fillCost := utils.CalculateFillCost(request.Amount, request.PricePerUnit, request.AdditionalCosts)

	var resp *dto.PondFillResponse
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		pondRepo := s.pondRepo.WithTx(tx)
		activePondRepo := s.activePondRepo.WithTx(tx)

		var newTotalCost, newNetResult decimal.Decimal
		var newTotalFish int
		var newFishTypes []string
		if activePond != nil {
			newTotalCost = activePond.TotalCost.Add(fillCost)
			newNetResult = activePond.TotalProfit.Sub(newTotalCost)
			newTotalFish = activePond.TotalFish + request.Amount
			newFishTypes = utils.AppendStringIfMissing(activePond.FishTypes, request.FishType)
		}

		// Mapping
		var newActivePond *model.ActivePond
		if activePond == nil {
			newActivePond = &model.ActivePond{
				PondId:    pondId,
				StartDate: activityDate,
				IsActive:  true,
				TotalCost: fillCost,
				NetResult: decimal.Zero.Sub(fillCost),
				TotalFish: request.Amount,
				FishTypes: []string{request.FishType},
			}
		}

		// Save
		if activePond == nil {
			if err := activePondRepo.Create(ctx, newActivePond); err != nil {
				return err
			}
			activePond = newActivePond
			if pond.Status == constants.FarmStatusMaintenance {
				pond.Status = constants.FarmStatusActive
				if err := pondRepo.Update(ctx, pond); err != nil {
					return err
				}
			}
		} else {
			activePond.TotalCost = newTotalCost
			activePond.NetResult = newNetResult
			activePond.TotalFish = newTotalFish
			activePond.FishTypes = newFishTypes
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
	if data.Pond.Status == constants.FarmStatusMaintenance {
		return errors.ErrPondInMaintenance
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

	sourceActive := sourceData.ActivePond
	destPond := destData.Pond
	destActive := destData.ActivePond

	// Calculate: price part = total fish weight * price per kg; split additional cost 50/50
	fishCost, additionalCost := utils.CalculateMoveCost(request.Amount, request.PricePerUnit, request.FishWeight, request.AdditionalCosts)
	halfAdditional := additionalCost.Div(decimal.NewFromInt(2))
	destMoveCost := fishCost.Add(halfAdditional)

	var destTotalCost, destNetResult decimal.Decimal
	var destTotalFish int
	var destFishTypes []string
	if destActive != nil {
		destTotalCost = destActive.TotalCost.Add(destMoveCost)
		destNetResult = destActive.TotalProfit.Sub(destTotalCost)
		destTotalFish = destActive.TotalFish + request.Amount
		destFishTypes = utils.AppendStringIfMissing(destActive.FishTypes, request.FishType)
	}

	sourceTotalCost := sourceActive.TotalCost.Add(halfAdditional)
	sourceTotalProfit := sourceActive.TotalProfit.Add(fishCost)
	sourceNetResult := sourceTotalProfit.Sub(sourceTotalCost)
	sourceTotalFish := max(sourceActive.TotalFish-request.Amount, 0)

	var resp *dto.PondMoveResponse
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		pondRepo := s.pondRepo.WithTx(tx)
		activePondRepo := s.activePondRepo.WithTx(tx)

		// Mapping
		var newDestActive *model.ActivePond
		if destActive == nil {
			newDestActive = &model.ActivePond{
				PondId:      request.ToPondId,
				StartDate:   activityDate,
				IsActive:    true,
				TotalCost:   destMoveCost,
				TotalProfit: decimal.Zero,
				NetResult:   decimal.Zero.Sub(destMoveCost),
				TotalFish:   request.Amount,
				FishTypes:   []string{request.FishType},
			}
		}

		// Save
		if destActive == nil {
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
			destActive.TotalCost = destTotalCost
			destActive.NetResult = destNetResult
			destActive.TotalFish = destTotalFish
			destActive.FishTypes = destFishTypes
			if err := activePondRepo.Update(ctx, destActive); err != nil {
				return err
			}
		}

		sourceActive.TotalCost = sourceTotalCost
		sourceActive.TotalProfit = sourceTotalProfit
		sourceActive.NetResult = sourceNetResult
		sourceActive.TotalFish = sourceTotalFish
		if request.MarkToClose {
			sourceActive.IsActive = false
			sourceActive.EndDate = &activityDate
		}
		if err := activePondRepo.Update(ctx, sourceActive); err != nil {
			return err
		}
		if request.MarkToClose {
			sourcePond := sourceData.Pond
			sourcePond.Status = constants.FarmStatusMaintenance
			if err := pondRepo.Update(ctx, sourcePond); err != nil {
				return err
			}
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

// validatePondForSell ensures data has pond, active cycle, and client (for sell flow).
func (s *pondService) validatePondForSell(data *repository.PondWithFarmAndActivePond) error {
	if data == nil || data.Pond == nil {
		return errors.ErrPondNotFound
	}
	if data.Pond.Status == constants.FarmStatusMaintenance {
		return errors.ErrPondInMaintenance
	}
	if data.ActivePond == nil {
		return errors.ErrPondNotActive
	}
	if data.ClientId == 0 {
		return errors.ErrFarmNotFound
	}
	return nil
}

// validateSellMerchantIfSet checks that merchantId exists when provided.
func (s *pondService) validateSellMerchantIfSet(merchantId *int) error {
	if merchantId == nil {
		return nil
	}
	merchant, err := s.merchantRepo.GetByID(*merchantId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if merchant == nil {
		return errors.ErrMerchantNotFound
	}
	return nil
}

func buildSellDetailModels(activityId int, details []dto.PondSellDetailItem) []*model.SellDetail {
	out := make([]*model.SellDetail, 0, len(details))
	for _, d := range details {
		out = append(out, &model.SellDetail{
			SellId:       activityId,
			FishType:     d.FishType,
			Size:         d.Size,
			Amount:       d.Amount,
			FishUnit:     d.FishUnit,
			PricePerUnit: d.PricePerUnit,
		})
	}
	return out
}

func (s *pondService) SellPond(ctx context.Context, pondId int, request dto.PondSellRequest, username string) (*dto.PondSellResponse, error) {
	data, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if err := s.validatePondForSell(data); err != nil {
		return nil, err
	}
	ok, err := utils.CanAccessClient(ctx, data.ClientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, errors.ErrAuthPermissionDenied
	}
	if err := s.validateSellMerchantIfSet(request.MerchantId); err != nil {
		return nil, err
	}
	activityDate, err := time.Parse("2006-01-02", request.ActivityDate)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}

	activePond := data.ActivePond
	pond := data.Pond

	var resp *dto.PondSellResponse
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		resp, err = s.executeSellTransaction(ctx, tx, activePond, pond, request, activityDate)
		return err
	})
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	return resp, nil
}

// executeSellTransaction creates sell activity + details, updates active pond, optionally closes pond.
func (s *pondService) executeSellTransaction(
	ctx context.Context,
	tx *gorm.DB,
	activePond *model.ActivePond,
	pond *model.Pond,
	request dto.PondSellRequest,
	activityDate time.Time,
) (*dto.PondSellResponse, error) {
	sellDetailRepo := s.sellDetailRepo.WithTx(tx)
	activePondRepo := s.activePondRepo.WithTx(tx)
	pondRepo := s.pondRepo.WithTx(tx)

	// Calculate
	revenue, additionalCostTotal := utils.CalculateSellTotals(request.Details, request.AdditionalCosts)
	newTotalCost := activePond.TotalCost
	if len(request.AdditionalCosts) > 0 {
		newTotalCost = newTotalCost.Add(additionalCostTotal)
	}
	newTotalProfit := activePond.TotalProfit.Add(revenue)
	newNetResult := newTotalProfit.Sub(newTotalCost)

	// Mapping
	activity := &model.Activity{
		ActivePondId: activePond.Id,
		Mode:         constants.ActivityModeSell,
		MerchantId:   request.MerchantId,
		ActivityDate: activityDate,
	}

	// Save
	if err := s.createActivityWithAdditionalCosts(ctx, tx, activity, request.AdditionalCosts); err != nil {
		return nil, err
	}
	sellDetails := buildSellDetailModels(activity.Id, request.Details)
	if err := sellDetailRepo.CreateBatch(ctx, sellDetails); err != nil {
		return nil, err
	}
	activePond.TotalCost = newTotalCost
	activePond.TotalProfit = newTotalProfit
	activePond.NetResult = newNetResult
	if request.MarkToClose {
		activePond.IsActive = false
		activePond.EndDate = &activityDate
	}
	if err := activePondRepo.Update(ctx, activePond); err != nil {
		return nil, err
	}
	if request.MarkToClose {
		pond.Status = constants.FarmStatusMaintenance
		if err := pondRepo.Update(ctx, pond); err != nil {
			return nil, err
		}
	}
	return &dto.PondSellResponse{
		ActivityId:   int64(activity.Id),
		ActivePondId: int64(activePond.Id),
	}, nil
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
			resp.StartDate = &ap.StartDate
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
