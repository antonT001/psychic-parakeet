package integer

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

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
		fns []func(int, Rules) error
	}
)

var vld = new(validator)

func init() {
	vld.fns = append(vld.fns, vld.validateMin)
	vld.fns = append(vld.fns, vld.validateMax)
	vld.fns = append(vld.fns, vld.validateIn)
}

func Validate(value int, rule string) error {
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
			if errors.Is(err, e.ErrInternal) {
				return err
			}
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (vo validator) validateMin(in int, r Rules) error {
	if v, exist := r.Conditions(constants.MinKey); exist && in < v.(int) {
		return errors.New("the value is less than the set limit")
	}

	return nil
}

func (vo validator) validateMax(in int, r Rules) error {
	if v, exist := r.Conditions(constants.MaxKey); exist && in > v.(int) {
		return errors.New("the value is greater than the set limit")
	}

	return nil
}

func (vo validator) validateIn(in int, r Rules) error {
	if v, exist := r.Conditions(constants.InKey); exist {
		value := v.([]string)
		result := make([]int, len(value))
		for i := range value {
			v, err := strconv.Atoi(value[i])
			if err != nil {
				return fmt.Errorf("%w: rule: %s: %w", e.ErrInternal, constants.InKey, err)
			}
			result[i] = v
		}
		_, success := slices.BinarySearch(result, in)
		if !success {
			return errors.New("the value is not included in the specified set")
		}
	}
	return nil
}
