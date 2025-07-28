package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"savebite/internal/app/ingredient_analyses/usecase"
	"savebite/internal/middlewares"
)

type AnalysisHandler struct {
	analysisUsecase usecase.AnalysisUsecaseItf
}

func NewAnalysisHandler(router fiber.Router, m middlewares.MiddlewareItf, analysisUsecase usecase.AnalysisUsecaseItf) {
	AnalysisHandler := AnalysisHandler{
		analysisUsecase: analysisUsecase,
	}

	router = router.Group("/analysis", m.RequireAuth)
	router.Post("/analyze", AnalysisHandler.Analyze)
}

func (h *AnalysisHandler) Analyze(c *fiber.Ctx) error {
	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "validation_error",
			"error":      "Missing required image",
			"message":    "The image file is required but was not provided",
		})
	}

	userID := c.Locals("userID").(string)
	userUUID := uuid.MustParse(userID)

	result, err := h.analysisUsecase.Analyze(image, userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      utils.StatusMessage(fiber.StatusInternalServerError),
			"message":    "Something went wrong on our end. Please try again later",
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
