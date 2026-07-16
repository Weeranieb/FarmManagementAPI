package service

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FeedCostCalculator --output=./mocks --outpkg=service --filename=feed_cost_calculator.go --structname=MockFeedCostCalculator --with-expecter=false

// FeedCostCalculator derives the total feed cost of an active pond cycle from
// its daily logs and the feed collections' price history. It is a pure
// read-side aggregation (no writes) shared by the pond read path (live for
// active cycles) and the close path (snapshotted onto active_ponds.feed_cost).
type FeedCostCalculator interface {
	// CalcCycleFeedCost returns the whole-cycle feed cost for one active pond.
	CalcCycleFeedCost(ctx context.Context, activePond *model.ActivePond) (decimal.Decimal, error)
	// CalcCycleFeedCostBatch returns feed cost keyed by active pond id for many
	// ponds using at most two queries total (logs + price history), so a farm
	// listing does not fan out into an N+1. Every non-nil input pond is present
	// in the result (zero when it has no priced feed logs).
	CalcCycleFeedCostBatch(ctx context.Context, activePonds []*model.ActivePond) (map[int]decimal.Decimal, error)
}

type feedCostCalculator struct {
	dailyLogRepo         repository.DailyLogRepository
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository
}

func NewFeedCostCalculator(
	dailyLogRepo repository.DailyLogRepository,
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository,
) FeedCostCalculator {
	return &feedCostCalculator{
		dailyLogRepo:         dailyLogRepo,
		feedPriceHistoryRepo: feedPriceHistoryRepo,
	}
}

func (c *feedCostCalculator) CalcCycleFeedCost(ctx context.Context, activePond *model.ActivePond) (decimal.Decimal, error) {
	if activePond == nil {
		return decimal.Zero, nil
	}
	byPond, err := c.CalcCycleFeedCostBatch(ctx, []*model.ActivePond{activePond})
	if err != nil {
		return decimal.Zero, err
	}
	return byPond[activePond.Id], nil
}

