package relation_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/relation"
)

func TestRelationV2_ToLogData(t *testing.T) {
	type fields struct {
		ID        string
		Object    relation.Object
		Subject   relation.Subject
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   relation.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				ID: "1",
				Object: relation.Object{
					ID:          "objID",
					NamespaceID: "objNS",
				},
				Subject: relation.Subject{
					ID:        "subID",
					Namespace: "subNS",
					RoleID:    "roleID",
				},
			},
			want: relation.LogData{
				Entity:           relation.AuditEntity,
				ID:               "id",
				ObjectID:         "objID",
				ObjectNamespace:  "objNS",
				SubjectID:        "subID",
				SubjectNamespace: "subNS",
				RoleID:           "roleID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relation := relation.RelationV2{
				ID:        tt.fields.ID,
				Object:    tt.fields.Object,
				Subject:   tt.fields.Subject,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if got := relation.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RelationV2.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSubjectLogData(t *testing.T) {
	type args struct {
		resourceType       string
		optionalResourceID string
	}
	tests := []struct {
		name string
		args args
		want relation.SubjectLogData
	}{

		{
			name: "should return log data",
			args: args{
				resourceType:       "type",
				optionalResourceID: "resID",
			},
			want: relation.SubjectLogData{
				Entity:             relation.AuditEntitySubject,
				ResourceType:       "type",
				OptionalResourceID: "resID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := relation.ToSubjectLogData(tt.args.resourceType, tt.args.optionalResourceID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSubjectLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
