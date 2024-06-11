package user_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/goto/shield/core/user"
	"github.com/goto/shield/pkg/metadata"
)

func TestUser_ToLogData(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		Email     string
		Metadata  metadata.Metadata
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   user.LogData
	}{
		{
			name: "should return log data",
			fields: fields{
				ID:    "id",
				Name:  "name",
				Email: "email",
			},
			want: user.LogData{
				Entity: user.AuditEntity,
				Name:   "name",
				Email:  "email",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := user.User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Email:     tt.fields.Email,
				Metadata:  tt.fields.Metadata,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if got := user.ToLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.ToLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserMetadataKey_ToMetadataKeyLogData(t *testing.T) {
	type fields struct {
		Key         string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   user.MetadataKeyLogData
	}{
		{
			name: "should return log data",
			fields: fields{
				Key:         "key",
				Description: "description",
			},
			want: user.MetadataKeyLogData{
				Entity:      user.AuditEntityMetadata,
				Key:         "key",
				Description: "description",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMetadataKey := user.UserMetadataKey{
				Key:         tt.fields.Key,
				Description: tt.fields.Description,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if got := userMetadataKey.ToMetadataKeyLogData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserMetadataKey.ToMetadataKeyLogData() = %v, want %v", got, tt.want)
			}
		})
	}
}
