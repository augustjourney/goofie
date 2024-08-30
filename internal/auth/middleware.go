package auth

import (
	"api/pkg/handler"
	"context"
	"github.com/gofiber/fiber/v2"
)

func ExtractUser(c *fiber.Ctx) error {
	token := handler.ExtractBearerToken(c)
	if token == "" {
		return c.Next()
	}

	claims, err := validateJWTToken(context.Background(), token)

	if err == nil {
		handler.PutUserIDToFiberContext(c, claims.UserID)
	}

	return c.Next()
}

func UserRequired(c *fiber.Ctx) error {
	var resp handler.Response
	userId, ok := handler.GetUserIDFromFiberContext(c)

	if !ok || userId == 0 {
		return resp.WithStatus(fiber.StatusUnauthorized).Do(c)
	}

	return c.Next()
}
