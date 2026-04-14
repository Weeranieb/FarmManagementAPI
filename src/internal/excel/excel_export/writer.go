package excel_export

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/xuri/excelize/v2"
)

const templateSheet = "Sheet1"

// Layout constants shared by both template variants.
const (
	monthHeaderRow = 1  // 1-indexed Excel row for the month header ("Feb-69")
	dayStartRow    = 5  // 1-indexed Excel row where day 1 data begins
	maxDays        = 31 // rows 5..35

	// Metadata positions (1-indexed rows) — mirror the template structure.
	metaFreshFeedIDRow  = 50
	metaPelletFeedIDRow = 51
	metaFreshFCRRow     = 43
	metaPelletFCRRow    = 44
	metaFeedIDCol       = "B" // column for feed collection ID values
	metaFCRValueCol     = "E" // column for FCR value
)

// templateLayout describes column offsets within one month block.
type templateLayout struct {
	blockDataCols int // number of data columns per block (excl. gap)
	stride        int // total columns per block including gap
	// Offsets within a block (0-based from block start column).
	freshMorning  int
	freshEvening  int
	pelletMorning int
	pelletEvening int
	deathFish     int
	touristCatch  int // -1 when not present
}

var noFishingLayout = templateLayout{
	blockDataCols: 5,
	stride:        7,
	freshMorning:  0,
	freshEvening:  1,
	pelletMorning: 2,
	pelletEvening: 3,
	deathFish:     4,
	touristCatch:  -1,
}

var fishingLayout = templateLayout{
	blockDataCols: 6,
	stride:        8,
	freshMorning:  0,
	freshEvening:  1,
	pelletMorning: 2,
	pelletEvening: 3,
	deathFish:     4,
	touristCatch:  5,
}

// monthKey groups daily logs by calendar month.
type monthKey struct {
	Year  int
	Month time.Month
}

// GenerateExport opens the template, creates one sheet per pond filled with daily-log data,
// and returns the resulting workbook as bytes.
func GenerateExport(templatePath string, ponds []PondExportData, withFishing bool) ([]byte, error) {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("open template %q: %w", templatePath, err)
	}
	defer func() { _ = f.Close() }()

	layout := noFishingLayout
	if withFishing {
		layout = fishingLayout
	}

	for i, pond := range ponds {
		sheetName := pond.PondName
		if err := cloneTemplateSheet(f, sheetName, i); err != nil {
			return nil, fmt.Errorf("clone sheet for %q: %w", sheetName, err)
		}
		if err := fillSheet(f, sheetName, pond, layout); err != nil {
			return nil, fmt.Errorf("fill sheet %q: %w", sheetName, err)
		}
	}

	// Remove the original template sheet (only if we created at least one pond sheet).
	if len(ponds) > 0 {
		idx, err := f.GetSheetIndex(templateSheet)
		if err == nil && idx >= 0 {
			_ = f.DeleteSheet(templateSheet)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write workbook: %w", err)
	}
	return buf.Bytes(), nil
}

// cloneTemplateSheet duplicates the template sheet and renames the copy.
func cloneTemplateSheet(f *excelize.File, name string, idx int) error {
	srcIdx, err := f.GetSheetIndex(templateSheet)
	if err != nil {
		return fmt.Errorf("template sheet lookup: %w", err)
	}
	if srcIdx < 0 {
		return fmt.Errorf("template sheet %q not found", templateSheet)
	}
	newIdx, err := f.NewSheet(name)
	if err != nil {
		return err
	}
	if err := f.CopySheet(srcIdx, newIdx); err != nil {
		return err
	}
	return nil
}

