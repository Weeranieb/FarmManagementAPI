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

//go:generate go run github.com/vektra/mockery/v2@latest --name=ActivityService --output=./mocks --outpkg=service --filename=activity_service.go --structname=MockActivityService --with-expecter=false
type ActivityService interface {
	// ListFeed returns the newest discrete events (fill / move / sell) across
	// every pond the caller's client owns, ordered by activity date desc.
	// limit <= 0 returns the full history.
	ListFeed(ctx context.Context, limit int) ([]*dto.ActivityFeedItem, error)
}

type activityService struct {
	activityRepo repository.ActivityRepository
}

func NewActivityService(activityRepo repository.ActivityRepository) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
	}
}

func (s *activityService) ListFeed(ctx context.Context, limit int) ([]*dto.ActivityFeedItem, error) {
	// Super admin (clientId == nil) sees all clients; regular users are
	// scoped to their own client. Missing clientId on a non-admin token is
	// an auth problem, not an empty feed.
	clientId, canAccess := utils.GetClientIdForAccess(ctx)
	if !canAccess {
		return nil, errors.ErrAuthPermissionDenied
	}

	rows, err := s.activityRepo.ListRecentByClientID(ctx, clientId, limit)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	// Same total rules as pondService.ListActivities: sell totals from
	// sell_details, fill/move totals folded with additional_costs.
	sellIds := make([]int, 0)
	fillMoveIds := make([]int, 0)
	for _, r := range rows {
		switch r.Mode {
		case constants.ActivityModeSell:
			sellIds = append(sellIds, r.Id)
		default:
			fillMoveIds = append(fillMoveIds, r.Id)
		}
	}

	type sellTotal struct {
		total  float64
		weight float64
	}
	sellTotals := map[int]sellTotal{}
	if len(sellIds) > 0 {
		totals, err := s.activityRepo.SumSellDetailsByActivityIDs(ctx, sellIds)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		for _, t := range totals {
			total, _ := t.Total.Float64()
			weight, _ := t.TotalWeight.Float64()
			sellTotals[t.SellId] = sellTotal{total: total, weight: weight}
		}
	}

	additionalCostTotals := map[int]float64{}
	if len(fillMoveIds) > 0 {
		totals, err := s.activityRepo.SumAdditionalCostsByActivityIDs(ctx, fillMoveIds)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		for _, t := range totals {
			f, _ := t.Total.Float64()
			additionalCostTotals[t.ActivityId] = f
		}
	}

	out := make([]*dto.ActivityFeedItem, 0, len(rows))
	for _, r := range rows {
		pricePerUnit, _ := r.PricePerUnit.Float64()
		fishWeight, _ := r.FishWeight.Float64()

		var total float64
		var totalWeight *float64
		switch r.Mode {
		case constants.ActivityModeSell:
			st := sellTotals[r.Id]
			total = st.total
			totalWeight = lo.ToPtr(st.weight)
		default:
			// amount × fish_weight × price_per_unit + Σ additional_costs —
			// mirrors utils.CalculateFillCost / CalculateMoveCost so the feed
			// matches the cost shown on submit.
			total = float64(r.Amount)*fishWeight*pricePerUnit + additionalCostTotals[r.Id]
		}

		out = append(out, &dto.ActivityFeedItem{
			Id:            r.Id,
			Mode:          r.Mode,
			ActivityDate:  r.ActivityDate,
			CreatedAt:     r.CreatedAt,
			CreatedBy:     r.CreatedBy,
			CreatedByName: r.CreatedByName,
			PondName:      r.PondName,
			ToPondName:    r.ToPondName,
			FishType:      r.FishType,
			Amount:        r.Amount,
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
