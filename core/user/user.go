package user

import (
	"context"
	"time"

	"github.com/goto/shield/pkg/metadata"
)

const (
	AuditEntityUser         = "user"
	AuditEntityUserMetadata = "user_metadata_key"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByIDs(ctx context.Context, userIds []string) ([]User, error)
	Create(ctx context.Context, user User) (User, error)
	List(ctx context.Context, flt Filter) ([]User, error)
	UpdateByID(ctx context.Context, toUpdate User) (User, error)
	UpdateByEmail(ctx context.Context, toUpdate User) (User, error)
	CreateMetadataKey(ctx context.Context, key UserMetadataKey) (UserMetadataKey, error)
}

type User struct {
	ID        string
	Name      string
	Email     string
	Metadata  metadata.Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserMetadataKey struct {
	Key         string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PagedUsers struct {
	Count int32
	Users []User
}

type UserLogData struct {
	Entity string `mapstructure:"entity"`
	Name   string `mapstructure:"name"`
	Email  string `mapstructure:"email"`
}

type UserMetadataKeyLogData struct {
	Entity      string `mapstructure:"entity"`
	Key         string `mapstructure:"key"`
	Description string `mapstructure:"description"`
}

func (user User) ToUserLogData() UserLogData {
	return UserLogData{
		Entity: AuditEntityUser,
		Name:   user.Name,
		Email:  user.Email,
	}
}

func (userMetadataKey UserMetadataKey) ToUserMetadataKeyLogData() UserMetadataKeyLogData {
	return UserMetadataKeyLogData{
		Entity:      AuditEntityUserMetadata,
		Key:         userMetadataKey.Key,
		Description: userMetadataKey.Description,
	}
}
