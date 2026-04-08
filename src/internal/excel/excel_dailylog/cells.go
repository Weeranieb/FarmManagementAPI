package excel_dailylog

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func parseDecimalCell(s string) (decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "#") {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(s)
}

func parseIntCell(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "#") {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func parseOptionalIntCell(s string) (*int, error) {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "#") {
		return nil, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseOptionalDecimalCell(s string) (*decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "#") {
		return nil, nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func cellStr(rows [][]string, row, col int) string {
	if row < 0 || row >= len(rows) || col < 0 {
		return ""
	}
	if col >= len(rows[row]) {
		return ""
	}
	return rows[row][col]
}

func compositeColHeader(rows [][]string, col, topRow, bottomRow int) string {
	var parts []string
	for r := topRow; r <= bottomRow && r < len(rows); r++ {
		s := strings.TrimSpace(cellStr(rows, r, col))
		if s != "" && !monthHeaderRe.MatchString(s) {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " ")
}

func parseDayOfMonth(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty day")
	}
	if strings.HasPrefix(s, "#") {
		return 0, fmt.Errorf("formula error in day cell")
	}
	day, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if day < 1 || day > 31 {
		return 0, fmt.Errorf("day out of range: %d", day)
	}
	return day, nil
}
