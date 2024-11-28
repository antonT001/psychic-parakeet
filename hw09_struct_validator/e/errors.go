package e

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrInternal             = errors.New("internal error")
	ErrUnsupportedFieldType = errors.New("unsupported field type")
	ErrUnsupportedItemType  = errors.New("unsupported item type")
	ErrNoRule               = errors.New("no rule")
	ErrParseRule            = errors.New("failed parse rule")
)

type (
	ValidationError struct {
		Field string
		Err   error
	}

	ValidationErrors []ValidationError
)

// for tests.
func VError(field string, err error) ValidationErrors {
	errs := []error{}
	errs = append(errs, err)
	return []ValidationError{
		{
			Field: field,
			Err:   errors.Join(errs...),
		},
	}
}

func (v ValidationErrors) Error() string {
	var buffer bytes.Buffer
	buffer.Grow(len(v))

	for _, err := range v {
		buffer.WriteString(fmt.Sprintf("field: %v, error: %v\n", err.Field, err.Err))
	}
	return fmt.Sprintf("failed validate: %s", buffer.String())
}
