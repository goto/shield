package relation

import (
	"context"
	"time"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/role"
)

const (
	AuditEntity        = "relation"
	AuditEntitySubject = "relation_subject"
)

type Repository interface {
	Get(ctx context.Context, id string) (RelationV2, error)
	Create(ctx context.Context, relation RelationV2) (RelationV2, error)
	List(ctx context.Context) ([]RelationV2, error)
	Update(ctx context.Context, toUpdate Relation) (Relation, error)
	DeleteByID(ctx context.Context, id string) error
	GetByFields(ctx context.Context, rel RelationV2) (RelationV2, error)
}

type AuthzRepository interface {
	Add(ctx context.Context, rel Relation) error
	Check(ctx context.Context, rel Relation, act action.Action) (bool, error)
	BulkCheck(ctx context.Context, rels []Relation, acts []action.Action) ([]Permission, error)
	DeleteV2(ctx context.Context, rel RelationV2) error
	DeleteSubjectRelations(ctx context.Context, resourceType, optionalResourceID string) error
	AddV2(ctx context.Context, rel RelationV2) error
	LookupResources(ctx context.Context, resourceType, permission, subjectType, subjectID string) ([]string, error)
	CheckIsPublic(ctx context.Context, rel Relation, act action.Action) (bool, error)
}

type Relation struct {
	ID                 string
	SubjectNamespace   namespace.Namespace
	SubjectNamespaceID string `json:"subject_namespace_id"`
	SubjectID          string `json:"subject_id"`
	SubjectRoleID      string `json:"subject_role_id"`
	ObjectNamespace    namespace.Namespace
	ObjectNamespaceID  string `json:"object_namespace_id"`
	ObjectID           string `json:"object_id"`
	Role               role.Role
	RoleID             string       `json:"role_id"`
	RelationType       RelationType `json:"role_type"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Object struct {
	ID          string
	NamespaceID string
}

type Subject struct {
	ID        string
	Namespace string
	RoleID    string
}

type RelationV2 struct {
	ID        string
	Object    Object
	Subject   Subject
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RelationType string

var RelationTypes = struct {
	Role      RelationType
	Namespace RelationType
}{
	Role:      "role",
	Namespace: "namespace",
}

type LogData struct {
	Entity           string `mapstructure:"entity"`
	ID               string `mapstructure:"id"`
	ObjectID         string `mapstructure:"object_id"`
	ObjectNamespace  string `mapstructure:"object_namespace"`
	SubjectID        string `mapstructure:"subject_id"`
	SubjectNamespace string `mapstructure:"subject_namespace"`
	RoleID           string `mapstructure:"role"`
}

type SubjectLogData struct {
	Entity             string `mapstructure:"entity"`
	ResourceType       string `mapstructure:"resource_type"`
	OptionalResourceID string `mapstructure:"optional_resource_id"`
}

type Permission struct {
	ObjectID        string `mapstructure:"object_id"`
	ObjectNamespace string `mapstructure:"object_namespace"`
	Permission      string `mapstructure:"permission"`
	Allowed         bool   `mapstructure:"allowed"`
}

func (relation RelationV2) ToLogData() LogData {
	return LogData{
		Entity:           AuditEntity,
		ID:               relation.ID,
		ObjectID:         relation.Object.ID,
		ObjectNamespace:  relation.Object.NamespaceID,
		SubjectID:        relation.Subject.ID,
		SubjectNamespace: relation.Subject.Namespace,
		RoleID:           relation.Subject.RoleID,
	}
}

func ToSubjectLogData(resourceType, optionalResourceID string) SubjectLogData {
	return SubjectLogData{
		Entity:             AuditEntitySubject,
		ResourceType:       resourceType,
		OptionalResourceID: optionalResourceID,
	}
}
