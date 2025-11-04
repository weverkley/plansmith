package smith

import (
	"encoding/json"
	"fmt"
	"plansmith/pkg/ai"
	"plansmith/pkg/logging"
)

type Agent struct {
	executor ai.AIExecutor
	loader   *Loader
}

func NewAgent(executor ai.AIExecutor) *Agent {
	logging.Info("Creating new Smith agent")
	return &Agent{
		executor: executor,
		loader:   NewLoader(),
	}
}

// Executor returns the AI executor
func (a *Agent) Executor() ai.AIExecutor {
	return a.executor
}

type VisionResponse struct {
	ProjectName   string   `json:"project_name"`
	ProductVision string   `json:"product_vision"`
	Epics         []Epic   `json:"epics"`
}

type Epic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (a *Agent) GenerateVision(markdownInput string) (*VisionResponse, error) {
	logging.Info("Generating vision from markdown input (length: %d)", len(markdownInput))

	prompt, err := a.loader.LoadInitPrompt(1, map[string]interface{}{
		"MarkdownInput": markdownInput,
	})
	if err != nil {
		logging.Error("Failed to load vision prompt: %v", err)
		return nil, fmt.Errorf("failed to load vision prompt: %w", err)
	}

	logging.Info("Loaded vision prompt (length: %d)", len(prompt))

	response, err := a.executor.ExecutePrompt(prompt)
	if err != nil {
		logging.Error("Failed to execute vision prompt: %v", err)
		return nil, fmt.Errorf("failed to execute vision prompt: %w", err)
	}

	logging.Info("Received vision response (length: %d)", len(response))

	var vision VisionResponse
	err = json.Unmarshal([]byte(response), &vision)
	if err != nil {
		logging.Error("Failed to parse vision response: %v", err)
		logging.Debug("Response content that failed to parse: %s", response)
		return nil, fmt.Errorf("failed to parse vision response: %w", err)
	}

	logging.Info("Successfully parsed vision response: project=%s, vision_length=%d, epics_count=%d",
		vision.ProjectName, len(vision.ProductVision), len(vision.Epics))

	return &vision, nil
}

type StoryResponse struct {
	UserStories []UserStory `json:"user_stories"`
}

type UserStory struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Story    string `json:"story"`
	Priority int    `json:"priority"`
	EpicID   int    `json:"epic_id"`
}

func (a *Agent) GenerateStories(vision, epic string) (*StoryResponse, error) {
	logging.Info("Generating stories for epic: %s", epic)

	prompt, err := a.loader.LoadInitPrompt(2, map[string]interface{}{
		"ProductVision": vision,
		"Epic":          epic,
	})
	if err != nil {
		logging.Error("Failed to load stories prompt: %v", err)
		return nil, fmt.Errorf("failed to load stories prompt: %w", err)
	}

	logging.Info("Loaded stories prompt (length: %d)", len(prompt))

	response, err := a.executor.ExecutePrompt(prompt)
	if err != nil {
		logging.Error("Failed to execute stories prompt: %v", err)
		return nil, fmt.Errorf("failed to execute stories prompt: %w", err)
	}

	logging.Info("Received stories response (length: %d)", len(response))

	var stories StoryResponse
	err = json.Unmarshal([]byte(response), &stories)
	if err != nil {
		logging.Error("Failed to parse stories response: %v", err)
		logging.Debug("Response content that failed to parse: %s", response)
		return nil, fmt.Errorf("failed to parse stories response: %w", err)
	}

	logging.Info("Successfully parsed stories response: stories_count=%d", len(stories.UserStories))

	return &stories, nil
}

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	ID           int   `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	StoryID      int   `json:"story_id"`
	Dependencies []int  `json:"dependencies"`
}

func (a *Agent) GenerateTasks(storyTitle, story string) (*TaskResponse, error) {
	logging.Info("Generating tasks for story: %s", storyTitle)

	prompt, err := a.loader.LoadInitPrompt(3, map[string]interface{}{
		"StoryTitle": storyTitle,
		"UserStory":  story,
	})
	if err != nil {
		logging.Error("Failed to load tasks prompt: %v", err)
		return nil, fmt.Errorf("failed to load tasks prompt: %w", err)
	}

	logging.Info("Loaded tasks prompt (length: %d)", len(prompt))

	response, err := a.executor.ExecutePrompt(prompt)
	if err != nil {
		logging.Error("Failed to execute tasks prompt: %v", err)
		return nil, fmt.Errorf("failed to execute tasks prompt: %w", err)
	}

	logging.Info("Received tasks response (length: %d)", len(response))

	var tasks TaskResponse
	err = json.Unmarshal([]byte(response), &tasks)
	if err != nil {
		logging.Error("Failed to parse tasks response: %v", err)
		logging.Debug("Response content that failed to parse: %s", response)
		return nil, fmt.Errorf("failed to parse tasks response: %w", err)
	}

	logging.Info("Successfully parsed tasks response: tasks_count=%d", len(tasks.Tasks))

	return &tasks, nil
}

type BundleResponse struct {
	FeatureBundle []FeatureBundle `json:"feature_bundle"`
}

type FeatureBundle struct {
	Title    string `json:"title"`
	Story    string `json:"story"`
	Priority int    `json:"priority"`
	Tasks    []Task `json:"tasks"`
}

func (a *Agent) GenerateAddBundle(currentPlanJSON string, markdownInput string) (*BundleResponse, error) {
	logging.Info("Generating feature bundle: plan_length=%d, input_length=%d", len(currentPlanJSON), len(markdownInput))

	prompt, err := a.loader.LoadAddPrompt(1, map[string]interface{}{
		"CurrentPlanJSON": currentPlanJSON,
		"MarkdownInput":   markdownInput,
	})
	if err != nil {
		logging.Error("Failed to load bundle prompt: %v", err)
		return nil, fmt.Errorf("failed to load bundle prompt: %w", err)
	}

	logging.Info("Loaded bundle prompt (length: %d)", len(prompt))

	response, err := a.executor.ExecutePrompt(prompt)
	if err != nil {
		logging.Error("Failed to execute bundle prompt: %v", err)
		return nil, fmt.Errorf("failed to execute bundle prompt: %w", err)
	}

	logging.Info("Received bundle response (length: %d)", len(response))

	var bundle BundleResponse
	err = json.Unmarshal([]byte(response), &bundle)
	if err != nil {
		logging.Error("Failed to parse bundle response: %v", err)
		logging.Debug("Response content that failed to parse: %s", response)
		return nil, fmt.Errorf("failed to parse bundle response: %w", err)
	}

	logging.Info("Successfully parsed bundle response: bundle_count=%d", len(bundle.FeatureBundle))

	return &bundle, nil
}
