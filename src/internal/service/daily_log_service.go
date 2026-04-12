package service

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	excel_dailylog "github.com/weeranieb/boonmafarm-backend/src/internal/excel/excel_dailylog"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/transaction"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyLogService --output=./mocks --outpkg=service --filename=daily_log_service.go --structname=MockDailyLogService --with-expecter=false
type DailyLogService interface {
	GetMonth(ctx context.Context, pondId int, month string) (*dto.DailyLogMonthResponse, error)
	BulkUpsert(ctx context.Context, pondId int, request dto.DailyLogBulkUpsertRequest, username string) error
	ImportFromExcelFile(ctx context.Context, pondId int, freshFcID, pelletFcID *int, month, filePath, username string) (int, error)
	ImportFromTemplate(ctx context.Context, farmId int, selectedPondIds []int, file []byte, username string) (*dto.DailyLogTemplateImportResponse, error)
}

type dailyLogService struct {
	dailyLogRepo         repository.DailyLogRepository
	activePondRepo       repository.ActivePondRepository
	feedCollectionRepo   repository.FeedCollectionRepository
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository
	pondRepo             repository.PondRepository
	farmRepo             repository.FarmRepository
	txManager            transaction.Manager
}

func NewDailyLogService(
	dailyLogRepo repository.DailyLogRepository,
	activePondRepo repository.ActivePondRepository,
	feedCollectionRepo repository.FeedCollectionRepository,
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository,
	pondRepo repository.PondRepository,
	farmRepo repository.FarmRepository,
	txManager transaction.Manager,
) DailyLogService {
	return &dailyLogService{
		dailyLogRepo:         dailyLogRepo,
		activePondRepo:       activePondRepo,
		feedCollectionRepo:   feedCollectionRepo,
		feedPriceHistoryRepo: feedPriceHistoryRepo,
		pondRepo:             pondRepo,
		farmRepo:             farmRepo,
		txManager:            txManager,
	}
}

// loadActivePondWithClientAccess loads the pond with farm client_id, enforces JWT client scope, and returns the active cycle row.
func (s *dailyLogService) loadActivePondWithClientAccess(ctx context.Context, pondId int) (*model.ActivePond, error) {
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
	if data.ActivePond == nil {
		return nil, errors.ErrPondNotActive
	}
	return data.ActivePond, nil
}

func (s *dailyLogService) ensureFarmTemplateImportAccess(ctx context.Context, farmId int) error {
	farm, err := s.farmRepo.GetByID(farmId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if farm == nil {
		return errors.ErrFarmNotFound
	}
	if farm.ClientId == 0 {
		return errors.ErrFarmNotFound
	}
	ok, err := utils.CanAccessClient(ctx, farm.ClientId)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return errors.ErrAuthPermissionDenied
	}
	return nil
}

