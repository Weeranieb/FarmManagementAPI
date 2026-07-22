package repository

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"

	"gorm.io/gorm"
)

// ActivityListRow is the flat join projection returned by ListByPondID.
// It is NOT a persisted model — only a query result.
//
// For "move" activities the same row appears in BOTH the source and the
// destination pond's history. `IsIncoming` distinguishes the two perspectives:
// when true, the requested pond was the destination (fish in); when false,
// it was the source (fish out, or a non-move row which is always outgoing).
type ActivityListRow struct {
	Id           int             `gorm:"column:id"`
	Mode         string          `gorm:"column:mode"`
	ActivityDate time.Time       `gorm:"column:activity_date"`
	FishType     string          `gorm:"column:fish_type"`
	Amount       int             `gorm:"column:amount"`
	FishWeight   decimal.Decimal `gorm:"column:fish_weight"`
	PricePerUnit decimal.Decimal `gorm:"column:price_per_unit"`
	MerchantName *string         `gorm:"column:merchant_name"`
	ToPondName   *string         `gorm:"column:to_pond_name"`
	FromPondName *string         `gorm:"column:from_pond_name"`
	IsIncoming   bool            `gorm:"column:is_incoming"`
}

// SellTotalRow is the grouped sum projection for sell totals keyed by activity id.
type SellTotalRow struct {
	SellId int             `gorm:"column:sell_id"`
	Total  decimal.Decimal `gorm:"column:total"`
	// TotalWeight is the summed sell_details.weight (kg) — the activity row
	// itself only stores an aggregate fish count, so weight must come from
	// the detail lines.
	TotalWeight decimal.Decimal `gorm:"column:total_weight"`
}

// ActivityFeedRow is the flat join projection returned by ListRecentByClientID.
// One row per activity (moves are NOT duplicated per pond perspective here —
// the feed shows "source → destination" in a single entry).
type ActivityFeedRow struct {
	Id            int             `gorm:"column:id"`
	Mode          string          `gorm:"column:mode"`
	ActivityDate  time.Time       `gorm:"column:activity_date"`
	CreatedAt     time.Time       `gorm:"column:created_at"`
	CreatedBy     string          `gorm:"column:created_by"`
	CreatedByName string          `gorm:"column:created_by_name"`
	PondName      string          `gorm:"column:pond_name"`
	FarmName      string          `gorm:"column:farm_name"`
	ToPondName    *string         `gorm:"column:to_pond_name"`
	FishType      string          `gorm:"column:fish_type"`
	Amount        int             `gorm:"column:amount"`
	FishWeight    decimal.Decimal `gorm:"column:fish_weight"`
	FishUnit      string          `gorm:"column:fish_unit"`
	PricePerUnit  decimal.Decimal `gorm:"column:price_per_unit"`
	MerchantName  *string         `gorm:"column:merchant_name"`
}

