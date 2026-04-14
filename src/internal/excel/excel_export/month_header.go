package excel_export

import (
	"fmt"
	"time"
)

var monthAbbrevs = [13]string{
	0:  "",
	1:  "Jan",
	2:  "Feb",
	3:  "Mar",
	4:  "Apr",
	5:  "May",
	6:  "Jun",
	7:  "Jul",
	8:  "Aug",
	9:  "Sep",
	10: "Oct",
	11: "Nov",
	12: "Dec",
}

// FormatMonthHeaderBE converts a Gregorian (year, month) to the "Mon-YY" Buddhist-era
// header used in the daily-log template. Example: (2026, February) -> "Feb-69".
func FormatMonthHeaderBE(year int, month time.Month) string {
	beYear := year + 543
	yy := beYear % 100
	return fmt.Sprintf("%s-%02d", monthAbbrevs[month], yy)
}
