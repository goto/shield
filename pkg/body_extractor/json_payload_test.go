package body_extractor

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestJSONPayloadHandler_Extract(t *testing.T) {
	defaultTestMessage := map[string]any{
		"k1": "v1",
		"nested_k1": map[string]any{
			"k1_1": "v1",
			"k2_2": 1,
		},
	}

	tests := []struct {
		name          string
		key           string
		testMessage   any
		want          any
		wantErrString string
	}{
		{
			name:          "should return error if field not exist",
			key:           "x",
			testMessage:   defaultTestMessage,
			wantErrString: "failed to find field: x",
		},
		{
			name:        "should return value if field found",
			key:         "nested_k1.k2_2",
			testMessage: defaultTestMessage,
			want:        float64(1),
		},
	}
	for _, tt := range tests {
		h := JSONPayloadHandler{}
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			msg, err := json.Marshal(tt.testMessage)
			assert.NoError(t, err)

			testReader := io.NopCloser(bytes.NewBuffer(msg))
			extractedData, err := h.Extract(&testReader, tt.key)

			if tt.wantErrString != "" {
				assert.Equal(t, tt.wantErrString, err.Error())
			} else {
				if diff := cmp.Diff(tt.want, extractedData); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}
