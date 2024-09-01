package auth

import (
	"api/pkg/handler"
	"context"
	"github.com/gofiber/fiber/v2"
)

// ExtractUser is a middleware that extracts JWT-token from authorization header.
// And if token is not empty, then it extracts user ID from that token.
// Gotten user ID will be put to fiber context locals.
// We can get that userID in http-request handler.
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

// UserRequired is a middleware that checks for user ID in fiber ctx locals
// which ended up there from middleware [ExtractUser].
// If user ID not found, returns unauthorized status.
func UserRequired(c *fiber.Ctx) error {
	var resp handler.Response
	userId, ok := handler.GetUserIDFromFiberContext(c)

	if !ok || userId == 0 {
		return resp.WithStatus(fiber.StatusUnauthorized).Do(c)
	}

	return c.Next()
}
