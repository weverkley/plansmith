package tui

import "github.com/charmbracelet/lipgloss"

const (
	white       = lipgloss.Color("#FFFFFF")
	black       = lipgloss.Color("#000000")
	purple      = lipgloss.Color("#7D56F4")
	lightPurple = lipgloss.Color("#AD58B4")
	darkGray    = lipgloss.Color("#333333")
	lightGray   = lipgloss.Color("#808080")
	blue        = lipgloss.Color("#00BFFF")
	green       = lipgloss.Color("#32CD32")
	yellow      = lipgloss.Color("#FFFF00")
	red         = lipgloss.Color("#FF0000")
	cyan        = lipgloss.Color("#00FFFF")
)

var (
	// General
	appStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(purple).
			Padding(1, 2).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lightGray).
			Italic(true).
			Padding(0, 1)

	// Input
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(0, 1)

	// Spinner
	spinnerStyle = lipgloss.NewStyle().
			Foreground(purple)

	// List
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(purple).
			Padding(0, 1)

	NormalTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#dddddd"}).
				Padding(0, 0, 0, 2)

	NormalDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
				Padding(0, 0, 0, 2)

	SelectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
				Padding(0, 0, 0, 1)

	SelectedDescStyle = SelectedTitleStyle.Copy().
				Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	DimmedTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
				Padding(0, 0, 0, 2)

	DimmedDescStyle = DimmedTitleStyle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	// Chat
	chatContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(purple).
				Padding(1)

	userMessageStyle = lipgloss.NewStyle().
				Foreground(blue).
				Padding(0, 1)
	userIcon = "👤"

	assistantMessageStyle = lipgloss.NewStyle().
				Foreground(green).
				Padding(0, 1)
	assistantIcon = "🤖"

	systemMessageStyle = lipgloss.NewStyle().
				Foreground(yellow).
				Padding(0, 1)
	systemIcon = "⚙️"

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(red).
				Bold(true).
				Padding(0, 1)
	errorIcon = "🔥"

	pathStyle = lipgloss.NewStyle().
			Foreground(cyan)

	fileStyle = lipgloss.NewStyle().
			Foreground(green)

	// Autocomplete
	suggestionStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	activeSuggestionStyle = lipgloss.NewStyle().
				Background(purple).
				Foreground(white).
				PaddingLeft(1)
)
