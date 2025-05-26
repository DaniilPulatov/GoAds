package usecaseerr

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrAccessDenied   = Error("access denied")
	ErrFileNotAllowed = Error("file type not allowed for upload")
	ErrAdNotFound     = Error("ad not found")

	ErrInvalidUserData   = Error("invalid user data provided")
	ErrUserAlreadyExists = Error("user already exists")
	ErrCheckUserExists   = Error("error checking if user exists")
	ErrUserNotFound      = Error("user not found")

	ErrTokenGeneration = Error("error generating token")
	ErrInvalidToken    = Error("invalid token provided")
	ErrTokenExpired    = Error("token has expired")
)
