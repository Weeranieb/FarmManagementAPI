package excel_export

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
)

// PondExportData holds everything needed to fill one sheet (one pond) in the export workbook.
type PondExportData struct {
	PondName               string
	Logs                   []*model.DailyLog // full cycle, ordered by feed_date ASC
	FreshFeedCollectionId  *int
	PelletFeedCollectionId *int
	FreshFCR               *decimal.Decimal
	PelletFCR              *decimal.Decimal
	TotalFish              int
	StartDate              time.Time
}
