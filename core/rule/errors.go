package rule

import "errors"

var (
	ErrUnknown                  = errors.New("undefined proxy rule")
	ErrUpsertConfigNotSupported = errors.New("upsert rule config is currently not supported")
	ErrInvalidRuleConfig        = errors.New("invalid rule config")
	ErrMarshal                  = errors.New("error while marshalling rule config")
)
