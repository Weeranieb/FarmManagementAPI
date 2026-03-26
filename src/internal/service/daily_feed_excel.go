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

var dailyFeedExcelDateHeaders = []string{"date", "วันที่"}
var dailyFeedExcelMorningHeaders = []string{"morning", "เช้า"}
var dailyFeedExcelEveningHeaders = []string{"evening", "เย็น"}

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

// parseDailyFeedExcelFile reads the first sheet of an Excel file and returns
// feed entries for the given month (YYYY-MM). Expects a header row with
// Date/วันที่, Morning/เช้า, Evening/เย็น columns.
func parseDailyFeedExcelFile(filePath string, month string) ([]dto.DailyFeedEntryInput, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("only .xlsx files are supported"))
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}
	defer func() { _ = f.Close() }()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("empty workbook"))
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}
	if len(rows) < 2 {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("excel has no data rows"))
	}

	headers := rows[0]
	dateCol := findColumnIndex(headers, dailyFeedExcelDateHeaders)
	morningCol := findColumnIndex(headers, dailyFeedExcelMorningHeaders)
	eveningCol := findColumnIndex(headers, dailyFeedExcelEveningHeaders)
	if dateCol < 0 || morningCol < 0 || eveningCol < 0 {
		return nil, errors.ErrValidationFailed.Wrap(
			fmt.Errorf("missing columns: need Date/วันที่, Morning/เช้า, Evening/เย็น"),
		)
	}

	maxCol := dateCol
	if morningCol > maxCol {
		maxCol = morningCol
	}
	if eveningCol > maxCol {
		maxCol = eveningCol
	}

	start, _, err := parseMonth(month)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(err)
	}
	year, mth := start.Year(), int(start.Month())
	daysInMonth := time.Date(year, time.Month(mth)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	byDay := make(map[int]dto.DailyFeedEntryInput)

	for _, row := range rows[1:] {
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

		morning, err := parseDecimalCell(row[morningCol])
		if err != nil {
			return nil, errors.ErrValidationFailed.Wrap(err)
		}
		evening, err := parseDecimalCell(row[eveningCol])
		if err != nil {
			return nil, errors.ErrValidationFailed.Wrap(err)
		}
		if morning.IsZero() && evening.IsZero() {
			continue
		}

		byDay[day] = dto.DailyFeedEntryInput{
			Day:     day,
			Morning: morning,
			Evening: evening,
		}
	}

	if len(byDay) == 0 {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("no valid rows for month %s", month))
	}

	out := make([]dto.DailyFeedEntryInput, 0, len(byDay))
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
