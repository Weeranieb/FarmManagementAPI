package service

import (
	"context"

	"github.com/samber/lo"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

type ActivityService interface {
	// ListFeed returns the newest discrete events (fill / move / sell) across
	// every pond the caller's client owns, ordered by activity date desc.
	// limit <= 0 returns the full history. Pass `before` to continue after the
	// last row of a previous page.
	ListFeed(ctx context.Context, limit int, before *repository.ActivityFeedCursor) ([]*dto.ActivityFeedItem, error)
	// ListSellDetails returns one sale's size-grade lines (smallest grade
	// first). Empty when the activity is not a sell, does not exist, or is not
	// the caller's — the repository query scopes by client, so a guessed id
	// cannot read another client's sale.
	ListSellDetails(ctx context.Context, activityId int) ([]*dto.SellDetailLine, error)
}

type activityService struct {
	activityRepo repository.ActivityRepository
}

func NewActivityService(activityRepo repository.ActivityRepository) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
	}
}

func (s *activityService) ListFeed(ctx context.Context, limit int, before *repository.ActivityFeedCursor) ([]*dto.ActivityFeedItem, error) {
	// Super admin (clientId == nil) sees all clients; regular users are
	// scoped to their own client. Missing clientId on a non-admin token is
	// an auth problem, not an empty feed.
	clientId, canAccess := utils.GetClientIdForAccess(ctx)
	if !canAccess {
		return nil, errors.ErrAuthPermissionDenied
	}

	rows, err := s.activityRepo.ListRecentByClientID(ctx, clientId, limit, before)
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

	out := make([]*dto.ActivityFeedItem, 0, len(rows))
	for _, r := range rows {
		pricePerUnit, _ := r.PricePerUnit.Float64()
		fishWeight, _ := r.FishWeight.Float64()
		amount := r.Amount

		var total float64
		var totalWeight *float64
		switch r.Mode {
		case constants.ActivityModeSell:
			st := sellTotals[r.Id]
			total = st.total
			// A sell row's own amount/price/weight columns stay empty — every
			// figure lives on its detail lines. Backfill from the summed details
			// and derive average ฿/kg here (guarding 0 kg), the same way the
			// per-pond timeline does, so one sale reads identically in both.
			amount = st.fishCount
			totalWeight = lo.ToPtr(st.weight)
			if st.weight > 0 {
				pricePerUnit = st.total / st.weight
			}
		default:
			total = fillMoveActivityTotal(r.Amount, fishWeight, pricePerUnit, additionalCostTotals[r.Id])
		}

		out = append(out, &dto.ActivityFeedItem{
			Id:            r.Id,
			Mode:          r.Mode,
			ActivityDate:  r.ActivityDate,
			CreatedAt:     r.CreatedAt,
			CreatedBy:     r.CreatedBy,
			CreatedByName: r.CreatedByName,
			PondId:        r.PondId,
			PondName:      r.PondName,
			FarmName:      r.FarmName,
			ToPondId:      r.ToPondId,
			ToPondName:    r.ToPondName,
			FishType:      r.FishType,
			Amount:        amount,
			FishWeight:    fishWeight,
			FishUnit:      r.FishUnit,
			PricePerUnit:  pricePerUnit,
			Total:         total,
			TotalWeight:   totalWeight,
			Merchant:      r.MerchantName,
		})
	}
	return out, nil
}

func (s *activityService) ListSellDetails(ctx context.Context, activityId int) ([]*dto.SellDetailLine, error) {
	clientId, canAccess := utils.GetClientIdForAccess(ctx)
	if !canAccess {
		return nil, errors.ErrAuthPermissionDenied
	}

	rows, err := s.activityRepo.ListSellDetailsByActivityID(ctx, activityId, clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	out := make([]*dto.SellDetailLine, 0, len(rows))
	for _, r := range rows {
		weight, _ := r.Weight.Float64()
		pricePerUnit, _ := r.PricePerUnit.Float64()
		// Multiply in decimal, then convert once — going through float64 first
		// would let rounding drift away from the summed headline total, which
		// SumSellDetailsByActivityIDs computes the same way in SQL.
		total, _ := r.Weight.Mul(r.PricePerUnit).Float64()
		out = append(out, &dto.SellDetailLine{
			FishSizeGradeId: r.FishSizeGradeId,
			SizeName:        r.SizeName,
			FishCount:       r.FishCount,
			Weight:          weight,
			PricePerUnit:    pricePerUnit,
			Total:           total,
		})
	}
	return out, nil
}
