package policy

import (
	"context"
)

const auditEntity = "policy"

type Repository interface {
	Get(ctx context.Context, id string) (Policy, error)
	List(ctx context.Context) ([]Policy, error)
	Upsert(ctx context.Context, pol Policy) (string, error)
	Update(ctx context.Context, pol Policy) (string, error)
}

type AuthzRepository interface {
	Add(ctx context.Context, policies []Policy) error
}

type Policy struct {
	ID          string
	RoleID      string
	NamespaceID string
	ActionID    string
}

type Filters struct {
	NamespaceID string
}

type PolicyLogData struct {
	Entity      string `mapstructure:"entity"`
	ID          string `mapstructure:"id"`
	RoleID      string `mapstructure:"role_id"`
	NamespaceID string `mapstructure:"namespace_id"`
	ActionID    string `mapstructure:"action_id"`
}

func (policy Policy) ToPolicyLogData(policyId string) PolicyLogData {
	return PolicyLogData{
		Entity:      auditEntity,
		ID:          policyId,
		RoleID:      policy.RoleID,
		NamespaceID: policy.NamespaceID,
		ActionID:    policy.ActionID,
	}
}
