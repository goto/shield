package schema_generator

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/goto/shield/internal/schema"

	"github.com/stretchr/testify/assert"
)

func makeDefnMap(s []string) map[string][]string {
	finalMap := make(map[string][]string)

	for _, v := range s {
		splitedConfigText := strings.Split(v, "\n")
		k := splitedConfigText[0]
		sort.Strings(splitedConfigText)
		finalMap[k] = splitedConfigText
	}

	return finalMap
}

// Test to check difference between predefined_schema.txt and schema defined in predefined.go
func TestPredefinedSchema(t *testing.T) {
	content, err := os.ReadFile("predefined_schema")
	assert.NoError(t, err)

	// slice and sort as GenerateSchema() generated the permissions and relations in random order
	schema.PreDefinedSystemNamespaceConfig[schema.ServiceDataKeyNamespace] = schema.ServiceDataKeyConfig
	actualPredefinedConfigs := makeDefnMap(GenerateSchema(schema.PreDefinedSystemNamespaceConfig))
	expectedPredefinedConfigs := makeDefnMap(strings.Split(string(content), "\n--\n"))
	assert.Equal(t, expectedPredefinedConfigs, actualPredefinedConfigs)
}
