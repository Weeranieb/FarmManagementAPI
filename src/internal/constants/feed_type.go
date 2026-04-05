package constants

import "slices"

const (
	FeedTypeFresh  = "fresh"
	FeedTypePellet = "pellet"

	FeedTypeLabelTHFresh  = "เหยื่อสด"
	FeedTypeLabelTHPellet = "อาหารเม็ด"
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
