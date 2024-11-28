package str

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/constants"
	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/e"
	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/rules"
)

type (
	Rules interface {
		Parse(in string) error
		Conditions(in constants.Key) (any, bool)
	}
	validator struct {
		fns []func(string, Rules) error
	}
)

var vld = new(validator)

func init() {
	vld.fns = append(vld.fns, vld.validateLen)
	vld.fns = append(vld.fns, vld.validateRegexp)
	vld.fns = append(vld.fns, vld.validateIn)
}

func Validate(value string, rule string) error {
	rules := rules.New()
	err := rules.Parse(rule)
	if err != nil {
		if errors.Is(err, e.ErrNoRule) {
			return nil
		}
		return err
	}

	errs := []error{}
	for _, fn := range vld.fns {
		err := fn(value, rules)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (vo validator) validateLen(in string, r Rules) error {
	if v, exist := r.Conditions(constants.LenKey); exist && len(in) != v.(int) {
		return errors.New("the value does not match the specified length")
	}

	return nil
}

func (vo validator) validateRegexp(in string, r Rules) error {
	if v, exist := r.Conditions(constants.RegexpKey); exist {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("%w: rule: %s: %w", e.ErrInternal, constants.RegexpKey, err)
		}

		if r.MatchString(in) {
			return nil
		}
		return errors.New("the value does not match the specified regular expression")
	}

	return nil
}

func (vo validator) validateIn(in string, r Rules) error {
	if v, exist := r.Conditions(constants.InKey); exist {
		for _, elem := range v.([]string) {
			if elem == in {
				return nil
			}
		}
		return errors.New("the value is not included in the specified set")
	}
	return nil
}
