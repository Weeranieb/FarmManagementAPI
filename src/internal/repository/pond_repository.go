package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository/projection"

	"gorm.io/gorm"
)

// PondWithFarmAndActivePond is the result of a single query joining pond, farm (client_id), and active pond.
type PondWithFarmAndActivePond struct {
	Pond               *model.Pond
	ClientId           int
	ActivePond         *model.ActivePond // nil when pond has no active cycle
	LatestActivityDate *time.Time        // max(activity_date) for this pond's activities
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=PondRepository --output=./mocks --outpkg=mocks --filename=pond_repository.go --structname=MockPondRepository --with-expecter=false
type PondRepository interface {
	Create(ctx context.Context, pond *model.Pond) error
	CreateBatch(ctx context.Context, ponds []*model.Pond) error
	GetByID(id int) (*model.Pond, error)
	GetByIDWithFarmAndActivePond(ctx context.Context, pondId int) (*PondWithFarmAndActivePond, error)
	GetByFarmIdAndName(farmId int, name string) (*model.Pond, error)
	Update(ctx context.Context, pond *model.Pond) error
	ListByFarmId(farmId int) ([]*model.Pond, error)
	ListByFarmIdWithActivePond(ctx context.Context, farmId int) ([]*PondWithFarmAndActivePond, error)
	Delete(id int) error
}

type pondRepository struct {
	db *gorm.DB
}

func NewPondRepository(db *gorm.DB) PondRepository {
	return &pondRepository{db: db}
}

func (r *pondRepository) Create(ctx context.Context, pond *model.Pond) error {
	return r.db.WithContext(ctx).Create(pond).Error
}

func (r *pondRepository) CreateBatch(ctx context.Context, ponds []*model.Pond) error {
	return r.db.WithContext(ctx).Create(ponds).Error
}

func (r *pondRepository) GetByID(id int) (*model.Pond, error) {
	var pond model.Pond
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pond).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pond, nil
}

const pondWithFarmAndActivePondQuery = `
SELECT
  p.id AS pond_id, p.farm_id AS pond_farm_id, p.name AS pond_name, p.status AS pond_status,
  p.deleted_at AS pond_deleted_at, p.created_at AS pond_created_at, p.created_by AS pond_created_by,
  p.updated_at AS pond_updated_at, p.updated_by AS pond_updated_by,
  f.client_id,
  ap.id AS ap_id, ap.pond_id AS ap_pond_id, ap.start_date AS ap_start_date, ap.end_date AS ap_end_date,
  ap.is_active AS ap_is_active, ap.total_cost AS ap_total_cost, ap.total_profit AS ap_total_profit,
  ap.net_result AS ap_net_result, ap.total_fish AS ap_total_fish, ap.fish_types AS ap_fish_types,
  ap.created_at AS ap_created_at, ap.created_by AS ap_created_by,
  ap.updated_at AS ap_updated_at, ap.updated_by AS ap_updated_by,
  (SELECT MAX(a.activity_date) FROM activities a INNER JOIN active_ponds ap2 ON a.active_pond_id = ap2.id AND ap2.deleted_at IS NULL WHERE ap2.pond_id = p.id AND a.deleted_at IS NULL) AS latest_activity_date
FROM ponds p
INNER JOIN farms f ON p.farm_id = f.id AND f.deleted_at IS NULL
LEFT JOIN active_ponds ap ON ap.pond_id = p.id AND ap.is_active = true AND ap.deleted_at IS NULL
WHERE p.id = ? AND p.deleted_at IS NULL`

const pondListWithActivePondQuery = `
SELECT
  p.id AS pond_id, p.farm_id AS pond_farm_id, p.name AS pond_name, p.status AS pond_status,
  p.deleted_at AS pond_deleted_at, p.created_at AS pond_created_at, p.created_by AS pond_created_by,
  p.updated_at AS pond_updated_at, p.updated_by AS pond_updated_by,
  f.client_id,
  ap.id AS ap_id, ap.pond_id AS ap_pond_id, ap.start_date AS ap_start_date, ap.end_date AS ap_end_date,
  ap.is_active AS ap_is_active, ap.total_cost AS ap_total_cost, ap.total_profit AS ap_total_profit,
  ap.net_result AS ap_net_result, ap.total_fish AS ap_total_fish, ap.fish_types AS ap_fish_types,
  ap.created_at AS ap_created_at, ap.created_by AS ap_created_by,
  ap.updated_at AS ap_updated_at, ap.updated_by AS ap_updated_by,
  (SELECT MAX(a.activity_date) FROM activities a INNER JOIN active_ponds ap2 ON a.active_pond_id = ap2.id AND ap2.deleted_at IS NULL WHERE ap2.pond_id = p.id AND a.deleted_at IS NULL) AS latest_activity_date
FROM ponds p
INNER JOIN farms f ON p.farm_id = f.id AND f.deleted_at IS NULL
LEFT JOIN active_ponds ap ON ap.pond_id = p.id AND ap.is_active = true AND ap.deleted_at IS NULL
WHERE p.farm_id = ? AND p.deleted_at IS NULL`

