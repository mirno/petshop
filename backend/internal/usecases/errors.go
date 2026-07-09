package usecases

import "errors"

var (
	ErrKeyNotFoundError  = errors.New("key not found")
	ErrInvalidValueError = errors.New("invalid value: value cannot be nil")
	ErrAgeDeniedError    = errors.New("invalid value: age cannot be less than 18")
)
