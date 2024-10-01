package blob

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/resource/config"
	"github.com/goto/shield/internal/schema"

	"github.com/pkg/errors"
	"gocloud.dev/blob"
)

type SchemaConfig struct {
	bucket Bucket
	config schema.NamespaceConfigMapType
}

func NewSchemaConfigRepository(b Bucket) *SchemaConfig {
	return &SchemaConfig{bucket: b}
}

func (s *SchemaConfig) GetSchema(ctx context.Context) (schema.NamespaceConfigMapType, error) {
	configMap := make(schema.NamespaceConfigMapType)
	if s.config != nil {
		return s.config, nil
	}

	configFromFiles, err := s.readYAMLFiles(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range configFromFiles {
		for k, v := range c {
			if v.Type == "resource_group" {
				configMap = schema.MergeNamespaceConfigMap(configMap, config.GetNamespacesForResourceGroup(k, v))
			} else {
				configMap = schema.MergeNamespaceConfigMap(config.GetNamespaceFromConfig(k, v.Roles, v.Permissions), configMap)
			}
		}
	}

	s.config = configMap

	return configMap, nil
}

func (s *SchemaConfig) readYAMLFiles(ctx context.Context) (config.ConfigYAML, error) {
	configYAMLs := make(config.ConfigYAML, 0)

	// iterate over bucket files, only read .yml & .yaml files
	it := s.bucket.List(&blob.ListOptions{})
	for {
		obj, err := it.Next(ctx)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if obj.IsDir {
			continue
		}
		if !(strings.HasSuffix(obj.Key, ".yaml") || strings.HasSuffix(obj.Key, ".yml")) {
			continue
		}
		fileBytes, err := s.bucket.ReadAll(ctx, obj.Key)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", "error in reading bucket object", err.Error())
		}

		configYAML, err := config.ParseConfigYaml(fileBytes)
		if err != nil {
			return nil, errors.Wrap(err, "yaml.Unmarshal: "+obj.Key)
		}
		if len(configYAML) == 0 {
			continue
		}

		configYAMLs = append(configYAMLs, configYAML)
	}

	return configYAMLs, nil
}

func (repo *SchemaConfig) UpsertResourceConfigs(ctx context.Context, name string, config schema.NamespaceConfigMapType) (resource.ResourceConfig, error) {
	// upsert resource config is not supported for BLOB storage type
	return resource.ResourceConfig{}, resource.ErrUpsertConfigNotSupported
}
