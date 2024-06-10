package policy_test

import (
	"reflect"
	"testing"

	"github.com/goto/shield/core/policy"
)

func TestPolicy_ToLogData(t *testing.T) {
	type fields struct {
		ID          string
		RoleID      string
		NamespaceID string
		ActionID    string
		PolicyID    string
	}
	tests := []struct {
		name   string
		fields fields
		want   policy.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				ID:          "1",
				RoleID:      "role",
				NamespaceID: "namespace",
				ActionID:    "action",
			},
			want: policy.LogData{
				Entity:      policy.AuditEntity,
				ID:          "1",
				RoleID:      "role",
				NamespaceID: "namespace",
				ActionID:    "action",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := policy.Policy{
				ID:          tt.fields.ID,
				RoleID:      tt.fields.RoleID,
				NamespaceID: tt.fields.NamespaceID,
				ActionID:    tt.fields.ActionID,
			}
			if got := policy.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Policy.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
