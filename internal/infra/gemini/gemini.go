package gemini

import (
	"context"
	"encoding/json"
	"google.golang.org/genai"
	"savebite/internal/domain/env"
	"savebite/pkg/log"
)

type AnalysisResult struct {
	DetectedItems       []string `json:"detected_items"`
	UsableIngredients   []string `json:"usable_ingredients"`
	UnusableIngredients []string `json:"unusable_ingredients"`
	Feedback            string   `json:"feedback"`
}

type GeminiItf interface {
	AnalyzeIngredientImage(fileData []byte, mimeType string) (AnalysisResult, error)
}

type GeminiStruct struct {
	client *genai.Client
}

var Gemini = getGemini()

func getGemini() GeminiItf {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: env.AppEnv.GeminiAPIKey,
	})
	if err != nil {
		log.Panic(log.LogInfo{
			"error": err.Error(),
		}, "[Gemini][getGemini] Failed to create genai client")
		return &GeminiStruct{}
	}

	return &GeminiStruct{client}
}

func (g *GeminiStruct) AnalyzeIngredientImage(fileData []byte, mimeType string) (AnalysisResult, error) {
	prompt := `
		You are a highly specialized assistant for food safety, culinary reuse, and sustainable waste management.

		You are given an image of a food item or dish, and your job is to analyze it thoroughly and return a structured, actionable assessment.

		TASKS:
		When given an image, perform the following steps in order

		1. Detect and Identify Ingredients
				- Scan the image and list all ingredients or food items as precisely as possible.
				- The number of ingredients may vary, and some might be partially visible.
		2. Analyze Quality of Each Ingredient
				- For each detected item, determine whether it is:
					- Usable (still safe, good condition, reusable)
					- Unusable (rotten, spoiled, contaminated, unsafe)
				- Use cues such as:
					- Color change, bruising, mold, texture, dryness, sogginess, etc.
		3. Provide Practical Feedback
				- Always start with "Hasil Analisis" formatted as title
				- Provide opening paragraph
				- Always begin the response with a concise summary of the detected ingredient conditions, clearly stating which items are still good and which are not.
			  - Case A: If there are usable ingredients
					- Provide a detailed recipe suggestion using only the usable ingredients.
					- Recipe must include:
						- Dish name
						- List of ingredients
						- Step-by-step instructions (simple and based on common household methods)
					- Format this clearly as Markdown.

			  - Case B: If there are only unusable ingredients
					- Provide detailed instructions (step by step) on sustainable disposal or processing.
					- For example: composting, upcycling, fermentation, stock making, etc.
					- Explain why and how each method helps the environment.
					- Format as Markdown.

			  - Case C: If there are both usable and unusable ingredients
					- Do both:
						- Recommend a recipe using only the usable ingredients (as above).
						- Suggest eco-friendly processing for the unusable ingredients (as above).
					- Provide transition paragraph between recipe recommendation and eco-friendly processing suggestion
		4. Translate all responses to Indonesian
		5. Ensure the feedback is properly formatted in valid markdown; if not, regenerate it until it meets the correct format.
		6. FORMAT YOUR RESPONSE AS A VALID JSON OBJECT with these fields:

		{
			\"detected_items\": [],
			\"usable_ingredients\": [],
			\"unusable_ingredients\": [],
			\"feedback\": \"[markdown]\"
		}

		IMPORTANT RULES
		- Only use usable ingredients in recipes.
		- Do not output any other explanations.
		- Your response must consist only of a single Markdown block and a JSON object as specified.
		- The markdown content in the feedback field must not contain any newline characters (\n).
		- The entire response must be written entirely in Bahasa Indonesia
	`

	parts := []*genai.Part{
		genai.NewPartFromText(prompt),
		{
			InlineData: &genai.Blob{
				MIMEType: mimeType,
				Data:     fileData,
			},
		},
	}

	temp := float32(0.1)
	topP := float32(0.8)
	topK := float32(40)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("You are a food quality and sustainability assistant	", genai.RoleUser),
		Temperature:       &temp,
		TopP:              &topP,
		TopK:              &topK,
		ResponseMIMEType:  "application/json",
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	response, _ := g.client.Models.GenerateContent(
		context.Background(),
		env.AppEnv.GeminiModel,
		contents,
		config,
	)

	result := AnalysisResult{}
	err := json.Unmarshal([]byte(response.Text()), &result)
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[Gemini][AnalyzeIngredientImage] Failed to unmarshal json response")
		return result, err
	}

	return result, nil
}
