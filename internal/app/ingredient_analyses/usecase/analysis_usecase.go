package usecase

import (
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"path/filepath"
	"savebite/internal/app/ingredient_analyses/repository"
	"savebite/internal/domain/dto"
	"savebite/internal/domain/entity"
	"savebite/internal/infra/gemini"
	"savebite/internal/infra/helper"
	"savebite/pkg/log"
	"savebite/pkg/markdown"
	"strings"
)

type AnalysisUsecaseItf interface {
	Analyze(image *multipart.FileHeader, userID uuid.UUID) (gemini.AnalysisResult, error)
	GetHistory(userID uuid.UUID) ([]dto.AnalysisResponse, error)
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

	for _, ingredient := range result.UsableIngredients {
		ingredients = append(ingredients, entity.Ingredient{
			AnalysisID: analysisUUID,
			Name:       ingredient,
			Status:     "usable",
		})
	}

	for _, ingredient := range result.UnusableIngredients {
		ingredients = append(ingredients, entity.Ingredient{
			AnalysisID: analysisUUID,
			Name:       ingredient,
			Status:     "unusable",
		})
	}

	analysis := &entity.Analysis{
		ID:          analysisUUID,
		UserID:      userID,
		Feedback:    result.Feedback,
		Ingredients: ingredients,
	}

	err = u.analysisRepo.Create(analysis)
	if err != nil {
		return gemini.AnalysisResult{}, err
	}

	return result, nil
}

func (u *AnalysisUsecase) GetHistory(userID uuid.UUID) ([]dto.AnalysisResponse, error) {
	analyses, err := u.analysisRepo.GetByUserID(userID)
	if err != nil {
		return []dto.AnalysisResponse{}, err
	}

	var res []dto.AnalysisResponse

	for _, analysis := range *analyses {
		usableIngredients := []string{}
		var unusableIngredients []string
		var items []string

		ingredients := analysis.Ingredients
		for _, ingredient := range ingredients {
			items = append(items, ingredient.Name)

			if ingredient.Status == "usable" {
				usableIngredients = append(usableIngredients, ingredient.Name)
			} else if ingredient.Status == "unusable" {
				unusableIngredients = append(unusableIngredients, ingredient.Name)
			}
		}

		items = helper.RemoveDuplicate(items)

		res = append(res, dto.AnalysisResponse{
			DetectedItems:       items,
			UsableIngredients:   usableIngredients,
			UnusableIngredients: unusableIngredients,
			Feedback:            analysis.Feedback,
		})
	}

	return res, nil
}
