package dto

import "time"

// ActivityFeedItem is one row of the farm-wide activity feed returned by
// GET /activity — the data source for the Home "กิจกรรมล่าสุด" section and
// the full "ประวัติกิจกรรม" history screen.
//
// The feed only carries discrete user-authored events (fill / move / sell).
// Daily-log saves are deliberately excluded — they happen every day on every
// pond and would bury the events users actually scroll back for.
//
// Total is computed server-side with the same rules as the per-pond
// timeline (GET /pond/:pondId/activities): sell totals come from
// sell_details (Σ weight × price_per_unit), fill/move totals are
// amount × fish_weight × price_per_unit + Σ additional_costs.
type ActivityFeedItem struct {
	Id   int    `json:"id"`
	Mode string `json:"mode"` // fill | move | sell
	// ActivityDate is the user-chosen event date (date-only; the time part
	// is always midnight). Group rows by this, not CreatedAt.
	ActivityDate time.Time `json:"activityDate"`
	// CreatedAt is when the record was saved — clients may show its clock
	// time when it falls on ActivityDate (same-day logging) and omit it for
	// backdated entries.
	CreatedAt time.Time `json:"createdAt"`
	// CreatedBy is the author's username (audit identity).
	CreatedBy string `json:"createdBy"`
	// CreatedByName is the author's display name (first name), falling back
	// to the username when the user record no longer exists.
	CreatedByName string  `json:"createdByName"`
	PondName      string  `json:"pondName"`
	// FarmName is the farm the source pond belongs to — shown as a secondary
	// label so a feed row is identifiable when pond names repeat across farms.
	FarmName   string  `json:"farmName"`
	ToPondName *string `json:"toPondName,omitempty"` // move only — destination pond
	FishType      string  `json:"fishType"`
	Amount        int     `json:"amount"`
	FishWeight    float64 `json:"fishWeight"`
	FishUnit      string  `json:"fishUnit"`
	PricePerUnit  float64 `json:"pricePerUnit"`
	Total         float64 `json:"total"`
	// TotalWeight is sell-only: Σ sell_details.weight (kg) — the headline
	// "ขายปลานิล 312 กก." figure.
	TotalWeight *float64 `json:"totalWeight,omitempty"`
	Merchant    *string  `json:"merchant,omitempty"`
}
