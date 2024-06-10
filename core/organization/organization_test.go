package organization_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/pkg/metadata"
)

func TestOrganization_ToLogData(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		Slug      string
		Metadata  metadata.Metadata
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   organization.LogData
	}{

		{
			name: "should return log data",
			fields: fields{
				ID:   "1",
				Name: "name",
				Slug: "slug",
				Metadata: metadata.Metadata{
					"k1": "v1",
				},
			},
			want: organization.LogData{
				Entity: organization.AuditEntity,
				ID:     "1",
				Name:   "name",
				Slug:   "slug",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			organization := organization.Organization{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Slug:      tt.fields.Slug,
				Metadata:  tt.fields.Metadata,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if got := organization.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Organization.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
