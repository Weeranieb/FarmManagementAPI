package constants

import "slices"

const (
	FeedTypeFresh  = "fresh"
	FeedTypePellet = "pellet"

	FeedTypeLabelTHFresh  = "เหยื่อสด"
	FeedTypeLabelTHPellet = "อาหารเม็ด"

	// Default kg per purchase pack when a feed has none specified — a fresh
	// crate (ลัง) ≈ 30 กก., a pellet bag (ถุง) ≈ 20 กก. pack_size_kg is NOT NULL,
	// so create must always resolve to one of these when the caller omits it.
	DefaultPackSizeKgFresh  = 30
	DefaultPackSizeKgPellet = 20
)

// FeedTypeLabelsTH maps canonical feed_type (API/DB) to Thai UI labels.
var FeedTypeLabelsTH = map[string]string{
	FeedTypeFresh:  FeedTypeLabelTHFresh,
	FeedTypePellet: FeedTypeLabelTHPellet,
}

// ValidFeedTypes returns allowed feed_type values for API/DB.
func ValidFeedTypes() []string {
	return []string{FeedTypeFresh, FeedTypePellet}
}

// IsValidFeedType reports whether s is a known feed type.
func IsValidFeedType(s string) bool {
	return slices.Contains(ValidFeedTypes(), s)
}

// FeedTypeLabelTH returns the Thai label for a canonical feed type, or empty if unknown.
func FeedTypeLabelTH(s string) string {
	if v, ok := FeedTypeLabelsTH[s]; ok {
		return v
	}
	return ""
}

// DefaultPackSizeKg returns the assumed pack size (กก.) for a feed type, used
// when a create request omits pack_size_kg (the column is NOT NULL).
func DefaultPackSizeKg(feedType string) float64 {
	if feedType == FeedTypeFresh {
		return DefaultPackSizeKgFresh
	}
	return DefaultPackSizeKgPellet
}
