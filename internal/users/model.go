package users

import (
	"api/internal/bookmarks"
	"api/internal/images"
	"gorm.io/gorm"
	"time"
)

// User stores data about user and represents table «users» in db
type User struct {
	gorm.Model
	Email            string  `gorm:"type:varchar(255);unique;not null;index;" json:"email"`
	FirstName        string  `gorm:"type:varchar(255);not null;" json:"first_name"`
	LastName         *string `gorm:"type:varchar(255);" json:"last_name"`
	Username         string  `gorm:"type:varchar(255);not null;" json:"username"`
	Password         string  `gorm:"type:varchar(255);not null;" json:"-"`
	ConfirmedEmailAt *time.Time
	Images           []images.Image       `gorm:"foreignKey:AuthorID"`
	Bookmarks        []bookmarks.Bookmark `gorm:"many2many:bookmarks;"`
}
