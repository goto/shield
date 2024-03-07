package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpression_IsEmpty(t *testing.T) {
	tests := []struct {
		name       string
		expression Expression
	}{
		{
			name:       "empty expression",
			expression: Expression{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.expression.IsEmpty())
		})
	}
}

func TestExpression_Evaluate(t *testing.T) {
	tests := []struct {
		name       string
		expression Expression
		wantOutput any
		wantNilErr bool
	}{
		{
			name: "comparison expression, A == A",
			expression: Expression{
				Attribute: "A",
				Operator:  "==",
				Value:     "A",
			},
			wantNilErr: true,
			wantOutput: true,
		},
		{
			name: "comparison expression, A == B",
			expression: Expression{
				Attribute: "A",
				Operator:  "==",
				Value:     "B",
			},
			wantNilErr: true,
			wantOutput: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.False(t, tt.expression.IsEmpty())
			output, err := tt.expression.Evaluate()
			if !tt.wantNilErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, output, tt.wantOutput)
		})
	}
}
