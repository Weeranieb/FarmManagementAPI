package service

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

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

		var total float64
		var totalWeight *float64
		switch r.Mode {
		case constants.ActivityModeSell:
			st := sellTotals[r.Id]
			total = st.total
			w := st.weight
			totalWeight = &w
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
			PondName:      r.PondName,
			FarmName:      r.FarmName,
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
