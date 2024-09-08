package errs

import "net/http"

var Statuses = map[error]int{
	ErrInternal:            http.StatusInternalServerError,
	ErrUserAlreadyExists:   http.StatusConflict,
	ErrEmptyUsername:       http.StatusBadRequest,
	ErrEmptyPassword:       http.StatusBadRequest,
	ErrShortPassword:       http.StatusBadRequest,
	ErrEmptyEmail:          http.StatusBadRequest,
	ErrEmptyFirstName:      http.StatusBadRequest,
	ErrUserNotFound:        http.StatusNotFound,
	ErrWrongCredentials:    http.StatusBadRequest,
	ErrNoImageForUploading: http.StatusBadRequest,
	ErrNoMultipartFormData: http.StatusBadRequest,
}
