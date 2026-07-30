package servicedata

import (
	"context"
	"fmt"
	"strings"
)

const auditEntityServiceDataKey = "service_data_key"

type Repository interface {
	Transactor
	CreateKey(ctx context.Context, key Key) (Key, error)
	Upsert(ctx context.Context, servicedata ServiceData) (ServiceData, error)
	GetKeyByURN(ctx context.Context, URN string) (Key, error)
	Get(ctx context.Context, filter Filter) ([]ServiceData, error)
	GetDistinctKeyValues(ctx context.Context, filter DistinctValueFilter) ([]any, error)
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context, err error) error
	Commit(ctx context.Context) error
}

type Key struct {
	ID          string
	URN         string
	ProjectID   string
	ProjectSlug string
	Name        string
	Description string
	ResourceID  string
}

type ServiceData struct {
	ID          string
	NamespaceID string
	EntityID    string
	Key         Key
	Value       any
}

type KeyLogData struct {
	Entity      string `mapstructure:"entity"`
	URN         string `mapstructure:"urn"`
	ProjectSlug string `mapstructure:"project_slug"`
	Key         string `mapstructure:"key"`
	Description string `mapstructure:"description"`
}

type Filter struct {
	ID        string
	Namespace string
	Entities  []string
	EntityIDs [][]string
	Project   string
}

// DistinctValueFilter describes a request for the distinct values found at a
// metadata path across all entities of a namespace.
//
// Path is dot-separated: the first segment is the service data key name and the
// remaining segments address a location inside that key's jsonb value, e.g.
// "functional_roles.roles" -> key "functional_roles", inner path "roles".
type DistinctValueFilter struct {
	Path      string
	Namespace string
	Project   string
}

// KeyName returns the service data key name (the first path segment).
func (f DistinctValueFilter) KeyName() string {
	key, _, _ := strings.Cut(f.Path, ".")
	return key
}

// JSONPath builds a SQL/JSON path expression from the segments after the key
// name. A trailing [*] flattens leaf arrays; lax mode (the default) auto-unwraps
// intermediate arrays. When there is no inner path it targets the whole value.
func (f DistinctValueFilter) JSONPath() string {
	_, rest, found := strings.Cut(f.Path, ".")
	var sb strings.Builder
	sb.WriteString("$")
	if found && rest != "" {
		for _, seg := range strings.Split(rest, ".") {
			sb.WriteString(`."`)
			// escape backslashes and quotes for the jsonpath string literal
			sb.WriteString(strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(seg))
			sb.WriteString(`"`)
		}
	}
	sb.WriteString("[*]")
	return sb.String()
}

func CreateURN(projectSlug, keyName string) string {
	return fmt.Sprintf("%s:servicedata_key:%s", projectSlug, keyName)
}

func (key Key) ToKeyLogData() KeyLogData {
	return KeyLogData{
		Entity:      auditEntityServiceDataKey,
		URN:         key.URN,
		ProjectSlug: key.ProjectSlug,
		Key:         key.Name,
		Description: key.Description,
	}
}
