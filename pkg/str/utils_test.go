package str

import "testing"

func TestDefaultStringIfEmpty(t *testing.T) {
	tests := []struct {
		name          string
		str           string
		defaultString string
		want          string
	}{
		{
			name:          "return str if not empty",
			str:           "the str",
			defaultString: "the default",
			want:          "the str",
		},
		{
			name:          "return defaultStr if empty",
			str:           "",
			defaultString: "the default",
			want:          "the default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultStringIfEmpty(tt.str, tt.defaultString); got != tt.want {
				t.Errorf("DefaultStringIfEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
