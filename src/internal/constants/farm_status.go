package constants

const (
	// FarmStatusActive - Farm is active and operational
	FarmStatusActive = "active"

	// FarmStatusMaintenance - Farm is under maintenance
	FarmStatusMaintenance = "maintenance"
)

// ValidFarmStatuses returns all valid farm status values
func ValidFarmStatuses() []string {
	return []string{
		FarmStatusActive,
		FarmStatusMaintenance,
	}
}

// IsValidFarmStatus checks if the provided status is valid
func IsValidFarmStatus(status string) bool {
	for _, validStatus := range ValidFarmStatuses() {
		if status == validStatus {
			return true
		}
	}
	return false
}
