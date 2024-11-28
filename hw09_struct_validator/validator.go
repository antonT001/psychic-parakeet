package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/constants"
	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/e"
	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/integer"
	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/str"
)

func Validate(v interface{}) error {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Struct:
		return validateStruct(t, reflect.ValueOf(v))
	default:
		return fmt.Errorf("%w: the validated object is not a structure", e.ErrInternal)
	}
}

func validateStruct(t reflect.Type, v reflect.Value) error {
	errs := make(e.ValidationErrors, 0)
	for i := 0; i < t.NumField(); i++ {
		err := validateField(t.Field(i), v.Field(i))
		if err != nil {
			if errors.Is(err, e.ErrInternal) {
				return fmt.Errorf("%w: field: %v", err, t.Field(i).Name)
			}

			errs = append(errs, e.ValidationError{
				Field: t.Field(i).Name,
				Err:   err,
			})
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateField(field reflect.StructField, value reflect.Value) error {
	switch field.Type.Kind() {
	case reflect.Int, reflect.String:
		return validateItem(field, value)

	case reflect.Array, reflect.Slice:
		errs := []error{}
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			err := validateItem(field, elem)
			if err != nil {
				if errors.Is(err, e.ErrInternal) {
					return err
				}
				errs = append(errs, err)
			}
		}
		if len(errs) == 0 {
			return nil
		}
		return errors.Join(errs...)

	default:
		return e.ErrUnsupportedFieldType
	}
}

func validateItem(field reflect.StructField, value reflect.Value) error {
	rule := field.Tag.Get(constants.ValidateTag)
	if len(rule) == 0 {
		return nil
	}
	switch value.Kind() {
	case reflect.Int:
		return integer.Validate(int(value.Int()), rule)
	case reflect.String:
		return str.Validate(value.String(), rule)
	default:
		return e.ErrUnsupportedItemType
	}
}
