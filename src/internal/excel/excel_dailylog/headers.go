package excel_dailylog

import "strings"

func normalizeHeader(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func headerLooksLikeFresh(nh string) bool {
	fresh := strings.Contains(nh, "เหยื่อ") || strings.Contains(nh, "fresh")
	session := strings.Contains(nh, "เช้า") || strings.Contains(nh, "เย็น") ||
		strings.Contains(nh, "morning") || strings.Contains(nh, "evening")
	return fresh && !session
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

func headerLooksLikeAvgBodyWeight(nh string) bool {
	return strings.Contains(nh, "นน.ตัว") || strings.Contains(nh, "น้ำหนักตัว")
}

func headerLooksLikeFishCount(nh string) bool {
	return strings.Contains(nh, "จำนวนปลา")
}

func isSummaryDayLabel(s string) bool {
	n := normalizeHeader(s)
	return strings.Contains(n, "รวม") || strings.Contains(n, "total") || strings.Contains(n, "average")
}
