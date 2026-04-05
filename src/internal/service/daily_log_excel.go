package service

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/xuri/excelize/v2"
)

func normalizeHeader(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func findColumnIndex(headers []string, candidates []string) int {
	for i, h := range headers {
		nh := normalizeHeader(h)
		for _, c := range candidates {
			if nh == normalizeHeader(c) {
				return i
			}
		}
	}
	return -1
}

func parseExcelDateCell(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return excelize.ExcelDateToTime(f, false)
	}
	return time.Time{}, fmt.Errorf("invalid date: %s", s)
}

func parseDecimalCell(s string) (decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(s)
}

func parseIntCell(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// parseOptionalIntCell returns nil for empty cells; non-empty must parse as int.
func parseOptionalIntCell(s string) (*int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func maxWidth(rows [][]string, maxRows int) int {
	w := 0
	for i := 0; i < len(rows) && i < maxRows; i++ {
		if len(rows[i]) > w {
			w = len(rows[i])
		}
	}
	return w
}

func compositeHeader(rows [][]string, headerRows, col int) string {
	var parts []string
	for r := 0; r < headerRows && r < len(rows); r++ {
		if col < len(rows[r]) {
			parts = append(parts, rows[r][col])
		}
	}
	return strings.Join(parts, " ")
}

func headerLooksLikeDate(nh string) bool {
	return strings.Contains(nh, "วันที่") || nh == "date" || strings.Contains(nh, "date")
}

func headerLooksLikeFreshMorning(nh string) bool {
	return (strings.Contains(nh, "เหยื่อ") || strings.Contains(nh, "fresh")) &&
		(strings.Contains(nh, "เช้า") || strings.Contains(nh, "morning"))
}

func headerLooksLikeFreshEvening(nh string) bool {
	return (strings.Contains(nh, "เหยื่อ") || strings.Contains(nh, "fresh")) &&
		(strings.Contains(nh, "เย็น") || strings.Contains(nh, "evening"))
}

func headerLooksLikePelletMorning(nh string) bool {
	pellet := strings.Contains(nh, "อาหาร") || strings.Contains(nh, "เม็ด") || strings.Contains(nh, "pellet")
	session := strings.Contains(nh, "เช้า") || strings.Contains(nh, "morning")
	return pellet && session
}

func headerLooksLikePelletEvening(nh string) bool {
	pellet := strings.Contains(nh, "อาหาร") || strings.Contains(nh, "เม็ด") || strings.Contains(nh, "pellet")
	session := strings.Contains(nh, "เย็น") || strings.Contains(nh, "evening")
	return pellet && session
}

func headerLooksLikeDeath(nh string) bool {
	return strings.Contains(nh, "ตาย") && !strings.Contains(nh, "%")
}

func headerLooksLikeTouristCatch(nh string) bool {
	return strings.Contains(nh, "นักท่อง") || strings.Contains(nh, "tourist") ||
		strings.Contains(nh, "จับปลา") || strings.Contains(nh, "ตกปลา")
}

// parseDailyLogExcelFile reads an .xlsx, finds a sheet with a date column and feed columns
// (Thai template: เหยื่อ+เช้า/เย็น, อาหาร+เช้า/เย็น, optional ตาย / tourist columns).
// Falls back to a simple Date + Morning + Evening row (values map to pellet only).
func parseDailyLogExcelFile(filePath string, month string) ([]dto.DailyLogEntryInput, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("only .xlsx files are supported"))
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}
	defer func() { _ = f.Close() }()

	start, _, err := parseMonth(month)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}
	year, mth := start.Year(), int(start.Month())
	daysInMonth := time.Date(year, time.Month(mth)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for _, sheet := range f.GetSheetList() {
		if sheet == "" {
			continue
		}
		rows, err := f.GetRows(sheet)
		if err != nil || len(rows) < 2 {
			continue
		}

		for headerRows := 2; headerRows >= 1; headerRows-- {
			if len(rows) < headerRows+1 {
				continue
			}
			w := maxWidth(rows, headerRows)
			if w == 0 {
				continue
			}
			composite := make([]string, w)
			for c := 0; c < w; c++ {
				composite[c] = compositeHeader(rows, headerRows, c)
			}

			dateCol := -1
			freshM, freshE, pelletM, pelletE := -1, -1, -1, -1
			deathCol, touristCol := -1, -1
			for c := 0; c < w; c++ {
				nh := normalizeHeader(composite[c])
				if nh == "" {
					continue
				}
				if headerLooksLikeDate(nh) {
					dateCol = c
					continue
				}
				if headerLooksLikeFreshMorning(nh) {
					freshM = c
					continue
				}
				if headerLooksLikeFreshEvening(nh) {
					freshE = c
					continue
				}
				if headerLooksLikePelletMorning(nh) {
					pelletM = c
					continue
				}
				if headerLooksLikePelletEvening(nh) {
					pelletE = c
					continue
				}
				if headerLooksLikeDeath(nh) {
					deathCol = c
					continue
				}
				if headerLooksLikeTouristCatch(nh) {
					touristCol = c
					continue
				}
			}

			simpleMorning, simpleEvening := -1, -1
			if dateCol >= 0 && freshM < 0 && pelletM < 0 {
				h0 := rows[0]
				simpleMorning = findColumnIndex(h0, []string{"morning", "เช้า"})
				simpleEvening = findColumnIndex(h0, []string{"evening", "เย็น"})
			}

			if dateCol < 0 {
				continue
			}
			if freshM < 0 && freshE < 0 && pelletM < 0 && pelletE < 0 {
				if simpleMorning < 0 || simpleEvening < 0 {
					continue
				}
			}

			maxCol := dateCol
			for _, col := range []int{freshM, freshE, pelletM, pelletE, deathCol, touristCol, simpleMorning, simpleEvening} {
				if col > maxCol {
					maxCol = col
				}
			}

			byDay := make(map[int]dto.DailyLogEntryInput)
			dataStart := headerRows
			for _, row := range rows[dataStart:] {
				if len(row) <= maxCol {
					continue
				}
				dateStr := row[dateCol]
				d, err := parseExcelDateCell(dateStr)
				if err != nil {
					continue
				}
				if d.Year() != year || int(d.Month()) != mth {
					continue
				}
				day := d.Day()
				if day < 1 || day > daysInMonth {
					continue
				}

				var fm, fe, pm, pe decimal.Decimal
				var dErr error
				if freshM >= 0 {
					fm, dErr = parseDecimalCell(row[freshM])
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
				}
				if freshE >= 0 {
					fe, dErr = parseDecimalCell(row[freshE])
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
				}
				if pelletM >= 0 {
					pm, dErr = parseDecimalCell(row[pelletM])
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
				}
				if pelletE >= 0 {
					pe, dErr = parseDecimalCell(row[pelletE])
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
				}
				if simpleMorning >= 0 && simpleEvening >= 0 {
					pm, dErr = parseDecimalCell(row[simpleMorning])
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
					pe, dErr = parseDecimalCell(row[simpleEvening])
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
				}

				deaths := 0
				var tourist *int
				if deathCol >= 0 {
					deaths, dErr = parseIntCell(row[deathCol])
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
				}
				if touristCol >= 0 {
					cell := ""
					if touristCol < len(row) {
						cell = row[touristCol]
					}
					tourist, dErr = parseOptionalIntCell(cell)
					if dErr != nil {
						return nil, errors.ErrValidationFailed.Wrap(dErr)
					}
				}

				touristIsZero := tourist == nil || *tourist == 0
				if fm.IsZero() && fe.IsZero() && pm.IsZero() && pe.IsZero() && deaths == 0 && touristIsZero {
					continue
				}

				byDay[day] = dto.DailyLogEntryInput{
					Day:               day,
					FreshMorning:      fm,
					FreshEvening:      fe,
					PelletMorning:     pm,
					PelletEvening:     pe,
					DeathFishCount:    deaths,
					TouristCatchCount: tourist,
				}
			}

			if len(byDay) == 0 {
				continue
			}

			out := make([]dto.DailyLogEntryInput, 0, len(byDay))
			for d := 1; d <= daysInMonth; d++ {
				if e, ok := byDay[d]; ok {
					out = append(out, e)
				}
			}
			if len(out) == 0 {
				return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("no entries to import"))
			}
			return out, nil
		}
	}

	return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("no sheet with required columns (date + feed columns)"))
}
