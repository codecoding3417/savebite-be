package dto

type AnalysisResponse struct {
	DetectedItems       []string `json:"detectedItems"`
	UsableIngredients   []string `json:"usableIngredients"`
	UnusableIngredients []string `json:"unusableIngredients"`
	Feedback            string   `json:"feedback"`
}
