package usecaseerr

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrAccessDenied = Error("access denied")
)
