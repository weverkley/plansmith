package tui

import (
	"fmt"
	"os"
	"path/filepath" // Added import
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"plansmith/pkg/logging"
	"plansmith/pkg/state"
	"plansmith/pkg/trello"
)

// Messages for async operations
type generationMsg struct {
	plan *state.Plan
	err  error
}

type trelloMsg struct {
	boardURL string
	err      error
}

type listModelsMsg struct {
	models []string
	err    error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logging.Debug("Main Update function received message of type %T", msg)
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case initializationMsg:
		m.chat.AddMessage("system", "Plansmith is here!")
		m.chat.AddMessage("system", "Like a blacksmith, but for crafting project plans. PlanSmith is an interactive, chat-like terminal application that 'crafts' raw project ideas (from Markdown) into fully-formed, actionable Kanban boards (in Trello), with human review at every step.")
		m.chat.AddMessage("system", "Version: v1.0")
		m.chat.AddMessage("system", "Would you like to start a new project or open an existing one? (new/existing)")
		m.conversationContext = ContextWaitingForNewOrExisting
		return m, nil

	case cursor.BlinkMsg,
		spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.chat.width = msg.Width
		m.filePicker.width = msg.Width
		m.filePicker.height = msg.Height - 6

		var inputView string
		if m.filePicker.visible {
			inputView = m.filePicker.View()
		} else {
			inputView = m.textInput.View()
		}
		help := helpStyle.Render("  ctrl+q: quit | @: browse files | arrow keys: scroll chat")
		otherComponentsHeight := lipgloss.Height(inputView) + lipgloss.Height(help)
		m.chat.viewport.Height = m.height - otherComponentsHeight

	case tea.KeyMsg:
		logging.Debug("KeyMsg received: %s", msg.String())
		if m.filePicker.visible {
			m.filePicker, cmd = m.filePicker.Update(msg)
			return m, cmd
		}

		if m.textInput.Focused() {
			switch msg.String() {
			case "tab":
				if m.autocomplete.Visible && len(m.autocomplete.Suggestions) > 0 {
					selected := m.autocomplete.Suggestions[m.autocomplete.Active]
					// Ensure we're replacing only the base part of the path being typed
					currentInput := m.textInput.Value()
					dir := filepath.Dir(currentInput)
					if dir == "." { // Handle current directory case
						dir = ""
					} else {
						dir += string(os.PathSeparator)
					}
					m.textInput.SetValue(dir + selected.Text)
					m.textInput.CursorEnd()
					m.autocomplete.Visible = false
					logging.Debug("Autocomplete: Tab pressed, selected '%s', new input '%s'", selected.Text, m.textInput.Value())
				}
			case "up":
				if m.autocomplete.Visible {
					m.autocomplete.Active--
					if m.autocomplete.Active < 0 {
						m.autocomplete.Active = len(m.autocomplete.Suggestions) - 1
					}
					logging.Debug("Autocomplete: Up pressed, active suggestion index: %d", m.autocomplete.Active)
				} else {
					m.chat, cmd = m.chat.Update(msg)
					cmds = append(cmds, cmd)
				}
			case "down":
				if m.autocomplete.Visible {
					m.autocomplete.Active++
					if m.autocomplete.Active >= len(m.autocomplete.Suggestions) {
						m.autocomplete.Active = 0
					}
					logging.Debug("Autocomplete: Down pressed, active suggestion index: %d", m.autocomplete.Active)
				} else {
					m.chat, cmd = m.chat.Update(msg)
					cmds = append(cmds, cmd)
				}
			case "esc":
				if m.autocomplete.Visible {
					m.autocomplete.Visible = false
					logging.Debug("Autocomplete: Esc pressed, hiding suggestions.")
					return m, nil
				}
			case "enter":
				if m.autocomplete.Visible && len(m.autocomplete.Suggestions) > 0 {
					selected := m.autocomplete.Suggestions[m.autocomplete.Active]
					newValue := selected.Text
					if selected.Description == "directory" {
						newValue += string(os.PathSeparator)
					}
					m.textInput.SetValue(newValue)
					m.textInput.CursorEnd()
					m.autocomplete.SetSuggestions(m.textInput.Value()) // Re-evaluate suggestions after completion
					m.autocomplete.Visible = true // Keep autocomplete visible if there are new suggestions
					logging.Debug("Autocomplete: Enter pressed, selected '%s', new input '%s'", selected.Text, m.textInput.Value())
					// After selecting, the input is now complete, so we fall through to submission logic
				}

				// Submission logic (moved outside the if/else for autocomplete selection)
				input := m.textInput.Value()
				logging.Info("User input: %s", input)
				if input == "/exit" {
					return m, tea.Quit
				}
				m.textInput.SetValue("")
				m.autocomplete.Visible = false // Hide autocomplete after enter
				m.chat.AddMessage("user", input)

				switch m.conversationContext {
				case ContextWaitingForNewOrExisting:
					if strings.ToLower(input) == "new" {
						m.chat.AddMessage("assistant", "Great! Please provide the path to your project's markdown file.")
						m.conversationContext = ContextWaitingForFilePath
						logging.Info("Conversation context changed to ContextWaitingForFilePath")
					} else if strings.ToLower(input) == "existing" {
						m.chat.AddMessage("assistant", "Great! Please provide the path to your project's .json file.")
						m.conversationContext = ContextWaitingForExistingFilePath
						logging.Info("Conversation context changed to ContextWaitingForExistingFilePath")
					} else {
						m.chat.AddMessage("assistant", "Sorry, I didn't understand that. Please type 'new' or 'existing'.")
						logging.Warn("Invalid input for new/existing project: %s", input)
					}
				case ContextWaitingForExistingFilePath:
					logging.Info("Attempting to load plan from path: %s", input)
					plan, err := m.stateManager.LoadPlanFromPath(input)
					if err != nil {
						m.chat.AddMessage("system", fmt.Sprintf("Error loading plan: %v", err))
						logging.Error("Failed to load plan from path %s: %v", input, err)
					} else {
						m.plan = plan
						formattedPlan := formatPlan(m.plan)
						m.chat.AddMessage("assistant", "Here is the loaded plan:\n"+formattedPlan)
						m.chat.AddMessage("assistant", "What would you like to do with this plan?")
						m.conversationContext = ContextNone // Or a new context for project dashboard
						logging.Info("Plan loaded successfully from %s", input)
					}
				case ContextWaitingForFilePath:
					m.markdownPath = input
					m.chat.AddMessage("assistant", fmt.Sprintf("Thanks! I'll start crafting a plan from '%s'.", input))
					m.chat.SetLoading(true)
					cmds = append(cmds, m.generatePlan())
					m.conversationContext = ContextNone
					logging.Info("Markdown path set to %s, starting plan generation.", input)
				case ContextWaitingForPlanConfirmation:
					if strings.ToLower(input) == "yes" {
						m.chat.AddMessage("assistant", "Great! What would you like to name the Trello board?")
						m.chat.AddMessage("assistant", fmt.Sprintf("I'll suggest a name for you: %s", m.plan.ProjectName))
						m.conversationContext = ContextWaitingForBoardName
						logging.Info("Plan confirmed, waiting for Trello board name.")
					} else if strings.ToLower(input) == "no" {
						m.chat.AddMessage("assistant", "Okay, the plan has been discarded. What would you like to do next?")
						m.conversationContext = ContextNone
						logging.Info("Plan discarded by user.")
					} else {
						m.chat.AddMessage("assistant", "Sorry, I didn't understand that. Please type 'yes' to confirm or 'no' to discard the plan.")
						logging.Warn("Invalid input for plan confirmation: %s", input)
					}
				case ContextWaitingForBoardName:
					boardName := input
					if boardName == "" {
						boardName = m.plan.ProjectName
						logging.Info("Using default board name: %s", boardName)
					}
					m.chat.AddMessage("assistant", fmt.Sprintf("Great! I'll proceed with creating the Trello board named '%s'.", boardName))
					m.chat.SetLoading(true)
					cmds = append(cmds, m.createTrelloBoard(boardName))
					m.conversationContext = ContextNone
					logging.Info("Creating Trello board with name: %s", boardName)
				default:
					logging.Debug("No specific conversation context, processing as general chat/command.")
					// Process as a command or general chat
					// cmds = append(cmds, m.processChatCommand(input))
				}
			default:
				m.textInput, cmd = m.textInput.Update(msg)
				cmds = append(cmds, cmd)
				// Only update autocomplete suggestions if the input value has changed
				if m.textInput.Value() != "" && (strings.HasSuffix(m.textInput.Value(), "/") || strings.HasSuffix(m.textInput.Value(), "./") || strings.HasSuffix(m.textInput.Value(), "../") || strings.Contains(m.textInput.Value(), string(os.PathSeparator))) {
					m.autocomplete.SetSuggestions(m.textInput.Value())
					m.autocomplete.Visible = true
					logging.Debug("Autocomplete suggestions updated for input: %s", m.textInput.Value())
				} else {
					m.autocomplete.Visible = false
				}
			}
		}

		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			logging.Info("Quit command received (ctrl+c or ctrl+q).")
			return m, tea.Quit
		case "@":
			m.filePicker.visible = true
			m.filePicker.updateFiles()
			logging.Info("File picker activated.")
			return m, nil
		}

		// If text input is not focused, or if it's an arrow key and autocomplete is not visible,
		// then it's likely for chat scrolling.
		if !m.textInput.Focused() || (msg.Type == tea.KeyRunes && (msg.String() == "up" || msg.String() == "down") && !m.autocomplete.Visible) {
			m.chat, cmd = m.chat.Update(msg)
			cmds = append(cmds, cmd)
		}

	case fileSelectedMsg:
		m.textInput.SetValue(m.textInput.Value() + msg.path)
		m.filePicker.visible = false

	case generationMsg:
		m.chat.SetLoading(false)
		if msg.err != nil {
			m.chat.AddMessage("system", fmt.Sprintf("Error generating plan: %v", msg.err))
		} else {
			m.plan = msg.plan
			formattedPlan := formatPlan(m.plan)
			m.chat.AddMessage("assistant", "Here's the plan I've crafted:\n"+formattedPlan)
			m.chat.AddMessage("assistant", "Please review it. Type 'yes' to confirm or 'no' to discard.")
			m.conversationContext = ContextWaitingForPlanConfirmation
		}

	case trelloMsg:
		m.chat.SetLoading(false)
		if msg.err != nil {
			m.chat.AddMessage("system", fmt.Sprintf("Error creating Trello board: %v", msg.err))
		} else {
			m.chat.AddMessage("assistant", fmt.Sprintf("I've created a Trello board for you at %s", msg.boardURL))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) generatePlan() tea.Cmd {
	return func() tea.Msg {
		logging.Info("Generating plan from file: %s", m.markdownPath)

		markdown, err := os.ReadFile(m.markdownPath)
		if err != nil {
			return generationMsg{err: fmt.Errorf("failed to read markdown file: %w", err)}
		}

		m.chat.AddMessage("assistant", "Generating vision...")
		vision, err := m.agent.GenerateVision(string(markdown))
		if err != nil {
			return generationMsg{err: fmt.Errorf("failed to generate vision: %w", err)}
		}

		plan := &state.Plan{
			ProjectName:   vision.ProjectName,
			ProductVision: vision.ProductVision,
			Epics:         []state.Epic{},
			UserStories:   []state.UserStory{},
			Tasks:         []state.Task{},
		}

		err = m.stateManager.SavePlan(plan)
		if err != nil {
			logging.Error("Failed to save plan: %v", err)
			return generationMsg{err: err}
		}

		for _, epicName := range vision.Epics {
			plan.Epics = append(plan.Epics, state.Epic{Name: epicName})
			m.chat.AddMessage("assistant", fmt.Sprintf("Generating stories for epic '%s'...", epicName))
			stories, err := m.agent.GenerateStories(vision.ProductVision, epicName)
			if err != nil {
				return generationMsg{err: fmt.Errorf("failed to generate stories for epic %s: %w", epicName, err)}
			}

			for _, story := range stories.UserStories {
				plan.UserStories = append(plan.UserStories, state.UserStory{Title: story.Title, Story: story.Story, Priority: story.Priority, Epic: story.Epic})
				m.chat.AddMessage("assistant", fmt.Sprintf("Generating tasks for story '%s'...", story.Title))
				tasks, err := m.agent.GenerateTasks(story.Title, story.Story)
				if err != nil {
					return generationMsg{err: fmt.Errorf("failed to generate tasks for story %s: %w", story.Title, err)}
				}

				for _, task := range tasks.Tasks {
					plan.Tasks = append(plan.Tasks, state.Task{Title: task.Title, Description: task.Description, StoryTitle: task.StoryTitle, Dependencies: task.Dependencies})
				}
				err = m.stateManager.SavePlan(plan)
				if err != nil {
					logging.Error("Failed to save plan: %v", err)
					return generationMsg{err: err}
				}
			}
		}

		return generationMsg{plan: plan}
	}
}

