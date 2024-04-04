package policy

import (
	"context"
)

const AuditEntity = "policy"

type Repository interface {
	Get(ctx context.Context, id string) (Policy, error)
	List(ctx context.Context) ([]Policy, error)
	Create(ctx context.Context, pol Policy) (string, error)
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

func (policy Policy) ToPolicyLogData(policyId string) map[string]string {
	return map[string]string{
		"entity":      AuditEntity,
		"id":          policyId,
		"roleId":      policy.RoleID,
		"namespaceId": policy.NamespaceID,
		"actionId":    policy.ActionID,
	}
}
