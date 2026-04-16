package config_test

import (
	"testing"

	"github.com/goto/shield/core/rule/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRulesetYaml_SkipReadBody(t *testing.T) {
	t.Run("should parse skip_read_body field", func(t *testing.T) {
		yamlContent := []byte(`
rules:
  - backends:
      - name: dex
        target: http://dex:80
        prefix: /dex
        frontends:
          - path: /dex/data_upload/upload
            method: POST
            skip_read_body: true
          - path: /dex/api/resource
            method: GET
`)
		ruleset, err := config.ParseRulesetYaml(yamlContent)
		require.NoError(t, err)

		require.Len(t, ruleset.Rules, 1)
		require.Len(t, ruleset.Rules[0].Backends, 1)
		require.Len(t, ruleset.Rules[0].Backends[0].Frontends, 2)

		uploadFrontend := ruleset.Rules[0].Backends[0].Frontends[0]
		assert.True(t, uploadFrontend.SkipReadBody, "skip_read_body should be true for upload route")
		assert.Equal(t, "/dex/data_upload/upload", uploadFrontend.Path)

		apiFrontend := ruleset.Rules[0].Backends[0].Frontends[1]
		assert.False(t, apiFrontend.SkipReadBody, "skip_read_body should default to false")
		assert.Equal(t, "/dex/api/resource", apiFrontend.Path)
	})
}
