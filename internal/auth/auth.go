package auth

import (
	"api/internal/users"
)

func New(usersService users.IService) (*Handler, *Service) {
	service := NewService(usersService)
	handler := NewHandler(service)
	return handler, service
}
