package constants

// ActivityType represents the type of activity.
type ActivityType string

const (
	FillType ActivityType = "FILL"
	MoveType ActivityType = "MOVE"
	SellType ActivityType = "SELL"
)
