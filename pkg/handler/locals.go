package handler

import "github.com/gofiber/fiber/v2"

// GetUserIDFromFiberContext gets user id from locals of fiber context.
func GetUserIDFromFiberContext(c *fiber.Ctx) (uint, bool) {
	userId, ok := c.Locals("user_id").(uint)
	return userId, ok
}

// PutUserIDToFiberContext puts user id to locals of fiber context.
func PutUserIDToFiberContext(c *fiber.Ctx, userId uint) {
	c.Locals("user_id", userId)
}
