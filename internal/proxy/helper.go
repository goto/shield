package proxy

import (
	"strings"

	"github.com/valyala/fasttemplate"
)

func ComposeAttribute(attribute string, attrs map[string]interface{}) string {
	if strings.Contains(attribute, "{") {
		template := fasttemplate.New(attribute, "{", "}")
		return template.ExecuteString(attrs)
	}
	return attribute
}
