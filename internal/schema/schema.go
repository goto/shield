package schema

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/policy"
	"github.com/goto/shield/core/role"
	"github.com/goto/shield/core/user"

	"github.com/goto/salt/log"
	"golang.org/x/exp/maps"
)

type NamespaceType string

var (
	SystemNamespace        NamespaceType = "system_namespace"
	ResourceGroupNamespace NamespaceType = "resource_group_namespace"

	ErrMigration     = errors.New("error in migrating authz schema")
	ErrInvalidDetail = errors.New("error in schema config")
)

const (
	RESOURCES_CONFIG_STORAGE_PG   = "postgres"
	RESOURCES_CONFIG_STORAGE_GS   = "gs"
	RESOURCES_CONFIG_STORAGE_FILE = "file"
	RESOURCES_CONFIG_STORAGE_MEM  = "mem"
)

type Config struct {
	ID        uint32
	Name      string
	Config    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AppConfig struct {
	ConfigStorage string
}

type InheritedNamespace struct {
	Name        string
	NamespaceId string
}

type NamespaceConfig struct {
	InheritedNamespaces []InheritedNamespace
	Type                NamespaceType
	Roles               map[string][]string
	Permissions         map[string][]string
}

type NamespaceConfigMapType map[string]NamespaceConfig

type NamespaceService interface {
	Upsert(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error)
}

type RoleService interface {
	Upsert(ctx context.Context, toCreate role.Role) (role.Role, error)
}

type PolicyService interface {
	Upsert(ctx context.Context, policy *policy.Policy) ([]policy.Policy, error)
}

type ActionService interface {
	Upsert(ctx context.Context, action action.Action) (action.Action, error)
}

type PGRepository interface {
	Transactor
	UpsertConfig(ctx context.Context, name string, config NamespaceConfigMapType) (Config, error)
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context, err error) error
	Commit(ctx context.Context) error
}

type FileService interface {
	GetSchema(ctx context.Context) (NamespaceConfigMapType, error)
}

type AuthzEngine interface {
	WriteSchema(ctx context.Context, schema NamespaceConfigMapType) error
}

