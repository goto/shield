package v1beta1

import (
	"context"

	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/internal/api"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"google.golang.org/grpc"
)

type RelationTransformer interface {
	TransformRelation(ctx context.Context, rlt relation.RelationV2) (relation.RelationV2, error)
}

type Handler struct {
	shieldv1beta1.UnimplementedShieldServiceServer
	orgService       OrganizationService
	projectService   ProjectService
	groupService     GroupService
	roleService      RoleService
	policyService    PolicyService
	userService      UserService
	namespaceService NamespaceService
	actionService    ActionService
	relationService  RelationService
	resourceService  ResourceService
	ruleService      RuleService
	relationAdapter  RelationTransformer
	checkAPILimit    int
}

func Register(ctx context.Context, s *grpc.Server, deps api.Deps, checkAPILimit int) error {
	s.RegisterService(
		&shieldv1beta1.ShieldService_ServiceDesc,
		&Handler{
			orgService:       deps.OrgService,
			projectService:   deps.ProjectService,
			groupService:     deps.GroupService,
			roleService:      deps.RoleService,
			policyService:    deps.PolicyService,
			userService:      deps.UserService,
			namespaceService: deps.NamespaceService,
			actionService:    deps.ActionService,
			relationService:  deps.RelationService,
			resourceService:  deps.ResourceService,
			ruleService:      deps.RuleService,
			relationAdapter:  deps.RelationAdapter,
			checkAPILimit:    checkAPILimit,
		},
	)

	return nil
}
