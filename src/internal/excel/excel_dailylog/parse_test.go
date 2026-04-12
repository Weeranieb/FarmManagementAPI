package excel_dailylog

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

// febBlockHeaders writes Thai headers for one month block starting at 1-based Excel col startCol.
func febBlockHeaders(t *testing.T, f *excelize.File, sheet string, startCol int, withTourist bool) {
	t.Helper()
	col := func(offset int) string {
		name, err := excelize.ColumnNumberToName(startCol + offset)
		require.NoError(t, err)
		return name
	}
	require.NoError(t, f.SetCellValue(sheet, col(0)+"2", "เหยื่อ"))
	require.NoError(t, f.SetCellValue(sheet, col(1)+"2", "เหยื่อ"))
	require.NoError(t, f.SetCellValue(sheet, col(2)+"2", "อาหาร"))
	require.NoError(t, f.SetCellValue(sheet, col(3)+"2", "อาหาร"))
	require.NoError(t, f.SetCellValue(sheet, col(4)+"2", "ตาย"))
	if withTourist {
		require.NoError(t, f.SetCellValue(sheet, col(5)+"2", "ตกปลา"))
	}
	require.NoError(t, f.SetCellValue(sheet, col(0)+"3", "เช้า"))
	require.NoError(t, f.SetCellValue(sheet, col(1)+"3", "เย็น"))
	require.NoError(t, f.SetCellValue(sheet, col(2)+"3", "เช้า"))
	require.NoError(t, f.SetCellValue(sheet, col(3)+"3", "เย็น"))
}

func TestParseEnglishMonthBEHeader(t *testing.T) {
	y, m, err := parseEnglishMonthBEHeader("Feb-69")
	require.NoError(t, err)
	require.Equal(t, 2026, y)
	require.Equal(t, time.February, m)
}

func TestParseSheet_BECalendarAndFeed(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "B1", "Feb-69"))
	febBlockHeaders(t, f, sheet, 2, false)
	require.NoError(t, f.SetCellValue(sheet, "A5", "1"))
	require.NoError(t, f.SetCellValue(sheet, "D5", "40"))
	require.NoError(t, f.SetCellValue(sheet, "E5", "5"))

	ref := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)
	ps, err := ParseSheetAt(f, sheet, ref)
	require.NoError(t, err)
	require.Equal(t, sheet, ps.PondName)
	require.Len(t, ps.Rows, 1)
	require.True(t, ps.Rows[0].FeedDate.Equal(time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)))
	require.True(t, ps.Rows[0].PelletMorning.Equal(decimal.NewFromInt(40)))
	require.True(t, ps.Rows[0].PelletEvening.Equal(decimal.NewFromInt(5)))
	require.Nil(t, ps.Rows[0].TouristCatchCount)
}

func TestParseSheet_TouristColumn(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "B1", "Feb-69"))
	febBlockHeaders(t, f, sheet, 2, true)
	require.NoError(t, f.SetCellValue(sheet, "A5", "1"))
	require.NoError(t, f.SetCellValue(sheet, "D5", "10"))
	require.NoError(t, f.SetCellValue(sheet, "G5", "3"))

	ref := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)
	ps, err := ParseSheetAt(f, sheet, ref)
	require.NoError(t, err)
	require.Len(t, ps.Rows, 1)
	require.True(t, ps.Rows[0].FeedDate.Equal(time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)))
	require.NotNil(t, ps.Rows[0].TouristCatchCount)
	require.Equal(t, 3, *ps.Rows[0].TouristCatchCount)
}

func TestParseSheet_FutureRowOmitted(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "B1", "Feb-69"))
	febBlockHeaders(t, f, sheet, 2, false)
	require.NoError(t, f.SetCellValue(sheet, "A5", "1"))
	require.NoError(t, f.SetCellValue(sheet, "D5", "1"))
	require.NoError(t, f.SetCellValue(sheet, "A6", "2"))
	require.NoError(t, f.SetCellValue(sheet, "D6", "2"))
	require.NoError(t, f.SetCellValue(sheet, "A7", "25"))
	require.NoError(t, f.SetCellValue(sheet, "D7", "99"))

	ref := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	ps, err := ParseSheetAt(f, sheet, ref)
	require.NoError(t, err)
	require.Len(t, ps.Rows, 2)
}

func TestParseSheet_FutureMonthBlockSkipped(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "B1", "Feb-69"))
	require.NoError(t, f.SetCellValue(sheet, "I1", "Apr-69"))
	febBlockHeaders(t, f, sheet, 2, false)
	febBlockHeaders(t, f, sheet, 9, false)
	require.NoError(t, f.SetCellValue(sheet, "A5", "1"))
	require.NoError(t, f.SetCellValue(sheet, "D5", "5"))
	require.NoError(t, f.SetCellValue(sheet, "J5", "9"))

	ref := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	ps, err := ParseSheetAt(f, sheet, ref)
	require.NoError(t, err)
	require.Len(t, ps.Rows, 1)
	require.Equal(t, time.February, ps.Rows[0].FeedDate.Month())
}

