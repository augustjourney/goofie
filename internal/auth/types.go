package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

type IHandler interface {
	Signup(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type IService interface {
	Signup(ctx context.Context, payload SignupDTO) (SignupResult, error)
	Login(ctx context.Context, payload LoginDTO) (LoginResult, error)
}
