package auth

import (
	"api/pkg/errs"
	"strings"
)

// SignupDTO stores data for signing up that came from request body
type SignupDTO struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// Validate checks required fields in [SignupDTO].
// Also makes email, username to lowercase.
// And trims space for email, username and password.
func (d *SignupDTO) Validate() error {
	if d.Username == "" {
		return errs.ErrEmptyUsername
	}
	if d.Email == "" {
		return errs.ErrEmptyEmail
	}
	if d.FirstName == "" {
		return errs.ErrEmptyFirstName
	}
	if d.Password == "" {
		return errs.ErrEmptyPassword
	}
	if len(d.Password) < 8 {
		return errs.ErrShortPassword
	}

	d.Email = strings.ToLower(strings.TrimSpace(d.Email))
	d.Username = strings.ToLower(strings.TrimSpace(d.Username))
	d.Password = strings.TrimSpace(d.Password)

	return nil
}

// SignupResult is a result data after signing up.
type SignupResult struct {
	// When user tries to sign up, they may already exist.
	// Returning this bool field we can show on client side a useful message.
	AlreadyExists bool `json:"already_exists"`
}

// LoginDTO stores data for login that came from request body
type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate checks required fields in [LoginDTO].
// Also makes email to lowercase.
// And trims space of email and password.
func (d *LoginDTO) Validate() error {
	if d.Email == "" {
		return errs.ErrEmptyEmail
	}
	if d.Password == "" {
		return errs.ErrEmptyPassword
	}
	d.Email = strings.ToLower(strings.TrimSpace(d.Email))
	d.Password = strings.TrimSpace(d.Password)
	return nil
}

// LoginResult is a result data after login.
type LoginResult struct {
	Token string `json:"token"`
}
