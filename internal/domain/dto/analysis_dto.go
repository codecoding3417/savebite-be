package dto

type AnalysisResponse struct {
	DetectedItems       []string `json:"detected_items"`
	UsableIngredients   []string `json:"usable_ingredients"`
	UnusableIngredients []string `json:"unusable_ingredients"`
	Feedback            string   `json:"feedback"`
}
