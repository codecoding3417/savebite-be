package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"savebite/internal/app/ingredient_analyses/usecase"
	"savebite/internal/middlewares"
)

type AnalysesHandler struct {
	analysesUsecase usecase.AnalysesUsecaseItf
}

func NewAnalysesHandler(router fiber.Router, m middlewares.MiddlewareItf, analysesUsecase usecase.AnalysesUsecaseItf) {
	AnalysesHandler := AnalysesHandler{
		analysesUsecase: analysesUsecase,
	}

	router = router.Group("/analyses")
	router.Post("/analyze", AnalysesHandler.Analyze)
}

func (h *AnalysesHandler) Analyze(c *fiber.Ctx) error {
	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "validation_error",
			"error":      "Missing required image",
			"message":    "The image file is required but was not provided",
		})
	}

	result, err := h.analysesUsecase.Analyze(image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      utils.StatusMessage(fiber.StatusInternalServerError),
			"message":    "Something went wrong on our end. Please try again later",
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
