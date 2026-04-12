package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartOfDayUTC(t *testing.T) {
	in := time.Date(2026, 3, 15, 14, 30, 0, 0, time.FixedZone("ICT", 7*3600))
	got := StartOfDayUTC(in)
	want := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	assert.True(t, got.Equal(want), "got %v want %v", got, want)
}
