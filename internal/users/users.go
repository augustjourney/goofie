package users

import (
	"gorm.io/gorm"
)

func New(db *gorm.DB) (*Handler, *Service) {
	repo := NewRepo(db)
	service := NewService(repo)
	handler := NewHandler()
	return handler, service
}
