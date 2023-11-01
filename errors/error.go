package errors

import (
	"errors"
	"net/http"
)

var (
	ErrUserExists     = errors.New("user already exists")
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrExpired        = errors.New("expired")
	ErrForbidden      = errors.New("forbidden")
	ErrDatabaseError  = errors.New("database error")
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("invalid input")
	// define more custom errors as needed
)

func MapErrorToStatusCode(err error) int {
	switch err {
	case ErrUserExists:
		return http.StatusConflict
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidInput:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrDatabaseError:
		return http.StatusInternalServerError
	case ErrExpired:
		return http.StatusForbidden
	case ErrInternalServer:
		return http.StatusInternalServerError
	case ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
