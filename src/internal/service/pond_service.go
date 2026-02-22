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
}

type PondServiceParams struct {
	dig.In

	PondRepo       repository.PondRepository
	FarmRepo       repository.FarmRepository
	ActivePondRepo repository.ActivePondRepository
	ActivityRepo   repository.ActivityRepository
	TxManager      transaction.Manager
}

type pondService struct {
	pondRepo       repository.PondRepository
	farmRepo       repository.FarmRepository
	activePondRepo repository.ActivePondRepository
	activityRepo   repository.ActivityRepository
	txManager      transaction.Manager
}

func NewPondService(params PondServiceParams) PondService {
	return &pondService{
		pondRepo:       params.PondRepo,
		farmRepo:       params.FarmRepo,
		activePondRepo: params.ActivePondRepo,
		activityRepo:   params.ActivityRepo,
		txManager:      params.TxManager,
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
	pond, err := s.pondRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if pond == nil {
		return nil, errors.ErrPondNotFound
	}

	return s.toPondResponse(pond), nil
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
	ponds, err := s.pondRepo.ListByFarmId(farmId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.PondResponse, 0, len(ponds))
	for _, pond := range ponds {
		responses = append(responses, s.toPondResponse(pond))
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

	additionalCosts := make([]decimal.Decimal, 0, len(request.AdditionalCosts))
	for _, item := range request.AdditionalCosts {
		additionalCosts = append(additionalCosts, item.Cost)
	}
	fillCost := utils.FillCost(request.Amount, request.PricePerUnit, additionalCosts)

	var resp *dto.PondFillResponse
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		if activePond == nil {
			activePond = &model.ActivePond{
				PondId:    pondId,
				StartDate: activityDate,
				IsActive:  true,
				TotalCost: fillCost,
				NetResult: decimal.Zero.Sub(fillCost), // TotalProfit - TotalCost; no sales yet
			}
			if err := tx.Create(activePond).Error; err != nil {
				return err
			}
			if pond.Status == constants.FarmStatusMaintenance {
				pond.Status = constants.FarmStatusActive
				if err := tx.Save(pond).Error; err != nil {
					return err
				}
			}
		} else {
			activePond.TotalCost = activePond.TotalCost.Add(fillCost)
			activePond.NetResult = activePond.TotalProfit.Sub(activePond.TotalCost)
			if err := tx.Save(activePond).Error; err != nil {
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
		if err := tx.Create(activity).Error; err != nil {
			return err
		}

		for _, item := range request.AdditionalCosts {
			ac := &model.AdditionalCost{
				ActivityId: activity.Id,
				Title:      item.Title,
				Cost:       item.Cost,
			}
			if err := tx.Create(ac).Error; err != nil {
				return err
			}
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

func (s *pondService) toPondResponse(pond *model.Pond) *dto.PondResponse {
	return &dto.PondResponse{
		Id:        pond.Id,
		FarmId:    pond.FarmId,
		Name:      pond.Name,
		Status:    pond.Status,
		CreatedAt: pond.CreatedAt,
		CreatedBy: pond.CreatedBy,
		UpdatedAt: pond.UpdatedAt,
		UpdatedBy: pond.UpdatedBy,
	}
}
