package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"plansmith/pkg/logging"
)

type QwenExecutor struct {
	apiKey string
	model  string
}

type QwenRequest struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type Input struct {
	Prompt string `json:"prompt"`
}

type Parameters struct {
	MaxTokens int `json:"max_tokens"`
}

type QwenResponse struct {
	Output    Output `json:"output"`
	RequestID string `json:"request_id"`
}

type Output struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

func NewQwenExecutor(apiKey string, model string) (*QwenExecutor, error) {
	logging.Info("Creating Qwen executor with API key (length: %d) and model: %s", len(apiKey), model)

	// If no model is specified, default to qwen-turbo
	if model == "" {
		model = "qwen-turbo"
	}

	return &QwenExecutor{
		apiKey: apiKey,
		model:  model,
	}, nil
}

func (e *QwenExecutor) ExecutePrompt(prompt string) (string, error) {
	logging.Info("Executing Qwen prompt (length: %d) with model: %s", len(prompt), e.model)
	logging.Debug("Prompt content: %s", prompt)

	// Create request
	request := QwenRequest{
		Model: e.model,
		Input: Input{
			Prompt: prompt,
		},
		Parameters: Parameters{
			MaxTokens: 2048,
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		logging.Error("Failed to marshal Qwen request: %v", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	logging.Debug("Qwen request JSON (length: %d): %s", len(jsonData), string(jsonData))

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation", bytes.NewBuffer(jsonData))
	if err != nil {
		logging.Error("Failed to create HTTP request for Qwen: %v", err)
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-SSE", "enable")

	logging.Info("Sending HTTP request to Qwen API")
	logging.Debug("Request headers: Authorization=%s, Content-Type=%s, X-DashScope-SSE=%s",
		"Bearer "+e.apiKey, req.Header.Get("Content-Type"), req.Header.Get("X-DashScope-SSE"))

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("Failed to execute HTTP request to Qwen: %v", err)
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Log response status
	logging.Info("Received HTTP response from Qwen API: status=%s (%d)", resp.Status, resp.StatusCode)
	logging.Debug("Response headers: %v", resp.Header)

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Read the response body for error details
		body, _ := io.ReadAll(resp.Body)
		logging.Error("Qwen API returned error status %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("HTTP request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var qwenResp QwenResponse
	err = json.NewDecoder(resp.Body).Decode(&qwenResp)
	if err != nil {
		logging.Error("Failed to decode Qwen response: %v", err)
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	logging.Info("Successfully received response from Qwen (length: %d)", len(qwenResp.Output.Text))
	logging.Debug("Response content: %s", qwenResp.Output.Text)

	return qwenResp.Output.Text, nil
}

// ListModels returns a list of available Qwen models
func (e *QwenExecutor) ListModels() ([]string, error) {
	// For Qwen, we'll return a static list of known models
	// In a real implementation, this would call an API endpoint to get the list
	models := []string{
		"qwen-turbo",
		"qwen-plus",
		"qwen-max",
	}

	return models, nil
}
