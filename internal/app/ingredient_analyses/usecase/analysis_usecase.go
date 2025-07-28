package usecase

import (
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"path/filepath"
	"savebite/internal/app/ingredient_analyses/repository"
	"savebite/internal/domain/entity"
	"savebite/internal/infra/gemini"
	"savebite/pkg/log"
	"savebite/pkg/markdown"
	"strings"
)

type AnalysisUsecaseItf interface {
	Analyze(image *multipart.FileHeader, userID uuid.UUID) (gemini.AnalysisResult, error)
}

type AnalysisUsecase struct {
	analysisRepo repository.AnalysisRepoItf
	gemini       gemini.GeminiItf
	md           markdown.MarkdownItf
}

func NewAnalysisUsecase(analysisRepo repository.AnalysisRepoItf, gemini gemini.GeminiItf, md markdown.MarkdownItf) AnalysisUsecaseItf {
	return &AnalysisUsecase{analysisRepo, gemini, md}
}

func (u *AnalysisUsecase) Analyze(imageFile *multipart.FileHeader, userID uuid.UUID) (gemini.AnalysisResult, error) {
	file, err := imageFile.Open()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisUsecase][Analyze] Failed to open image file")
		return gemini.AnalysisResult{}, err
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[AnalysisUsecase][Analyze] Failed to read image file")
		return gemini.AnalysisResult{}, err
	}

	mimeType := imageFile.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"

		filename := imageFile.Filename
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case "png":
			mimeType = "image/png"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		}
	}

	result, err := u.gemini.AnalyzeIngredientImage(fileData, mimeType)
	if err != nil {
		return gemini.AnalysisResult{}, err
	}

	result.Feedback = string(u.md.MdToHTML([]byte(result.Feedback)))

	result.Feedback = strings.ReplaceAll(result.Feedback, "\n", "")

	result.Feedback = strings.ReplaceAll(result.Feedback, "\t", "")

	analysisUUID := uuid.New()

	var ingredients []entity.Ingredient

	for _, ingredient := range result.UseableIngredients {
		ingredients = append(ingredients, entity.Ingredient{
			AnalysisID: analysisUUID,
			Name:       ingredient,
			Status:     "useable",
		})
	}

	for _, ingredient := range result.UnuseableIngredients {
		ingredients = append(ingredients, entity.Ingredient{
			AnalysisID: analysisUUID,
			Name:       ingredient,
			Status:     "unuseable",
		})
	}

	analysis := &entity.Analysis{
		ID:          analysisUUID,
		UserID:      userID,
		Feedback:    result.Feedback,
		Ingredients: ingredients,
	}

	err = u.analysisRepo.CreateWithIngredients(analysis)
	if err != nil {
		return gemini.AnalysisResult{}, err
	}

	return result, nil
}