func (s *dailyLogService) resolvePrices(feedCollectionId int, dates []time.Time) (map[time.Time]*decimal.Decimal, error) {
	history, err := s.feedPriceHistoryRepo.ListByFeedCollectionId(feedCollectionId)
	if err != nil {
		return nil, err
	}

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

// resolveFeedCollection loads and validates a feed collection when id is set; returns (nil, nil) when id is nil or non-positive.
func (s *dailyLogService) resolveFeedCollection(id *int, wantType string) (*model.FeedCollection, error) {
	if id == nil || *id <= 0 {
		return nil, nil
	}
	fc, err := s.feedCollectionRepo.GetByID(*id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if fc == nil {
		return nil, errors.ErrFeedCollectionNotFound
	}
	if fc.FeedType != wantType {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("feed collection %d must be type %s", *id, wantType))
	}
	return fc, nil
}

func (s *dailyLogService) validateFeedCollectionOptional(id *int, wantType string) error {
	_, err := s.resolveFeedCollection(id, wantType)
	return err
}

func (s *dailyLogService) validateBulkIDs(req *dto.DailyLogBulkUpsertRequest) error {
	if err := s.validateFeedCollectionOptional(req.FreshFeedCollectionId, constants.FeedTypeFresh); err != nil {
		return err
	}
	if err := s.validateFeedCollectionOptional(req.PelletFeedCollectionId, constants.FeedTypePellet); err != nil {
		return err
	}

	for _, e := range req.Entries {
		if !e.FreshMorning.IsZero() || !e.FreshEvening.IsZero() {
			if req.FreshFeedCollectionId == nil || *req.FreshFeedCollectionId <= 0 {
				return errors.ErrValidationFailed.Wrap(fmt.Errorf("freshFeedCollectionId is required when logging fresh feed amounts"))
			}
		}
		if !e.PelletMorning.IsZero() || !e.PelletEvening.IsZero() {
			if req.PelletFeedCollectionId == nil || *req.PelletFeedCollectionId <= 0 {
				return errors.ErrValidationFailed.Wrap(fmt.Errorf("pelletFeedCollectionId is required when logging pellet feed amounts"))
			}
		}
	}
	return nil
}

func (s *dailyLogService) GetMonth(ctx context.Context, pondId int, month string) (*dto.DailyLogMonthResponse, error) {
	ap, err := s.loadActivePondWithClientAccess(ctx, pondId)
	if err != nil {
		return nil, err
	}
	activePondId := ap.Id

	start, end, err := parseMonth(month)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}

	logs, err := s.dailyLogRepo.ListByActivePondAndMonth(ctx, activePondId, start, end)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	freshFc, err := s.resolveFeedCollection(ap.FreshFeedCollectionId, constants.FeedTypeFresh)
	if err != nil {
		return nil, err
	}
	pelletFc, err := s.resolveFeedCollection(ap.PelletFeedCollectionId, constants.FeedTypePellet)
	if err != nil {
		return nil, err
	}

	var freshID, pelletID *int
	if freshFc != nil {
		v := freshFc.Id
		freshID = &v
	}
	if pelletFc != nil {
		v := pelletFc.Id
		pelletID = &v
	}

	out := &dto.DailyLogMonthResponse{
		FreshFeedCollectionId:  freshID,
		PelletFeedCollectionId: pelletID,
		Entries:                []dto.DailyLogEntryResponse{},
	}

	if freshFc != nil {
		out.FreshFeedCollectionName = freshFc.Name
		out.FreshUnit = freshFc.Unit
	}
	if pelletFc != nil {
		out.PelletFeedCollectionName = pelletFc.Name
		out.PelletUnit = pelletFc.Unit
	}

	var freshPriceMap, pelletPriceMap map[time.Time]*decimal.Decimal
	if freshFc != nil {
		dates := make([]time.Time, 0, len(logs))
		for _, e := range logs {
			dates = append(dates, e.FeedDate)
		}
		freshPriceMap, err = s.resolvePrices(freshFc.Id, dates)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
	}
	if pelletFc != nil {
		dates := make([]time.Time, 0, len(logs))
		for _, e := range logs {
			dates = append(dates, e.FeedDate)
		}
		pelletPriceMap, err = s.resolvePrices(pelletFc.Id, dates)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
	}

	for _, e := range logs {
		er := dto.DailyLogEntryResponse{
			Id:                e.Id,
			Day:               utils.CalendarDay(e.FeedDate),
			FreshMorning:      e.FreshMorning,
			FreshEvening:      e.FreshEvening,
			PelletMorning:     e.PelletMorning,
			PelletEvening:     e.PelletEvening,
			DeathFishCount:    e.DeathFishCount,
			TouristCatchCount: e.TouristCatchCount,
		}
		if freshPriceMap != nil {
			er.FreshUnitPrice = freshPriceMap[e.FeedDate]
		}
		if pelletPriceMap != nil {
			er.PelletUnitPrice = pelletPriceMap[e.FeedDate]
		}
		out.Entries = append(out.Entries, er)
	}

	return out, nil
}

