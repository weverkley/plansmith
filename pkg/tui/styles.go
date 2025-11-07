package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(4).
			PaddingRight(4).
			MarginBottom(1)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2)

	selectedMenuItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				PaddingLeft(2).
				PaddingRight(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Background(lipgloss.Color("#FFFFFF")).
			Bold(true)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	// Styles for the confirmation list
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

	// Styles for chat messages
	userMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00BFFF"))
	userIcon = "👤"

	assistantMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#32CD32"))
	assistantIcon = "🤖"

	systemMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFF00"))
	systemIcon = "⚙️"

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000"))
	errorIcon = "🔥"

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF"))

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	helpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Italic(true)

	// Styles for autocomplete
	suggestionStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	activeSuggestionStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FFFFFF")).
				PaddingLeft(1)
)