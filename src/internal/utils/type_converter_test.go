package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

func TestConvertRepeatedFormInts(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got, err := utils.ConvertRepeatedFormInts("ids", []string{"1", " 2 ", "3"})
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("required", func(t *testing.T) {
		_, err := utils.ConvertRepeatedFormInts("ids", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ids is required")

		_, err = utils.ConvertRepeatedFormInts("ids", []string{})
		require.Error(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := utils.ConvertRepeatedFormInts("ids", []string{"1", "x"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ids")
	})
}