type UserRepository interface {
	Create(ctx context.Context, usr user.User) (user.User, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
}

type SchemaMigrationConfig struct {
	DefaultSystemEmail      string
	BootstrapServiceDataKey bool
}

type SchemaService struct {
	logger                log.Logger
	appConfig             AppConfig
	schemaConfig          FileService
	pgRepository          PGRepository
	namespaceService      NamespaceService
	roleService           RoleService
	actionService         ActionService
	policyService         PolicyService
	authzEngine           AuthzEngine
	userRepository        UserRepository
	schemaMigrationConfig SchemaMigrationConfig
}

func NewSchemaMigrationService(
	logger log.Logger,
	appConfig AppConfig,
	schemaConfig FileService,
	pgRepository PGRepository,
	namespaceService NamespaceService,
	roleService RoleService,
	actionService ActionService,
	policyService PolicyService,
	authzEngine AuthzEngine,
	userRepository UserRepository,
	schemaMigrationConfig SchemaMigrationConfig,
) *SchemaService {
	return &SchemaService{
		logger:                logger,
		appConfig:             appConfig,
		schemaConfig:          schemaConfig,
		pgRepository:          pgRepository,
		namespaceService:      namespaceService,
		roleService:           roleService,
		actionService:         actionService,
		policyService:         policyService,
		authzEngine:           authzEngine,
		userRepository:        userRepository,
		schemaMigrationConfig: schemaMigrationConfig,
	}
}

func (s SchemaService) RunMigrations(ctx context.Context) error {
	defaultUser := user.User{
		Name:  s.schemaMigrationConfig.DefaultSystemEmail,
		Email: s.schemaMigrationConfig.DefaultSystemEmail,
	}

	if _, err := s.userRepository.GetByEmail(ctx, defaultUser.Email); err != nil {
		// creating predefined user for log activity if user not exist
		if err == user.ErrNotExist {
			if _, err := s.userRepository.Create(ctx, defaultUser); err != nil {
				return err
			}
		} else {
			// if other error occurred, return the error
			return err
		}
	}

	namespaceConfigMap, err := s.schemaConfig.GetSchema(ctx)
	if err != nil {
		return err
	}

	// setting context with predefined user email
	ctx = user.SetContextWithEmail(ctx, defaultUser.Email)

	// append service data key schema if configured to be bootstrapped automatically
	if s.schemaMigrationConfig.BootstrapServiceDataKey {
		PreDefinedSystemNamespaceConfig[ServiceDataKeyNamespace] = ServiceDataKeyConfig
	}

	// combining predefined and configured namespaces
	namespaceConfigMap = MergeNamespaceConfigMap(namespaceConfigMap, PreDefinedSystemNamespaceConfig)

	// adding predefined roles and permissions for resource group namespaces
	for n, nc := range namespaceConfigMap {
		if nc.Type == ResourceGroupNamespace {
			namespaceConfigMap = MergeNamespaceConfigMap(namespaceConfigMap, NamespaceConfigMapType{
				n: PreDefinedResourceGroupNamespaceConfig,
			})
		}
	}

	// spiceDBSchema := GenerateSchema(namespaceConfigMap)

	// iterate over namespace
	for namespaceId, v := range namespaceConfigMap {
		// create namespace
		backend := ""
		resourceType := ""
		if v.Type == ResourceGroupNamespace {
			st := strings.Split(namespaceId, "/")
			backend = st[0]
			resourceType = st[1]
		}
		s.logger.Info(fmt.Sprintf("create namespace %s", namespaceId))
		_, err := s.namespaceService.Upsert(ctx, namespace.Namespace{
			ID:           namespaceId,
			Name:         namespaceId,
			Backend:      backend,
			ResourceType: resourceType,
		})
		if err != nil {
			return fmt.Errorf("%w: %s", ErrMigration, err.Error())
		}

		// create roles
		for roleId, principals := range v.Roles {
			s.logger.Info(fmt.Sprintf("create role %s with principals %s under namespace %s", roleId, principals, namespaceId))
			_, err := s.roleService.Upsert(ctx, role.Role{
				ID:          fmt.Sprintf("%s:%s", namespaceId, roleId),
				Name:        roleId,
				Types:       principals,
				NamespaceID: namespaceId,
			})
			if err != nil {
				return fmt.Errorf("%w: %s", ErrMigration, err.Error())
			}
		}

		// create role for inherited namespaces
		for _, ins := range v.InheritedNamespaces {
			s.logger.Info(fmt.Sprintf("create role %s from inherited namespace %s under namespace %s", ins.Name, ins.NamespaceId, namespaceId))
			_, err := s.roleService.Upsert(ctx, role.Role{
				ID:          fmt.Sprintf("%s:%s", namespaceId, ins.Name),
				Name:        ins.Name,
				Types:       []string{ins.NamespaceId},
				NamespaceID: namespaceId,
			})
			if err != nil {
				return fmt.Errorf("%w: %s", ErrMigration, err.Error())
			}
		}

		// create actions
		// IMP: we should depreciate actions with principals
		for actionId := range v.Permissions {
			s.logger.Info(fmt.Sprintf("create action %s under namespace %s", actionId, namespaceId))
			_, err := s.actionService.Upsert(ctx, action.Action{
				ID:          fmt.Sprintf("%s.%s", actionId, namespaceId),
				Name:        actionId,
				NamespaceID: namespaceId,
			})
			if err != nil {
				return fmt.Errorf("%w: %s", ErrMigration, err.Error())
			}
		}
	}

	for namespaceId, v := range namespaceConfigMap {
		// create policies
		for actionId, roles := range v.Permissions {
			for _, r := range roles {
				transformedRole, err := getRoleAndPrincipal(r, namespaceId)
				if err != nil {
					return fmt.Errorf("%w: %s", ErrMigration, err.Error())
				}

				if _, ok := namespaceConfigMap[GetNamespace(transformedRole.NamespaceID)].Roles[transformedRole.ID]; !ok {
					return fmt.Errorf("%w: role %s not associated with namespace: %s", ErrInvalidDetail, transformedRole.ID, transformedRole.NamespaceID)
				}

				roleId := GetRoleID(GetNamespace(transformedRole.NamespaceID), transformedRole.ID)
				actionId := fmt.Sprintf("%s.%s", actionId, namespaceId)
				s.logger.Info(fmt.Sprintf("create policy for role %s on namespace %s with action %s", roleId, namespaceId, actionId))
				_, err = s.policyService.Upsert(ctx, &policy.Policy{
					RoleID:      roleId,
					NamespaceID: namespaceId,
					ActionID:    actionId,
				})
				if err != nil {
					return fmt.Errorf("%w: %s", ErrMigration, err.Error())
				}
			}
		}
	}

	if err = s.authzEngine.WriteSchema(ctx, namespaceConfigMap); err != nil {
		return fmt.Errorf("%w: %s", ErrMigration, err.Error())
	}

	return nil
}

func (s SchemaService) UpsertConfig(ctx context.Context, name string, config string) (Config, error) {
	if strings.TrimSpace(name) == "" {
		return Config{}, ErrInvalidDetail
	}

	if strings.TrimSpace(config) == "" {
		return Config{}, ErrInvalidDetail
	}

	resourceConfig, err := ParseConfigYaml([]byte(config))
	if err != nil {
		return Config{}, ErrInvalidDetail
	}

	configMap := make(NamespaceConfigMapType)
	for k, v := range resourceConfig {
		if v.Type == "resource_group" {
			configMap = MergeNamespaceConfigMap(configMap, GetNamespacesForResourceGroup(k, v))
		} else {
			configMap = MergeNamespaceConfigMap(GetNamespaceFromConfig(k, v.Roles, v.Permissions), configMap)
		}
	}

	ctx = s.pgRepository.WithTransaction(ctx)

	res, err := s.pgRepository.UpsertConfig(ctx, name, configMap)
	if err != nil {
		if txErr := s.pgRepository.Rollback(ctx, err); txErr != nil {
			return Config{}, err
		}
		return Config{}, err
	}

	if s.appConfig.ConfigStorage == RESOURCES_CONFIG_STORAGE_PG {
		if err := s.RunMigrations(ctx); err != nil {
			if txErr := s.pgRepository.Rollback(ctx, err); txErr != nil {
				return Config{}, err
			}
			return Config{}, err
		}
	}

	err = s.pgRepository.Commit(ctx)
	if err != nil {
		return Config{}, err
	}

	return res, nil
}

func MergeNamespaceConfigMap(smallMap, largeMap NamespaceConfigMapType) NamespaceConfigMapType {
	combinedMap := make(NamespaceConfigMapType)
	maps.Copy(combinedMap, smallMap)
	for namespaceName, namespaceConfig := range largeMap {
		if _, ok := combinedMap[namespaceName]; !ok {
			combinedMap[namespaceName] = NamespaceConfig{
				Roles:       make(map[string][]string),
				Permissions: make(map[string][]string),
			}
		}

		for roleName := range namespaceConfig.Roles {
			if _, ok := combinedMap[namespaceName].Roles[roleName]; !ok {
				combinedMap[namespaceName].Roles[roleName] = namespaceConfig.Roles[roleName]
			} else {
				combinedMap[namespaceName].Roles[roleName] = AppendIfUnique(namespaceConfig.Roles[roleName], combinedMap[namespaceName].Roles[roleName])
			}
		}

		for permissionName := range namespaceConfig.Permissions {
			combinedMap[namespaceName].Permissions[permissionName] = AppendIfUnique(namespaceConfig.Permissions[permissionName], combinedMap[namespaceName].Permissions[permissionName])
		}

		if value, ok := combinedMap[namespaceName]; ok {
			value.Type = namespaceConfig.Type
			value.InheritedNamespaces = AppendIfUnique(value.InheritedNamespaces, namespaceConfig.InheritedNamespaces)
			combinedMap[namespaceName] = value
		}
	}

	return combinedMap
}

func NewSchemaMigrationConfig(defaultSystemEmail string, bootstrapServiceDataKey bool) SchemaMigrationConfig {
	return SchemaMigrationConfig{
		DefaultSystemEmail:      defaultSystemEmail,
		BootstrapServiceDataKey: bootstrapServiceDataKey,
	}
}
