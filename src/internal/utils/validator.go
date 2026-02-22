package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
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
}

// ValidateStruct validates a struct using go-validator v10
func ValidateStruct(s any) error {
	return validate.Struct(s)
}
