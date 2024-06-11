package role_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/role"
	"github.com/goto/shield/pkg/metadata"
)

func TestRole_ToLogData(t *testing.T) {
	type fields struct {
		ID          string
		Name        string
		Types       []string
		NamespaceID string
		Metadata    metadata.Metadata
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   role.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				ID:          "id",
				Name:        "name",
				Types:       []string{"types"},
				NamespaceID: "namespace_id",
			},
			want: role.LogData{
				Entity:      role.AuditEntity,
				ID:          "id",
				Name:        "name",
				Types:       []string{"types"},
				NamespaceID: "namespace_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := role.Role{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				Types:       tt.fields.Types,
				NamespaceID: tt.fields.NamespaceID,
				Metadata:    tt.fields.Metadata,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if got := role.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Role.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
