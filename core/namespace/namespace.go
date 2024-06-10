package namespace

import (
	"context"
	"fmt"
	"time"
)

const auditEntity = "namespace"

type Repository interface {
	Get(ctx context.Context, id string) (Namespace, error)
	Upsert(ctx context.Context, ns Namespace) (Namespace, error)
	List(ctx context.Context) ([]Namespace, error)
	Update(ctx context.Context, ns Namespace) (Namespace, error)
}

type Namespace struct {
	ID           string
	Name         string
	Backend      string
	ResourceType string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type NamspaceLogData struct {
	Entity       string `mapstructure:"entity"`
	ID           string `mapstructure:"id"`
	Name         string `mapstructure:"name"`
	Backend      string `mapstructure:"backend"`
	ResourceType string `mapstructure:"resource_type"`
}

func (namespace Namespace) ToNameSpaceLogData() NamspaceLogData {
	return NamspaceLogData{
		Entity:       auditEntity,
		ID:           namespace.ID,
		Name:         namespace.Name,
		Backend:      namespace.Backend,
		ResourceType: namespace.ResourceType,
	}
}

func strListHas(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsSystemNamespaceID(nsID string) bool {
	return strListHas(systemIdsDefinition, nsID)
}

func CreateID(backend, resourceType string) string {
	if resourceType == "" {
		return backend
	}

	return fmt.Sprintf("%s/%s", backend, resourceType)
}
