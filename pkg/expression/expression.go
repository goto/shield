package expression

import (
	"fmt"
	"reflect"

	"github.com/antonmedv/expr"
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
		code := fmt.Sprintf("Key %s Value", e.Comparison.Operator)

		program, err := expr.Compile(code, expr.Env(Comparison{}))
		if err != nil {
			return nil, err
		}

		output, err := expr.Run(program, e.Comparison)
		if err != nil {
			return nil, err
		}

		return output, nil
	}
	return nil, nil
}

func (e Expression) IsEmpty() bool {
	return reflect.ValueOf(e).IsZero()
}
