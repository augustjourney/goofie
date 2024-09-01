package images

import "gorm.io/gorm"

// New creates a new Images Module returning images [Handler] and [Service]
func New(db *gorm.DB) (*Handler, *Service) {
	repo := NewRepo(db)
	service := NewService(repo)
	handler := NewHandler(service)
	return handler, service
}
