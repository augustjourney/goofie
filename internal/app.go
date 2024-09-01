package app

import (
	"api/internal/auth"
	"api/internal/images"
	"api/internal/users"
	"api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// NewApp creates a new fiber app, setting all the necessary middlewares and routes.
func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Use(middleware.RequestLogger)
	app.Use(auth.ExtractUser)

	_, usersService := users.New(db)
	authHandler, _ := auth.New(usersService)
	imagesHandler, _ := images.New(db)

	app.Post("/auth/signup", authHandler.Signup)
	app.Post("/auth/login", authHandler.Login)
	app.Post("/images", auth.UserRequired, imagesHandler.Create)

	return app
}