// AdditionalCostTotalRow is the grouped sum projection for additional_costs
// keyed by activity id. Used to fold extra costs into fill/move totals.
type AdditionalCostTotalRow struct {
	ActivityId int             `gorm:"column:activity_id"`
	Total      decimal.Decimal `gorm:"column:total"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=ActivityRepository --output=./mocks --outpkg=mocks --filename=activity_repository.go --structname=MockActivityRepository --with-expecter=false
type ActivityRepository interface {
	WithTx(tx *gorm.DB) ActivityRepository
	Create(ctx context.Context, activity *model.Activity) error
	ListByPondID(ctx context.Context, pondId int) ([]ActivityListRow, error)
	ListRecentByClientID(ctx context.Context, clientId *int, limit int) ([]ActivityFeedRow, error)
	SumSellDetailsByActivityIDs(ctx context.Context, activityIds []int) ([]SellTotalRow, error)
	SumAdditionalCostsByActivityIDs(ctx context.Context, activityIds []int) ([]AdditionalCostTotalRow, error)
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) WithTx(tx *gorm.DB) ActivityRepository {
	return &activityRepository{db: tx}
}

func (r *activityRepository) Create(ctx context.Context, activity *model.Activity) error {
	return r.db.WithContext(ctx).Create(activity).Error
}

// Lists fill/move/sell rows for a pond. The pond may appear as the source
// (active_pond_id) — fill/sell rows are always sourced — or as the
// destination of a move (to_active_pond_id). The `is_incoming` CASE flags
// the destination perspective so the service can populate FromPondName /
// ToPondName accordingly. Without the OR on tap.pond_id, the destination
// pond would never see incoming moves in its history.
const activityListByPondQuery = `
SELECT
  a.id AS id,
  a.mode AS mode,
  a.activity_date AS activity_date,
  a.fish_type AS fish_type,
  a.amount AS amount,
  a.fish_weight AS fish_weight,
  a.price_per_unit AS price_per_unit,
  m.name AS merchant_name,
  CASE WHEN tap.pond_id = ? THEN NULL ELSE tp.name END AS to_pond_name,
  CASE WHEN tap.pond_id = ? THEN fp.name ELSE NULL END AS from_pond_name,
  CASE WHEN tap.pond_id = ? THEN TRUE ELSE FALSE END AS is_incoming
FROM activities a
INNER JOIN active_ponds ap ON a.active_pond_id = ap.id AND ap.deleted_at IS NULL
LEFT JOIN ponds fp ON ap.pond_id = fp.id AND fp.deleted_at IS NULL
LEFT JOIN merchants m ON a.merchant_id = m.id AND m.deleted_at IS NULL
LEFT JOIN active_ponds tap ON a.to_active_pond_id = tap.id AND tap.deleted_at IS NULL
LEFT JOIN ponds tp ON tap.pond_id = tp.id AND tp.deleted_at IS NULL
WHERE (ap.pond_id = ? OR tap.pond_id = ?) AND a.deleted_at IS NULL
ORDER BY a.activity_date DESC, a.id DESC`

func (r *activityRepository) ListByPondID(ctx context.Context, pondId int) ([]ActivityListRow, error) {
	var rows []ActivityListRow
	err := r.db.WithContext(ctx).
		Raw(activityListByPondQuery, pondId, pondId, pondId, pondId, pondId).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Lists the newest fill/move/sell rows across every pond the client owns
// (activities → active_ponds → ponds → farms → client_id). Unlike
// ListByPondID, each move appears exactly once — the feed renders it as
// "source → destination" instead of two per-pond perspectives. created_by
// stores a username; the users join resolves it to a display name for the
// "โดย…" meta line (falls back to the raw username for system/legacy rows).
const activityFeedQuery = `
SELECT
  a.id AS id,
  a.mode AS mode,
  a.activity_date AS activity_date,
  a.created_at AS created_at,
  a.created_by AS created_by,
  COALESCE(u.first_name, a.created_by) AS created_by_name,
  fp.name AS pond_name,
  f.name AS farm_name,
  tp.name AS to_pond_name,
  a.fish_type AS fish_type,
  a.amount AS amount,
  a.fish_weight AS fish_weight,
  a.fish_unit AS fish_unit,
  a.price_per_unit AS price_per_unit,
  m.name AS merchant_name
FROM activities a
INNER JOIN active_ponds ap ON a.active_pond_id = ap.id AND ap.deleted_at IS NULL
INNER JOIN ponds fp ON ap.pond_id = fp.id AND fp.deleted_at IS NULL
INNER JOIN farms f ON fp.farm_id = f.id AND f.deleted_at IS NULL
LEFT JOIN merchants m ON a.merchant_id = m.id AND m.deleted_at IS NULL
LEFT JOIN active_ponds tap ON a.to_active_pond_id = tap.id AND tap.deleted_at IS NULL
LEFT JOIN ponds tp ON tap.pond_id = tp.id AND tp.deleted_at IS NULL
LEFT JOIN users u ON a.created_by = u.username AND u.deleted_at IS NULL
WHERE a.deleted_at IS NULL`

func (r *activityRepository) ListRecentByClientID(ctx context.Context, clientId *int, limit int) ([]ActivityFeedRow, error) {
	query := activityFeedQuery
	args := make([]interface{}, 0, 2)
	if clientId != nil {
		query += ` AND f.client_id = ?`
		args = append(args, *clientId)
	}
	query += ` ORDER BY a.activity_date DESC, a.id DESC`
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	var rows []ActivityFeedRow
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *activityRepository) SumSellDetailsByActivityIDs(ctx context.Context, activityIds []int) ([]SellTotalRow, error) {
	if len(activityIds) == 0 {
		return nil, nil
	}
	var rows []SellTotalRow
	err := r.db.WithContext(ctx).Raw(`
SELECT sd.sell_id AS sell_id, SUM(sd.weight * sd.price_per_unit) AS total, SUM(sd.weight) AS total_weight
FROM sell_details sd
WHERE sd.sell_id IN ? AND sd.deleted_at IS NULL
GROUP BY sd.sell_id`, activityIds).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *activityRepository) SumAdditionalCostsByActivityIDs(ctx context.Context, activityIds []int) ([]AdditionalCostTotalRow, error) {
	if len(activityIds) == 0 {
		return nil, nil
	}
	var rows []AdditionalCostTotalRow
	err := r.db.WithContext(ctx).Raw(`
SELECT ac.activity_id AS activity_id, SUM(ac.cost) AS total
FROM additional_costs ac
WHERE ac.activity_id IN ? AND ac.deleted_at IS NULL
GROUP BY ac.activity_id`, activityIds).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
