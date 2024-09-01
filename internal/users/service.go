package users

import (
	"api/pkg/errs"
	"context"
	"errors"
)

// Service stores methods of users actions.
type Service struct {
	repo IRepo
}

// GetOneByEmail finds a user in DB by email and returns a found user and error
// Errors: ErrInternal, ErrUserNotFound
func (s *Service) GetOneByEmail(ctx context.Context, email string) (User, error) {
	return s.repo.GetOneByEmail(ctx, email)
}

// Create creates a new user in DB and returns created user, already exist bool and error
// Errors: ErrInternal
func (s *Service) Create(ctx context.Context, user User) (User, bool, error) {
	// check if user already exists
	foundUser, err := s.GetOneByEmail(ctx, user.Email)

	var alreadyExists bool

	if err != nil {
		if !errors.Is(err, errs.ErrUserNotFound) {
			return foundUser, alreadyExists, err
		}
	}

	if foundUser.ID != 0 {
		alreadyExists = true
		return foundUser, alreadyExists, nil
	}

	// create if user doesn't exist
	createdUser, err := s.repo.Create(ctx, user)
	return createdUser, alreadyExists, err
}

// NewService creates and returns a new users [Service] instance.
func NewService(repo IRepo) *Service {
	return &Service{repo: repo}
}
