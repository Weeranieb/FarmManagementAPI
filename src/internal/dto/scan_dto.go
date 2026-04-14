package dto

type ScanDailyLogRequest struct {
	Month string `json:"month" form:"month" validate:"required"`
}

type ScanEntry struct {
	Day            int      `json:"day"`
	FreshMorning   *float64 `json:"freshMorning"`
	FreshEvening   *float64 `json:"freshEvening"`
	PelletMorning  *float64 `json:"pelletMorning"`
	PelletEvening  *float64 `json:"pelletEvening"`
	DeathFishCount *int     `json:"deathFishCount"`
}

type ScanConfidence struct {
	Day            int     `json:"day"`
	FreshMorning   float64 `json:"freshMorning"`
	FreshEvening   float64 `json:"freshEvening"`
	PelletMorning  float64 `json:"pelletMorning"`
	PelletEvening  float64 `json:"pelletEvening"`
	DeathFishCount float64 `json:"deathFishCount"`
}

type ScanDailyLogResponse struct {
	ScanLogId  int              `json:"scanLogId"`
	Month      string           `json:"month"`
	Entries    []ScanEntry      `json:"entries"`
	Confidence []ScanConfidence `json:"confidence"`
	ImageUrls  []string         `json:"imageUrls"`
	Notes      string           `json:"notes"`
}
