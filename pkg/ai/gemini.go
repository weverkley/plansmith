package ai

import (
	"context"
	"fmt"
	"plansmith/pkg/logging"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiExecutor struct {
	client    *genai.Client
	modelName string
}

func NewGeminiExecutor(apiKey string, model string) (*GeminiExecutor, error) {
	// Log the attempt to create a Gemini client
	logging.Info("Creating Gemini client with API key (length: %d) and model: %s", len(apiKey), model)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logging.Error("Failed to create Gemini client: %v", err)
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// If no model is specified, default to gemini-pro
	if model == "" {
		model = "gemini-flash"
	}

	return &GeminiExecutor{
		client:    client,
		modelName: model,
	}, nil
}

func (e *GeminiExecutor) ExecutePrompt(prompt string) (string, error) {
	logging.Info("Executing Gemini prompt (length: %d) with model: %s", len(prompt), e.modelName)
	logging.Debug("Prompt content: %s", prompt)

	ctx := context.Background()
	model := e.client.GenerativeModel(e.modelName)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		logging.Error("Failed to generate content with Gemini: %v", err)
		return "", fmt.Errorf("failed to execute Gemini prompt: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		logging.Warn("No candidates returned from Gemini")
		return "", fmt.Errorf("no response from Gemini")
	}

	result := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		result += fmt.Sprintf("%v", part)
	}

	logging.Info("Successfully received response from Gemini (length: %d)", len(result))
	logging.Debug("Response content: %s", result)

	return cleanResponse(result), nil
}

func cleanResponse(response string) string {
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	}
	return strings.TrimSpace(response)
}

// ListModels returns a list of available Gemini models
func (e *GeminiExecutor) ListModels() ([]string, error) {
	// For Gemini, we'll return a static list of known models
	// In a real implementation, this might call an API endpoint to get the list
	models := []string{
		"gemini-pro",
		"gemini-pro-vision",
		"gemini-ultra",
	}

	return models, nil
}
