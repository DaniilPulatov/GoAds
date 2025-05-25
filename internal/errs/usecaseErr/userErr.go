package usecaseerr

import "errors"

var (
	ErrAccessDenied = errors.New("access denied")
)
