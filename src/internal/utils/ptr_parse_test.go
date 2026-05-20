package utils_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

func TestDecimalFromStringPtr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, utils.DecimalFromStringPtr(nil))
	})

	t.Run("empty", func(t *testing.T) {
		s := ""
		assert.Nil(t, utils.DecimalFromStringPtr(&s))
	})

	t.Run("invalid", func(t *testing.T) {
		s := "not-a-decimal"
		assert.Nil(t, utils.DecimalFromStringPtr(&s))
	})

	t.Run("valid", func(t *testing.T) {
		s := "12.5"
		got := utils.DecimalFromStringPtr(&s)
		require.NotNil(t, got)
		assert.True(t, got.Equal(decimal.RequireFromString("12.5")))
	})
}

func TestDecimalFromStringPtrOrZero(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.True(t, utils.DecimalFromStringPtrOrZero(nil).IsZero())
	})

	t.Run("invalid", func(t *testing.T) {
		s := "x"
		assert.True(t, utils.DecimalFromStringPtrOrZero(&s).IsZero())
	})

	t.Run("valid", func(t *testing.T) {
		s := "99.99"
		got := utils.DecimalFromStringPtrOrZero(&s)
		assert.True(t, got.Equal(decimal.RequireFromString("99.99")))
	})
}

func TestStringSliceFromJSONPtr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, utils.StringSliceFromJSONPtr(nil))
	})

	t.Run("empty", func(t *testing.T) {
		s := ""
		assert.Nil(t, utils.StringSliceFromJSONPtr(&s))
	})

	t.Run("invalid", func(t *testing.T) {
		s := "not-json"
		assert.Nil(t, utils.StringSliceFromJSONPtr(&s))
	})

	t.Run("valid", func(t *testing.T) {
		s := `["tilapia","catfish"]`
		got := utils.StringSliceFromJSONPtr(&s)
		assert.Equal(t, []string{"tilapia", "catfish"}, got)
	})
}
