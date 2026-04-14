package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	excel_export "github.com/weeranieb/boonmafarm-backend/src/internal/excel/excel_export"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

const (
	templateDir             = "src/assets/templates"
	templateNoFishingFile   = "template_insert_pond_no_fishing.xlsx"
	templateWithFishingFile = "template_insert_pond_fishing.xlsx"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=ExportService --output=./mocks --outpkg=service --filename=export_service.go --structname=MockExportService --with-expecter=false
type ExportService interface {
	ExportDailyLogs(ctx context.Context, farmId int, withFishing bool) ([]byte, string, error)
}

type exportService struct {
	farmRepo           repository.FarmRepository
	pondRepo           repository.PondRepository
	dailyLogRepo       repository.DailyLogRepository
	feedCollectionRepo repository.FeedCollectionRepository
}

func NewExportService(
	farmRepo repository.FarmRepository,
	pondRepo repository.PondRepository,
	dailyLogRepo repository.DailyLogRepository,
	feedCollectionRepo repository.FeedCollectionRepository,
) ExportService {
	return &exportService{
		farmRepo:           farmRepo,
		pondRepo:           pondRepo,
		dailyLogRepo:       dailyLogRepo,
		feedCollectionRepo: feedCollectionRepo,
	}
}

func (s *exportService) ExportDailyLogs(ctx context.Context, farmId int, withFishing bool) ([]byte, string, error) {
	farm, err := s.farmRepo.GetByID(farmId)
	if err != nil {
		return nil, "", errors.ErrGeneric.Wrap(err)
	}
	if farm == nil {
		return nil, "", errors.ErrFarmNotFound
	}

	ok, err := utils.CanAccessClient(ctx, farm.ClientId)
	if err != nil {
		return nil, "", errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, "", errors.ErrAuthPermissionDenied
	}

	pondData, err := s.buildPondExportData(ctx, farmId)
	if err != nil {
		return nil, "", err
	}
	if len(pondData) == 0 {
		return nil, "", errors.ErrValidationFailed.Wrap(fmt.Errorf("no active ponds with data to export"))
	}

	templatePath := s.resolveTemplatePath(withFishing)
	buf, err := excel_export.GenerateExport(templatePath, pondData, withFishing)
	if err != nil {
		return nil, "", errors.ErrGeneric.Wrap(err)
	}

	filename := fmt.Sprintf("ฟาร์ม %s.xlsx", farm.Name)
	return buf, filename, nil
}

func (s *exportService) buildPondExportData(ctx context.Context, farmId int) ([]excel_export.PondExportData, error) {
	ponds, err := s.pondRepo.ListByFarmIdWithActivePond(ctx, farmId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var result []excel_export.PondExportData
	for _, p := range ponds {
		if p.ActivePond == nil {
			continue
		}
		ap := p.ActivePond

		logs, err := s.dailyLogRepo.ListByActivePondAndMonth(ctx, ap.Id, ap.StartDate, today)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}

		ped := excel_export.PondExportData{
			PondName:               p.Pond.Name,
			Logs:                   logs,
			FreshFeedCollectionId:  ap.FreshFeedCollectionId,
			PelletFeedCollectionId: ap.PelletFeedCollectionId,
			TotalFish:              ap.TotalFish,
			StartDate:              ap.StartDate,
		}

		ped.FreshFCR, err = s.resolveFCR(ap.FreshFeedCollectionId)
		if err != nil {
			return nil, err
		}
		ped.PelletFCR, err = s.resolveFCR(ap.PelletFeedCollectionId)
		if err != nil {
			return nil, err
		}

		result = append(result, ped)
	}
	return result, nil
}

func (s *exportService) resolveFCR(feedCollectionId *int) (*decimal.Decimal, error) {
	if feedCollectionId == nil || *feedCollectionId <= 0 {
		return nil, nil
	}
	fc, err := s.feedCollectionRepo.GetByID(*feedCollectionId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if fc == nil {
		return nil, nil
	}
	return fc.Fcr, nil
}

func (s *exportService) resolveTemplatePath(withFishing bool) string {
	name := templateNoFishingFile
	if withFishing {
		name = templateWithFishingFile
	}
	return filepath.Join(templateDir, name)
}
