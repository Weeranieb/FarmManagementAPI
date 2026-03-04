package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendStringIfMissing(t *testing.T) {
	t.Run("nil slice returns single element", func(t *testing.T) {
		got := AppendStringIfMissing(nil, "a")
		assert.Equal(t, []string{"a"}, got)
	})
	t.Run("missing element is appended", func(t *testing.T) {
		got := AppendStringIfMissing([]string{"a", "b"}, "c")
		assert.Equal(t, []string{"a", "b", "c"}, got)
	})
	t.Run("existing element returns slice unchanged", func(t *testing.T) {
		slice := []string{"a", "b"}
		got := AppendStringIfMissing(slice, "a")
		assert.Equal(t, []string{"a", "b"}, got)
	})
}
