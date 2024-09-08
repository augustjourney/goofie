package bookmarks

import (
	"gorm.io/gorm"
	"time"
)

type Bookmark struct {
	gorm.Model
	UserID    uint
	ImageID   uint
	CreatedAt time.Time
}