func (c *feedCostCalculator) CalcCycleFeedCostBatch(ctx context.Context, activePonds []*model.ActivePond) (map[int]decimal.Decimal, error) {
	result := make(map[int]decimal.Decimal, len(activePonds))

	apIds := make([]int, 0, len(activePonds))
	freshByAp := make(map[int]int, len(activePonds))
	pelletByAp := make(map[int]int, len(activePonds))
	fcIdSet := make(map[int]struct{})
	for _, ap := range activePonds {
		if ap == nil {
			continue
		}
		apIds = append(apIds, ap.Id)
		result[ap.Id] = decimal.Zero
		if ap.FreshFeedCollectionId != nil && *ap.FreshFeedCollectionId > 0 {
			freshByAp[ap.Id] = *ap.FreshFeedCollectionId
			fcIdSet[*ap.FreshFeedCollectionId] = struct{}{}
		}
		if ap.PelletFeedCollectionId != nil && *ap.PelletFeedCollectionId > 0 {
			pelletByAp[ap.Id] = *ap.PelletFeedCollectionId
			fcIdSet[*ap.PelletFeedCollectionId] = struct{}{}
		}
	}
	if len(apIds) == 0 {
		return result, nil
	}

	logs, err := c.dailyLogRepo.ListByActivePondIds(ctx, apIds)
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return result, nil
	}

	fcIds := make([]int, 0, len(fcIdSet))
	for id := range fcIdSet {
		fcIds = append(fcIds, id)
	}
	histories, err := c.feedPriceHistoryRepo.ListByFeedCollectionIds(fcIds)
	if err != nil {
		return nil, err
	}
	histByFc := groupPriceHistoryAscending(histories)

	// Warn once per (active pond, feed collection) when feed was logged but the
	// collection has no price history at all — that feed is silently counted as
	// zero (resolveFeedPrice can only fall back when some price exists), which
	// would otherwise understate cost with no signal.
	warned := make(map[[2]int]struct{})
	warnMissingPrice := func(activePondId, feedCollectionId int) {
		key := [2]int{activePondId, feedCollectionId}
		if _, done := warned[key]; done {
			return
		}
		warned[key] = struct{}{}
		slog.WarnContext(ctx, "feed cost not counted: feed collection has no price history",
			"active_pond_id", activePondId, "feed_collection_id", feedCollectionId)
	}

	// Warn once per (active pond, feed kind) when feed was logged but the cycle
	// has no feed collection configured for that kind at all — the amount is
	// counted as zero and would otherwise understate cost with no signal.
	warnedNoCollection := make(map[noCollectionKey]struct{})
	warnMissingCollection := func(activePondId int, kind string) {
		key := noCollectionKey{activePondId: activePondId, kind: kind}
		if _, done := warnedNoCollection[key]; done {
			return
		}
		warnedNoCollection[key] = struct{}{}
		slog.WarnContext(ctx, "feed cost not counted: no feed collection configured for logged feed",
			"active_pond_id", activePondId, "feed_kind", kind)
	}

	for _, log := range logs {
		cost := decimal.Zero
		if !log.Fresh.IsZero() {
			if fcId, ok := freshByAp[log.ActivePondId]; ok {
				if p := resolveFeedPrice(histByFc[fcId], log.FeedDate); p != nil {
					cost = cost.Add(log.Fresh.Mul(*p))
				} else {
					warnMissingPrice(log.ActivePondId, fcId)
				}
			} else {
				warnMissingCollection(log.ActivePondId, "fresh")
			}
		}
		pellet := log.PelletMorning.Add(log.PelletEvening)
		if !pellet.IsZero() {
			if fcId, ok := pelletByAp[log.ActivePondId]; ok {
				if p := resolveFeedPrice(histByFc[fcId], log.FeedDate); p != nil {
					cost = cost.Add(pellet.Mul(*p))
				} else {
					warnMissingPrice(log.ActivePondId, fcId)
				}
			} else {
				warnMissingCollection(log.ActivePondId, "pellet")
			}
		}
		result[log.ActivePondId] = result[log.ActivePondId].Add(cost)
	}
	return result, nil
}

// noCollectionKey dedups "feed logged but no collection configured" warnings per
// active pond and feed kind ("fresh"/"pellet").
type noCollectionKey struct {
	activePondId int
	kind         string
}

// groupPriceHistoryAscending buckets price history rows by feed collection id,
// each bucket sorted ascending by PriceUpdatedDate (oldest first).
func groupPriceHistoryAscending(histories []*model.FeedPriceHistory) map[int][]*model.FeedPriceHistory {
	byFc := make(map[int][]*model.FeedPriceHistory)
	for _, h := range histories {
		byFc[h.FeedCollectionId] = append(byFc[h.FeedCollectionId], h)
	}
	for id := range byFc {
		rows := byFc[id]
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].PriceUpdatedDate.Before(rows[j].PriceUpdatedDate)
		})
	}
	return byFc
}

// resolveFeedPrice returns the feed unit price to use for a log dated `date`,
// given that collection's price history sorted ascending by PriceUpdatedDate:
//  1. the latest price whose PriceUpdatedDate is on or before `date` (the price
//     actually in effect that day); otherwise
//  2. the earliest recorded price (nearest available in the future) as a
//     fallback when `date` predates every recorded price; otherwise
//  3. nil when there is no price history at all.
func resolveFeedPrice(historyAsc []*model.FeedPriceHistory, date time.Time) *decimal.Decimal {
	if len(historyAsc) == 0 {
		return nil
	}
	for i := len(historyAsc) - 1; i >= 0; i-- {
		if !historyAsc[i].PriceUpdatedDate.After(date) {
			p := historyAsc[i].Price
			return &p
		}
	}
	p := historyAsc[0].Price
	return &p
}
