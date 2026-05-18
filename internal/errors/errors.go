package errors

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.Code == t.Code
}

func New(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, target *Error) error {
	return &Error{
		Code:    target.Code,
		Message: target.Message,
		Err:     err,
	}
}

func WrapMessage(err error, target *Error, message string) error {
	return &Error{
		Code:    target.Code,
		Message: message,
		Err:     err,
	}
}
