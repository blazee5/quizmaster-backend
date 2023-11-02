package http_errors

import "errors"

var (
	PermissionDenied = errors.New("permission denied")
)
