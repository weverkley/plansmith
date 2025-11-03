package tui

// AppState represents the current state of the application
type AppState int

const (
	StateChat AppState = iota
)

// ConversationContext represents the current context of the conversation
type ConversationContext int

const (
	ContextNone ConversationContext = iota
	ContextWaitingForNewOrExisting
	ContextWaitingForFilePath
	ContextWaitingForPlanConfirmation
	ContextWaitingForExistingFilePath
	ContextWaitingForBoardName
)
