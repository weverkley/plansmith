package tui

import (
	"plansmith/pkg/logging"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	logging.Debug("Rendering view")

	// Calculate total frame sizes for appStyle
	appHorizontalFrameSize := lipgloss.Width(appStyle.Render(""))
	appVerticalFrameSize := lipgloss.Height(appStyle.Render(""))

	// Calculate available width and height for the main content area (inside appStyle)
	chatAvailableWidth := m.width - appHorizontalFrameSize
	chatAvailableHeight := m.height - appVerticalFrameSize

	// Set component dimensions based on available space
	// Title
	title := titleStyle.Width(chatAvailableWidth).Render("PlanSmith AI Agent")
	titleHeight := lipgloss.Height(title)

	// Help
	helpText := "  ctrl+q: quit | @: browse files | arrow keys: scroll chat"
	if m.conversationContext == ContextWaitingForNewOrExisting ||
		m.conversationContext == ContextWaitingForVisionConfirmation ||
		m.conversationContext == ContextWaitingForStoriesConfirmation ||
		m.conversationContext == ContextWaitingForBoardCreationConfirmation ||
		m.conversationContext == ContextWaitingForPlanConfirmation ||
		m.conversationContext == ContextWaitingForFeatureConfirmation ||
		m.conversationContext == ContextWaitingForBoardSelection {
		helpText = "  enter: select | up/down: navigate"
	} else if m.autocomplete.Visible {
		helpText = "  esc: close autocomplete | tab: select | up/down: navigate"
	}
	help := helpStyle.Render(helpText)
	helpHeight := lipgloss.Height(help)

	// Input
	var bottomView string
	inputContentWidth := chatAvailableWidth - inputStyle.GetHorizontalPadding()*2 - inputStyle.GetHorizontalBorderSize()*2
	if inputContentWidth < 0 {
		inputContentWidth = 0
	}
	m.textInput.Width = inputContentWidth // Set text input's internal width
	inputView := inputStyle.Width(chatAvailableWidth).Render(m.textInput.View())
	inputHeight := lipgloss.Height(inputView)

	if m.conversationContext == ContextWaitingForNewOrExisting ||
		m.conversationContext == ContextWaitingForVisionConfirmation ||
		m.conversationContext == ContextWaitingForStoriesConfirmation ||
		m.conversationContext == ContextWaitingForBoardCreationConfirmation ||
		m.conversationContext == ContextWaitingForPlanConfirmation ||
		m.conversationContext == ContextWaitingForFeatureConfirmation ||
		m.conversationContext == ContextWaitingForBoardSelection {
		m.confirmationList.SetWidth(chatAvailableWidth)
		bottomView = m.confirmationList.View()
	} else if m.filePicker.visible {
		m.filePicker.width = chatAvailableWidth
		m.filePicker.height = chatAvailableHeight - titleHeight - helpHeight - inputHeight // Remaining height for file picker
		bottomView = m.filePicker.View()
	} else {
		autocompleteView := m.autocomplete.View()
		if m.autocomplete.Visible {
			bottomView = lipgloss.JoinVertical(lipgloss.Left, inputView, autocompleteView)
		} else {
			bottomView = inputView
		}
	}

	// Chat
	m.chat.width = chatAvailableWidth
	m.chat.height = chatAvailableHeight - titleHeight - helpHeight - lipgloss.Height(bottomView) // Remaining height for chat
	chatView := m.chat.View()

	// Assemble the final view
	mainContent := lipgloss.JoinVertical(lipgloss.Left,
		title,
		chatView,
		bottomView,
		help,
	)

	return appStyle.Width(m.width).Height(m.height).Render(mainContent)
}
