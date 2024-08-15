package users

import (
	"api/pkg/errs"
	"context"
	"errors"
)

type Service struct {
	storage IRepo
}

// GetOneByEmail finds a user in DB by email and returns a found user and error
// Errors: ErrInternal, ErrUserNotFound
func (s *Service) GetOneByEmail(ctx context.Context, email string) (User, error) {
	return s.storage.GetOneByEmail(ctx, email)
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
	createdUser, err := s.storage.Create(ctx, user)
	return createdUser, alreadyExists, err
}

func NewService(storage IRepo) *Service {
	return &Service{storage: storage}
}
