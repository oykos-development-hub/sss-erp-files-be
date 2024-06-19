package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFoundField = errors.New("field not found")
)

type AppError struct {
	Code int
	Err  error
}

func (e AppError) Unwrap() error {
	return e.Err
}

func (e AppError) Error() string {
	return e.Err.Error()
}

func (e AppError) HTTPStatusCode() int {
	return httpStatusCode(e.Code)
}

func (e AppError) PrettyMsg() string {
	return prettyMsg(e.Code)
}

func IsErr(err error, code int) bool {
	var e AppError

	if errors.As(err, &e) {
		return e.Code == code
	}

	return false
}

func New(format string, args ...interface{}) error {
	return AppError{
		Code: InternalCode,
		Err:  fmt.Errorf(format, args...),
	}
}

func NewWithCode(code int, format string, args ...interface{}) error {
	return AppError{
		Code: code,
		Err:  fmt.Errorf(format, args...),
	}
}

func Wrap(err error, message string) error {
	code := InternalCode

	var e AppError

	if errors.As(err, &e) {
		code = e.Code
	}

	return AppError{
		Code: code,
		Err:  fmt.Errorf("%s: %w", message, err),
	}
}

func NewNotFoundError(message string, args ...interface{}) error {
	return AppError{
		Code: NotFoundCode,
		Err:  fmt.Errorf(message, args...),
	}
}

func NewBadRequestError(message string, args ...interface{}) error {
	return AppError{
		Code: BadRequestCode,
		Err:  fmt.Errorf(message, args...),
	}
}

func WrapBadRequestError(err error, message string, args ...interface{}) error {
	return AppError{
		Code: BadRequestCode,
		Err:  fmt.Errorf("%s: %w", fmt.Sprintf(message, args...), err),
	}
}

func WrapNotFoundError(err error, message string, args ...any) error {
	return AppError{
		Code: NotFoundCode,
		Err:  fmt.Errorf("%s: %w", fmt.Sprintf(message, args...), err),
	}
}

func NewInternalServerError(message string, args ...interface{}) error {
	return AppError{
		Code: InternalCode,
		Err:  fmt.Errorf(message, args...),
	}
}

func WrapInternalServerError(err error, message string, args ...interface{}) error {
	return AppError{
		Code: InternalCode,
		Err:  fmt.Errorf("%s: %w", fmt.Sprintf(message, args...), err),
	}
}
