package users

import "context"

type IRepo interface {
	Create(ctx context.Context, user User) (User, error)
	GetOneByEmail(ctx context.Context, email string) (User, error)
	UpdatePassword(ctx context.Context, userId int, password string) error
}

type IService interface {
	GetOneByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, bool, error)
}