func (s *dailyLogService) BulkUpsert(ctx context.Context, pondId int, request dto.DailyLogBulkUpsertRequest, username string) error {
	ap, err := s.loadActivePondWithClientAccess(ctx, pondId)
	if err != nil {
		return err
	}
	activePondId := ap.Id

	if request.FreshFeedCollectionId == nil && ap.FreshFeedCollectionId != nil {
		v := *ap.FreshFeedCollectionId
		request.FreshFeedCollectionId = &v
	}
	if request.PelletFeedCollectionId == nil && ap.PelletFeedCollectionId != nil {
		v := *ap.PelletFeedCollectionId
		request.PelletFeedCollectionId = &v
	}

	start, _, err := parseMonth(request.Month)
	if err != nil {
		return errors.ErrValidationFailed.Wrap(err)
	}

	if err := s.validateBulkIDs(&request); err != nil {
		return err
	}

	var models []*model.DailyLog
	for _, e := range request.Entries {
		feedDate := time.Date(start.Year(), start.Month(), e.Day, 0, 0, 0, 0, time.UTC)
		if feedDate.Month() != start.Month() {
			continue
		}
		models = append(models, &model.DailyLog{
			ActivePondId:      activePondId,
			FeedDate:          feedDate,
			FreshMorning:      e.FreshMorning,
			FreshEvening:      e.FreshEvening,
			PelletMorning:     e.PelletMorning,
			PelletEvening:     e.PelletEvening,
			DeathFishCount:    e.DeathFishCount,
			TouristCatchCount: e.TouristCatchCount,
		})
	}

	var deleteDates []time.Time
	for _, d := range request.DeleteDays {
		dt := time.Date(start.Year(), start.Month(), d, 0, 0, 0, 0, time.UTC)
		if dt.Month() == start.Month() {
			deleteDates = append(deleteDates, dt)
		}
	}

	if len(models) == 0 && len(deleteDates) == 0 {
		return nil
	}

	return s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		dr := s.dailyLogRepo.WithTx(tx)
		if err := dr.Upsert(ctx, models); err != nil {
			return err
		}
		if err := dr.HardDeleteByActivePondAndDates(ctx, activePondId, deleteDates); err != nil {
			return err
		}
		updated := false
		if request.FreshFeedCollectionId != nil {
			v := *request.FreshFeedCollectionId
			ap.FreshFeedCollectionId = &v
			updated = true
		}
		if request.PelletFeedCollectionId != nil {
			v := *request.PelletFeedCollectionId
			ap.PelletFeedCollectionId = &v
			updated = true
		}
		if !updated {
			return nil
		}
		return s.activePondRepo.WithTx(tx).Update(ctx, ap)
	})
}

func (s *dailyLogService) ImportFromExcelFile(ctx context.Context, pondId int, freshFcID, pelletFcID *int, month, filePath, username string) (int, error) {
	entries, err := parseDailyLogExcelFile(filePath, month)
	if err != nil {
		return 0, err
	}
	req := dto.DailyLogBulkUpsertRequest{
		Month:                  month,
		FreshFeedCollectionId:  freshFcID,
		PelletFeedCollectionId: pelletFcID,
		Entries:                entries,
	}
	if err := s.BulkUpsert(ctx, pondId, req, username); err != nil {
		return 0, err
	}
	return len(entries), nil
}

