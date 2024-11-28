package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/constants"
	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/e"
)

type (
	rules map[constants.Key]any
)

func New() rules {
	return make(rules)
}

func (r rules) Parse(in string) error {
	if len(in) == 0 {
		return e.ErrNoRule
	}

	for _, v := range strings.Split(in, "|") {
		rule := strings.Split(v, ":")
		if len(rule) != 2 {
			return fmt.Errorf("%w: more than one separator [:]", e.ErrInternal)
		}

		value := strings.Split(rule[1], ",")
		err := r.set(constants.Key(rule[0]), value)
		if err != nil {
			return fmt.Errorf("%w: %v: %w", e.ErrInternal, e.ErrParseRule, err)
		}
	}
	return nil
}

func (r rules) Conditions(in constants.Key) (any, bool) {
	out, exist := r[in]
	return out, exist
}

func (r rules) set(key constants.Key, value []string) error {
	switch key {
	case constants.MinKey,
		constants.MaxKey,
		constants.LenKey:
		if len(value) != 1 {
			return fmt.Errorf("more than one value for a rule %s", key)
		}
		v, err := strconv.Atoi(value[0])
		if err != nil {
			return fmt.Errorf("rule: %s: %w", key, err)
		}
		return r.setValue(key, v)

	case constants.InKey:
		if len(value) < 2 {
			return fmt.Errorf("less than two values for the rule %s", key)
		}
		return r.setValue(key, value)

	case constants.RegexpKey:
		if len(value) != 1 {
			return fmt.Errorf("more than one value for a rule %s", key)
		}
		return r.setValue(key, value[0])

	default:
		return fmt.Errorf("unsupported rule [%s]", key)
	}
}

func (r rules) setValue(key constants.Key, in any) error {
	if _, ok := r[key]; ok {
		return fmt.Errorf("rule %s is set", key)
	}
	r[key] = in
	return nil
}
