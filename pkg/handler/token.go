package handler

import (
	"api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ExtractBearerToken(c *fiber.Ctx) string {
	authHeader := c.Get(fiber.HeaderAuthorization)
	if authHeader == "" {
		logger.Warn(logger.Record{
			Message: "trying to extract bearer token but auth header is empty",
			Data: map[string]interface{}{
				"url":        c.OriginalURL(),
				"referer":    string(c.Context().Referer()),
				"user_agent": string(c.Context().UserAgent()),
			},
		})
		return ""
	}

	authHeaderData := strings.SplitAfterN(authHeader, "Bearer ", 2)

	if len(authHeaderData) < 2 {
		logger.Warn(logger.Record{
			Message: "unable to split auth header",
			Data: map[string]interface{}{
				"authHeader": authHeader,
			},
		})
		return ""
	}

	return authHeaderData[1]
}
