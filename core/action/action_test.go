package action_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/action"
)

func TestAction_ToActionLogData(t *testing.T) {
	timeNow := time.Now()
	type fields struct {
		ID          string
		Name        string
		NamespaceID string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   action.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				ID:          "1",
				Name:        "name",
				NamespaceID: "id",
				CreatedAt:   timeNow,
				UpdatedAt:   timeNow,
			},
			want: action.LogData{
				Entity:      action.AuditEntity,
				ID:          "1",
				Name:        "name",
				NamespaceID: "id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := action.Action{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				NamespaceID: tt.fields.NamespaceID,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if got := action.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Action.ToActionLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
