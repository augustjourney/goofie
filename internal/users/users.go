package users

import (
	"gorm.io/gorm"
)

// New creates a new Users Module returning users [Handler] and [Service]
func New(db *gorm.DB) (*Handler, *Service) {
	repo := NewRepo(db)
	service := NewService(repo)
	handler := NewHandler()
	return handler, service
}
