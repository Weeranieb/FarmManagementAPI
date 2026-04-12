package excel_dailylog

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

const scanMaxRow = 50
const headerBottomRow = 5
const dayScanMaxRow = 40

const (
	headerThaiMorning    = "เช้า"
	headerThaiEvening    = "เย็น"
	headerEnglishMorning = "morning"
	headerEnglishEvening = "evening"
	headerThaiFresh      = "เหยื่อ"
	headerEnglishFresh   = "fresh"
	headerThaiPellet     = "อาหาร"
	headerThaiPelletAlt  = "เม็ด"
	headerEnglishPellet  = "pellet"
)

func todayUTC(ref time.Time) time.Time {
	if ref.IsZero() {
		ref = time.Now().UTC()
	}
	y, m, d := ref.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func maxRowWidth(rows [][]string) int {
	max := 0
	// Scan a bounded header/data window; templates are expected to keep all month blocks near the top.
	n := min(len(rows), scanMaxRow)
	for r := range n {
		if len(rows[r]) > max {
			max = len(rows[r])
		}
	}
	return max
}

func findBlockStarts(row0 []string) ([]int, error) {
	if len(row0) == 0 {
		return nil, fmt.Errorf("empty first row")
	}
	var starts []int
	for c := range row0 {
		v := strings.TrimSpace(row0[c])
		if monthHeaderRe.MatchString(v) {
			starts = append(starts, c)
		}
	}
	if len(starts) == 0 {
		return nil, fmt.Errorf("no month headers (Mon-yy) found on row 1")
	}
	return starts, nil
}

func blockEndCol(starts []int, idx int, maxCol int) int {
	if idx+1 < len(starts) {
		return starts[idx+1] - 1
	}
	return maxCol - 1
}

type columnMap struct {
	freshMorning      int
	freshEvening      int
	pelletMorning     int
	pelletEvening     int
	deathFishCount    int
	touristCatchCount int
	avgBodyWeight     int
	fishCount         int
}

func newColumnMap() columnMap {
	return columnMap{
		freshMorning:      -1,
		freshEvening:      -1,
		pelletMorning:     -1,
		pelletEvening:     -1,
		deathFishCount:    -1,
		touristCatchCount: -1,
		avgBodyWeight:     -1,
		fishCount:         -1,
	}
}

func mapBlockColumns(rows [][]string, start, end int) (columnMap, error) {
	cm := newColumnMap()
	top, bot := 0, headerBottomRow
	if len(rows)-1 < bot {
		bot = len(rows) - 1
	}
	containsAny := func(text string, tokens ...string) bool {
		for _, token := range tokens {
			if strings.Contains(text, token) {
				return true
			}
		}
		return false
	}
	withPreviousGroupIfNeeded := func(c int, normalizedHeader string) string {
		hasSession := containsAny(normalizedHeader,
			headerThaiMorning, headerThaiEvening, headerEnglishMorning, headerEnglishEvening,
		)
		hasFeedGroup := containsAny(normalizedHeader,
			headerThaiFresh, headerEnglishFresh, headerThaiPellet, headerThaiPelletAlt, headerEnglishPellet,
		)
		if !hasSession || hasFeedGroup || c <= start {
			return normalizedHeader
		}
		previousHeader := normalizeHeader(compositeColHeader(rows, c-1, top, bot))
		if previousHeader == "" {
			return normalizedHeader
		}
		return normalizeHeader(previousHeader + " " + normalizedHeader)
	}
	for c := start; c <= end; c++ {
		comp := compositeColHeader(rows, c, top, bot)
		normalizedHeader := withPreviousGroupIfNeeded(c, normalizeHeader(comp))
		if normalizedHeader == "" {
			continue
		}
		if headerLooksLikeFreshMorning(normalizedHeader) && cm.freshMorning < 0 {
			cm.freshMorning = c
			continue
		}
		if headerLooksLikeFreshEvening(normalizedHeader) && cm.freshEvening < 0 {
			cm.freshEvening = c
			continue
		}
		if headerLooksLikePelletMorning(normalizedHeader) && cm.pelletMorning < 0 {
			cm.pelletMorning = c
			continue
		}
		if headerLooksLikePelletEvening(normalizedHeader) && cm.pelletEvening < 0 {
			cm.pelletEvening = c
			continue
		}
		if headerLooksLikeDeath(normalizedHeader) && cm.deathFishCount < 0 {
			cm.deathFishCount = c
			continue
		}
		if headerLooksLikeAvgBodyWeight(normalizedHeader) && cm.avgBodyWeight < 0 {
			cm.avgBodyWeight = c
			continue
		}
		if headerLooksLikeFishCount(normalizedHeader) && cm.fishCount < 0 {
			cm.fishCount = c
			continue
		}
	}
	if cm.freshMorning < 0 || cm.freshEvening < 0 || cm.pelletMorning < 0 || cm.pelletEvening < 0 {
		return cm, fmt.Errorf("missing bait/feed morning/evening headers")
	}
	if cm.deathFishCount >= 0 && cm.deathFishCount+1 <= end {
		nhAdj := normalizeHeader(compositeColHeader(rows, cm.deathFishCount+1, top, bot))
		if headerLooksLikeTouristCatch(nhAdj) {
			cm.touristCatchCount = cm.deathFishCount + 1
		}
	}
	if cm.touristCatchCount < 0 {
		for c := start; c <= end; c++ {
			nh := normalizeHeader(compositeColHeader(rows, c, top, bot))
			if headerLooksLikeTouristCatch(nh) {
				cm.touristCatchCount = c
				break
			}
		}
	}
	return cm, nil
}

func findDayStartRow(rows [][]string) (int, error) {
	lim := min(len(rows), dayScanMaxRow)
	for r := range lim {
		a := cellStr(rows, r, 0)
		if isSummaryDayLabel(a) {
			continue
		}
		day, err := parseDayOfMonth(a)
		if err == nil && day == 1 {
			return r, nil
		}
	}
	return 0, fmt.Errorf("no day-1 row in column A (scanned first %d rows)", lim)
}

func parseSheetFeedIDs(rows [][]string) (freshID, pelletID *int) {
	parseCellID := func(row, col int) *int {
		raw := strings.TrimSpace(cellStr(rows, row, col))
		if raw == "" || strings.HasPrefix(raw, "#") {
			return nil
		}
		id, err := strconv.Atoi(raw)
		if err != nil {
			return nil
		}
		return &id
	}

	// Hardcoded metadata positions from current template:
	// B42 = fresh feed collection ID, B43 = pellet feed collection ID.
	return parseCellID(41, 1), parseCellID(42, 1)
}

// rowHasAnySignal is true when the row has feed, deaths, or (if the template has that column) tourist catch.
// Optional weight/fish-count columns alone do not count as signal.
func rowHasAnySignal(e ExtractedDailyLogRow, touristColPresent bool) bool {
	if anyNonZeroFeed(e) {
		return true // Any non-zero fresh or pellet (morning or evening).
	}
	if e.DeathFishCount != 0 {
		return true // Non-zero death count.
	}
	if touristColPresent && e.TouristCatchCount != nil && *e.TouristCatchCount != 0 {
		return true // Sheet has tourist column and catch count is present and non-zero.
	}
	return false // No feed, death, or tourist signal; skip row (e.g. weight-only or all zeros).
}

func anyNonZeroFeed(e ExtractedDailyLogRow) bool {
	return lo.SomeBy([]decimal.Decimal{
		e.FreshMorning,
		e.FreshEvening,
		e.PelletMorning,
		e.PelletEvening,
	}, func(d decimal.Decimal) bool { return !d.IsZero() })
}

func hasFeed(e ExtractedDailyLogRow) bool {
	return anyNonZeroFeed(e)
}

func extractBlock(
	rows [][]string,
	year int, month time.Month,
	cm columnMap,
	dayStart int,
	today time.Time,
) ([]ExtractedDailyLogRow, error) {
	touristPresent := cm.touristCatchCount >= 0
	var out []ExtractedDailyLogRow
	for r := dayStart; r < len(rows); r++ {
		a := cellStr(rows, r, 0)
		if isSummaryDayLabel(a) {
			break
		}
		day, err := parseDayOfMonth(a)
		if err != nil {
			break
		}
		feedDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		if feedDate.Month() != month || feedDate.Day() != day {
			continue
		}
		if feedDate.After(today) {
			continue
		}

		freshMorning, err := parseDecimalCell(cellStr(rows, r, cm.freshMorning))
		if err != nil {
			return nil, err
		}
		freshEvening, err := parseDecimalCell(cellStr(rows, r, cm.freshEvening))
		if err != nil {
			return nil, err
		}
		pelletMorning, err := parseDecimalCell(cellStr(rows, r, cm.pelletMorning))
		if err != nil {
			return nil, err
		}
		pelletEvening, err := parseDecimalCell(cellStr(rows, r, cm.pelletEvening))
		if err != nil {
			return nil, err
		}
		deaths := 0
		if cm.deathFishCount >= 0 {
			deaths, err = parseIntCell(cellStr(rows, r, cm.deathFishCount))
			if err != nil {
				return nil, err
			}
		}
		var tourist *int
		if cm.touristCatchCount >= 0 {
			tourist, err = parseOptionalIntCell(cellStr(rows, r, cm.touristCatchCount))
			if err != nil {
				return nil, err
			}
		}
		var weight *decimal.Decimal
		if cm.avgBodyWeight >= 0 {
			weight, err = parseOptionalDecimalCell(cellStr(rows, r, cm.avgBodyWeight))
			if err != nil {
				return nil, err
			}
		}
		var fish *int
		if cm.fishCount >= 0 {
			fish, err = parseOptionalIntCell(cellStr(rows, r, cm.fishCount))
			if err != nil {
				return nil, err
			}
		}
		e := ExtractedDailyLogRow{
			FeedDate:          feedDate,
			FreshMorning:      freshMorning,
			FreshEvening:      freshEvening,
			PelletMorning:     pelletMorning,
			PelletEvening:     pelletEvening,
			DeathFishCount:    deaths,
			TouristCatchCount: tourist,
			AvgBodyWeight:     weight,
			FishCount:         fish,
		}
		if !rowHasAnySignal(e, touristPresent) {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}

func blockMonthAfterToday(year int, month time.Month, today time.Time) bool {
	ty, tm, _ := today.Date()
	if year > ty {
		return true
	}
	if year == ty && month > tm {
		return true
	}
	return false
}

// ParseSheet parses one worksheet using time.Now (UTC date) as "today" for future cutoffs.
func ParseSheet(f *excelize.File, sheetName string) (*ParsedSheet, error) {
	return ParseSheetAt(f, sheetName, time.Time{})
}

// ParseSheetAt is like ParseSheet but uses ref (UTC date only) as "today"; zero means now.
func ParseSheetAt(f *excelize.File, sheetName string, ref time.Time) (*ParsedSheet, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("get rows %q: %w", sheetName, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty sheet %q", sheetName)
	}
	today := todayUTC(ref)
	maxCol := maxRowWidth(rows)
	if maxCol < 2 {
		return nil, fmt.Errorf("sheet too narrow")
	}
	starts, err := findBlockStarts(rows[0])
	if err != nil {
		return nil, err
	}
	dayStart, err := findDayStartRow(rows)
	if err != nil {
		return nil, err
	}
	var all []ExtractedDailyLogRow
	for i, start := range starts {
		end := blockEndCol(starts, i, maxCol)
		if end < start {
			return nil, fmt.Errorf("invalid block bounds at column %d", start)
		}
		headerVal := strings.TrimSpace(cellStr(rows, 0, start))
		gy, gm, err := parseEnglishMonthBEHeader(headerVal)
		if err != nil {
			break
		}
		if blockMonthAfterToday(gy, gm, today) {
			break
		}
		cm, err := mapBlockColumns(rows, start, end)
		if err != nil {
			return nil, fmt.Errorf("block cols %d-%d: %w", start, end, err)
		}
		blockRows, err := extractBlock(rows, gy, gm, cm, dayStart, today)
		if err != nil {
			return nil, err
		}
		hasFeedEmitted := slices.ContainsFunc(blockRows, hasFeed)
		if i >= 1 && !hasFeedEmitted {
			break
		}
		all = append(all, blockRows...)
	}
	freshFeedCollectionID, pelletFeedCollectionID := parseSheetFeedIDs(rows)
	return &ParsedSheet{
		PondName:               sheetName,
		Rows:                   all,
		FreshFeedCollectionId:  freshFeedCollectionID,
		PelletFeedCollectionId: pelletFeedCollectionID,
	}, nil
}

// ParseFile opens path and parses the given sheet (ref semantics match ParseSheetAt).
func ParseFile(path, sheetName string, ref time.Time) (*ParsedSheet, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	return ParseSheetAt(f, sheetName, ref)
}

// parseAllSheets is the shared implementation for all-sheets parsing.
// Sheets that are not valid daily-log templates (blank tabs, helpers, etc.) are skipped; no error is returned for those.
func parseAllSheets(f *excelize.File, ref time.Time) (map[string]*ParsedSheet, error) {
	out := make(map[string]*ParsedSheet)
	for _, name := range f.GetSheetList() {
		ps, err := ParseSheetAt(f, name, ref)
		if err != nil {
			continue
		}
		out[name] = ps
	}
	return out, nil
}

// ParseFileAllSheets parses every sheet from a file on disk.
// Unreadable or non-template sheets are omitted from the map; the returned error is only for open/read failures.
func ParseFileAllSheets(path string, ref time.Time) (map[string]*ParsedSheet, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	return parseAllSheets(f, ref)
}

// ParseReaderAllSheets parses every sheet from an io.Reader (e.g. an in-memory upload).
// Non-template sheets are skipped; errors are only from opening the workbook.
func ParseReaderAllSheets(r io.Reader, ref time.Time) (map[string]*ParsedSheet, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("open reader: %w", err)
	}
	defer func() { _ = f.Close() }()
	return parseAllSheets(f, ref)
}
