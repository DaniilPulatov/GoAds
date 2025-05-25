package usecaseerr

import "errors"

var (
	ErrAccessDenied = errors.New("access denied")
	ErrFileNotAllowed = errors.New("file type not allowed for upload")
	ErrAdNotFound = errors.New("ad not found")
)
