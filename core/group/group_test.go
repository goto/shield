package group_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/pkg/metadata"
)

func TestGroup_ToLogData(t *testing.T) {
	timeNow := time.Now()
	type fields struct {
		ID             string
		Name           string
		Slug           string
		OrganizationID string
		Metadata       metadata.Metadata
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   group.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				ID:             "1",
				Name:           "name",
				Slug:           "slug",
				OrganizationID: "id",
				Metadata: metadata.Metadata{
					"k1": "v1",
				},
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			want: group.LogData{
				Entity:         group.AuditEntity,
				ID:             "1",
				Name:           "name",
				Slug:           "slug",
				OrganizationID: "id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := group.Group{
				ID:             tt.fields.ID,
				Name:           tt.fields.Name,
				Slug:           tt.fields.Slug,
				OrganizationID: tt.fields.OrganizationID,
				Metadata:       tt.fields.Metadata,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			if got := group.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
