package expression

import (
	"fmt"
	"reflect"
)

type Expression struct {
	Attribute any    `yaml:"attribute" mapstructure:"attribute"`
	Operator  string `yaml:"operator" mapstructure:"operator"`
	Value     any    `yaml:"value" mapstructure:"value"`
}

func (e Expression) Evaluate() (bool, error) {
	output := false
	var err error

	switch e.Operator {
	case "==":
		output = reflect.DeepEqual(e.Attribute, e.Value)
		err = nil
	default:
		err = fmt.Errorf("unsupported operator %s", e.Operator)
	}

	return output, err
}

func (e Expression) IsEmpty() bool {
	return reflect.ValueOf(e).IsZero()
}
