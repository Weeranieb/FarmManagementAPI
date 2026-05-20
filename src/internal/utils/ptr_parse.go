package utils

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// DecimalFromStringPtr parses *string into *decimal.Decimal.
// Returns nil when s is nil, empty, or not a valid decimal.
func DecimalFromStringPtr(s *string) *decimal.Decimal {
	if s == nil || *s == "" {
		return nil
	}
	d, err := decimal.NewFromString(*s)
	if err != nil {
		return nil
	}
	return &d
}

// DecimalFromStringPtrOrZero parses *string into decimal.Decimal.
// Returns decimal.Zero when s is nil, empty, or not a valid decimal.
func DecimalFromStringPtrOrZero(s *string) decimal.Decimal {
	d := DecimalFromStringPtr(s)
	if d == nil {
		return decimal.Zero
	}
	return *d
}

// StringSliceFromJSONPtr unmarshals a JSON string array from *string.
// Returns nil when s is nil, empty, or not valid JSON for []string.
func StringSliceFromJSONPtr(s *string) []string {
	if s == nil || *s == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(*s), &out); err != nil {
		return nil
	}
	return out
}
