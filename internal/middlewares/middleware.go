package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"savebite/pkg/jwt"
)

type MiddlewareItf interface {
	RequireAuth(c *fiber.Ctx) error
}

type Middleware struct {
	jwt jwt.CustomJWTItf
}

func NewMiddleware(jwt jwt.CustomJWTItf) MiddlewareItf {
	return &Middleware{jwt}
}
