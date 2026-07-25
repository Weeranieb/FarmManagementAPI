package service

import (
	"context"
	"slices"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

type sellActivityTotal struct {
	total     float64
	weight    float64
	fishCount int
}

type activityIDMode struct {
	id   int
	mode string
}

func partitionActivityIDs(rows []activityIDMode) (sellIds, fillMoveIds []int) {
	for _, r := range rows {
		switch r.mode {
		case constants.ActivityModeSell:
			sellIds = append(sellIds, r.id)
		default:
			fillMoveIds = append(fillMoveIds, r.id)
		}
	}
	return sellIds, fillMoveIds
}

func loadActivityCostTotals(
	ctx context.Context,
	repo repository.ActivityRepository,
	sellIds, fillMoveIds []int,
) (map[int]sellActivityTotal, map[int]float64, error) {
	sellTotals := map[int]sellActivityTotal{}
	if len(sellIds) > 0 {
		totals, err := repo.SumSellDetailsByActivityIDs(ctx, sellIds)
		if err != nil {
			return nil, nil, errors.ErrGeneric.Wrap(err)
		}
		for _, t := range totals {
			total, _ := t.Total.Float64()
			weight, _ := t.TotalWeight.Float64()
			sellTotals[t.SellId] = sellActivityTotal{
				total:     total,
				weight:    weight,
				fishCount: t.TotalFishCount,
			}
		}
	}

	// Additional costs are summed for every mode. Fill/move fold them into the
	// activity total (see fillMoveActivityTotal); a sell keeps its total at gross
	// revenue and reports the cost separately, so the two never double-count.
	costIds := slices.Concat(sellIds, fillMoveIds)
	additionalCostTotals := map[int]float64{}
	if len(costIds) > 0 {
		totals, err := repo.SumAdditionalCostsByActivityIDs(ctx, costIds)
		if err != nil {
			return nil, nil, errors.ErrGeneric.Wrap(err)
		}
		for _, t := range totals {
			f, _ := t.Total.Float64()
			additionalCostTotals[t.ActivityId] = f
		}
	}
	return sellTotals, additionalCostTotals, nil
}

func fillMoveActivityTotal(amount int, fishWeight, pricePerUnit, additionalCost float64) float64 {
	return float64(amount)*fishWeight*pricePerUnit + additionalCost
}
