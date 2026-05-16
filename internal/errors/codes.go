package errors

var (
	ErrValidation = New("validation_error", "validation failed")
	ErrNotFound   = New("not_found", "resource not found")
	ErrConflict   = New("conflict", "resource already exists")
	ErrTimeout    = New("timeout", "operation timeout")
	ErrCanceled   = New("canceled", "operation canceled")
	ErrInternal   = New("internal_error", "internal server error")
)