func (s *dailyLogService) ImportFromTemplate(ctx context.Context, farmId int, selectedPondIds []int, file []byte, username string) (*dto.DailyLogTemplateImportResponse, error) {
	if err := s.ensureFarmTemplateImportAccess(ctx, farmId); err != nil {
		return nil, err
	}

	sheets, parseErr := excel_dailylog.ParseReaderAllSheets(bytes.NewReader(file), time.Now())
	if parseErr != nil && len(sheets) == 0 {
		return nil, errors.ErrGeneric.Wrap(parseErr)
	}

	ponds, err := s.pondRepo.ListByFarmId(farmId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	pondByName := make(map[string]*model.Pond, len(ponds))
	for _, p := range ponds {
		pondByName[strings.TrimSpace(p.Name)] = p
	}

	selectedSet := make(map[int]bool, len(selectedPondIds))
	for _, id := range selectedPondIds {
		selectedSet[id] = true
	}

	var results []dto.DailyLogTemplateImportResult
	var skipped []string

	for sheetName, ps := range sheets {
		pond, ok := pondByName[strings.TrimSpace(ps.PondName)]
		if !ok || !selectedSet[pond.Id] {
			skipped = append(skipped, sheetName)
			continue
		}

		activePond, err := s.activePondRepo.GetActiveByPondID(ctx, pond.Id)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		if activePond == nil {
			skipped = append(skipped, sheetName)
			continue
		}

		logs := make([]*model.DailyLog, 0, len(ps.Rows))
		for _, row := range ps.Rows {
			dl := row.ToDailyLog(activePond.Id, username)
			logs = append(logs, &dl)
		}

		if err := s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
			repo := s.dailyLogRepo.WithTx(tx)
			if len(logs) > 0 {
				importKeys, _ := templateImportDateKeys(logs)
				minD, maxD := templateImportReconcileDateRangeUTC(activePond, logs)
				existing, err := repo.ListIDAndFeedDateByActivePondRange(ctx, activePond.Id, minD, maxD)
				if err != nil {
					return err
				}
				deleteIDs := staleDailyLogIDsForTemplateImport(existing, importKeys)
				if err := repo.HardDeleteByIDs(ctx, deleteIDs); err != nil {
					return err
				}
			}
			if err := repo.Upsert(ctx, logs); err != nil {
				return err
			}
			apr := s.activePondRepo.WithTx(tx)
			updated := false
			if ps.FreshFeedCollectionId != nil {
				v := *ps.FreshFeedCollectionId
				activePond.FreshFeedCollectionId = &v
				updated = true
			}
			if ps.PelletFeedCollectionId != nil {
				v := *ps.PelletFeedCollectionId
				activePond.PelletFeedCollectionId = &v
				updated = true
			}
			if updated {
				return apr.Update(ctx, activePond)
			}
			return nil
		}); err != nil {
			return nil, errors.ErrGeneric.Wrap(fmt.Errorf("pond %q: %w", pond.Name, err))
		}

		results = append(results, dto.DailyLogTemplateImportResult{
			PondId:       pond.Id,
			PondName:     pond.Name,
			RowsImported: len(logs),
		})
	}

	return &dto.DailyLogTemplateImportResponse{
		Results: results,
		Skipped: skipped,
	}, nil
}

func dailyLogUTCDate(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// templateImportDateKeys returns distinct calendar feed dates present in the import (UTC, YYYY-MM-DD keys).
func templateImportDateKeys(logs []*model.DailyLog) (map[string]struct{}, bool) {
	if len(logs) == 0 {
		return nil, false
	}
	keys := make(map[string]struct{}, len(logs))
	for _, l := range logs {
		nd := dailyLogUTCDate(l.FeedDate)
		keys[nd.Format("2006-01-02")] = struct{}{}
	}
	return keys, true
}

// templateImportReconcileDateRangeUTC is the window for listing existing rows to reconcile: from active pond
// start (UTC date) through today (UTC date). If StartDate is zero, uses the earliest feed_date in the import.
func templateImportReconcileDateRangeUTC(activePond *model.ActivePond, logs []*model.DailyLog) (minD, maxD time.Time) {
	maxD = dailyLogUTCDate(time.Now())
	if activePond != nil && !activePond.StartDate.IsZero() {
		minD = dailyLogUTCDate(activePond.StartDate)
	} else {
		minD = dailyLogUTCDate(logs[0].FeedDate)
		for _, l := range logs[1:] {
			d := dailyLogUTCDate(l.FeedDate)
			if d.Before(minD) {
				minD = d
			}
		}
	}
	if maxD.Before(minD) {
		minD = maxD
	}
	return minD, maxD
}

func staleDailyLogIDsForTemplateImport(existing []repository.DailyLogIDFeedDate, importDateKeys map[string]struct{}) []int {
	var out []int
	for _, row := range existing {
		k := dailyLogUTCDate(row.FeedDate).Format("2006-01-02")
		if _, ok := importDateKeys[k]; !ok {
			out = append(out, row.Id)
		}
	}
	return out
}

func parseMonth(month string) (start, end time.Time, err error) {
	start, err = time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month format, expected YYYY-MM: %w", err)
	}
	end = start.AddDate(0, 1, -1)
	return start, end, nil
}
