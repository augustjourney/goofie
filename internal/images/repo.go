package images

import (
	"context"
	"gorm.io/gorm"
)

// Repo stores methods for getting, creating, updating and deleting images data in the database.
type Repo struct {
	db *gorm.DB
}

// Create saves image data to the database.
func (r *Repo) Create(ctx context.Context, image *Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}

// NewRepo creates and returns a new [Repo] instance.
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}
