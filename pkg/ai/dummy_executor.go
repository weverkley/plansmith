package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DummyExecutor implements the AIExecutor interface for testing purposes.
// It returns mock data without making any actual AI calls.
type DummyExecutor struct{}

func (d *DummyExecutor) ExecutePrompt(prompt string) (string, error) {
	if strings.Contains(prompt, "Project Markdown") {
		return d.mockVisionResponse(), nil
	} else if strings.Contains(prompt, "Your overall project vision is") {
		return d.mockStoriesResponse(), nil
	} else if strings.Contains(prompt, "The User Story you are breaking down is") {
		return d.mockTasksResponse(), nil
	}
	return "", fmt.Errorf("unknown prompt type for dummy executor")
}

func (d *DummyExecutor) mockVisionResponse() string {
	response := map[string]interface{}{
		"project_name":   "Dummy Project",
		"product_vision": "A dummy vision for a dummy project.",
		"epics":          []map[string]string{{"name": "Dummy Epic 1"}, {"name": "Dummy Epic 2"}},
	}
	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse)
}

func (d *DummyExecutor) mockStoriesResponse() string {
	response := map[string]interface{}{
		"user_stories": []map[string]interface{}{
			{
				"title":    "Dummy Story 1",
				"story":    "As a dummy user, I want a dummy feature, so that I can have a dummy benefit.",
				"priority": 10,
			},
			{
				"title":    "Dummy Story 2",
				"story":    "As another dummy user, I want another dummy feature, so that I can have another dummy benefit.",
				"priority": 8,
			},
		},
	}
	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse)
}

func (d *DummyExecutor) mockTasksResponse() string {
	response := map[string]interface{}{
		"tasks": []map[string]interface{}{
			{
				"title":        "Dummy Task 1",
				"description":  "A dummy task for a dummy story.",
				"dependencies": []string{},
				"labels":       []string{"Type: Bug", "Severity: High"},
			},
			{
				"title":        "Dummy Task 2",
				"description":  "Another dummy task for a dummy story.",
				"dependencies": []string{"Dummy Task 1"},
				"labels":       []string{"Type: Feature", "Priority: Medium"},
			},
		},
	}
	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse)
}

func (d *DummyExecutor) ListModels() ([]string, error) {
	return []string{"dummy-model-1", "dummy-model-2"}, nil
}
