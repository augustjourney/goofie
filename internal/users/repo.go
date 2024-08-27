package users

import (
	"api/pkg/errs"
	"api/pkg/logger"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

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

func (r *Repo) UpdatePassword(ctx context.Context, userId int, password string) error {
	return nil
}

func NewRepo(db *gorm.DB) IRepo {
	return &Repo{
		db,
	}
}
