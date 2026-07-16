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

// NullDecimalFromStringPtr parses *string into a decimal.NullDecimal.
// Returns an invalid (NULL) NullDecimal when s is nil, empty, or not a valid
// decimal; a valid one otherwise.
func NullDecimalFromStringPtr(s *string) decimal.NullDecimal {
	return NullDecimalFromDecimalPtr(DecimalFromStringPtr(s))
}

// NullDecimalFromDecimalPtr converts a *decimal.Decimal (the shape carried by
// DTO request/response fields) into a decimal.NullDecimal (the shape persisted
// on models). Invalid (NULL) when the pointer is nil.
func NullDecimalFromDecimalPtr(d *decimal.Decimal) decimal.NullDecimal {
	if d == nil {
		return decimal.NullDecimal{}
	}
	return decimal.NullDecimal{Decimal: *d, Valid: true}
}

// DecimalPtrFromNullDecimal converts a decimal.NullDecimal (persisted on models)
// back into a *decimal.Decimal for DTO fields. Returns nil when invalid (NULL).
func DecimalPtrFromNullDecimal(d decimal.NullDecimal) *decimal.Decimal {
	if !d.Valid {
		return nil
	}
	v := d.Decimal
	return &v
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
