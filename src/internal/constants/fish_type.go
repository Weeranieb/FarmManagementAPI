package constants

import "slices"

const (
	// FishTypeNil - Nile tilapia (ปลานิล)
	FishTypeNil = "nil"

	// FishTypeKaphong - Barramundi (ปลากะพง)
	FishTypeKaphong = "kaphong"

	// FishTypeKang - Kang (ปลาคัง)
	FishTypeKang = "kang"

	// FishTypeDuk - Catfish (ปลาดุก)
	FishTypeDuk = "duk"

	// FishUnitKg - Default unit for fish weight (kilograms)
	FishUnitKg = "kg"
)

// ValidFishTypes returns all valid fish type values (lowercase, for API/DB).
func ValidFishTypes() []string {
	return []string{
		FishTypeNil,
		FishTypeKaphong,
		FishTypeKang,
		FishTypeDuk,
	}
}

// IsValidFishType checks if the provided fish type is valid.
func IsValidFishType(fishType string) bool {
	return slices.Contains(ValidFishTypes(), fishType)
}
