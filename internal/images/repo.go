package images

import (
	"context"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func (r *Repo) Create(ctx context.Context, image *Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}
