package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyFeedService --output=./mocks --outpkg=service --filename=daily_feed_service.go --structname=MockDailyFeedService --with-expecter=false
type DailyFeedService interface {
	GetMonth(ctx context.Context, pondId int, month string) ([]*dto.DailyFeedTableResponse, error)
	BulkUpsert(ctx context.Context, pondId int, request dto.DailyFeedBulkUpsertRequest, username string) error
	DeleteTable(ctx context.Context, pondId int, feedCollectionId int) error
}

type dailyFeedService struct {
	dailyFeedRepo        repository.DailyFeedRepository
	activePondRepo       repository.ActivePondRepository
	feedCollectionRepo   repository.FeedCollectionRepository
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository
}

func NewDailyFeedService(
	dailyFeedRepo repository.DailyFeedRepository,
	activePondRepo repository.ActivePondRepository,
	feedCollectionRepo repository.FeedCollectionRepository,
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository,
) DailyFeedService {
	return &dailyFeedService{
		dailyFeedRepo:        dailyFeedRepo,
		activePondRepo:       activePondRepo,
		feedCollectionRepo:   feedCollectionRepo,
		feedPriceHistoryRepo: feedPriceHistoryRepo,
	}
}

func parseMonth(month string) (start, end time.Time, err error) {
	start, err = time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month format, expected YYYY-MM: %w", err)
	}
	end = start.AddDate(0, 1, -1)
	return start, end, nil
}

func (s *dailyFeedService) getActivePondId(ctx context.Context, pondId int) (int, error) {
	ap, err := s.activePondRepo.GetActiveByPondID(ctx, pondId)
	if err != nil {
		return 0, errors.ErrGeneric.Wrap(err)
	}
	if ap == nil {
		return 0, errors.ErrPondNotActive
	}
	return ap.Id, nil
}

// resolvePrices returns the effective unit price for each feed_date, given a feed_collection_id.
// It loads the full price history for that feed and finds the latest price <= each date.
func (s *dailyFeedService) resolvePrices(feedCollectionId int, dates []time.Time) (map[time.Time]*decimal.Decimal, error) {
	history, err := s.feedPriceHistoryRepo.ListByFeedCollectionId(feedCollectionId)
	if err != nil {
		return nil, err
	}

	// history is already sorted DESC by price_updated_date
	sort.Slice(history, func(i, j int) bool {
		return history[i].PriceUpdatedDate.Before(history[j].PriceUpdatedDate)
	})

	result := make(map[time.Time]*decimal.Decimal, len(dates))
	for _, d := range dates {
		var found *decimal.Decimal
		for i := len(history) - 1; i >= 0; i-- {
			if !history[i].PriceUpdatedDate.After(d) {
				p := history[i].Price
				found = &p
				break
			}
		}
		result[d] = found
	}
	return result, nil
}

func (s *dailyFeedService) GetMonth(ctx context.Context, pondId int, month string) ([]*dto.DailyFeedTableResponse, error) {
	activePondId, err := s.getActivePondId(ctx, pondId)
	if err != nil {
		return nil, err
	}

	start, end, err := parseMonth(month)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}

	feeds, err := s.dailyFeedRepo.ListByActivePondAndMonth(activePondId, start, end)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	// Group by feed_collection_id
	grouped := make(map[int][]*model.DailyFeed)
	fcIds := make(map[int]bool)
	for _, f := range feeds {
		grouped[f.FeedCollectionId] = append(grouped[f.FeedCollectionId], f)
		fcIds[f.FeedCollectionId] = true
	}

	// Also include feed collections with no entries this month but that have entries in other months
	allFcIds, err := s.dailyFeedRepo.ListFeedCollectionIdsByActivePond(activePondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	for _, id := range allFcIds {
		fcIds[id] = true
	}

	var tables []*dto.DailyFeedTableResponse
	for fcId := range fcIds {
		fc, err := s.feedCollectionRepo.GetByID(fcId)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		if fc == nil {
			continue
		}

		entries := grouped[fcId]
		var dates []time.Time
		for _, e := range entries {
			dates = append(dates, e.FeedDate)
		}

		priceMap, err := s.resolvePrices(fcId, dates)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}

		var dtoEntries []dto.DailyFeedEntryResponse
		for _, e := range entries {
			dtoEntries = append(dtoEntries, dto.DailyFeedEntryResponse{
				Id:        e.Id,
				Day:       e.FeedDate.Day(),
				Morning:   e.MorningAmount,
				Evening:   e.EveningAmount,
				UnitPrice: priceMap[e.FeedDate],
			})
		}
		if dtoEntries == nil {
			dtoEntries = []dto.DailyFeedEntryResponse{}
		}

		tables = append(tables, &dto.DailyFeedTableResponse{
			FeedCollectionId:   fcId,
			FeedCollectionName: fc.Name,
			FeedUnit:           fc.Unit,
			Entries:            dtoEntries,
		})
	}

	if tables == nil {
		tables = []*dto.DailyFeedTableResponse{}
	}
	return tables, nil
}

func (s *dailyFeedService) BulkUpsert(ctx context.Context, pondId int, request dto.DailyFeedBulkUpsertRequest, username string) error {
	activePondId, err := s.getActivePondId(ctx, pondId)
	if err != nil {
		return err
	}

	start, _, err := parseMonth(request.Month)
	if err != nil {
		return errors.ErrValidationFailed.Wrap(err)
	}

	// Verify feed collection exists
	fc, err := s.feedCollectionRepo.GetByID(request.FeedCollectionId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if fc == nil {
		return errors.ErrFeedCollectionNotFound
	}

	var models []*model.DailyFeed
	for _, e := range request.Entries {
		feedDate := time.Date(start.Year(), start.Month(), e.Day, 0, 0, 0, 0, time.UTC)
		if feedDate.Month() != start.Month() {
			continue // skip invalid day numbers for the month
		}
		models = append(models, &model.DailyFeed{
			ActivePondId:     activePondId,
			FeedCollectionId: request.FeedCollectionId,
			FeedDate:         feedDate,
			MorningAmount:    e.Morning,
			EveningAmount:    e.Evening,
		})
	}

	return s.dailyFeedRepo.Upsert(ctx, models)
}

func (s *dailyFeedService) DeleteTable(ctx context.Context, pondId int, feedCollectionId int) error {
	activePondId, err := s.getActivePondId(ctx, pondId)
	if err != nil {
		return err
	}
	return s.dailyFeedRepo.SoftDeleteByActivePondAndFeedCollection(ctx, activePondId, feedCollectionId)
}
