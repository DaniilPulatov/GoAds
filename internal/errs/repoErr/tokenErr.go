package repoerr

import "errors"

var (
	ErrSavingToken = errors.New("Refresh token not saved")
)
