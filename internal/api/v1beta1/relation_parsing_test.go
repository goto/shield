package v1beta1

import "testing"

func TestSubjectParsingCompatibility(t *testing.T) {
	for _, tc := range []struct{ input, principal, id string }{
		{"user:person@example.test", "user", "person@example.test"},
		{"user:existing:extra", "user", "existing"},
		{"user:", "user", ""},
		{"missing-delimiter", "", ""},
	} {
		principal, id := extractSubjectFromPrincipal(tc.input)
		if principal != tc.principal || id != tc.id {
			t.Fatalf("unexpected subject parse for %q", tc.input)
		}
	}
}
