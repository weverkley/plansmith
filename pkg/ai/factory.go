package ai

import (
	"fmt"
)

func NewExecutor(provider, apiKey, model string) (AIExecutor, error) {
	switch provider {
	case "gemini":
		return NewGeminiExecutor(apiKey, model)
	case "openai":
		return NewOpenAIExecutor(apiKey, model)
	case "qwen":
		return NewQwenExecutor(apiKey, model)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", provider)
	}
}
