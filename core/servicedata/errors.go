package servicedata

import "errors"

var (
	ErrInvalidDetail = errors.New("invalid service data detail")
	ErrConflict      = errors.New("key already exist")
	ErrNotExist      = errors.New("service data not exist")
	ErrLogActivity   = errors.New("error while logging activity")
)
