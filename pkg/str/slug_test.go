package str

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		options SlugifyOptions
		want    string
	}{
		{
			name: "should slugify all parameters",
			str:  "a word to-slugify#fragment:colon",
			want: "a_word_to_slugify_fragment_colon",
		},
		{
			name: "should slugify all parameters except hyphen",
			str:  "a word to-slugify#fragment:colon",
			want: "a_word_to-slugify_fragment_colon",
			options: SlugifyOptions{
				KeepHyphen: true,
			},
		},
		{
			name: "should slugify all parameters except colon",
			str:  "a word to-slugify#fragment:colon",
			want: "a_word_to_slugify_fragment:colon",
			options: SlugifyOptions{
				KeepColon: true,
			},
		},
		{
			name: "should slugify all parameters except hash",
			str:  "a word to-slugify#fragment:colon",
			want: "a_word_to_slugify#fragment_colon",
			options: SlugifyOptions{
				KeepHash: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.str, tt.options); got != tt.want {
				t.Errorf("Slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}
