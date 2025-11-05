package tui

import (
	"plansmith/pkg/ai"
	"plansmith/pkg/logging"
	"plansmith/pkg/smith"
	"plansmith/pkg/state"
	"plansmith/pkg/trello"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type ConversationContext int

const (
	ContextNone ConversationContext = iota
	ContextWaitingForNewOrExisting
	ContextWaitingForFilePath
	ContextWaitingForExistingFilePath
	ContextWaitingForPlanConfirmation
	ContextWaitingForBoardName
	ContextWaitingForVisionConfirmation
	ContextWaitingForStoriesConfirmation
	ContextWaitingForBoardCreationConfirmation
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type initializationMsg struct{}

type Model struct {
	conversationContext ConversationContext
	textInput           textinput.Model
	filePicker          FilePicker
	autocomplete        Autocomplete
	chat                *Chat
	viewport            viewport.Model
	spinner             spinner.Model
	stateManager        *state.Manager
	width               int
	height              int
	err                 error
	generatingMsg       string
	markdownPath        string
	confirmationList    list.Model

	// AI and services
	agent        *smith.Agent
	trelloClient *trello.Client
	plan         *state.Plan
	aiProvider   string
	aiModel      string
}

func NewModel(dummy bool) Model {
	logging.Info("Creating new TUI model")

	ti := textinput.New()
	ti.Placeholder = "Type your command or @ to browse files..."
	ti.Focus()
	logging.Debug("TextInput initialized and focused.")

	fp := NewFilePicker()
	logging.Debug("FilePicker initialized.")
	ac := NewAutocomplete()
	logging.Debug("Autocomplete initialized.")
	chat := NewChat()
	logging.Debug("Chat initialized.")

	vp := viewport.New(80, 20)
	logging.Debug("Viewport initialized.")

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	logging.Debug("Spinner initialized.")

	// Initialize AI executor
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.plansmith")
	viper.AddConfigPath(".")
	logging.Debug("Viper config paths set.")

	var executor ai.AIExecutor
	var trelloClient *trello.Client
	aiProvider := "gemini"
	aiModel := ""

	if dummy {
		executor = &ai.DummyExecutor{}
	} else {
		if err := viper.ReadInConfig(); err == nil {
			logging.Info("Successfully loaded configuration file")

			// AI configuration
			aiProvider = viper.GetString("ai.default_provider")
			var apiKey string
			if aiProvider == "gemini" {
				apiKey = viper.GetString("ai.keys.gemini")
				aiModel = viper.GetString("ai.models.gemini")
			} else if aiProvider == "openai" {
				apiKey = viper.GetString("ai.keys.openai")
				aiModel = viper.GetString("ai.models.openai")
			} else if aiProvider == "qwen" {
				apiKey = viper.GetString("ai.keys.qwen")
				aiModel = viper.GetString("ai.models.qwen")
			}

			logging.Info("AI provider: %s, Model: %s", aiProvider, aiModel)

			if apiKey != "" {
				var err error
				executor, err = ai.NewExecutor(aiProvider, apiKey, aiModel)
				if err != nil {
					logging.Error("Failed to create AI executor: %v", err)
				} else {
					logging.Info("Successfully created AI executor")
				}
			} else {
				logging.Warn("No API key found for provider: %s", aiProvider)
			}

		} else {
			logging.Warn("Failed to load configuration file: %v", err)
		}
	}

	// Trello configuration
	key := viper.GetString("trello.key")
	token := viper.GetString("trello.token")
	if key != "" && token != "" {
		trelloClient = trello.NewClient(key, token)
		logging.Info("Successfully created Trello client")
	} else {
		logging.Warn("Trello key or token not found in config")
	}

	// If we can't load config or create executor, create a dummy one
	if executor == nil {
		logging.Info("Using dummy executor")
		executor = &ai.DummyExecutor{}
	}

	agent := smith.NewAgent(executor)
	logging.Debug("Smith Agent initialized.")

		// Initialize confirmation list
	confirmationList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	confirmationList.SetShowTitle(false)
	confirmationList.SetShowHelp(false)
	confirmationList.SetFilteringEnabled(false)
	confirmationList.SetHeight(2)
	confirmationList.SetWidth(10)
	logging.Debug("Confirmation list initialized.")

	return Model{
		conversationContext: ContextWaitingForNewOrExisting,
		textInput:           ti,
		filePicker:          fp,
		autocomplete:        ac,
		chat:                &chat,
		viewport:            vp,
		spinner:             sp,
		stateManager:        state.NewManager(),
		agent:               agent,
		trelloClient:        trelloClient,
		plan:                &state.Plan{},
		aiProvider:          aiProvider,
		aiModel:             aiModel,
		confirmationList:    confirmationList,
	}
}

func (m Model) Init() tea.Cmd {
	logging.Info("Initializing TUI model")
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			return initializationMsg{}
		},
	)
}

// GetChatMessages returns the current chat messages.
func (m Model) GetChatMessages() []ChatMessage {
	return m.chat.messages
}
