package errs

import (
	"errors"
)

var (
	// general

	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("not found")

	// users

	ErrUserAlreadyExists = errors.New("you already have an account")
	ErrUserNotFound      = errors.New("user not found")

	// auth

	ErrEmptyUsername    = errors.New("username cannot be empty")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrShortPassword    = errors.New("password should be at least 8 characters")
	ErrEmptyEmail       = errors.New("email cannot be empty")
	ErrEmptyFirstName   = errors.New("first name cannot be empty")
	ErrWrongCredentials = errors.New("wrong login or password")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidToken     = errors.New("invalid token")

	// images

	ErrNoImageForUploading      = errors.New("no image provided for upload")
	ErrWrongQualityValue        = errors.New("wrong quality value")
	ErrWrongHeightValue         = errors.New("wrong height value")
	ErrWrongWidthValue          = errors.New("wrong width value")
	ErrOutputFormatNotSupported = errors.New("output format not supported")
	ErrUnableToUploadImage      = errors.New("unable to upload image")
	// handler

	ErrNoMultipartFormData = errors.New("file in multipart form data not provided")
)
