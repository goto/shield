package config

import (
	"fmt"
	"maps"

	"github.com/goto/shield/internal/schema"
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

type Config struct {
	Type string `yaml:"type" json:"type"`

	ResourceTypes []ResourceTypeConfig `yaml:"resource_types" json:"resource_types,omitempty"`

	Roles       []RoleConfig        `yaml:"roles" json:"roles,omitempty"`
	Permissions []PermissionsConfig `yaml:"permissions" json:"permissions,omitempty"`
}

type ConfigYAML []map[string]Config

func ParseConfigYaml(fileBytes []byte) (map[string]Config, error) {
	var config map[string]Config
	if err := yaml.Unmarshal(fileBytes, &config); err != nil {
		return map[string]Config{}, err
	}
	return config, nil
}

func GetNamespacesForResourceGroup(name string, c Config) schema.NamespaceConfigMapType {
	namespaceConfig := schema.NamespaceConfigMapType{}

	for _, v := range c.ResourceTypes {
		maps.Copy(namespaceConfig, GetNamespaceFromConfig(name, v.Roles, v.Permissions, v.Name))
	}

	return namespaceConfig
}

func GetNamespaceFromConfig(name string, rolesConfigs []RoleConfig, permissionConfigs []PermissionsConfig, resourceType ...string) schema.NamespaceConfigMapType {
	tnc := schema.NamespaceConfig{
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
		tnc.Type = schema.SystemNamespace
	} else {
		tnc.Type = schema.ResourceGroupNamespace
		name = fmt.Sprintf("%s/%s", name, resourceType[0])
	}

	return schema.NamespaceConfigMapType{name: tnc}
}
