package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"savebite/internal/app/analysis/usecase"
	"savebite/internal/domain/dto"
	"savebite/internal/middlewares"
)

type AnalysisHandler struct {
	analysisUsecase usecase.AnalysisUsecaseItf
}

func NewAnalysisHandler(router fiber.Router, m middlewares.MiddlewareItf, analysisUsecase usecase.AnalysisUsecaseItf) {
	AnalysisHandler := AnalysisHandler{
		analysisUsecase: analysisUsecase,
	}

	router = router.Group("/", m.RequireAuth)
	router.Get("/me/analyses", AnalysisHandler.GetHistory)

	router = router.Group("/analyses")
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": result,
	})
}

func (h *AnalysisHandler) GetHistory(c *fiber.Ctx) error {
	pagination := dto.PaginationRequest{}
	if err := c.QueryParser(&pagination); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "bad_request",
			"error":      utils.StatusMessage(fiber.StatusBadRequest),
			"message":    "Failed to parse query",
		})
	}

	userID := c.Locals("userID").(string)
	userUUID := uuid.MustParse(userID)

	result, meta, err := h.analysisUsecase.GetHistory(userUUID, pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      utils.StatusMessage(fiber.StatusInternalServerError),
			"message":    "Something went wrong on our end. Please try again later",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": result,
		"meta":    meta,
	})
}
