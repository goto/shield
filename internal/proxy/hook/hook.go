package hook

import (
	"net/http"

	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/proxy/middleware"
)

type Service interface {
	Info() Info
	ServeHook(res *http.Response, err error) (*http.Response, error)
}

type Info struct {
	Name        string
	Description string
}

func ExtractHook(r *http.Request, name string) (rule.HookSpec, bool) {
	rl, ok := ExtractRule(r)
	if !ok {
		return rule.HookSpec{}, false
	}
	return rl.Hooks.Get(name)
}

func ExtractRule(r *http.Request) (*rule.Rule, bool) {
	rl, ok := middleware.ExtractRule(r)
	if !ok {
		return nil, false
	}

	return rl, true
}

type Hook struct{}

func New() Hook {
	return Hook{}
}

func (h Hook) Info() Info {
	return Info{}
}

func (h Hook) ServeHook(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
		// TODO: clear or add error body as well
	}

	return res, nil
}
