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
type ActivityListRow struct {
	Id           int             `gorm:"column:id"`
	Mode         string          `gorm:"column:mode"`
	ActivityDate time.Time       `gorm:"column:activity_date"`
	FishType     string          `gorm:"column:fish_type"`
	Amount       int             `gorm:"column:amount"`
	PricePerUnit decimal.Decimal `gorm:"column:price_per_unit"`
	MerchantName *string         `gorm:"column:merchant_name"`
	ToPondName   *string         `gorm:"column:to_pond_name"`
}

// SellTotalRow is the grouped sum projection for sell totals keyed by activity id.
type SellTotalRow struct {
	SellId int             `gorm:"column:sell_id"`
	Total  decimal.Decimal `gorm:"column:total"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=ActivityRepository --output=./mocks --outpkg=mocks --filename=activity_repository.go --structname=MockActivityRepository --with-expecter=false
type ActivityRepository interface {
	WithTx(tx *gorm.DB) ActivityRepository
	Create(ctx context.Context, activity *model.Activity) error
	ListByPondID(ctx context.Context, pondId int) ([]ActivityListRow, error)
	SumSellDetailsByActivityIDs(ctx context.Context, activityIds []int) ([]SellTotalRow, error)
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

const activityListByPondQuery = `
SELECT
  a.id AS id,
  a.mode AS mode,
  a.activity_date AS activity_date,
  a.fish_type AS fish_type,
  a.amount AS amount,
  a.price_per_unit AS price_per_unit,
  m.name AS merchant_name,
  tp.name AS to_pond_name
FROM activities a
INNER JOIN active_ponds ap ON a.active_pond_id = ap.id AND ap.deleted_at IS NULL
LEFT JOIN merchants m ON a.merchant_id = m.id AND m.deleted_at IS NULL
LEFT JOIN active_ponds tap ON a.to_active_pond_id = tap.id AND tap.deleted_at IS NULL
LEFT JOIN ponds tp ON tap.pond_id = tp.id AND tp.deleted_at IS NULL
WHERE ap.pond_id = ? AND a.deleted_at IS NULL
ORDER BY a.activity_date DESC, a.id DESC`

func (r *activityRepository) ListByPondID(ctx context.Context, pondId int) ([]ActivityListRow, error) {
	var rows []ActivityListRow
	err := r.db.WithContext(ctx).Raw(activityListByPondQuery, pondId).Scan(&rows).Error
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
SELECT sd.sell_id AS sell_id, SUM(sd.weight * sd.price_per_unit) AS total
FROM sell_details sd
WHERE sd.sell_id IN ? AND sd.deleted_at IS NULL
GROUP BY sd.sell_id`, activityIds).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
