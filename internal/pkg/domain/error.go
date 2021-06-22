package domain

type Error struct {
	Message string `json:"message"`
}

func NewError(message string) error {
	return &Error{
		Message: message,
	}
}

func (err Error) Error() string {
	return err.Message
}

var (
	NotFoundError = NewError("not found")
	ConflictError = NewError("conflict error")
	ParentError   = NewError("no parent in this thread")
)