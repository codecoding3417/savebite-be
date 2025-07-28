package dto

type AnalysisResponse struct {
	DetectedItems       []string `json:"detected_items"`
	UsableIngredients   []string `json:"usable_ingredients"`
	UnusableIngredients []string `json:"unusable_ingredients"`
	Feedback            string   `json:"feedback"`
	ImageURL            string   `json:"image_url"`
}

type PaginationRequest struct {
	Page   int `query:"page"`
	Limit  int `query:"limit"`
	Offset int
}

type Meta struct {
	TotalData int `json:"total_items,omitempty"`
	TotalPage int `json:"total_pages,omitempty"`
	Page      int `json:"page,omitempty"`
	Limit     int `json:"limit,omitempty"`
}
