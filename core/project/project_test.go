package project_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/pkg/metadata"
)

func TestProject_ToLogData(t *testing.T) {
	type fields struct {
		ID           string
		Name         string
		Slug         string
		Organization organization.Organization
		Metadata     metadata.Metadata
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   project.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				ID:   "1",
				Name: "name",
				Slug: "slug",
				Organization: organization.Organization{
					ID: "orgID",
				},
			},
			want: project.LogData{
				Entity:         project.AuditEntity,
				ID:             "1",
				Name:           "name",
				Slug:           "slug",
				OrganizationID: "orgID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := project.Project{
				ID:           tt.fields.ID,
				Name:         tt.fields.Name,
				Slug:         tt.fields.Slug,
				Organization: tt.fields.Organization,
				Metadata:     tt.fields.Metadata,
				CreatedAt:    tt.fields.CreatedAt,
				UpdatedAt:    tt.fields.UpdatedAt,
			}
			if got := project.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
