package utils

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
)

// DeriveFarmStatusFromPonds returns FarmStatusActive if at least one non-nil pond
// has status active; otherwise FarmStatusMaintenance (including no ponds).
// Callers pass only non-deleted ponds (same as pondRepo.ListByFarmId).
func DeriveFarmStatusFromPonds(ponds []*model.Pond) string {
	for _, p := range ponds {
		if p != nil && p.Status == constants.FarmStatusActive {
			return constants.FarmStatusActive
		}
	}
	return constants.FarmStatusMaintenance
}
