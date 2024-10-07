package schema

import (
	"fmt"
	"maps"

	"gopkg.in/yaml.v2"
)

type RoleConfig struct {
	Name       string   `yaml:"name" json:"name"`
	Principals []string `yaml:"principals" json:"principals"`
}

type PermissionsConfig struct {
	Name  string   `yaml:"name" json:"name"`
	Roles []string `yaml:"roles" json:"roles"`
}

type ResourceTypeConfig struct {
	Name        string              `yaml:"name" json:"name"`
	Roles       []RoleConfig        `yaml:"roles" json:"roles"`
	Permissions []PermissionsConfig `yaml:"permissions" json:"permissions"`
}

type ResourceConfig struct {
	Type string `yaml:"type" json:"type"`

	ResourceTypes []ResourceTypeConfig `yaml:"resource_types" json:"resource_types,omitempty"`

	Roles       []RoleConfig        `yaml:"roles" json:"roles,omitempty"`
	Permissions []PermissionsConfig `yaml:"permissions" json:"permissions,omitempty"`
}

type ConfigYAML []map[string]ResourceConfig

func ParseConfigYaml(fileBytes []byte) (map[string]ResourceConfig, error) {
	var config map[string]ResourceConfig
	if err := yaml.Unmarshal(fileBytes, &config); err != nil {
		return map[string]ResourceConfig{}, err
	}
	return config, nil
}

func GetNamespacesForResourceGroup(name string, c ResourceConfig) NamespaceConfigMapType {
	namespaceConfig := NamespaceConfigMapType{}

	for _, v := range c.ResourceTypes {
		maps.Copy(namespaceConfig, GetNamespaceFromConfig(name, v.Roles, v.Permissions, v.Name))
	}

	return namespaceConfig
}

func GetNamespaceFromConfig(name string, rolesConfigs []RoleConfig, permissionConfigs []PermissionsConfig, resourceType ...string) NamespaceConfigMapType {
	tnc := NamespaceConfig{
		Roles:       make(map[string][]string),
		Permissions: make(map[string][]string),
	}

	for _, v1 := range rolesConfigs {
		tnc.Roles[v1.Name] = v1.Principals
	}

	for _, v2 := range permissionConfigs {
		tnc.Permissions[v2.Name] = v2.Roles
	}

	if len(resourceType) == 0 {
		tnc.Type = SystemNamespace
	} else {
		tnc.Type = ResourceGroupNamespace
		name = fmt.Sprintf("%s/%s", name, resourceType[0])
	}

	return NamespaceConfigMapType{name: tnc}
}
