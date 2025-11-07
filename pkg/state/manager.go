package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Manager struct {
	planPath  string
	statePath string
}

func NewManager() *Manager {
	// Use local directory for development
	// In production, this would be in the user's home directory
	dir := ".plansmith"

	// Create directory if it doesn't exist
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		// If we can't create the directory, use current directory
		dir = "."
	}

	return &Manager{
		statePath: filepath.Join(dir, "state.json"),
		planPath:  filepath.Join(dir, "plan.json"),
	}
}

func (m *Manager) SavePlan(plan *Plan) error {
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plan: %w", err)
	}

	err = os.WriteFile(m.planPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write plan file: %w", err)
	}

	return nil
}

func (m *Manager) LoadPlan() (*Plan, error) {
	// Check if file exists
	if _, err := os.Stat(m.planPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plan file does not exist")
	}

	data, err := os.ReadFile(m.planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	var plan Plan
	err = json.Unmarshal(data, &plan)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
	}

	return &plan, nil
}

func (m *Manager) LoadPlanFromPath(path string) (*Plan, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("plan file does not exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	var plan Plan
	err = json.Unmarshal(data, &plan)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
	}

	return &plan, nil
}

func (m *Manager) SaveState(state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	err = os.WriteFile(m.statePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func (m *Manager) LoadState() (*State, error) {
	// Check if file exists
	if _, err := os.Stat(m.statePath); os.IsNotExist(err) {
		// Return default state if file doesn't exist
		return &State{}, nil
	}

	data, err := os.ReadFile(m.statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &state, nil
}

func (m *Manager) PlanExists() bool {
	_, err := os.Stat(m.planPath)
	return !os.IsNotExist(err)
}
