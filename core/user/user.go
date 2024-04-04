package user

import (
	"context"
	"time"

	"github.com/goto/shield/pkg/metadata"
	"golang.org/x/exp/maps"
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

func (user User) ToUserAuditData() (map[string]string, error) {
	logData := map[string]string{
		"entity": AuditEntityUser,
		"name":   user.Name,
		"email":  user.Email,
	}
	userMetadata, err := user.Metadata.ToStringValueMap()
	if err != nil {
		return logData, err
	}

	maps.Copy(logData, userMetadata)
	return logData, nil
}

func (userMetadataKey UserMetadataKey) ToUserMetadataKey() map[string]string {
	return map[string]string{
		"entity":      AuditEntityUserMetadata,
		"key":         userMetadataKey.Key,
		"description": userMetadataKey.Description,
	}
}
