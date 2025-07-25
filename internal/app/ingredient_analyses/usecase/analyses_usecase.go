package usecase

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"savebite/internal/infra/gemini"
	"savebite/pkg/log"
	"savebite/pkg/markdown"
	"strings"
)

type AnalysesUsecaseItf interface {
	Analyze(image *multipart.FileHeader) (gemini.AnalysisResult, error)
}

type AnalysesUsecase struct {
	gemini gemini.GeminiItf
	md     markdown.MarkdownItf
}

func NewAnalysesUsecase(gemini gemini.GeminiItf, md markdown.MarkdownItf) AnalysesUsecaseItf {
	return &AnalysesUsecase{gemini, md}
}

func (u *AnalysesUsecase) Analyze(imageFile *multipart.FileHeader) (gemini.AnalysisResult, error) {
	file, err := imageFile.Open()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[Gemini][AnalyzeIngredientImage] Failed to open image file")
		return gemini.AnalysisResult{}, err
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[Gemini][AnalyzeIngredientImage] Failed to read image file")
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

	return result, nil
}
