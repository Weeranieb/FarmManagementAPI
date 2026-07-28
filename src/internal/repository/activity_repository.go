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
	// TotalFishCount is the summed sell_details.fish_count — a sell activity
	// row leaves `amount` unset (every head count lives on the detail lines),
	// so this is the only source for "how many fish were sold".
	TotalFishCount int `gorm:"column:total_fish_count"`
}

// ActivityFeedCursor points at the last row a caller has already seen. It is
// the feed's sort key — (activity_date, id) — and nothing else, so paging is
// stable no matter how many activities are recorded mid-scroll.
type ActivityFeedCursor struct {
	ActivityDate time.Time
	Id           int
}

// ActivityFeedRow is the flat join projection returned by ListRecentByClientID.
// One row per activity (moves are NOT duplicated per pond perspective here —
// the feed shows "source → destination" in a single entry).
type ActivityFeedRow struct {
	Id            int       `gorm:"column:id"`
	Mode          string    `gorm:"column:mode"`
	ActivityDate  time.Time `gorm:"column:activity_date"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	CreatedBy     string    `gorm:"column:created_by"`
	CreatedByName string    `gorm:"column:created_by_name"`
	// PondId is the ponds.id of the source pond — the identifier the pond
	// screens key off (not the active_ponds.id the activity itself points at),
	// so a feed row can link straight to its pond.
	PondId       int             `gorm:"column:pond_id"`
	PondName     string          `gorm:"column:pond_name"`
	FarmName     string          `gorm:"column:farm_name"`
	ToPondId     *int            `gorm:"column:to_pond_id"`
	ToPondName   *string         `gorm:"column:to_pond_name"`
	FishType     string          `gorm:"column:fish_type"`
	Amount       int             `gorm:"column:amount"`
	FishWeight   decimal.Decimal `gorm:"column:fish_weight"`
	FishUnit     string          `gorm:"column:fish_unit"`
	PricePerUnit decimal.Decimal `gorm:"column:price_per_unit"`
	MerchantName *string         `gorm:"column:merchant_name"`
}

// SellDetailRow is one size-grade line of a sale, joined to its grade name.
// Returned by ListSellDetailsByActivityID for the per-size breakdown.
type SellDetailRow struct {
	FishSizeGradeId int    `gorm:"column:fish_size_grade_id"`
	SizeName        string `gorm:"column:size_name"`
	// FishCount is nullable on legacy rows written before the column became
	// required, so it stays a pointer rather than silently reporting 0 fish.
	FishCount    *int            `gorm:"column:fish_count"`
	Weight       decimal.Decimal `gorm:"column:weight"`
	PricePerUnit decimal.Decimal `gorm:"column:price_per_unit"`
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
	ListRecentByClientID(ctx context.Context, clientId *int, limit int, before *ActivityFeedCursor) ([]ActivityFeedRow, error)
	SumSellDetailsByActivityIDs(ctx context.Context, activityIds []int) ([]SellTotalRow, error)
	SumAdditionalCostsByActivityIDs(ctx context.Context, activityIds []int) ([]AdditionalCostTotalRow, error)
	ListSellDetailsByActivityID(ctx context.Context, activityId int, clientId *int) ([]SellDetailRow, error)
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
  fp.id AS pond_id,
  fp.name AS pond_name,
  f.name AS farm_name,
  tp.id AS to_pond_id,
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

func (r *activityRepository) ListRecentByClientID(ctx context.Context, clientId *int, limit int, before *ActivityFeedCursor) ([]ActivityFeedRow, error) {
	query := activityFeedQuery
	args := make([]interface{}, 0, 4)
	if clientId != nil {
		query += ` AND f.client_id = ?`
		args = append(args, *clientId)
	}
	// Row-wise comparison against the feed's own sort key. A plain OFFSET would
	// shift under the caller: this feed is append-heavy (workers record sells
	// while an owner scrolls history), and every insert lands at the top, so
	// page 2 would repeat rows page 1 already showed.
	if before != nil {
		query += ` AND (a.activity_date, a.id) < (?, ?)`
		args = append(args, before.ActivityDate, before.Id)
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
SELECT sd.sell_id AS sell_id, SUM(sd.weight * sd.price_per_unit) AS total, SUM(sd.weight) AS total_weight, COALESCE(SUM(sd.fish_count), 0) AS total_fish_count
FROM sell_details sd
WHERE sd.sell_id IN ? AND sd.deleted_at IS NULL
GROUP BY sd.sell_id`, activityIds).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Lists one sale's size-grade lines, smallest grade first (fish_size_grades
// carries an explicit sort_index for exactly this).
//
// Client scoping is enforced inside the query rather than by a separate
// ownership check: the activity → active_ponds → ponds → farms chain is joined
// so another client's sale cannot be read even with a guessed id. That also
// collapses "no such activity", "not a sell" and "not yours" into the same
// empty result, so the endpoint is not an existence oracle.
const sellDetailsByActivityQuery = `
SELECT
  sd.fish_size_grade_id AS fish_size_grade_id,
  g.name AS size_name,
  sd.fish_count AS fish_count,
  sd.weight AS weight,
  sd.price_per_unit AS price_per_unit
FROM sell_details sd
INNER JOIN activities a ON sd.sell_id = a.id AND a.deleted_at IS NULL
INNER JOIN active_ponds ap ON a.active_pond_id = ap.id AND ap.deleted_at IS NULL
INNER JOIN ponds p ON ap.pond_id = p.id AND p.deleted_at IS NULL
INNER JOIN farms f ON p.farm_id = f.id AND f.deleted_at IS NULL
LEFT JOIN fish_size_grades g ON sd.fish_size_grade_id = g.id AND g.deleted_at IS NULL
WHERE sd.sell_id = ? AND sd.deleted_at IS NULL`

func (r *activityRepository) ListSellDetailsByActivityID(ctx context.Context, activityId int, clientId *int) ([]SellDetailRow, error) {
	query := sellDetailsByActivityQuery
	args := []interface{}{activityId}
	if clientId != nil {
		query += ` AND f.client_id = ?`
		args = append(args, *clientId)
	}
	query += ` ORDER BY g.sort_index ASC, sd.id ASC`

	var rows []SellDetailRow
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
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
