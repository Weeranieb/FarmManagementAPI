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
	CreatedByName string `json:"createdByName"`
	// PondId is the source pond's ponds.id — lets a client open the pond a feed
	// row came from (GET /pond/:pondId/… all key off this id, not the
	// active_ponds.id the activity row itself references).
	PondId   int    `json:"pondId"`
	PondName string `json:"pondName"`
	// FarmName is the farm the source pond belongs to — shown as a secondary
	// label so a feed row is identifiable when pond names repeat across farms.
	FarmName string `json:"farmName"`
	// ToPondId / ToPondName are the move destination — set for moves only.
	ToPondId   *int    `json:"toPondId,omitempty"`
	ToPondName *string `json:"toPondName,omitempty"`
	// FishType is empty for a sell: species is recorded per size-grade on the
	// detail lines, so a sale has no single species to report.
	FishType string `json:"fishType"`
	// Amount is the head count. For a sell it comes from Σ sell_details.fish_count
	// (the activity row's own column is unset) — see ActivityFeedItem's doc.
	Amount     int     `json:"amount"`
	FishWeight float64 `json:"fishWeight"` // kg per fish; 0 for a sell
	FishUnit   string  `json:"fishUnit"`
	// PricePerUnit is ฿ per kg as entered; for a sell it is the derived average
	// (Σ revenue ÷ Σ weight), so clients never have to divide.
	PricePerUnit float64 `json:"pricePerUnit"`
	Total        float64 `json:"total"`
	// TotalWeight is sell-only: Σ sell_details.weight (kg) — the headline
	// "ขายปลานิล 312 กก." figure.
	TotalWeight *float64 `json:"totalWeight,omitempty"`
	Merchant    *string  `json:"merchant,omitempty"`
}

// SellDetailLine is one size-grade line of a sale, returned by
// GET /activity/:activityId/sell-details. A sale's money is priced per grade,
// so this is the only place the "how many baht per size" question is answerable
// — the feed row and ActivityFeedItem only carry the summed totals.
//
// Ordered smallest grade first. An empty list means the activity has no sell
// lines: it does not exist, is not a sell, or belongs to another client.
type SellDetailLine struct {
	FishSizeGradeId int `json:"fishSizeGradeId"`
	// SizeName is the grade's display name ("ไซส์ 1", "ใหญ่", …). Empty when the
	// grade row has since been soft-deleted; clients should fall back to the id.
	SizeName string `json:"sizeName"`
	// FishCount is omitted on legacy lines written before the column was
	// required — absent means "not recorded", not "zero fish".
	FishCount    *int    `json:"fishCount,omitempty"`
	Weight       float64 `json:"weight"`       // kg sold in this grade
	PricePerUnit float64 `json:"pricePerUnit"` // ฿ per kg for this grade
	// Total is weight × pricePerUnit, computed here so every client reports the
	// same figure as the summed headline rather than re-deriving it.
	Total float64 `json:"total"`
}
