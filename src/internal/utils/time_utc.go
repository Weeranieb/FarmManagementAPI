package utils

import "time"

// StartOfDayUTC returns midnight UTC on the civil calendar date of t in the UTC zone
// (year/month/day from t.UTC(), time 00:00:00 UTC).
func StartOfDayUTC(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
