package excel_dailylog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var monthHeaderRe = regexp.MustCompile(`(?i)^(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)-(\d{2})$`)

var monthAbbrevs = map[string]time.Month{
	"jan": time.January, "feb": time.February, "mar": time.March, "apr": time.April,
	"may": time.May, "jun": time.June, "jul": time.July, "aug": time.August,
	"sep": time.September, "oct": time.October, "nov": time.November, "dec": time.December,
}

// parseEnglishMonthBEHeader parses "Feb-69" as February พ.ศ. 2569 → Gregorian 2026.
func parseEnglishMonthBEHeader(raw string) (year int, month time.Month, err error) {
	raw = strings.TrimSpace(raw)
	m := monthHeaderRe.FindStringSubmatch(raw)
	if m == nil {
		return 0, 0, fmt.Errorf("not a month header: %q", raw)
	}
	mon, ok := monthAbbrevs[strings.ToLower(m[1])]
	if !ok {
		return 0, 0, fmt.Errorf("unknown month: %s", m[1])
	}
	yy, err := strconv.Atoi(m[2])
	if err != nil {
		return 0, 0, fmt.Errorf("parse BE year suffix %q: %w", m[2], err)
	}
	beYear := 2500 + yy
	gy := beYear - 543
	return gy, mon, nil
}
