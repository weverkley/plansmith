package ai

import (
	"context"
	"fmt"
	"plansmith/pkg/logging"

	"github.com/sashabaranov/go-openai"
)

type OpenAIExecutor struct {
	client *openai.Client
	model  string
}

func NewOpenAIExecutor(apiKey string, model string) (*OpenAIExecutor, error) {
	logging.Info("Creating OpenAI client with API key (length: %d) and model: %s", len(apiKey), model)

	client := openai.NewClient(apiKey)

	// If no model is specified, default to gpt-3.5-turbo
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	return &OpenAIExecutor{
		client: client,
		model:  model,
	}, nil
}

func (e *OpenAIExecutor) ExecutePrompt(prompt string) (string, error) {
	logging.Info("Executing OpenAI prompt (length: %d) with model: %s", len(prompt), e.model)
	logging.Debug("Prompt content: %s", prompt)

	ctx := context.Background()

	resp, err := e.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: e.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		logging.Error("Failed to create chat completion with OpenAI: %v", err)
		return "", fmt.Errorf("failed to execute OpenAI prompt: %w", err)
	}

	if len(resp.Choices) == 0 {
		logging.Warn("No choices returned from OpenAI")
		return "", fmt.Errorf("no response from OpenAI")
	}

	result := resp.Choices[0].Message.Content
	logging.Info("Successfully received response from OpenAI (length: %d)", len(result))
	logging.Debug("Response content: %s", result)

	return result, nil
}

// ListModels returns a list of available OpenAI models
func (e *OpenAIExecutor) ListModels() ([]string, error) {
	ctx := context.Background()

	// Get the list of models from OpenAI
	models, err := e.client.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list OpenAI models: %w", err)
	}

	// Extract model IDs
	var modelIDs []string
	for _, model := range models.Models {
		modelIDs = append(modelIDs, model.ID)
	}

	return modelIDs, nil
}
