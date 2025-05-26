package usecaseerr

var (
	ErrFileNotAllowed = Error("file type not allowed for upload")
	ErrAdNotFound     = Error("ad not found")
	ErrGettingUser    = Error("error getting user from database")
	ErrUserNotFound   = Error("user not found")
)
