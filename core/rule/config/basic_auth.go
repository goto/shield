package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

var bcryptEncoding = regexp.MustCompile(`^\$2[aby]\$[0-9]{2}\$[./A-Za-z0-9]{53}$`)

// ValidateMiddlewareConfig validates security-sensitive configuration without
// including credentials or the complete configuration in returned errors.
func ValidateMiddlewareConfig(name string, raw map[string]interface{}) error {
	if name != "basic_auth" {
		return nil
	}
	var cfg struct {
		Users []struct {
			User         string
			Password     string
			Capabilities []string
		}
		UserDB string
		Scope  struct {
			Action     string
			Attributes map[string]interface{}
		}
	}
	if err := mapstructure.Decode(raw, &cfg); err != nil {
		return errors.New("invalid basic_auth configuration")
	}
	if cfg.UserDB != "" {
		return errors.New("basic_auth userdb is not supported; configure bcrypt users")
	}
	if len(cfg.Users) == 0 {
		return errors.New("basic_auth requires at least one bcrypt user")
	}
	seen := map[string]bool{}
	for i, u := range cfg.Users {
		if strings.TrimSpace(u.User) == "" || seen[u.User] {
			return fmt.Errorf("basic_auth user %d has an empty or duplicate name", i)
		}
		seen[u.User] = true
		cost, err := bcrypt.Cost([]byte(u.Password))
		if err != nil || !bcryptEncoding.MatchString(u.Password) || cost < bcrypt.DefaultCost {
			return fmt.Errorf("basic_auth user %d requires bcrypt with cost at least %d", i, bcrypt.DefaultCost)
		}
		for _, capability := range u.Capabilities {
			if strings.HasPrefix(capability, "r#") {
				if _, err := regexp.Compile(strings.TrimPrefix(capability, "r#")); err != nil {
					return fmt.Errorf("basic_auth user %d has an invalid capability expression", i)
				}
			}
		}
	}
	if cfg.Scope.Action != "" {
		if _, err := template.New("basic_auth").Parse(cfg.Scope.Action); err != nil {
			return errors.New("invalid basic_auth scope action template")
		}
	}
	return nil
}
