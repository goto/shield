package activity

import "errors"

var (
	ErrInvalidUUID = errors.New("invalid syntax of uuid")
	ErrInvalidData = errors.New("invalid log data")
)