package apperr

type Kind string

const (
	Invalid      Kind = "invalid"
	NotFound     Kind = "not-found"
	Unauthorized Kind = "unauthorized"
	Forbidden    Kind = "forbidden"
	Unavailable  Kind = "unavailable"
	Internal     Kind = "internal"
)

type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return e.Message + ": " + e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(kind Kind, message string) *Error {
	return &Error{Kind: kind, Message: message}
}

func Wrap(kind Kind, message string, err error) *Error {
	return &Error{Kind: kind, Message: message, Err: err}
}
