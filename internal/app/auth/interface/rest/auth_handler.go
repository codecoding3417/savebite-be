package rest

import (
	"savebite/internal/app/auth/usecase"
	"savebite/internal/domain/dto"
	"savebite/internal/domain/env"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Validator   *validator.Validate
	AuthUsecase usecase.AuthUsecaseItf
}

func NewAuthHandler(r fiber.Router, v *validator.Validate, u usecase.AuthUsecaseItf) {
	AuthHandler := AuthHandler{
		Validator:   v,
		AuthUsecase: u,
	}

	r = r.Group("/auth/google/oauth")
	r.Get("/redirect", AuthHandler.Redirect)
	r.Post("/callback", AuthHandler.Callback)
}

func (h *AuthHandler) Redirect(c *fiber.Ctx) error {
	state := "rhegergeigiei3r8723ry2f3ufwjef4g3487ygf"

	c.Cookie(&fiber.Cookie{
		Name:     "oauth2state",
		Value:    state,
		MaxAge:   0,
		Expires:  time.Now().Add(5 * time.Minute),
		Secure:   env.AppEnv.AppEnv != "development",
		HTTPOnly: true,
	})

	redirectURL := h.AuthUsecase.Redirect(state)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"redirect_url": redirectURL,
	})
}

func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	req := dto.OAuthCallbackRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "bad_request",
			"error":      "Bad Request",
			"message":    "Invalid Body Request",
		})
	}

	if err := h.Validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "validation_error",
			"error":      err.Error(),
			"message":    "Validation Error",
		})
	}

	token, err := h.AuthUsecase.Callback(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      "Internal Server Error",
			"message":    "Something went wrong on our end. Please try again later",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}
