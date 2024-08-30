package handler

import "github.com/gofiber/fiber/v2"

func GetUserIDFromFiberContext(c *fiber.Ctx) (uint, bool) {
	userId, ok := c.Locals("user_id").(uint)
	return userId, ok
}

func PutUserIDToFiberContext(c *fiber.Ctx, userId uint) {
	c.Locals("user_id", userId)
}
