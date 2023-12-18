package http_errors

import "errors"

var (
	PermissionDenied = errors.New("permission denied")
	WrongArgument    = errors.New("wrong argument")
	ErrInvalidImage  = errors.New("invalid image")
)
