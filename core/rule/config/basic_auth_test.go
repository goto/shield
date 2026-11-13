package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/core/rule/config"
	"golang.org/x/crypto/bcrypt"
)

func TestBasicAuthValidation(t *testing.T) {
	strong, err := bcrypt.GenerateFromPassword([]byte("synthetic-test-only"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	weak, err := bcrypt.GenerateFromPassword([]byte("synthetic-test-only"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, hash string
		valid      bool
	}{
		{"bcrypt", string(strong), true},
		{"low-cost-bcrypt", string(weak), false},
		{"legacy-sha1", "{SHA}synthetic-not-a-real-credential", false},
		{"legacy-md5", "$apr1$synthetic$not-a-real-credential", false},
		{"plain", "synthetic-not-a-real-credential", false},
		{"truncated", string(strong[:30]), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw := map[string]interface{}{"users": []map[string]interface{}{{"user": "test-user", "password": tc.hash}}}
			err := config.ValidateMiddlewareConfig("basic_auth", raw)
			if (err == nil) != tc.valid {
				t.Fatalf("unexpected validation result: %v", err)
			}
			if err != nil && strings.Contains(err.Error(), tc.hash) {
				t.Fatal("credential leaked in error")
			}
			source := fmt.Sprintf("rules:\n- backends:\n  - frontends:\n    - middlewares:\n      - name: basic_auth\n        config:\n          users:\n          - user: test-user\n            password: %q\n", tc.hash)
			_, err = config.ParseRulesetYaml([]byte(source))
			if (err == nil) != tc.valid {
				t.Fatalf("YAML load result: %v", err)
			}
			err = rule.ValidateRuleset(rule.Ruleset{Rules: []rule.Rule{{Middlewares: rule.MiddlewareSpecs{{Name: "basic_auth", Config: raw}}}}})
			if (err == nil) != tc.valid {
				t.Fatalf("persisted rule validation result: %v", err)
			}
		})
	}
	if err := config.ValidateMiddlewareConfig("authz", nil); err != nil {
		t.Fatal("unrelated middleware rejected")
	}
}

func TestBasicAuthMalformedConfiguration(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("synthetic-test-only"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	user := map[string]interface{}{"user": "test-user", "password": string(hash)}
	for _, raw := range []map[string]interface{}{
		nil,
		{"users": "invalid"},
		{"users": []interface{}{user, user}},
		{"users": []interface{}{user}, "userdb": "unused"},
		{"users": []interface{}{user}, "scope": map[string]interface{}{"action": "{{"}},
		{"users": []interface{}{map[string]interface{}{"user": "test", "password": string(hash), "capabilities": []string{"r#["}}}},
	} {
		if err := config.ValidateMiddlewareConfig("basic_auth", raw); err == nil {
			t.Fatal("invalid config accepted")
		}
	}
}
