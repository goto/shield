package user

import (
	"context"
	"time"

	"github.com/goto/shield/pkg/metadata"
)

const (
	AuditEntity         = "user"
	AuditEntityMetadata = "user_metadata_key"
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
	DeleteByEmail(ctx context.Context, email string, emailTag string) error
	DeleteById(ctx context.Context, id string, emailTag string) error
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

type Config struct {
	InactiveEmailTag string
}

type PagedUsers struct {
	Count int32
	Users []User
}

type LogData struct {
	Entity string `mapstructure:"entity"`
	Name   string `mapstructure:"name"`
	Email  string `mapstructure:"email"`
}

type MetadataKeyLogData struct {
	Entity      string `mapstructure:"entity"`
	Key         string `mapstructure:"key"`
	Description string `mapstructure:"description"`
}

func (user User) ToLogData() LogData {
	return LogData{
		Entity: AuditEntity,
		Name:   user.Name,
		Email:  user.Email,
	}
}

func (userMetadataKey UserMetadataKey) ToMetadataKeyLogData() MetadataKeyLogData {
	return MetadataKeyLogData{
		Entity:      AuditEntityMetadata,
		Key:         userMetadataKey.Key,
		Description: userMetadataKey.Description,
	}
}
