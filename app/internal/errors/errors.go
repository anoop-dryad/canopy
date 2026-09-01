package errors

import "errors"

var (
	ErrDeviceNotFound = errors.New("device not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDuplicate      = errors.New("device already exists")
)
