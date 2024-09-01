package users

import (
	"api/pkg/errs"
	"api/pkg/logger"
	"context"
	"errors"
	"gorm.io/gorm"
)

// Repo stores methods for getting, creating, updating and deleting users data in the database.
type Repo struct {
	db *gorm.DB
}

// GetOneByEmail finds a user by email in the database
func (r *Repo) GetOneByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := r.db.WithContext(ctx).Where(&User{Email: email}).First(&user).Error
	if err == nil {
		return user, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errs.ErrUserNotFound
	}

	logger.Error(logger.Record{
		Context: ctx,
		Message: "[Repo.GetOneByEmail]: unable to get a user by email",
		Error:   err,
		Data:    map[string]interface{}{"email": email},
	})

	return user, errs.ErrInternal
}

// Create saves user data to the database.
func (r *Repo) Create(ctx context.Context, user User) (User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		logger.Error(logger.Record{
			Context: ctx,
			Message: "[Repo.Create]: unable to create a new user",
			Error:   err,
			Data:    map[string]interface{}{"user": user},
		})
		return user, errs.ErrInternal
	}
	return user, nil
}

// UpdatePassword TODO
func (r *Repo) UpdatePassword(ctx context.Context, userId int, password string) error {
	return nil
}

// NewRepo creates and returns a new [Repo] instance.
func NewRepo(db *gorm.DB) IRepo {
	return &Repo{
		db,
	}
}
