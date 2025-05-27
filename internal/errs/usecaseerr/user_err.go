package usecaseerr

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrInvalidUserData   = Error("invalid user data provided")
	ErrUserAlreadyExists = Error("user already exists")
	ErrCheckUserExists   = Error("error checking if user exists")
	ErrUserNotHaveAds    = Error("user does not have any ads")

	ErrTokenGeneration = Error("error generating token")
	ErrInvalidToken    = Error("invalid token provided")
	ErrTokenExpired    = Error("token has expired")
	ErrInvalidTokenDuration = Error("invalid token duration provided")

	ErrFileNotAllowed = Error("file type not allowed for upload")
	ErrAdNotFound     = Error("ad not found")
	ErrGettingUser    = Error("error getting user from database")
	ErrUserNotFound   = Error("user not found")
	ErrInvalidParams  = Error("invalid parameters")
)
