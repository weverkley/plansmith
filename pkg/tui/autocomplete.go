package tui

import (
	"fmt" // Added import
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"plansmith/pkg/logging"
)

type Suggestion struct {
	Text        string
	Description string
}

type Autocomplete struct {
	Suggestions []Suggestion
	Active      int
	Visible     bool
}

func NewAutocomplete() Autocomplete {
	return Autocomplete{}
}

func (a *Autocomplete) SetSuggestions(input string) {
	logging.Debug("Autocomplete: SetSuggestions called with input: %s", input)
	a.Suggestions = []Suggestion{}
	a.Active = 0 // Reset active selection
	a.Visible = false

	if input == "" {
		logging.Debug("Autocomplete: Empty input, hiding suggestions.")
		return
	}

	dir := filepath.Dir(input)
	base := filepath.Base(input)
	logging.Debug("Autocomplete: Dir: %s, Base: %s", dir, base)

	// In a real application, you might want to cache this or limit the scan depth
	// to avoid performance issues on very large codebases.
	var allPaths []string
	// Always scan the current directory for suggestions, regardless of path indicator.
	// This allows for more "fuzzy" matching like powerlevel10k.
	if dir == "." || dir == "./" || dir == "/" || dir == string(os.PathSeparator) { // Normalize to current directory
		dir = "."
	}
	if absDir, err := filepath.Abs(dir); err == nil {
		scanResult, err := a._scanDirectory(absDir, "", 3) // Scan up to 3 levels deep
		if err != nil {
			logging.Error("Autocomplete: Failed to scan directory %s: %v", absDir, err)
		}
		allPaths = append(allPaths, scanResult...)
	} else {
		// If the input has a directory part, scan that directory
		if absDir, err := filepath.Abs(dir); err == nil {
			scanResult, err := a._scanDirectory(absDir, "", 3) // Scan up to 3 levels deep
			if err != nil {
				logging.Error("Autocomplete: Failed to scan directory %s: %v", absDir, err)
			}
			allPaths = append(allPaths, scanResult...)
		}
	}

	for _, p := range allPaths {
		matchBase := strings.ToLower(filepath.Base(p))
		matchFull := strings.ToLower(p)
		lowerBase := strings.ToLower(base)

		// Fuzzy match: either starts with base or contains base
		if strings.HasPrefix(matchBase, lowerBase) || strings.Contains(matchFull, lowerBase) {
			info, err := os.Stat(p)
			description := "file"
			if err == nil && info.IsDir() {
				description = "directory"
			}
			a.Suggestions = append(a.Suggestions, Suggestion{Text: p, Description: description})
		}
	}

	if len(a.Suggestions) > 0 {
		a.Visible = true
		logging.Debug("Autocomplete: %d suggestions found, visible: true", len(a.Suggestions))
	} else {
		logging.Debug("Autocomplete: No suggestions found.")
	}
}

// _scanDirectory recursively scans a directory up to a given depth and returns
// a list of relative paths from the original call.
func (a *Autocomplete) _scanDirectory(root, currentPath string, depth int) ([]string, error) {
	if depth < 0 {
		return nil, nil // Stop recursion if depth limit reached
	}

	var results []string
	fullPath := filepath.Join(root, currentPath) // Construct the absolute path for scanning
	
	// Check if the directory exists and is accessible
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err // Return error if directory not found or inaccessible
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", fullPath)
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// Ignore hidden files/directories and common VCS/build folders
		if strings.HasPrefix(entry.Name(), ".") || entry.Name() == "node_modules" || entry.Name() == "vendor" || entry.Name() == "target" || entry.Name() == "build" {
			continue
		}

		entryPath := filepath.Join(currentPath, entry.Name())
		results = append(results, entryPath)

		if entry.IsDir() {
			subResults, err := a._scanDirectory(root, entryPath, depth-1)
			if err != nil {
				logging.Warn("Autocomplete: Error scanning subdirectory %s: %v", filepath.Join(root, entryPath), err)
				continue
			}
			results = append(results, subResults...)
		}
	}
	return results, nil
}

func (a *Autocomplete) View() string {
	if !a.Visible || len(a.Suggestions) == 0 {
		return ""
	}

	var s strings.Builder
	for i, suggestion := range a.Suggestions {
		style := suggestionStyle
		if i == a.Active {
			style = activeSuggestionStyle
		}
		s.WriteString(style.Render(suggestion.Text))
		s.WriteString("\n")
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		PaddingLeft(1).
		PaddingRight(1).
		Render(s.String())
}
