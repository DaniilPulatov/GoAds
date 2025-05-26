package repoerr

var (
	ErrUserExists         = Error("user already exists")
	ErrUserNotFound       = Error("user not found")
	ErrInvalidCredentials = Error("invalid credentials")
	ErrUserSelectFailed   = Error("failed to select user")
	ErrUserInsertFailed   = Error("failed to insert user")

	ErrCreatingToken      = Error("failed to create refresh token")
	ErrTokenNotFound      = Error("refresh token not found")
	ErrTokenUpdateFailed  = Error("failed to update refresh token")
	ErrTokenDeleteFailed  = Error("failed to delete refresh token")
	ErrTokenSelectFailed  = Error("failed to select refresh token")
	ErrTokenAlreadyExists = Error("refresh token already exists")
)
