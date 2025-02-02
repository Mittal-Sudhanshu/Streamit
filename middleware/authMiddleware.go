package middleware

import (
	"streamit/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorization header",
		})
	}

	// Example: Extract userID from token (mocked, replace with actual decoding logic)
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	// Assume the token is just the user ID for this example
	userID, _ := utils.DecodeJWT(tokenParts[1])

	// Store userID in Fiber's context locals
	c.Locals("userID", userID["userId"])

	// Proceed to the next middleware/handler
	return c.Next()
}
