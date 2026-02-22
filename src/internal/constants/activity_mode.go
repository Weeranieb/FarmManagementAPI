package constants

import "slices"

const (
	// ActivityModeFill - Add fish to a pond
	ActivityModeFill = "fill"

	// ActivityModeMove - Transfer fish from one pond to another
	ActivityModeMove = "move"

	// ActivityModeSell - Record a sell
	ActivityModeSell = "sell"
)

// ValidActivityModes returns all valid activity mode values (for API/DB).
func ValidActivityModes() []string {
	return []string{
		ActivityModeFill,
		ActivityModeMove,
		ActivityModeSell,
	}
}

// IsValidActivityMode checks if the provided mode is valid.
func IsValidActivityMode(mode string) bool {
	return slices.Contains(ValidActivityModes(), mode)
}