func (m Model) createTrelloBoard(boardName string) tea.Cmd {
	return func() tea.Msg {
		logging.Info("Creating Trello board for plan: %s", m.plan.ProjectName)
		if m.trelloClient == nil {
			return trelloMsg{err: fmt.Errorf("Trello client not initialized")}
		}

		board, err := m.trelloClient.CreateProjectBoard(boardName)
		if err != nil {
			return trelloMsg{err: fmt.Errorf("failed to create Trello board: %w", err)}
		}

		trelloPlan := &trello.Plan{
			ProjectName:   m.plan.ProjectName,
			ProductVision: m.plan.ProductVision,
			Epics:         []trello.Epic{},
			UserStories:   []trello.UserStory{},
			Tasks:         []trello.Task{},
		}

		for _, epic := range m.plan.Epics {
			trelloPlan.Epics = append(trelloPlan.Epics, trello.Epic{Name: epic.Name})
		}

		for _, story := range m.plan.UserStories {
			trelloPlan.UserStories = append(trelloPlan.UserStories, trello.UserStory{Title: story.Title, Story: story.Story, Priority: story.Priority, Epic: story.Epic})
		}

		for _, task := range m.plan.Tasks {
			trelloPlan.Tasks = append(trelloPlan.Tasks, trello.Task{Title: task.Title, Description: task.Description, StoryTitle: task.StoryTitle, Dependencies: task.Dependencies})
		}

		err = m.trelloClient.PopulateBoard(board.ID, trelloPlan)
		if err != nil {
			return trelloMsg{err: fmt.Errorf("failed to populate Trello board: %w", err)}
		}

		err = m.trelloClient.LinkCards(board.ID, trelloPlan)
		if err != nil {
			return trelloMsg{err: fmt.Errorf("failed to link cards on Trello board: %w", err)}
		}

		return trelloMsg{boardURL: board.URL}
	}
}

func formatPlan(plan *state.Plan) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Project Name: %s\n", plan.ProjectName))
	b.WriteString(fmt.Sprintf("Product Vision: %s\n", plan.ProductVision))
	b.WriteString("Epics:\n")
	for _, epic := range plan.Epics {
		b.WriteString(fmt.Sprintf("- %s\n", epic.Name))
	}
	b.WriteString("User Stories:\n")
	for _, story := range plan.UserStories {
		b.WriteString(fmt.Sprintf("- %s (Epic: %s)\n", story.Title, story.Epic))
	}
	b.WriteString("Tasks:\n")
	for _, task := range plan.Tasks {
		b.WriteString(fmt.Sprintf("- %s (Story: %s): %s\n", task.Title, task.StoryTitle, task.Description))
	}
	return b.String()
}
