package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/universeh2h/report/pkg/config"
)

var authHelpers = NewAuthHelpers()

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		if token := c.Cookies("vazzaccess"); token != "" {
			tokenString = token                                // ✅ First assign the token
			fmt.Printf("token from cookie: %s\n", tokenString) // ✅ Then print it
		} else {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Missing authorization token",
				})
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Invalid authorization format",
				})
			}

			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			fmt.Printf("token from header: %s\n", tokenString)
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Missing token",
			})
		}

		claims, err := config.ValidateToken(tokenString)
		if err != nil {
			fmt.Printf("Token validation failed: %s, Error: %v\n", tokenString, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token",
			})
		}

		// Set user information in context
		c.Locals("username", claims["username"])

		return c.Next()
	}
}
