package tui

import (
	"fmt"
	"os"
	"path/filepath" // Added import
	"strings"

	"github.com/charmbracelet/bubbles/list"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"plansmith/pkg/logging"
	"plansmith/pkg/smith"
	"plansmith/pkg/state"
	"plansmith/pkg/trello"
)

// Messages for async operations
type visionGeneratedMsg struct {
	vision *smith.VisionResponse
	err    error
}

type storiesForEpicGeneratedMsg struct {
	stories *smith.StoryResponse
	err     error
}

type allStoriesGeneratedMsg struct {
	stories []state.UserStory
	err     error
}

type tasksForStoryGeneratedMsg struct {
	tasks *smith.TaskResponse
	err   error
}

type allTasksGeneratedMsg struct {
	tasks []state.Task
	err   error
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
		m.chat.AddMessage("system", "Would you like to start a new project or open an existing one?")

		items := []list.Item{
			item{title: "new", desc: "Start a new project"},
			item{title: "existing", desc: "Open an existing project"},
		}
		m.confirmationList.SetItems(items)
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
		logging.Debug("KeyMsg received: %s, Context: %v", msg.String(), m.conversationContext)


		if m.conversationContext == ContextWaitingForNewOrExisting ||
			m.conversationContext == ContextWaitingForVisionConfirmation ||
			m.conversationContext == ContextWaitingForStoriesConfirmation ||
			m.conversationContext == ContextWaitingForBoardCreationConfirmation ||
			m.conversationContext == ContextWaitingForPlanConfirmation {
			
			var listCmd tea.Cmd
			m.confirmationList, listCmd = m.confirmationList.Update(msg)
			cmds = append(cmds, listCmd)

			switch msg.String() {
			case "enter":
				selectedItem := m.confirmationList.SelectedItem().(item)
				input := selectedItem.title
				m.chat.AddMessage("user", input)
				m.textInput.Focus()
				
				switch m.conversationContext {
				case ContextWaitingForNewOrExisting:
					if input == "new" {
						m.chat.AddMessage("assistant", "Great! Please provide the path to your project's markdown file.")
						m.conversationContext = ContextWaitingForFilePath
						logging.Info("Conversation context changed to ContextWaitingForFilePath")
					} else if input == "existing" {
						m.chat.AddMessage("assistant", "Great! Please provide the path to your project's .json file.")
						m.conversationContext = ContextWaitingForExistingFilePath
						logging.Info("Conversation context changed to ContextWaitingForExistingFilePath")
					}
				case ContextWaitingForVisionConfirmation:
					if input == "yes" {
						m.chat.AddMessage("assistant", "Generating stories...")
						m.chat.SetLoading(true)
						cmds = append(cmds, m.generateStoriesCmd())
						m.conversationContext = ContextNone // Reset context after starting generation
						logging.Info("Vision confirmed, starting story generation.")
					} else if input == "no" {
						m.chat.AddMessage("assistant", "Vision discarded. Please provide a new markdown file to start over.")
						m.conversationContext = ContextWaitingForFilePath // Allow user to provide a new file
						logging.Info("Vision discarded by user.")
					}
				case ContextWaitingForStoriesConfirmation:
					if input == "yes" {
						m.chat.AddMessage("assistant", "Generating tasks...")
						m.chat.SetLoading(true)
						cmds = append(cmds, m.generateTasksCmd())
						m.conversationContext = ContextNone // Reset context after starting generation
						logging.Info("Stories confirmed, starting task generation.")
					} else if input == "no" {
						m.chat.AddMessage("assistant", "Stories discarded. Please provide a new markdown file to start over.")
						m.conversationContext = ContextWaitingForFilePath // Allow user to provide a new file
						logging.Info("Stories discarded by user.")
					}
				case ContextWaitingForBoardCreationConfirmation:
					if input == "yes" {
						m.chat.AddMessage("assistant", fmt.Sprintf("Great! I'll proceed with creating the Trello board named '%s'.", m.plan.ProjectName))
						m.chat.SetLoading(true)
						cmds = append(cmds, m.createTrelloBoard(m.plan.ProjectName))
						m.conversationContext = ContextNone // Reset context after starting creation
						logging.Info("Board creation confirmed, starting Trello board creation.")
					} else if input == "no" {
						m.chat.AddMessage("assistant", "Trello board creation cancelled. What would you like to do next?")
						m.conversationContext = ContextNone // Or a new context for project dashboard
						logging.Info("Trello board creation cancelled by user.")
					}
				case ContextWaitingForPlanConfirmation:
					if input == "yes" {
						m.chat.AddMessage("assistant", "Great! What would you like to name the Trello board?")
						m.chat.AddMessage("assistant", fmt.Sprintf("I'll suggest a name for you: %s", m.plan.ProjectName))
						m.conversationContext = ContextWaitingForBoardName
						logging.Info("Plan confirmed, waiting for Trello board name.")
					} else if input == "no" {
						m.chat.AddMessage("assistant", "Okay, the plan has been discarded. What would you like to do next?")
						m.conversationContext = ContextNone
						logging.Info("Plan discarded by user.")
					}
				}
				return m, tea.Batch(cmds...)
			}
		}

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
					m.autocomplete.Visible = true                       // Keep autocomplete visible if there are new suggestions
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
					// Check if the file exists at the given path. If not, try to find it.
					if _, err := os.Stat(input); os.IsNotExist(err) {
						logging.Info("File not found at '%s', searching for it...", input)
						filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if !info.IsDir() && info.Name() == input {
								logging.Info("Found file at '%s'", path)
								input = path
							}
							return nil
						})
					}
					m.markdownPath = input
					m.chat.AddMessage("assistant", fmt.Sprintf("Thanks! I'll start crafting a plan from '%s'.", input))
					m.chat.SetLoading(true)
					cmds = append(cmds, m.generateVisionCmd())
					m.conversationContext = ContextNone
					logging.Info("Markdown path set to %s, starting plan generation.", input)
				case ContextWaitingForBoardName:
					boardName := input
					if boardName == "" {
						boardName = m.plan.ProjectName
						logging.Info("Using default board name: %s", boardName)
					}
					m.plan.ProjectName = boardName // Update plan with confirmed board name
					m.chat.AddMessage("assistant", fmt.Sprintf("You chose to name the Trello board '%s'. Do you want to create it now? (yes/no)", boardName))
					items := []list.Item{
						item{title: "yes", desc: "Create Trello board"},
						item{title: "no", desc: "Cancel Trello board creation"},
					}
					m.confirmationList.SetItems(items)
					m.conversationContext = ContextWaitingForBoardCreationConfirmation
					return m, nil
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
		// When a file is selected, we replace the entire input with the file path.
		m.textInput.SetValue(msg.path)
		m.filePicker.visible = false
		m.textInput.CursorEnd()

	case visionGeneratedMsg:
		m.chat.SetLoading(false)
		if msg.err != nil {
			m.chat.AddMessage("system", fmt.Sprintf("Error generating vision: %v", msg.err))
		} else {
			m.plan = &state.Plan{
				ProjectName:   msg.vision.ProjectName,
				ProductVision: msg.vision.ProductVision,
			}
			for _, epic := range msg.vision.Epics {
				m.plan.Epics = append(m.plan.Epics, state.Epic{ID: epic.ID, Name: epic.Name})
			}
			formattedPlan := formatPlan(m.plan)
			m.chat.AddMessage("assistant", "Here's the generated vision:\n"+formattedPlan)
			m.chat.AddMessage("assistant", "Do you want to proceed with generating stories based on this vision?")
			
			items := []list.Item{
				item{title: "yes", desc: "Proceed with this vision"},
				item{title: "no", desc: "Discard and restart vision generation"},
			}
			m.confirmationList.SetItems(items)
			m.conversationContext = ContextWaitingForVisionConfirmation
		}
	case storiesForEpicGeneratedMsg:
		if msg.err != nil {
			m.chat.AddMessage("system", fmt.Sprintf("Error generating stories for epic: %v", msg.err))
		} else {
			if len(m.plan.Epics) > 0 {
				cmds = append(cmds, m.generateStoriesCmd())
			} else {
				formattedPlan := formatPlan(m.plan)
				m.chat.AddMessage("assistant", "All stories generated! Here's the plan so far:\n"+formattedPlan)
				m.chat.AddMessage("assistant", "Do you want to proceed with generating tasks based on these stories? (yes/no)")
				items := []list.Item{
					item{title: "yes", desc: "Proceed with these stories"},
					item{title: "no", desc: "Discard and restart story generation"},
				}
				m.confirmationList.SetItems(items)
				m.conversationContext = ContextWaitingForStoriesConfirmation
			}
		}

	case allStoriesGeneratedMsg:
		m.chat.SetLoading(false)
		if msg.err != nil {
			m.chat.AddMessage("system", fmt.Sprintf("Error generating stories: %v", msg.err))
		} else {
			formattedPlan := formatPlan(m.plan)
			m.chat.AddMessage("assistant", "All stories generated! Here's the plan so far:\n"+formattedPlan)
			m.chat.AddMessage("assistant", "Do you want to proceed with generating tasks based on these stories? (yes/no)")
			items := []list.Item{
				item{title: "yes", desc: "Proceed with these stories"},
				item{title: "no", desc: "Discard and restart story generation"},
			}
			m.confirmationList.SetItems(items)
			m.conversationContext = ContextWaitingForStoriesConfirmation
		}
	case tasksForStoryGeneratedMsg:
		if msg.err != nil {
			m.chat.AddMessage("system", fmt.Sprintf("Error generating tasks for story: %v", msg.err))
		} else {
			if len(m.plan.UserStories) > 0 {
				cmds = append(cmds, m.generateTasksCmd())
			} else {
				m.chat.AddMessage("assistant", "Plan generated successfully!")
				formattedPlan := formatPlan(m.plan)
				m.chat.AddMessage("assistant", "Here's the plan I've crafted:\n"+formattedPlan)
				m.chat.AddMessage("assistant", "Please review it. Type 'yes' to confirm or 'no' to discard.")
				m.conversationContext = ContextWaitingForPlanConfirmation
			}
		}

	case allTasksGeneratedMsg:
				m.chat.SetLoading(false)
				if msg.err != nil {
					m.chat.AddMessage("system", fmt.Sprintf("Error generating tasks: %v", msg.err))
				} else {
					m.chat.AddMessage("assistant", "Plan generated successfully!")
					// Save the plan to a file
					err := m.stateManager.SavePlan(m.plan)
					if err != nil {
						m.chat.AddMessage("system", fmt.Sprintf("Error writing plan to file: %v", err))
					} else {
						m.chat.AddMessage("assistant", fmt.Sprintf("I've saved the plan to '.plansmith/plan.json'."))
					}
					m.chat.AddMessage("assistant", "Please review the plan. Do you want to confirm it? (yes/no)")
					items := []list.Item{
						item{title: "yes", desc: "Confirm the plan"},
						item{title: "no", desc: "Discard the plan"},
					}
					m.confirmationList.SetItems(items)
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

func (m *Model) generateVisionCmd() tea.Cmd {
	return func() tea.Msg {
		logging.Info("Generating vision from file: %s", m.markdownPath)

		markdown, err := os.ReadFile(m.markdownPath)
		if err != nil {
			return visionGeneratedMsg{err: fmt.Errorf("failed to read markdown file: %w", err)}
		}

		m.chat.AddMessage("assistant", "Generating vision...")
		vision, err := m.agent.GenerateVision(string(markdown))
		if err != nil {
			return visionGeneratedMsg{err: fmt.Errorf("failed to generate vision: %w", err)}
		}

		for i := range vision.Epics {
			vision.Epics[i].ID = smith.GenerateID("EPIC", i+1)
		}

		return visionGeneratedMsg{vision: vision}
	}
}

func (m *Model) generateStoriesCmd() tea.Cmd {
	return func() tea.Msg {
		if len(m.plan.Epics) == 0 {
			return allStoriesGeneratedMsg{stories: m.plan.UserStories}
		}

		epic := m.plan.Epics[0]
		m.plan.Epics = m.plan.Epics[1:]

		m.chat.AddMessage("assistant", fmt.Sprintf("Generating stories for epic '%s'...", epic.Name))
		stories, err := m.agent.GenerateStories(m.plan.ProductVision, epic.Name, epic.ID)
		if err != nil {
			return storiesForEpicGeneratedMsg{err: fmt.Errorf("failed to generate stories for epic %s: %w", epic.Name, err)}
		}
		for i := range stories.UserStories {
			stories.UserStories[i].ID = smith.GenerateID("STORY", len(m.plan.UserStories) + i + 1)
			stories.UserStories[i].EpicID = epic.ID
			m.plan.UserStories = append(m.plan.UserStories, state.UserStory{ID: stories.UserStories[i].ID, Title: stories.UserStories[i].Title, Story: stories.UserStories[i].Story, Priority: stories.UserStories[i].Priority, EpicID: stories.UserStories[i].EpicID})
		}

		return storiesForEpicGeneratedMsg{stories: stories}
	}
}

func (m *Model) generateTasksCmd() tea.Cmd {

	return func() tea.Msg {

		if len(m.plan.UserStories) == 0 {

			return allTasksGeneratedMsg{tasks: m.plan.Tasks}

		}



		story := m.plan.UserStories[0]

		m.plan.UserStories = m.plan.UserStories[1:]



		m.chat.AddMessage("assistant", fmt.Sprintf("Generating tasks for story '%s'...", story.Title))

		tasks, err := m.agent.GenerateTasks(story.Title, story.Story, story.ID)

		if err != nil {

			return tasksForStoryGeneratedMsg{err: fmt.Errorf("failed to generate tasks for story %s: %w", story.Title, err)}

		}

		for i := range tasks.Tasks {

			tasks.Tasks[i].ID = smith.GenerateID("TASK", len(m.plan.Tasks) + i + 1)

			tasks.Tasks[i].StoryID = story.ID

			m.plan.Tasks = append(m.plan.Tasks, state.Task{ID: tasks.Tasks[i].ID, Title: tasks.Tasks[i].Title, Description: tasks.Tasks[i].Description, StoryID: tasks.Tasks[i].StoryID, Dependencies: tasks.Tasks[i].Dependencies, Labels: tasks.Tasks[i].Labels})

		}



		return tasksForStoryGeneratedMsg{tasks: tasks}

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
			trelloPlan.Epics = append(trelloPlan.Epics, trello.Epic{ID: epic.ID, Name: epic.Name})
		}

		for _, story := range m.plan.UserStories {
			trelloPlan.UserStories = append(trelloPlan.UserStories, trello.UserStory{ID: story.ID, Title: story.Title, Story: story.Story, Priority: story.Priority, EpicID: story.EpicID})
		}

		for _, task := range m.plan.Tasks {
			trelloPlan.Tasks = append(trelloPlan.Tasks, trello.Task{ID: task.ID, Title: task.Title, Description: task.Description, StoryID: task.StoryID, Dependencies: task.Dependencies, Labels: task.Labels})
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
		b.WriteString(fmt.Sprintf("- [%s] %s\n", epic.ID, epic.Name))
	}
	b.WriteString("User Stories:\n")
	for _, story := range plan.UserStories {
		b.WriteString(fmt.Sprintf("- [%s] %s (Epic: %s)\n", story.ID, story.Title, story.EpicID))
	}
	b.WriteString("Tasks:\n")
	for _, task := range plan.Tasks {
		b.WriteString(fmt.Sprintf("- [%s] %s (Story: %s): %s\n", task.ID, task.Title, task.StoryID, task.Description))
	}
	return b.String()
}
