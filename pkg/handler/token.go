package handler

import (
	"api/pkg/logger"
	"context"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// ExtractBearerToken extracts token from bearer authorization.
// Returns empty string if authorization header is empty or invalid.
func ExtractBearerToken(c *fiber.Ctx) string {
	authHeader := c.Get(fiber.HeaderAuthorization)
	if authHeader == "" {
		logger.Warn(context.TODO(), "trying to extract bearer token but auth header is empty",
			"url", c.OriginalURL(), "referer", string(c.Context().Referer()), "user_agent", string(c.Context().UserAgent()))
		return ""
	}

	authHeaderData := strings.SplitAfterN(authHeader, "Bearer ", 2)

	if len(authHeaderData) < 2 {
		logger.Warn(context.TODO(), "unable to split auth header", "authHeader", authHeader)
		return ""
	}

	return authHeaderData[1]
}
