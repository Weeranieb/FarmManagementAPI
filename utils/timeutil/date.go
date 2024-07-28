package timeutil

import "time"

func DaysInMonth(year int, month time.Month) int {
	// Create a time.Time object for the first day of the next month
	firstDayOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)

	// Subtract one day to get the last day of the current month
	lastDayOfMonth := firstDayOfNextMonth.AddDate(0, 0, -1)

	// Return the day of the last day of the month, which is the total number of days
	return lastDayOfMonth.Day()
}

func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

var ThaiMonths = [...]string{
	"",      // Padding to align month numbers with indexes
	"ม.ค.",  // January
	"ก.พ.",  // February
	"มี.ค.", // March
	"เม.ย.", // April
	"พ.ค.",  // May
	"มิ.ย.", // June
	"ก.ค.",  // July
	"ส.ค.",  // August
	"ก.ย.",  // September
	"ต.ค.",  // October
	"พ.ย.",  // November
	"ธ.ค.",  // December
}

var FullThaiMonths = [...]string{
	"",           // Padding to align month numbers with indexes
	"มกราคม",     // January
	"กุมภาพันธ์", // February
	"มีนาคม",     // March
	"เมษายน",     // April
	"พฤษภาคม",    // May
	"มิถุนายน",   // June
	"กรกฎาคม",    // July
	"สิงหาคม",    // August
	"กันยายน",    // September
	"ตุลาคม",     // October
	"พฤศจิกายน",  // November
	"ธันวาคม",    // December
}

var FullThaiMonthToTime = map[string]time.Month{
	"มกราคม":     time.January,
	"กุมภาพันธ์": time.February,
	"มีนาคม":     time.March,
	"เมษายน":     time.April,
	"พฤษภาคม":    time.May,
	"มิถุนายน":   time.June,
	"กรกฎาคม":    time.July,
	"สิงหาคม":    time.August,
	"กันยายน":    time.September,
	"ตุลาคม":     time.October,
	"พฤศจิกายน":  time.November,
	"ธันวาคม":    time.December,
}
