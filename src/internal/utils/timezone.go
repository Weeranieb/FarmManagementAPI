package utils

import "time"

// ThailandLocation is the Asia/Bangkok timezone, used for converting
// PostgreSQL DATE values (which may be scanned as UTC instants) back
// to the intended civil calendar date.
var ThailandLocation *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic("utils: load Asia/Bangkok: " + err.Error())
	}
	ThailandLocation = loc
}

// CalendarDay returns the civil day-of-month for t in Thailand time.
func CalendarDay(t time.Time) int {
	return t.In(ThailandLocation).Day()
}