func rowToPondWithFarmAndActivePond(row *projection.PondFillQueryRow) *PondWithFarmAndActivePond {
	pond := &model.Pond{
		Id:     row.PondId,
		FarmId: row.PondFarmId,
		Name:   row.PondName,
		Status: row.PondStatus,
		BaseModel: model.BaseModel{
			DeletedAt: row.PondDeletedAt,
			CreatedAt: row.PondCreatedAt,
			CreatedBy: row.PondCreatedBy,
			UpdatedAt: row.PondUpdatedAt,
			UpdatedBy: row.PondUpdatedBy,
		},
	}
	out := &PondWithFarmAndActivePond{Pond: pond, ClientId: row.ClientId, LatestActivityDate: row.LatestActivityDate}
	if row.ApId != nil {
		totalCost := decimal.Zero
		if row.ApTotalCost != nil {
			if d, err := decimal.NewFromString(*row.ApTotalCost); err == nil {
				totalCost = d
			}
		}
		totalProfit := decimal.Zero
		if row.ApTotalProfit != nil {
			if d, err := decimal.NewFromString(*row.ApTotalProfit); err == nil {
				totalProfit = d
			}
		}
		netResult := decimal.Zero
		if row.ApNetResult != nil {
			if d, err := decimal.NewFromString(*row.ApNetResult); err == nil {
				netResult = d
			}
		}
		fishTypes := parseFishTypesJSON(row.ApFishTypes)
		ap := &model.ActivePond{
			Id:          *row.ApId,
			PondId:      *row.ApPondId,
			StartDate:   *row.ApStartDate,
			EndDate:     row.ApEndDate,
			IsActive:    *row.ApIsActive,
			TotalCost:   totalCost,
			TotalProfit: totalProfit,
			NetResult:   netResult,
			TotalFish:   ptrToInt(row.ApTotalFish),
			FishTypes:   fishTypes,
		}
		if row.ApCreatedAt != nil && row.ApCreatedBy != nil && row.ApUpdatedAt != nil && row.ApUpdatedBy != nil {
			ap.BaseModel = model.BaseModel{
				CreatedAt: *row.ApCreatedAt,
				CreatedBy: *row.ApCreatedBy,
				UpdatedAt: *row.ApUpdatedAt,
				UpdatedBy: *row.ApUpdatedBy,
			}
		}
		out.ActivePond = ap
	}
	return out
}

func (r *pondRepository) GetByIDWithFarmAndActivePond(ctx context.Context, pondId int) (*PondWithFarmAndActivePond, error) {
	var row projection.PondFillQueryRow
	err := r.db.WithContext(ctx).Raw(pondWithFarmAndActivePondQuery, pondId).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.PondId == 0 {
		return nil, nil
	}
	return rowToPondWithFarmAndActivePond(&row), nil
}

func (r *pondRepository) ListByFarmIdWithActivePond(ctx context.Context, farmId int) ([]*PondWithFarmAndActivePond, error) {
	var rows []projection.PondFillQueryRow
	err := r.db.WithContext(ctx).Raw(pondListWithActivePondQuery, farmId).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*PondWithFarmAndActivePond, 0, len(rows))
	for i := range rows {
		out = append(out, rowToPondWithFarmAndActivePond(&rows[i]))
	}
	return out, nil
}

func ptrToInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func parseFishTypesJSON(s *string) []string {
	if s == nil || *s == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(*s), &out); err != nil {
		return nil
	}
	return out
}

func (r *pondRepository) GetByFarmIdAndName(farmId int, name string) (*model.Pond, error) {
	var pond model.Pond
	err := r.db.Where("farm_id = ? AND name = ? AND deleted_at IS NULL", farmId, name).First(&pond).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pond, nil
}

func (r *pondRepository) Update(ctx context.Context, pond *model.Pond) error {
	return r.db.WithContext(ctx).Save(pond).Error
}

func (r *pondRepository) ListByFarmId(farmId int) ([]*model.Pond, error) {
	var ponds []*model.Pond
	err := r.db.Where("farm_id = ? AND deleted_at IS NULL", farmId).Find(&ponds).Error
	return ponds, err
}

func (r *pondRepository) Delete(id int) error {
	return r.db.Delete(&model.Pond{}, id).Error
}
