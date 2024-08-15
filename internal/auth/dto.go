package auth

import (
	"api/pkg/errs"
	"strings"
)

type SignupDTO struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

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

type SignupResult struct {
	AlreadyExists bool `json:"already_exists"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

type LoginResult struct {
	Token string `json:"token"`
}
