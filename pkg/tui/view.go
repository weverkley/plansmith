package tui

import (
	"plansmith/pkg/logging"

	"github.com/charmbracelet/lipgloss"
)

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
)

func (m Model) View() string {
	logging.Debug("Rendering view")

	var inputView string
	if m.filePicker.visible {
		inputView = m.filePicker.View()
	} else {
		inputView = m.textInput.View()
	}

	autocompleteView := m.autocomplete.View()

	var help string
	if m.autocomplete.Visible {
		help = helpStyle.Render("  esc: close autocomplete | tab: select | up/down: navigate")
	} else {
		help = helpStyle.Render("  ctrl+q: quit | @: browse files | arrow keys: scroll chat")
	}

	chatView := m.chat.View()

	return lipgloss.JoinVertical(lipgloss.Left,
		chatView,
		inputView,
		autocompleteView,
		help,
	)
}
