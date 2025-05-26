package utilserr

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrTitleRequired    = Error("title is required")
	ErrLocationRequired = Error("location is required")
	ErrCategoryRequired = Error("category is required")
)
