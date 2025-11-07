package tui

import (
	"plansmith/pkg/logging"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	logging.Debug("Rendering view")

	var bottomView string
	var help string

	if m.conversationContext == ContextWaitingForNewOrExisting ||
		m.conversationContext == ContextWaitingForVisionConfirmation ||
		m.conversationContext == ContextWaitingForStoriesConfirmation ||
		m.conversationContext == ContextWaitingForBoardCreationConfirmation ||
		m.conversationContext == ContextWaitingForPlanConfirmation ||
		m.conversationContext == ContextWaitingForBoardSelection {
		bottomView = m.confirmationList.View()
		help = helpStyle.Render("  enter: select | up/down: navigate")
	} else if m.filePicker.visible {
		bottomView = m.filePicker.View()
		help = helpStyle.Render("  ctrl+q: quit | @: browse files | arrow keys: scroll chat")
	} else {
		bottomView = m.textInput.View()
		autocompleteView := m.autocomplete.View()
		if m.autocomplete.Visible {
			help = helpStyle.Render("  esc: close autocomplete | tab: select | up/down: navigate")
			bottomView = lipgloss.JoinVertical(lipgloss.Left, bottomView, autocompleteView)
		} else {
			help = helpStyle.Render("  ctrl+q: quit | @: browse files | arrow keys: scroll chat")
		}
	}

	chatView := m.chat.View()

	return lipgloss.JoinVertical(lipgloss.Left,
		chatView,
		bottomView,
		help,
	)
}
