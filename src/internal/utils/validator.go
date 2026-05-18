package utils

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var validate *validator.Validate

// passwordRE: ≥8 ASCII alphanumerics, with at least one upper and one lower case letter.
var passwordRE = regexp.MustCompile(`^[A-Za-z0-9]{8,}$`)

func init() {
	validate = validator.New()

	// Report the JSON name (camelCase) instead of the Go struct field name so
	// validation errors say `freshFeedCollectionId`, not `FreshFeedCollectionId`.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	// decimal_gt0: decimal.Decimal must be > 0 (use with required or omitempty as needed)
	_ = validate.RegisterValidation("decimal_gt0", func(fl validator.FieldLevel) bool {
		f := fl.Field()
		if f.Kind() != reflect.Struct {
			return true
		}
		d, ok := f.Interface().(decimal.Decimal)
		if !ok {
			return true
		}
		return d.GreaterThan(decimal.Zero)
	})
	// decimal_gte0: decimal.Decimal must be >= 0
	_ = validate.RegisterValidation("decimal_gte0", func(fl validator.FieldLevel) bool {
		f := fl.Field()
		if f.Kind() != reflect.Struct {
			return true
		}
		d, ok := f.Interface().(decimal.Decimal)
		if !ok {
			return true
		}
		return d.GreaterThanOrEqual(decimal.Zero)
	})
	// password: ≥8 ASCII alphanumerics, must contain upper and lower case.
	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if !passwordRE.MatchString(s) {
			return false
		}
		hasUpper := false
		hasLower := false
		for _, r := range s {
			if r >= 'A' && r <= 'Z' {
				hasUpper = true
			} else if r >= 'a' && r <= 'z' {
				hasLower = true
			}
			if hasUpper && hasLower {
				return true
			}
		}
		return false
	})
}

// ValidateStruct validates a struct using go-validator v10
func ValidateStruct(s any) error {
	return validate.Struct(s)
}
