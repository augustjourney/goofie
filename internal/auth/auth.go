package auth

import (
	"api/internal/users"
)

// New creates a new Auth Module returning [Handler] and [Service]
func New(usersService users.IService) (*Handler, *Service) {
	service := NewService(usersService)
	handler := NewHandler(service)
	return handler, service
}
