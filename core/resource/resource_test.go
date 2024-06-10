package resource_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/resource"
)

func TestResource_ToLogData(t *testing.T) {
	type fields struct {
		Idxa           string
		URN            string
		Name           string
		ProjectID      string
		OrganizationID string
		NamespaceID    string
		UserID         string
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   resource.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				URN:            "urn",
				Name:           "name",
				OrganizationID: "organization_id",
				ProjectID:      "project_id",
				NamespaceID:    "namespace_id",
				UserID:         "user_id",
			},
			want: resource.LogData{
				Entity:         resource.AuditEntity,
				URN:            "urn",
				Name:           "name",
				OrganizationID: "organization_id",
				ProjectID:      "project_id",
				NamespaceID:    "namespace_id",
				UserID:         "user_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := resource.Resource{
				Idxa:           tt.fields.Idxa,
				URN:            tt.fields.URN,
				Name:           tt.fields.Name,
				ProjectID:      tt.fields.ProjectID,
				OrganizationID: tt.fields.OrganizationID,
				NamespaceID:    tt.fields.NamespaceID,
				UserID:         tt.fields.UserID,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			if got := resource.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
