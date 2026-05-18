package http

import (
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldError is the per-field detail attached to a validation failure response.
// Surfaces the JSON path of the offending field plus a human-readable cause
// so clients can render inline form errors without having to parse strings.
type FieldError struct {
	Field   string `json:"field" example:"freshFeedCollectionId"`
	Tag     string `json:"tag" example:"required"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message" example:"is required"`
}

// BuildFieldErrors converts a go-playground/validator error chain into a slice
// of structured FieldError entries. Returns nil if err isn't a ValidationErrors
// — callers fall back to the plain wrapped message in that case.
func BuildFieldErrors(err error) []FieldError {
	if err == nil {
		return nil
	}
	var ves validator.ValidationErrors
	if !stderrors.As(err, &ves) {
		return nil
	}
	out := make([]FieldError, 0, len(ves))
	for _, fe := range ves {
		out = append(out, FieldError{
			Field:   fieldPath(fe),
			Tag:     fe.Tag(),
			Value:   formatValue(fe.Value()),
			Message: humanize(fe),
		})
	}
	return out
}

// SummarizeFieldErrors joins per-field messages into a single line so older
// clients that only read the `details` string still see something actionable.
func SummarizeFieldErrors(fields []FieldError) string {
	if len(fields) == 0 {
		return ""
	}
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		parts = append(parts, fmt.Sprintf("%s %s", f.Field, f.Message))
	}
	return strings.Join(parts, "; ")
}

// fieldPath returns the JSON namespace of the offending field with the root
// struct name stripped — `RegisterDTO.username` → `username`,
// `Body.items[0].quantity` → `items[0].quantity`.
func fieldPath(fe validator.FieldError) string {
	ns := fe.Namespace()
	if i := strings.IndexByte(ns, '.'); i >= 0 {
		return ns[i+1:]
	}
	return ns
}

func formatValue(v any) string {
	if v == nil {
		return ""
	}
	s := fmt.Sprintf("%v", v)
	if len(s) > 64 {
		return s[:61] + "..."
	}
	return s
}

// humanize turns validator tag codes into a human-readable phrase. Falls back
// to the tag itself for tags we don't have a specific phrase for, which is
// still better than the raw error string from validator.
func humanize(fe validator.FieldError) string {
	param := fe.Param()
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "url":
		return "must be a valid URL"
	case "uuid":
		return "must be a valid UUID"
	case "min":
		return fmt.Sprintf("must be at least %s", param)
	case "max":
		return fmt.Sprintf("must be at most %s", param)
	case "gt":
		return fmt.Sprintf("must be greater than %s", param)
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", param)
	case "lt":
		return fmt.Sprintf("must be less than %s", param)
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", param)
	case "len":
		return fmt.Sprintf("must have length %s", param)
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", param)
	case "eqfield":
		return fmt.Sprintf("must equal field %s", param)
	case "nefield":
		return fmt.Sprintf("must not equal field %s", param)
	case "decimal_gt0":
		return "must be greater than 0"
	case "decimal_gte0":
		return "must be greater than or equal to 0"
	case "password":
		return "must be at least 8 alphanumeric characters with at least one upper and one lower case letter"
	default:
		if param != "" {
			return fmt.Sprintf("failed %q (%s)", fe.Tag(), param)
		}
		return fmt.Sprintf("failed %q", fe.Tag())
	}
}
