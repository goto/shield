package attribute

import (
	"strings"

	"github.com/valyala/fasttemplate"
)

const (
	TypeJSONPayload AttributeType = "json_payload"
	TypeGRPCPayload AttributeType = "grpc_payload"
	TypeQuery       AttributeType = "query"
	TypeHeader      AttributeType = "header"
	TypePathParam   AttributeType = "path_param"
	TypeConstant    AttributeType = "constant"
	TypeComposite   AttributeType = "composite"

	SourceRequest  AttributeType = "request"
	SourceResponse AttributeType = "response"
)

type AttributeType string

type Attribute struct {
	Key    string        `yaml:"key" mapstructure:"key"`
	Type   AttributeType `yaml:"type" mapstructure:"type"`
	Index  string        `yaml:"index" mapstructure:"index"` // proto index
	Path   string        `yaml:"path" mapstructure:"path"`
	Params []string      `yaml:"params" mapstructure:"params"`
	Source string        `yaml:"source" mapstructure:"source"`
	Value  string        `yaml:"value" mapstructure:"value"`
}

func Compose(attribute string, attrs map[string]interface{}) string {
	if strings.Contains(attribute, "${") {
		template := fasttemplate.New(attribute, "${", "}")
		return template.ExecuteString(attrs)
	}
	return attribute
}
