package http_errors

import "errors"

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrWrongArgument    = errors.New("wrong argument")
	ErrInvalidImage     = errors.New("invalid image")
	ErrCodeExpired      = errors.New("code is expired")
)
