package servicedata_test

import (
	"testing"

	"github.com/goto/shield/core/servicedata"
	"github.com/stretchr/testify/assert"
)

func TestDistinctValueFilter_KeyName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "KeyOnly", path: "functional_roles", want: "functional_roles"},
		{name: "KeyWithSingleSubPath", path: "functional_roles.roles", want: "functional_roles"},
		{name: "KeyWithNestedSubPath", path: "functional_roles.parent_organizations.organization_id", want: "functional_roles"},
		{name: "EmptyPath", path: "", want: ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := servicedata.DistinctValueFilter{Path: tt.path}
			assert.Equal(t, tt.want, f.KeyName())
		})
	}
}

func TestDistinctValueFilter_JSONPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "KeyOnlyTargetsWholeValue", path: "functional_roles", want: `$[*]`},
		{name: "SingleSubPath", path: "functional_roles.roles", want: `$."roles"[*]`},
		{name: "NestedSubPath", path: "functional_roles.parent_organizations.organization_id", want: `$."parent_organizations"."organization_id"[*]`},
		{name: "EscapesQuotesAndBackslashes", path: `key.we"ird\seg`, want: `$."we\"ird\\seg"[*]`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := servicedata.DistinctValueFilter{Path: tt.path}
			assert.Equal(t, tt.want, f.JSONPath())
		})
	}
}