func fillSheet(f *excelize.File, sheet string, pond PondExportData, layout templateLayout) error {
	grouped := groupLogsByMonth(pond.Logs)
	months := sortedMonthKeys(grouped)

	for mi, mk := range months {
		blockStartCol := 2 + (mi * layout.stride) // 2 = column B (1-indexed)
		header := FormatMonthHeaderBE(mk.Year, mk.Month)
		cell, _ := excelize.CoordinatesToCellName(blockStartCol, monthHeaderRow)
		_ = f.SetCellValue(sheet, cell, header)

		for _, log := range grouped[mk] {
			day := log.FeedDate.Day()
			if day < 1 || day > maxDays {
				continue
			}
			row := dayStartRow + day - 1 // day 1 -> row 5
			writeLogRow(f, sheet, row, blockStartCol, layout, log)
		}
	}

	writeMetadata(f, sheet, pond)
	return nil
}

func writeLogRow(f *excelize.File, sheet string, row, blockStartCol int, layout templateLayout, log *model.DailyLog) {
	setDecimal(f, sheet, row, blockStartCol+layout.freshMorning, log.FreshMorning)
	setDecimal(f, sheet, row, blockStartCol+layout.freshEvening, log.FreshEvening)
	setDecimal(f, sheet, row, blockStartCol+layout.pelletMorning, log.PelletMorning)
	setDecimal(f, sheet, row, blockStartCol+layout.pelletEvening, log.PelletEvening)
	setInt(f, sheet, row, blockStartCol+layout.deathFish, log.DeathFishCount)
	if layout.touristCatch >= 0 && log.TouristCatchCount != nil {
		setInt(f, sheet, row, blockStartCol+layout.touristCatch, *log.TouristCatchCount)
	}
}

func writeMetadata(f *excelize.File, sheet string, pond PondExportData) {
	if pond.FreshFeedCollectionId != nil {
		cell, _ := excelize.CoordinatesToCellName(colIndex(metaFeedIDCol), metaFreshFeedIDRow)
		_ = f.SetCellValue(sheet, cell, *pond.FreshFeedCollectionId)
	}
	if pond.PelletFeedCollectionId != nil {
		cell, _ := excelize.CoordinatesToCellName(colIndex(metaFeedIDCol), metaPelletFeedIDRow)
		_ = f.SetCellValue(sheet, cell, *pond.PelletFeedCollectionId)
	}
	if pond.FreshFCR != nil {
		cell, _ := excelize.CoordinatesToCellName(colIndex(metaFCRValueCol), metaFreshFCRRow)
		v, _ := pond.FreshFCR.Float64()
		_ = f.SetCellValue(sheet, cell, v)
	}
	if pond.PelletFCR != nil {
		cell, _ := excelize.CoordinatesToCellName(colIndex(metaFCRValueCol), metaPelletFCRRow)
		v, _ := pond.PelletFCR.Float64()
		_ = f.SetCellValue(sheet, cell, v)
	}
}

// --- helpers ---

func groupLogsByMonth(logs []*model.DailyLog) map[monthKey][]*model.DailyLog {
	m := make(map[monthKey][]*model.DailyLog)
	for _, l := range logs {
		mk := monthKey{Year: l.FeedDate.Year(), Month: l.FeedDate.Month()}
		m[mk] = append(m[mk], l)
	}
	return m
}

func sortedMonthKeys(m map[monthKey][]*model.DailyLog) []monthKey {
	keys := make([]monthKey, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Year != keys[j].Year {
			return keys[i].Year < keys[j].Year
		}
		return keys[i].Month < keys[j].Month
	})
	return keys
}

func setDecimal(f *excelize.File, sheet string, row, col int, d decimal.Decimal) {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	v, _ := d.Float64()
	_ = f.SetCellValue(sheet, cell, v)
}

func setInt(f *excelize.File, sheet string, row, col, v int) {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	_ = f.SetCellValue(sheet, cell, v)
}

// colIndex converts "A" -> 1, "B" -> 2, etc. Only single-letter columns supported for metadata.
func colIndex(letter string) int {
	if len(letter) != 1 {
		return 1
	}
	return int(letter[0]-'A') + 1
}
