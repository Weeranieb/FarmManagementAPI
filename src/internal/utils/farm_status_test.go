package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
)

func TestDeriveFarmStatusFromPonds(t *testing.T) {
	t.Parallel()

	t.Run("nil slice", func(t *testing.T) {
		assert.Equal(t, constants.FarmStatusMaintenance, DeriveFarmStatusFromPonds(nil))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, constants.FarmStatusMaintenance, DeriveFarmStatusFromPonds([]*model.Pond{}))
	})

	t.Run("nil entry skipped", func(t *testing.T) {
		ponds := []*model.Pond{nil, {Status: constants.FarmStatusMaintenance}}
		assert.Equal(t, constants.FarmStatusMaintenance, DeriveFarmStatusFromPonds(ponds))
	})

	t.Run("all maintenance", func(t *testing.T) {
		ponds := []*model.Pond{
			{Status: constants.FarmStatusMaintenance},
			{Status: constants.FarmStatusMaintenance},
		}
		assert.Equal(t, constants.FarmStatusMaintenance, DeriveFarmStatusFromPonds(ponds))
	})

	t.Run("one active", func(t *testing.T) {
		ponds := []*model.Pond{
			{Status: constants.FarmStatusMaintenance},
			{Status: constants.FarmStatusActive},
		}
		assert.Equal(t, constants.FarmStatusActive, DeriveFarmStatusFromPonds(ponds))
	})

	t.Run("all active", func(t *testing.T) {
		ponds := []*model.Pond{
			{Status: constants.FarmStatusActive},
			{Status: constants.FarmStatusActive},
		}
		assert.Equal(t, constants.FarmStatusActive, DeriveFarmStatusFromPonds(ponds))
	})

	t.Run("unknown status treated as not active", func(t *testing.T) {
		ponds := []*model.Pond{{Status: "other"}}
		assert.Equal(t, constants.FarmStatusMaintenance, DeriveFarmStatusFromPonds(ponds))
	})
}
