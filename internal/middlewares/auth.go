package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"savebite/pkg/jwt"
	"strings"
)

func (m *Middleware) RequireAuth(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if len(token) < 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error_code": "unauthorized",
			"error":      "missing authentication token",
			"message":    "Authentication token is required",
		})
	}

	parts := strings.Split(token, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error_code": "unauthorized",
			"error":      "invalid authentication format",
			"message":    "Authorization token format must be Bearer <token>",
		})
	}

	claims := &jwt.Claims{}

	err := m.jwt.Decode(parts[1], claims)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error_code": "unauthorized",
			"error":      "invalid or expired authentication token",
			"message":    "Invalid or expired token",
		})
	}

	c.Locals("id", claims.ID.String())
	c.Locals("email", claims.Email)
	c.Locals("name", claims.Name)

	return c.Next()
}
