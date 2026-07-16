package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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

const bulkImportMaxPonds = 5000

//go:generate go run github.com/vektra/mockery/v2@latest --name=PondService --output=./mocks --outpkg=service --filename=pond_service.go --structname=MockPondService --with-expecter=false
type PondService interface {
	CreatePonds(ctx context.Context, request dto.CreatePondsRequest) error
	Get(ctx context.Context, id int) (*dto.PondResponse, error)
	Update(ctx context.Context, request dto.UpdatePondRequest) error
	GetList(ctx context.Context, farmId int) ([]*dto.PondResponse, error)
	Delete(ctx context.Context, id int) error
	ListActivities(ctx context.Context, pondId int) ([]*dto.ActivityResponse, error)
	ListCycles(ctx context.Context, pondId int) ([]*dto.PondCycleResponse, error)
	FillPond(ctx context.Context, pondId int, request dto.PondFillRequest, username string) (*dto.PondFillResponse, error)
	MovePond(ctx context.Context, sourcePondId int, request dto.PondMoveRequest, username string) (*dto.PondMoveResponse, error)
	SellPond(ctx context.Context, pondId int, request dto.PondSellRequest, username string) (*dto.PondSellResponse, error)
	PreviewFillPond(ctx context.Context, pondId int, request dto.PondFillRequest) (*dto.PondFillPreviewResponse, error)
	PreviewMovePond(ctx context.Context, sourcePondId int, request dto.PondMoveRequest) (*dto.PondMovePreviewResponse, error)
	PreviewSellPond(ctx context.Context, pondId int, request dto.PondSellRequest) (*dto.PondSellPreviewResponse, error)
	CalcFillPond(ctx context.Context, request dto.PondFillCalcRequest) *dto.PondFillCalcResponse
	CalcMovePond(ctx context.Context, request dto.PondMoveCalcRequest) *dto.PondMoveCalcResponse
	CalcSellPond(ctx context.Context, request dto.PondSellCalcRequest) *dto.PondSellCalcResponse
	BulkImportFarmPond(ctx context.Context, clientId int, request dto.BulkImportFarmPondRequest) (*dto.BulkImportFarmPondResponse, error)
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
	FishSizeGradeRepo  repository.FishSizeGradeRepository
	FeedCostCalc       FeedCostCalculator
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
	fishSizeGradeRepo  repository.FishSizeGradeRepository
	feedCostCalc       FeedCostCalculator
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
		fishSizeGradeRepo:  params.FishSizeGradeRepo,
		feedCostCalc:       params.FeedCostCalc,
		txManager:          params.TxManager,
	}
}