func TestParseSheet_FirstBlockEmptySecondHasFeed(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "B1", "Feb-69"))
	require.NoError(t, f.SetCellValue(sheet, "I1", "Mar-69"))
	febBlockHeaders(t, f, sheet, 2, false)
	febBlockHeaders(t, f, sheet, 9, false)
	require.NoError(t, f.SetCellValue(sheet, "A5", "1"))
	require.NoError(t, f.SetCellValue(sheet, "K5", "7"))

	ref := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	ps, err := ParseSheetAt(f, sheet, ref)
	require.NoError(t, err)
	require.Len(t, ps.Rows, 1)
	require.Equal(t, time.March, ps.Rows[0].FeedDate.Month())
	require.True(t, ps.Rows[0].PelletMorning.Equal(decimal.NewFromInt(7)))
}

func TestParseSheet_SecondBlockNoFeedStopsRest(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	sheet := f.GetSheetName(0)
	require.NoError(t, f.SetCellValue(sheet, "B1", "Feb-69"))
	require.NoError(t, f.SetCellValue(sheet, "I1", "Mar-69"))
	require.NoError(t, f.SetCellValue(sheet, "P1", "Apr-69"))
	febBlockHeaders(t, f, sheet, 2, false)
	febBlockHeaders(t, f, sheet, 9, false)
	febBlockHeaders(t, f, sheet, 16, false)
	require.NoError(t, f.SetCellValue(sheet, "A5", "1"))
	require.NoError(t, f.SetCellValue(sheet, "D5", "3"))
	require.NoError(t, f.SetCellValue(sheet, "F5", "2"))
	require.NoError(t, f.SetCellValue(sheet, "I5", "0"))
	require.NoError(t, f.SetCellValue(sheet, "J5", "0"))
	require.NoError(t, f.SetCellValue(sheet, "K5", "0"))
	require.NoError(t, f.SetCellValue(sheet, "L5", "0"))
	require.NoError(t, f.SetCellValue(sheet, "M5", "5"))
	require.NoError(t, f.SetCellValue(sheet, "R5", "9"))

	ref := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	ps, err := ParseSheetAt(f, sheet, ref)
	require.NoError(t, err)
	require.Len(t, ps.Rows, 1)
	require.Equal(t, time.February, ps.Rows[0].FeedDate.Month())
}

func TestToDailyLog(t *testing.T) {
	e := ExtractedDailyLogRow{
		FeedDate:       time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		FreshMorning:   decimal.NewFromInt(1),
		PelletEvening:  decimal.NewFromInt(2),
		DeathFishCount: 1,
	}
	m := e.ToDailyLog(99, "alice")
	require.Equal(t, 99, m.ActivePondId)
	require.Equal(t, "alice", m.CreatedBy)
	require.Equal(t, "alice", m.UpdatedBy)
}

func TestParseFile_NoFishing(t *testing.T) {
	ref := time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC)
	ps, err := ParseFile("test_no_fishing.xlsx", "1 ซ้าย", ref)
	require.NoError(t, err)

	require.Equal(t, "1 ซ้าย", ps.PondName)
	require.NotNil(t, ps.FreshFeedCollectionId)
	require.NotNil(t, ps.PelletFeedCollectionId)
	require.Equal(t, 1, *ps.FreshFeedCollectionId)
	require.Equal(t, 2, *ps.PelletFeedCollectionId)

	require.Len(t, ps.Rows, 3)

	day5 := ps.Rows[0]
	require.True(t, day5.FeedDate.Equal(time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)))
	require.True(t, day5.FreshMorning.Equal(decimal.NewFromInt(1)))
	require.True(t, day5.FreshEvening.Equal(decimal.NewFromInt(1)))
	require.True(t, day5.PelletMorning.Equal(decimal.Zero))
	require.True(t, day5.PelletEvening.Equal(decimal.Zero))
	require.Equal(t, 0, day5.DeathFishCount)
	require.Nil(t, day5.TouristCatchCount)

	day6 := ps.Rows[1]
	require.True(t, day6.FeedDate.Equal(time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)))
	require.True(t, day6.PelletMorning.Equal(decimal.NewFromInt(1)))
	require.Equal(t, 0, day6.DeathFishCount)
	require.Nil(t, day6.TouristCatchCount)

	day7 := ps.Rows[2]
	require.True(t, day7.FeedDate.Equal(time.Date(2026, 3, 7, 0, 0, 0, 0, time.UTC)))
	require.True(t, day7.PelletEvening.Equal(decimal.NewFromInt(1)))
	require.Equal(t, 0, day7.DeathFishCount)
	require.Nil(t, day7.TouristCatchCount)
}
