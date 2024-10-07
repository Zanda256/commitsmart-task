// Package validate contains the support for validating models.
// validate holds the settings and caches for validating request struct values.
package validate

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var validate *validator.Validate

func init() {

	// Instantiate a validator.
	validate = validator.New()

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Check validates the provided model against it's declared tags.
func Check(val any) error {
	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Error(),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
