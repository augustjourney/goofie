package app

import (
	"api/internal/auth"
	"api/internal/images"
	"api/internal/users"
	"api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Use(middleware.RequestLogger)
	app.Use(auth.ExtractUser)

	_, usersService := users.New(db)
	authHandler, _ := auth.New(usersService)
	imagesHandler := images.NewHandler(images.NewService(images.NewRepo(db)))

	app.Post("/auth/signup", authHandler.Signup)
	app.Post("/auth/login", authHandler.Login)
	app.Post("/images", auth.UserRequired, imagesHandler.Create)

	return app
}
