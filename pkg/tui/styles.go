package tui

import "github.com/charmbracelet/lipgloss"

var (
	// General styles
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	// Message styles
	userMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4CAF50")).Bold(true)
	assistantMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2196F3")).Bold(true)
	systemMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9800")).Bold(true)
	errorMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	loadingMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	// Icons
	userIcon      = "👤"
	assistantIcon = "🤖"
	systemIcon    = "⚙️"
	errorIcon       = "❌"
	loadingIcon   = "⏳"

	// Path and file highlighting
	pathStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8")).Underline(true)
	fileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8"))

	// Autocomplete styles
	suggestionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	activeSuggestionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4"))

	// Input field style
	inputFieldStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	// Help text style
	helpTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true)
)
