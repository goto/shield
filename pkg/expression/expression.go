package expression

import (
	"errors"
	"fmt"
	"reflect"
)

type Expression struct {
	Comparison Comparison `yaml:"comparison" mapstructure:"comparison"`
}

type Comparison struct {
	Key      any    `yaml:"key" mapstructure:"key"`
	Operator string `yaml:"operator" mapstructure:"operator"`
	Value    any    `yaml:"value" mapstructure:"value"`
}

func (e Expression) Evaluate() (any, error) {
	if !reflect.ValueOf(e.Comparison).IsZero() {
		var output any
		var err error
		switch e.Comparison.Operator {
		case "==":
			output = reflect.DeepEqual(e.Comparison.Key, e.Comparison.Value)
			err = nil
		default:
			err = errors.New(fmt.Sprintf("unsupported comparison operator %s", e.Comparison.Operator))
			output = nil
		}

		return output, err
	}
	return nil, errors.New("no supported operators configured")
}

func (e Expression) IsEmpty() bool {
	return reflect.ValueOf(e).IsZero()
}
