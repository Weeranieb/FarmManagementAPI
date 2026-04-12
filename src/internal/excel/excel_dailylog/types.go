package excel_dailylog

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
)

// ParsedSheet is the result of parsing one worksheet in the horizontal monthly template.
type ParsedSheet struct {
	PondName               string
	Rows                   []ExtractedDailyLogRow
	FreshFeedCollectionId  *int
	PelletFeedCollectionId *int
}

// ExtractedDailyLogRow is one logical day for one month block on the sheet.
// AvgBodyWeight and FishCount are optional Excel-only columns (not on model.DailyLog).
type ExtractedDailyLogRow struct {
	FeedDate          time.Time
	FreshMorning      decimal.Decimal
	FreshEvening      decimal.Decimal
	PelletMorning     decimal.Decimal
	PelletEvening     decimal.Decimal
	DeathFishCount    int
	TouristCatchCount *int
	AvgBodyWeight     *decimal.Decimal
	FishCount         *int
}

// ToDailyLog builds a GORM model row. Feed collection IDs live on active_ponds, not daily_logs.
func (e ExtractedDailyLogRow) ToDailyLog(activePondId int, createdBy string) model.DailyLog {
	return model.DailyLog{
		ActivePondId:      activePondId,
		FeedDate:          e.FeedDate,
		FreshMorning:      e.FreshMorning,
		FreshEvening:      e.FreshEvening,
		PelletMorning:     e.PelletMorning,
		PelletEvening:     e.PelletEvening,
		DeathFishCount:    e.DeathFishCount,
		TouristCatchCount: e.TouristCatchCount,
		BaseModel: model.BaseModel{
			CreatedBy: createdBy,
			UpdatedBy: createdBy,
		},
	}
}
