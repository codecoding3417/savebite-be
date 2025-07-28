package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"savebite/internal/app/user/usecase"
	"savebite/internal/middlewares"
)

type UserHandler struct {
	useCase usecase.UserUsecaseItf
}

func NewUserHandler(r fiber.Router, u usecase.UserUsecaseItf, m middlewares.MiddlewareItf) {
	UserHandler := UserHandler{
		useCase: u,
	}

	r = r.Group("/", m.RequireAuth)
	r.Get("/me", UserHandler.GetUserProfile)
	r = r.Group("/users")
}

func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	userUUID := uuid.MustParse(userId)
	res, err := h.useCase.GetProfile(userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      "Internal Server Error",
			"message":    "Something went wrong on our end. Please try again later",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": res,
	})
}
