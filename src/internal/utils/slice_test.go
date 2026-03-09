package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendStringIfMissing(t *testing.T) {
	t.Run("nil slice returns single element", func(t *testing.T) {
		// GIVEN — nil slice and element "a"
		// WHEN — AppendStringIfMissing is called
		got := AppendStringIfMissing(nil, "a")
		// THEN — result is ["a"]
		assert.Equal(t, []string{"a"}, got)
	})
	t.Run("missing element is appended", func(t *testing.T) {
		// GIVEN — slice ["a","b"] and element "c"
		// WHEN — AppendStringIfMissing is called
		got := AppendStringIfMissing([]string{"a", "b"}, "c")
		// THEN — result is ["a","b","c"]
		assert.Equal(t, []string{"a", "b", "c"}, got)
	})
	t.Run("existing element returns slice unchanged", func(t *testing.T) {
		// GIVEN — slice ["a","b"] and element "a"
		// WHEN — AppendStringIfMissing is called
		slice := []string{"a", "b"}
		got := AppendStringIfMissing(slice, "a")
		// THEN — result is unchanged ["a","b"]
		assert.Equal(t, []string{"a", "b"}, got)
	})
}
