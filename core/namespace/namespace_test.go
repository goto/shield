package namespace_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/namespace"
)

func TestNamespace_ToLogData(t *testing.T) {
	type fields struct {
		ID           string
		Name         string
		Backend      string
		ResourceType string
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   namespace.LogData
	}{

		{
			name: "should return log data",
			fields: fields{
				ID:           "1",
				Name:         "name",
				Backend:      "backend",
				ResourceType: "type",
			},
			want: namespace.LogData{
				Entity:       namespace.AuditEntity,
				ID:           "1",
				Name:         "name",
				Backend:      "backend",
				ResourceType: "type",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespace := namespace.Namespace{
				ID:           tt.fields.ID,
				Name:         tt.fields.Name,
				Backend:      tt.fields.Backend,
				ResourceType: tt.fields.ResourceType,
				CreatedAt:    tt.fields.CreatedAt,
				UpdatedAt:    tt.fields.UpdatedAt,
			}
			if got := namespace.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Namespace.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