// syncFarmStatusFromPonds updates farms.status from current ponds using pondRepo.WithTx(tx) and
// farmRepo.WithTx(tx). tx must be the active GORM transaction from txManager.WithTransaction.
func (s *pondService) syncFarmStatusFromPonds(ctx context.Context, tx *gorm.DB, farmId int) error {
	if tx == nil {
		return errors.ErrGeneric.Wrap(fmt.Errorf("syncFarmStatusFromPonds: transaction required"))
	}
	if farmId == 0 {
		return nil
	}
	pondRepo := s.pondRepo.WithTx(tx)
	farmRepo := s.farmRepo.WithTx(tx)
	ponds, err := pondRepo.ListByFarmId(farmId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	desired := utils.DeriveFarmStatusFromPonds(ponds)
	farm, err := farmRepo.GetByID(farmId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if farm == nil {
		return nil
	}
	if farm.Status == desired {
		return nil
	}
	farm.Status = desired
	return farmRepo.Update(ctx, farm)
}

func (s *pondService) CreatePonds(ctx context.Context, request dto.CreatePondsRequest) error {
	// Verify the farm exists and the caller can access its owning client.
	// The handler only checks admin-level; per-client scoping lives here so a
	// client admin can't create ponds in another client's farm by passing a
	// foreign farmId.
	farm, err := s.farmRepo.GetByID(request.FarmId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if farm == nil {
		return errors.ErrFarmNotFound
	}
	ok, err := utils.CanAccessClient(ctx, farm.ClientId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return errors.ErrAuthPermissionDenied
	}

	newPonds := make([]*model.Pond, 0, len(request.Ponds))
	for _, item := range request.Ponds {
		newPonds = append(newPonds, &model.Pond{
			FarmId: request.FarmId,
			Name:   utils.NormalizePondNameForStore(item.Name),
			Status: constants.FarmStatusMaintenance,
			Area:   utils.NullDecimalFromDecimalPtr(item.Area),
		})
	}

	return s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		pondRepo := s.pondRepo.WithTx(tx)
		for _, pond := range newPonds {
			checkPond, err := pondRepo.GetByFarmIdAndName(request.FarmId, pond.Name)
			if err != nil {
				return errors.ErrGeneric.Wrap(err)
			}
			if checkPond != nil {
				return errors.ErrPondAlreadyExists
			}
		}
		// CreatedBy/UpdatedBy set via BaseModel hook from ctx
		if err := pondRepo.CreateBatch(ctx, newPonds); err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		return s.syncFarmStatusFromPonds(ctx, tx, request.FarmId)
	})
}

func (s *pondService) Get(ctx context.Context, id int) (*dto.PondResponse, error) {
	pa, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if pa == nil {
		return nil, errors.ErrPondNotFound
	}
	resp := s.toPondResponseFromPondWithActive(pa)
	if pa.ActivePond != nil {
		feedCost, err := s.feedCostCalc.CalcCycleFeedCost(ctx, pa.ActivePond)
		if err != nil {
			// Feed cost is derived from separate tables (daily logs + price
			// history). If that read fails, still return the pond with nil
			// financials rather than failing the whole detail request.
			slog.ErrorContext(ctx, "feed cost calc failed; returning pond without financials",
				"pond_id", id, "error", err)
		} else {
			setPondFinancials(resp, pa.ActivePond, feedCost)
		}
	}
	return resp, nil
}

// cyclePLFloats maps a cycle's decimal P&L to the float64s the response layer
// uses: accumulated cost and revenue, the given feed cost, and the live net
// result (revenue − cost − feed). Single source of the net formula so the
// pond-detail and cycle-list read paths can't diverge.
func cyclePLFloats(ap *model.ActivePond, feedCost decimal.Decimal) (totalCost, totalRevenue, feed, net float64) {
	totalCost, _ = ap.TotalCost.Float64()
	totalRevenue, _ = ap.TotalProfit.Float64()
	feed, _ = feedCost.Float64()
	net, _ = ap.TotalProfit.Sub(ap.TotalCost).Sub(feedCost).Float64()
	return
}

// setPondFinancials fills the cycle P&L fields on resp from the active pond's
// stored transactional figures plus a derived (or snapshotted) feed cost.
// NetResult is computed live as revenue − cost − feed so an active cycle's
// figure always reflects current feed consumption.
func setPondFinancials(resp *dto.PondResponse, ap *model.ActivePond, feedCost decimal.Decimal) {
	totalCost, totalRevenue, feed, net := cyclePLFloats(ap, feedCost)
	resp.TotalCost = &totalCost
	resp.TotalRevenue = &totalRevenue
	resp.FeedCost = &feed
	resp.NetResult = &net
}

// applyCloseFeedSnapshot freezes the cycle's derived feed cost onto ap and folds
// it into the final net result. Shared by the sell-close and move-close paths so
// both close semantics stay identical. Reads ap.TotalProfit/TotalCost, which the
// caller must have already set to their final values.
func (s *pondService) applyCloseFeedSnapshot(ctx context.Context, ap *model.ActivePond) error {
	feedCost, err := s.feedCostCalc.CalcCycleFeedCost(ctx, ap)
	if err != nil {
		return err
	}
	ap.FeedCost = decimal.NullDecimal{Decimal: feedCost, Valid: true}
	ap.NetResult = ap.TotalProfit.Sub(ap.TotalCost).Sub(feedCost)
	return nil
}

// ListActivities returns the fill/move/sell activity timeline for a pond,
// ordered by date desc. Sell rows have their total computed from
// sell_details (sum of weight * price_per_unit) since the parent activity
// row only stores `amount` (an aggregate fish count).
func (s *pondService) ListActivities(ctx context.Context, pondId int) ([]*dto.ActivityResponse, error) {
	pond, err := s.pondRepo.GetByID(pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if pond == nil {
		return nil, errors.ErrPondNotFound
	}

	rows, err := s.activityRepo.ListByPondID(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	modes := make([]activityIDMode, 0, len(rows))
	for _, r := range rows {
		modes = append(modes, activityIDMode{id: r.Id, mode: r.Mode})
	}
	sellIds, fillMoveIds := partitionActivityIDs(modes)
	sellTotals, additionalCostTotals, err := loadActivityCostTotals(ctx, s.activityRepo, sellIds, fillMoveIds)
	if err != nil {
		return nil, err
	}

	out := make([]*dto.ActivityResponse, 0, len(rows))
	for _, r := range rows {
		pricePerUnit, _ := r.PricePerUnit.Float64()
		fishWeight, _ := r.FishWeight.Float64()
		var total float64
		switch r.Mode {
		case constants.ActivityModeSell:
			total = sellTotals[r.Id].total
		default:
			total = fillMoveActivityTotal(r.Amount, fishWeight, pricePerUnit, additionalCostTotals[r.Id])
		}
		direction := "out"
		if r.IsIncoming {
			direction = "in"
		}
		out = append(out, &dto.ActivityResponse{
			Id:           r.Id,
			Mode:         r.Mode,
			Direction:    direction,
			ActivityDate: r.ActivityDate,
			FishType:     r.FishType,
			Amount:       r.Amount,
			PricePerUnit: pricePerUnit,
			Total:        total,
			Merchant:     r.MerchantName,
			ToPondName:   r.ToPondName,
			FromPondName: r.FromPondName,
		})
	}
	return out, nil
}

// ListCycles returns every production cycle of a pond (active + closed), newest
// first, each with its P&L. Feed cost for the active cycle is derived live;
// closed cycles use the value frozen at close (nil for legacy cycles closed
// before feed-cost accounting). Access is scoped to the pond's client since the
// response exposes financial figures.
func (s *pondService) ListCycles(ctx context.Context, pondId int) ([]*dto.PondCycleResponse, error) {
	data, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if data == nil || data.Pond == nil {
		return nil, errors.ErrPondNotFound
	}
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

	cycles, err := s.activePondRepo.ListByPondID(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	// Feed cost is derived only for still-active cycles (closed ones carry a
	// snapshot); batch the active ones so there is no per-cycle query.
	activeCycles := make([]*model.ActivePond, 0, len(cycles))
	for _, c := range cycles {
		if c.IsActive {
			activeCycles = append(activeCycles, c)
		}
	}
	feedByAp, err := s.feedCostCalc.CalcCycleFeedCostBatch(ctx, activeCycles)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	out := make([]*dto.PondCycleResponse, 0, len(cycles))
	for _, c := range cycles {
		out = append(out, buildPondCycleResponse(c, feedByAp))
	}
	return out, nil
}

// buildPondCycleResponse maps one cycle to its P&L response. For an active cycle
// feed cost is taken from the derived batch and net = revenue − cost − feed; for
// a closed cycle the frozen snapshot is used as-is (feed_cost may be nil for
// legacy cycles, in which case net_result is the pre-feed figure it was stored
// with).
func buildPondCycleResponse(c *model.ActivePond, feedByAp map[int]decimal.Decimal) *dto.PondCycleResponse {
	totalCost, _ := c.TotalCost.Float64()
	totalRevenue, _ := c.TotalProfit.Float64()
	resp := &dto.PondCycleResponse{
		Id:           c.Id,
		StartDate:    c.StartDate,
		EndDate:      c.EndDate,
		IsActive:     c.IsActive,
		TotalFish:    c.TotalFish,
		FishTypes:    c.FishTypes,
		TotalCost:    totalCost,
		TotalRevenue: totalRevenue,
	}
	if c.IsActive {
		_, _, feed, net := cyclePLFloats(c, feedByAp[c.Id])
		resp.FeedCost = &feed
		resp.NetResult = net
		return resp
	}
	if c.FeedCost.Valid {
		feed, _ := c.FeedCost.Decimal.Float64()
		resp.FeedCost = &feed
	}
	net, _ := c.NetResult.Float64()
	resp.NetResult = net
	return resp
}

func (s *pondService) Update(ctx context.Context, req dto.UpdatePondRequest) error {
	existing, err := s.pondRepo.GetByID(req.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existing == nil {
		return errors.ErrPondNotFound
	}
	oldFarmId := existing.FarmId

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
	if req.Area != nil {
		existing.Area = utils.NullDecimalFromDecimalPtr(req.Area)
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

	return s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		pondRepo := s.pondRepo.WithTx(tx)
		// UpdatedBy set via BaseModel hook from ctx
		if err := pondRepo.Update(ctx, existing); err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		if err := s.syncFarmStatusFromPonds(ctx, tx, oldFarmId); err != nil {
			return err
		}
		if existing.FarmId != oldFarmId {
			if err := s.syncFarmStatusFromPonds(ctx, tx, existing.FarmId); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *pondService) GetList(ctx context.Context, farmId int) ([]*dto.PondResponse, error) {
	list, err := s.pondRepo.ListByFarmIdWithActivePond(ctx, farmId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	responses := make([]*dto.PondResponse, 0, len(list))
	activePonds := make([]*model.ActivePond, 0, len(list))
	for _, pa := range list {
		responses = append(responses, s.toPondResponseFromPondWithActive(pa))
		if pa != nil && pa.ActivePond != nil {
			activePonds = append(activePonds, pa.ActivePond)
		}
	}
	// Derive feed cost for every active cycle in one batch (no N+1) so the
	// listing shows a net result that already accounts for feed. Degrade
	// gracefully: if the feed-data read fails, return the listing without
	// financials rather than failing the whole farm view.
	feedByAp, err := s.feedCostCalc.CalcCycleFeedCostBatch(ctx, activePonds)
	if err != nil {
		slog.ErrorContext(ctx, "batch feed cost calc failed; listing without financials",
			"farm_id", farmId, "error", err)
		return responses, nil
	}
	for i, pa := range list {
		if pa != nil && pa.ActivePond != nil {
			setPondFinancials(responses[i], pa.ActivePond, feedByAp[pa.ActivePond.Id])
		}
	}
	return responses, nil
}

func (s *pondService) Delete(ctx context.Context, id int) error {
	pond, err := s.pondRepo.GetByID(id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if pond == nil {
		return errors.ErrPondNotFound
	}
	farmId := pond.FarmId
	return s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		pondRepo := s.pondRepo.WithTx(tx)
		if err := pondRepo.Delete(ctx, id); err != nil {
			return errors.ErrGeneric.Wrap(err)
		}
		return s.syncFarmStatusFromPonds(ctx, tx, farmId)
	})
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
	// Calculate: amount × fishWeight × pricePerUnit + additionalCosts (price is per kg)
	fillCost := utils.CalculateFillCost(request.Amount, request.PricePerUnit, request.FishWeight, request.AdditionalCosts)

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
		if err := s.syncFarmStatusFromPonds(ctx, tx, data.Pond.FarmId); err != nil {
			return err
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
	if request.Amount <= 0 {
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
	// Refuse moves that would drain the source past zero. Without this the
	// caller would see TotalFish silently clamped to 0 (see sourceTotalFish
	// below) and the off-by-N would only show up as data drift.
	if request.Amount > sourceActive.TotalFish {
		return nil, errors.ErrPondInsufficientFish
	}
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
			// Closing ends the cycle and empties the pond (same semantics as a
			// sell that closes): any fish not moved out are written off, so a
			// closed cycle never reports residual stock.
			sourceActive.TotalFish = 0
			// Freeze the source cycle's derived feed cost and fold it into the
			// final net result (same close semantics as a sell that closes).
			if err := s.applyCloseFeedSnapshot(ctx, sourceActive); err != nil {
				return err
			}
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
		srcFarmId := sourceData.Pond.FarmId
		dstFarmId := destData.Pond.FarmId
		if err := s.syncFarmStatusFromPonds(ctx, tx, srcFarmId); err != nil {
			return err
		}
		if dstFarmId != srcFarmId {
			if err := s.syncFarmStatusFromPonds(ctx, tx, dstFarmId); err != nil {
				return err
			}
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

// sumSellFishCount totals the head count across sell detail lines. FishCount is
// required at the DTO layer; a nil is treated as 0 defensively.
func sumSellFishCount(details []dto.PondSellDetailItem) int {
	total := 0
	for _, d := range details {
		if d.FishCount != nil {
			total += *d.FishCount
		}
	}
	return total
}

func buildSellDetailModels(activityId int, details []dto.PondSellDetailItem) []*model.SellDetail {
	out := make([]*model.SellDetail, 0, len(details))
	for _, d := range details {
		out = append(out, &model.SellDetail{
			SellId:          activityId,
			FishSizeGradeId: d.FishSizeGradeId,
			Weight:          d.Weight,
			PricePerUnit:    d.PricePerUnit,
			FishCount:       d.FishCount,
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
	if err := s.validateSellGradeIDs(request.Details); err != nil {
		return nil, err
	}
	activityDate, err := time.Parse("2006-01-02", request.ActivityDate)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}

	activePond := data.ActivePond
	pond := data.Pond

	// Stock is tracked by head count; refuse selling more fish than the pond
	// holds. This pre-transaction check gives the common case a clean 422; the
	// authoritative guard is the row-locked re-check inside the transaction
	// (executeSellTransaction), which serializes concurrent sells.
	if sumSellFishCount(request.Details) > activePond.TotalFish {
		return nil, errors.ErrPondInsufficientFish
	}

	var resp *dto.PondSellResponse
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		resp, err = s.executeSellTransaction(ctx, tx, activePond, pond, request, activityDate)
		return err
	})
	if err != nil {
		// Preserve the specific insufficient-fish error (surfaced by the locked
		// re-check) so a concurrent oversell still returns 422, not 500. The tx
		// manager returns the callback error verbatim, so the sentinel is intact.
		if err == errors.ErrPondInsufficientFish {
			return nil, errors.ErrPondInsufficientFish
		}
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

	// Re-read the cycle under a row lock so a concurrent sell/move on the same
	// cycle serializes behind this lock. Adopt the locked row's authoritative
	// figures (head count + running totals) onto the working copy so a
	// concurrent commit isn't lost, then re-check stock — this rejects an
	// oversell the pre-transaction check couldn't see (both sales read the same
	// stale count before either committed).
	locked, err := activePondRepo.GetByIDForUpdate(ctx, activePond.Id)
	if err != nil {
		return nil, err
	}
	if locked == nil {
		return nil, errors.ErrPondNotFound
	}
	activePond.TotalFish = locked.TotalFish
	activePond.TotalCost = locked.TotalCost
	activePond.TotalProfit = locked.TotalProfit
	if sumSellFishCount(request.Details) > activePond.TotalFish {
		return nil, errors.ErrPondInsufficientFish
	}

	// Calculate
	revenue, additionalCostTotal := utils.CalculateSellTotals(request.Details, request.AdditionalCosts)
	newTotalCost := activePond.TotalCost
	if len(request.AdditionalCosts) > 0 {
		newTotalCost = newTotalCost.Add(additionalCostTotal)
	}
	newTotalProfit := activePond.TotalProfit.Add(revenue)
	newNetResult := newTotalProfit.Sub(newTotalCost)

	// Stock is tracked by head count: a sale removes the sold fish; closing
	// empties the pond entirely.
	newTotalFish := max(activePond.TotalFish-sumSellFishCount(request.Details), 0)
	if request.MarkToClose {
		newTotalFish = 0
	}

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
	activePond.TotalFish = newTotalFish
	// On close, freeze feed cost and fold it into net result (active cycles keep
	// net = profit − cost and derive feed cost on read, so feed_cost stays NULL).
	if request.MarkToClose {
		if err := s.applyCloseFeedSnapshot(ctx, activePond); err != nil {
			return nil, err
		}
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
	if err := s.syncFarmStatusFromPonds(ctx, tx, pond.FarmId); err != nil {
		return nil, err
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

// --- Preview (Review & Confirm) methods ---

func buildAdditionalCostLines(costs []dto.AdditionalCostItem) []dto.AdditionalCostLine {
	lines := make([]dto.AdditionalCostLine, 0, len(costs))
	for _, c := range costs {
		f, _ := c.Cost.Float64()
		lines = append(lines, dto.AdditionalCostLine{Title: c.Title, Cost: f})
	}
	return lines
}

func (s *pondService) PreviewFillPond(ctx context.Context, pondId int, request dto.PondFillRequest) (*dto.PondFillPreviewResponse, error) {
	data, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if data == nil || data.Pond == nil {
		return &dto.PondFillPreviewResponse{Valid: false, ValidationError: errors.ErrPondNotFound.Message}, nil
	}
	if data.ClientId == 0 {
		return &dto.PondFillPreviewResponse{Valid: false, ValidationError: errors.ErrFarmNotFound.Message}, nil
	}
	ok, err := utils.CanAccessClient(ctx, data.ClientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, errors.ErrAuthPermissionDenied
	}
	if !constants.IsValidFishType(request.FishType) {
		return &dto.PondFillPreviewResponse{Valid: false, ValidationError: errors.ErrInvalidFishType.Message}, nil
	}

	stockBefore := 0
	if data.ActivePond != nil {
		stockBefore = data.ActivePond.TotalFish
	}

	// Reuse same calculation as FillPond
	fillCost := utils.CalculateFillCost(request.Amount, request.PricePerUnit, request.FishWeight, request.AdditionalCosts)
	additionalTotal := utils.CalculateAdditionalCostsTotal(request.AdditionalCosts)
	baseCost := fillCost.Sub(additionalTotal)
	totalCost, _ := fillCost.Float64()
	baseCostF, _ := baseCost.Float64()
	pricePerUnit, _ := request.PricePerUnit.Float64()
	fishWeight, _ := request.FishWeight.Float64()
	totalWeight := float64(request.Amount) * fishWeight
	additionalLines := buildAdditionalCostLines(request.AdditionalCosts)

	return &dto.PondFillPreviewResponse{
		Valid:           true,
		Species:         request.FishType,
		Quantity:        request.Amount,
		AvgWeightKg:     fishWeight,
		TotalWeight:     totalWeight,
		CostPerUnit:     pricePerUnit,
		BaseStockCost:   baseCostF,
		AdditionalCosts: additionalLines,
		TotalCost:       totalCost,
		StockBefore:     stockBefore,
		StockAfter:      stockBefore + request.Amount,
		StockDelta:      request.Amount,
	}, nil
}

func (s *pondService) PreviewMovePond(ctx context.Context, sourcePondId int, request dto.PondMoveRequest) (*dto.PondMovePreviewResponse, error) {
	sourceData, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, sourcePondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if err := s.validatePondWithFarmAndActivePondSource(sourceData); err != nil {
		return &dto.PondMovePreviewResponse{Valid: false, ValidationError: err.Error()}, nil
	}
	ok, err := utils.CanAccessClient(ctx, sourceData.ClientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, errors.ErrAuthPermissionDenied
	}
	if sourcePondId == request.ToPondId {
		return &dto.PondMovePreviewResponse{Valid: false, ValidationError: errors.ErrPondInvalidInput.Message}, nil
	}
	if !constants.IsValidFishType(request.FishType) {
		return &dto.PondMovePreviewResponse{Valid: false, ValidationError: errors.ErrInvalidFishType.Message}, nil
	}

	stockBefore := sourceData.ActivePond.TotalFish
	if request.Amount <= 0 {
		return &dto.PondMovePreviewResponse{Valid: false, ValidationError: errors.ErrPondInvalidInput.Message}, nil
	}
	if request.Amount > stockBefore {
		return &dto.PondMovePreviewResponse{Valid: false, ValidationError: errors.ErrPondInsufficientFish.Message}, nil
	}

	fishCost, additionalCost := utils.CalculateMoveCost(request.Amount, request.PricePerUnit, request.FishWeight, request.AdditionalCosts)
	halfAdditional := additionalCost.Div(decimal.NewFromInt(2))
	destTotal := fishCost.Add(halfAdditional)
	sourceNet := fishCost.Sub(halfAdditional)

	baseCost, _ := fishCost.Float64()
	additionalTotal, _ := additionalCost.Float64()
	halfAdditionalF, _ := halfAdditional.Float64()
	destTotalF, _ := destTotal.Float64()
	sourceNetF, _ := sourceNet.Float64()
	pricePerUnit, _ := request.PricePerUnit.Float64()
	fishWeight, _ := request.FishWeight.Float64()
	totalWeight := float64(request.Amount) * fishWeight
	additionalLines := buildAdditionalCostLines(request.AdditionalCosts)

	return &dto.PondMovePreviewResponse{
		Valid:                true,
		Species:              request.FishType,
		Quantity:             request.Amount,
		AvgWeightKg:          fishWeight,
		TotalWeight:          totalWeight,
		CostPerUnit:          pricePerUnit,
		BaseTransferCost:     baseCost,
		AdditionalCosts:      additionalLines,
		AdditionalCostsTotal: additionalTotal,
		TotalCost:            baseCost + additionalTotal,
		SourceFishRevenue:    baseCost,
		SourceAdditionalCost: halfAdditionalF,
		SourceNetEffect:      sourceNetF,
		DestFishCost:         baseCost,
		DestAdditionalCost:   halfAdditionalF,
		DestTotalCost:        destTotalF,
		StockBefore:          stockBefore,
		StockAfter:           max(stockBefore-request.Amount, 0),
		StockDelta:           -request.Amount,
	}, nil
}

func (s *pondService) PreviewSellPond(ctx context.Context, pondId int, request dto.PondSellRequest) (*dto.PondSellPreviewResponse, error) {
	data, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if err := s.validatePondForSell(data); err != nil {
		return &dto.PondSellPreviewResponse{Valid: false, ValidationError: err.Error()}, nil
	}
	ok, err := utils.CanAccessClient(ctx, data.ClientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, errors.ErrAuthPermissionDenied
	}
	if err := s.validateSellMerchantIfSet(request.MerchantId); err != nil {
		return &dto.PondSellPreviewResponse{Valid: false, ValidationError: err.Error()}, nil
	}

	gradeMap, err := s.buildGradeNameMap(request.Details)
	if err != nil {
		return &dto.PondSellPreviewResponse{Valid: false, ValidationError: errors.ErrFishSizeGradeNotFound.Message}, nil
	}

	detailLines := utils.CalculateSellDetailLines(request.Details)
	items := make([]dto.PondSellPreviewItem, 0, len(detailLines))
	var totalRevenue, totalWeight float64
	for _, line := range detailLines {
		items = append(items, dto.PondSellPreviewItem{
			FishSizeGradeId:   line.FishSizeGradeId,
			FishSizeGradeName: gradeMap[line.FishSizeGradeId],
			Weight:            line.Weight,
			PricePerKg:        line.PricePerUnit,
			Subtotal:          line.Subtotal,
			FishCount:         line.FishCount,
		})
		totalRevenue += line.Subtotal
		totalWeight += line.Weight
	}

	return &dto.PondSellPreviewResponse{
		Valid:        true,
		Items:        items,
		TotalRevenue: totalRevenue,
		TotalWeight:  totalWeight,
	}, nil
}

// toAdditionalCostItems converts the relaxed calc-request shape into the
// stricter shape consumed by utils.* and buildAdditionalCostLines.
func toAdditionalCostItems(calcs []dto.AdditionalCostCalcItem) []dto.AdditionalCostItem {
	if len(calcs) == 0 {
		return nil
	}
	out := make([]dto.AdditionalCostItem, 0, len(calcs))
	for _, c := range calcs {
		out = append(out, dto.AdditionalCostItem(c))
	}
	return out
}

// CalcFillPond returns live cost/weight totals for the fill form. Pure math,
// no DB access — used by the stock-action modal to display running totals.
func (s *pondService) CalcFillPond(_ context.Context, request dto.PondFillCalcRequest) *dto.PondFillCalcResponse {
	addCosts := toAdditionalCostItems(request.AdditionalCosts)
	fillCost := utils.CalculateFillCost(request.Amount, request.PricePerUnit, request.FishWeight, addCosts)
	additionalTotal := utils.CalculateAdditionalCostsTotal(addCosts)
	baseCost := fillCost.Sub(additionalTotal)
	totalCostF, _ := fillCost.Float64()
	additionalTotalF, _ := additionalTotal.Float64()
	baseCostF, _ := baseCost.Float64()
	pricePerUnit, _ := request.PricePerUnit.Float64()
	fishWeight, _ := request.FishWeight.Float64()
	return &dto.PondFillCalcResponse{
		Quantity:             request.Amount,
		AvgWeightKg:          fishWeight,
		TotalWeight:          float64(request.Amount) * fishWeight,
		CostPerUnit:          pricePerUnit,
		BaseStockCost:        baseCostF,
		AdditionalCosts:      buildAdditionalCostLines(addCosts),
		AdditionalCostsTotal: additionalTotalF,
		TotalCost:            totalCostF,
	}
}

// CalcMovePond returns live cost/weight totals for the move form. Pure math.
// Move cost formula: amount × fishWeight × pricePerUnit + additionalCosts.
// Returns a per-side split so the review UI can show source (treated as a
// sale) and destination (treated as a purchase) impacts separately.
func (s *pondService) CalcMovePond(_ context.Context, request dto.PondMoveCalcRequest) *dto.PondMoveCalcResponse {
	addCosts := toAdditionalCostItems(request.AdditionalCosts)
	fishCost, additionalCost := utils.CalculateMoveCost(request.Amount, request.PricePerUnit, request.FishWeight, addCosts)
	halfAdditional := additionalCost.Div(decimal.NewFromInt(2))
	totalCost := fishCost.Add(additionalCost)
	destTotal := fishCost.Add(halfAdditional)
	sourceNet := fishCost.Sub(halfAdditional)

	fishCostF, _ := fishCost.Float64()
	additionalCostF, _ := additionalCost.Float64()
	halfAdditionalF, _ := halfAdditional.Float64()
	totalCostF, _ := totalCost.Float64()
	destTotalF, _ := destTotal.Float64()
	sourceNetF, _ := sourceNet.Float64()
	pricePerUnit, _ := request.PricePerUnit.Float64()
	fishWeight, _ := request.FishWeight.Float64()
	return &dto.PondMoveCalcResponse{
		Quantity:             request.Amount,
		AvgWeightKg:          fishWeight,
		TotalWeight:          float64(request.Amount) * fishWeight,
		CostPerUnit:          pricePerUnit,
		BaseTransferCost:     fishCostF,
		AdditionalCosts:      buildAdditionalCostLines(addCosts),
		AdditionalCostsTotal: additionalCostF,
		TotalCost:            totalCostF,
		SourceFishRevenue:    fishCostF,
		SourceAdditionalCost: halfAdditionalF,
		SourceNetEffect:      sourceNetF,
		DestFishCost:         fishCostF,
		DestAdditionalCost:   halfAdditionalF,
		DestTotalCost:        destTotalF,
	}
}

// CalcSellPond returns live revenue/weight totals for the sell form. Pure math.
// Skips fish-size-grade lookup so partial in-progress rows can be sent.
func (s *pondService) CalcSellPond(_ context.Context, request dto.PondSellCalcRequest) *dto.PondSellCalcResponse {
	items := make([]dto.PondSellCalcLine, 0, len(request.Details))
	var totalRevenue, totalWeight float64
	for _, d := range request.Details {
		w, _ := d.Weight.Float64()
		ppu, _ := d.PricePerUnit.Float64()
		subtotal := w * ppu
		items = append(items, dto.PondSellCalcLine{
			FishSizeGradeId: d.FishSizeGradeId,
			Weight:          w,
			PricePerKg:      ppu,
			Subtotal:        subtotal,
			FishCount:       d.FishCount,
		})
		totalRevenue += subtotal
		totalWeight += w
	}
	addCosts := toAdditionalCostItems(request.AdditionalCosts)
	additionalTotal := utils.CalculateAdditionalCostsTotal(addCosts)
	additionalTotalF, _ := additionalTotal.Float64()
	return &dto.PondSellCalcResponse{
		Items:                items,
		TotalWeight:          totalWeight,
		TotalRevenue:         totalRevenue,
		AdditionalCosts:      buildAdditionalCostLines(addCosts),
		AdditionalCostsTotal: additionalTotalF,
		NetTotal:             totalRevenue - additionalTotalF,
	}
}

// validateSellGradeIDs checks that all FishSizeGradeId values in the details exist.
func (s *pondService) validateSellGradeIDs(details []dto.PondSellDetailItem) error {
	ids := collectGradeIDs(details)
	grades, err := s.fishSizeGradeRepo.GetByIDs(ids)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if len(grades) != len(ids) {
		return errors.ErrFishSizeGradeNotFound
	}
	return nil
}

// buildGradeNameMap returns a map of gradeId -> gradeName for preview display.
func (s *pondService) buildGradeNameMap(details []dto.PondSellDetailItem) (map[int]string, error) {
	ids := collectGradeIDs(details)
	grades, err := s.fishSizeGradeRepo.GetByIDs(ids)
	if err != nil {
		return nil, err
	}
	if len(grades) != len(ids) {
		return nil, errors.ErrFishSizeGradeNotFound
	}
	m := make(map[int]string, len(grades))
	for _, g := range grades {
		m[g.Id] = g.Name
	}
	return m, nil
}

func collectGradeIDs(details []dto.PondSellDetailItem) []int {
	seen := make(map[int]struct{}, len(details))
	ids := make([]int, 0, len(details))
	for _, d := range details {
		if _, ok := seen[d.FishSizeGradeId]; !ok {
			seen[d.FishSizeGradeId] = struct{}{}
			ids = append(ids, d.FishSizeGradeId)
		}
	}
	return ids
}

// validateBulkImportRequest runs every data-quality check against the parsed
// payload before any DB work. It collects *all* issues found (joined with
// "; ") so the user can correct the whole file in one pass instead of
// hitting the API once per fix.
//
// Validations performed here:
//   - Total pond count (across all farms) does not exceed bulkImportMaxPonds.
//   - Farm name is non-empty after normalize and ≤ 100 chars.
//   - Pond name is non-empty after normalize and ≤ 100 chars.
//   - Area, when provided, is not negative (defense-in-depth against the
//     `decimal_gte0` struct tag — that tag is what `validateAndParse` runs
//     in the handler, but we don't trust callers to go through it).
//   - No duplicate (farmName, pondName) pair, case-insensitive after the
//     normalize step, within the same request.
func (s *pondService) validateBulkImportRequest(request dto.BulkImportFarmPondRequest) error {
	var issues []string
	totalPonds := 0
	// Key: "<lower(farmNormalized)>\x00<lower(pondNormalized)>".
	// Value: 1-based ordinal of where the entry was first seen, so the dup
	// message can point at the original occurrence.
	seenFarmPond := make(map[string]int)

	for fi, f := range request.Farms {
		farmName := utils.NormalizeFarmNameForStore(f.Name)
		if farmName == "" {
			issues = append(issues, fmt.Sprintf("farm #%d: empty name after normalize", fi+1))
			continue
		}
		if len(farmName) > 100 {
			issues = append(issues, fmt.Sprintf("farm %q: name exceeds 100 chars", farmName))
		}
		farmKey := strings.ToLower(farmName)

		for pi, p := range f.Ponds {
			pondName := utils.NormalizePondNameForStore(p.Name)
			if pondName == "" {
				issues = append(issues, fmt.Sprintf("farm %q pond #%d: empty pond name after normalize", farmName, pi+1))
				continue
			}
			if len(pondName) > 100 {
				issues = append(issues, fmt.Sprintf("farm %q pond %q: name exceeds 100 chars", farmName, pondName))
				continue
			}
			if p.Area != nil && p.Area.IsNegative() {
				issues = append(issues, fmt.Sprintf("farm %q pond %q: area must be >= 0", farmName, pondName))
			}

			totalPonds++
			key := farmKey + "\x00" + strings.ToLower(pondName)
			if firstOrdinal, dup := seenFarmPond[key]; dup {
				issues = append(issues, fmt.Sprintf("duplicate pond %q in farm %q (also at item #%d)", pondName, farmName, firstOrdinal))
				continue
			}
			seenFarmPond[key] = pi + 1
		}
	}

	if totalPonds > bulkImportMaxPonds {
		issues = append(issues, fmt.Sprintf("too many ponds: got %d, max %d", totalPonds, bulkImportMaxPonds))
	}

	if len(issues) > 0 {
		return errors.ErrValidationFailed.Wrap(fmt.Errorf("%s", strings.Join(issues, "; ")))
	}
	return nil
}

// BulkImportFarmPond upserts farms and ponds for a client in one transaction.
//   - Missing farms are created (status maintenance; will be re-synced from ponds).
//   - Missing ponds are created (status maintenance).
//   - Existing ponds get their area updated when an area is provided.
//   - Nothing is ever deleted.
func (s *pondService) BulkImportFarmPond(ctx context.Context, clientId int, request dto.BulkImportFarmPondRequest) (*dto.BulkImportFarmPondResponse, error) {
	if err := s.validateBulkImportRequest(request); err != nil {
		return nil, err
	}

	resp := &dto.BulkImportFarmPondResponse{
		Farms: make([]dto.BulkImportFarmResult, 0, len(request.Farms)),
	}

	err := s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		farmRepo := s.farmRepo.WithTx(tx)
		pondRepo := s.pondRepo.WithTx(tx)
		touchedFarmIds := make(map[int]struct{}, len(request.Farms))

		for _, f := range request.Farms {
			farmName := utils.NormalizeFarmNameForStore(f.Name)

			existingFarm, err := farmRepo.GetByNameAndClientId(farmName, clientId)
			if err != nil {
				return errors.ErrGeneric.Wrap(err)
			}

			farmResult := dto.BulkImportFarmResult{Name: farmName}
			var targetFarm *model.Farm
			if existingFarm == nil {
				targetFarm = &model.Farm{
					ClientId: clientId,
					Name:     farmName,
					Status:   constants.FarmStatusMaintenance,
				}
				if err := farmRepo.Create(ctx, targetFarm); err != nil {
					return errors.ErrGeneric.Wrap(err)
				}
				resp.FarmsCreated++
				farmResult.IsNew = true
			} else {
				targetFarm = existingFarm
				resp.FarmsExisting++
			}
			touchedFarmIds[targetFarm.Id] = struct{}{}

			for _, p := range f.Ponds {
				pondName := utils.NormalizePondNameForStore(p.Name)

				existingPond, err := pondRepo.GetByFarmIdAndName(targetFarm.Id, pondName)
				if err != nil {
					return errors.ErrGeneric.Wrap(err)
				}
				if existingPond == nil {
					newPond := &model.Pond{
						FarmId: targetFarm.Id,
						Name:   pondName,
						Status: constants.FarmStatusMaintenance,
						Area:   utils.NullDecimalFromDecimalPtr(p.Area),
					}
					if err := pondRepo.Create(ctx, newPond); err != nil {
						return errors.ErrGeneric.Wrap(err)
					}
					resp.PondsCreated++
					farmResult.PondsCreated++
				} else {
					if p.Area != nil {
						existingPond.Area = utils.NullDecimalFromDecimalPtr(p.Area)
						if err := pondRepo.Update(ctx, existingPond); err != nil {
							return errors.ErrGeneric.Wrap(err)
						}
						resp.PondsUpdated++
						farmResult.PondsUpdated++
					} else {
						// Pond matched by name but no area to apply — no DB write.
						resp.PondsUnchanged++
						farmResult.PondsUnchanged++
					}
				}
			}

			resp.Farms = append(resp.Farms, farmResult)
		}

		for farmId := range touchedFarmIds {
			if err := s.syncFarmStatusFromPonds(ctx, tx, farmId); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
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
		Area:      utils.DecimalPtrFromNullDecimal(pond.Area),
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
