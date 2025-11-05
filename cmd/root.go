package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plansmith/pkg/ai"
	"plansmith/pkg/logging"
	"plansmith/pkg/smith"
	"plansmith/pkg/state"
	"plansmith/pkg/trello"
	"plansmith/pkg/tui"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logLevel string
var plansmithDir string // Global variable to store the resolved .plansmith directory

var dummy bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "plansmith",
	Short: "PlanSmith is an AI-powered project planning tool.",
	Long: `PlanSmith is an interactive, chat-like terminal application that 'crafts'
raw project ideas (from Markdown) into fully-formed, actionable Kanban boards (in Trello),
with human review at every step.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize logger after viper has read the config and plansmithDir is resolved
		parsedLogLevel, err := logging.ParseLogLevel(logLevel)
		if err != nil {
			log.Fatalf("Invalid log level: %v", err)
		}
		err = logging.InitGlobalLogger(parsedLogLevel, filepath.Join(plansmithDir, "logs"))
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
		defer logging.CloseGlobalLogger()

		logging.Info("Starting PlanSmith application")
		if len(args) == 0 {
			StartTUI(dummy)
		} else {
			runNonInteractive(args)
		}
		logging.Info("PlanSmith application finished")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra flags can be global and persistent across subcommands.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.plansmith/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set the logging level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&dummy, "dummy", false, "Use a dummy AI executor for testing")
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Determine the .plansmith directory
	cwd, err := os.Getwd()
	cobra.CheckErr(err)
	cwdPlansmithDir := filepath.Join(cwd, ".plansmith")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	homePlansmithDir := filepath.Join(home, ".plansmith")

	if _, err := os.Stat(cwdPlansmithDir); err == nil {
		plansmithDir = cwdPlansmithDir
	} else if _, err := os.Stat(homePlansmithDir); err == nil {
		plansmithDir = homePlansmithDir
	} else {
		// If neither exists, create in CWD
		err := os.MkdirAll(cwdPlansmithDir, 0755)
		cobra.CheckErr(err)
		plansmithDir = cwdPlansmithDir
		fmt.Fprintln(os.Stderr, "Created .plansmith directory in current working directory:", plansmithDir)
	}

	// Set viper config paths
	viper.AddConfigPath(plansmithDir)
	viper.SetConfigName("config")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		// Handle errors reading the config file, but don\'t fatal if it\'s just not found
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
			fmt.Fprintln(os.Stderr, "No config file found, using defaults and environment variables.")
		} else {
			// Some other error occurred
			fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		}
	}
}

func StartTUI(dummy bool) {
	logging.Info("Initializing TUI")

	model := tui.NewModel(dummy)

	logging.Info("Starting bubbletea program")
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		logging.Error("Failed to start bubbletea program: %v", err)
		panic(err)
	}

	// Check if chat messages should be saved on exit
	if viper.GetBool("save_chat_on_exit") {
		if m, ok := finalModel.(tui.Model); ok {
			saveChatHistory(m.GetChatMessages(), plansmithDir)
		}
	}

	logging.Info("TUI program finished")
}

func saveChatHistory(messages []tui.ChatMessage, plansmithDir string) {
	logDir := filepath.Join(plansmithDir, "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		logging.Error("Failed to create logs directory for chat history: %v", err)
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := filepath.Join(logDir, fmt.Sprintf("chat_history_%s.log", timestamp))

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logging.Error("Failed to create chat history file: %v", err)
		return
	}
	defer file.Close()

	for _, msg := range messages {
		_, err := file.WriteString(fmt.Sprintf("[%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), strings.ToUpper(msg.Role), msg.Content))
		if err != nil {
			logging.Error("Failed to write chat message to history file: %v", err)
			return
		}
	}
	logging.Info("Chat history saved to %s", logFileName)
}

func runNonInteractive(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: plansmith <markdown_file> [trello_board_name]")
		return
	}

	markdownFile := args[0]
	var boardName string
	if len(args) > 1 {
		boardName = args[1]
	}

	// Initialize AI executor
	var executor ai.AIExecutor
	if dummy {
		executor = &ai.DummyExecutor{}
	} else {
		aiProvider := viper.GetString("ai.default_provider")
		var apiKey string
		if aiProvider == "gemini" {
			apiKey = viper.GetString("ai.keys.gemini")
		} else if aiProvider == "openai" {
			apiKey = viper.GetString("ai.keys.openai")
		} else if aiProvider == "qwen" {
			apiKey = viper.GetString("ai.keys.qwen")
		}

		if apiKey != "" {
			var err error
			executor, err = ai.NewExecutor(aiProvider, apiKey, "")
			if err != nil {
				log.Fatalf("Failed to create AI executor: %v", err)
			}
		} else {
			log.Fatalf("No API key found for provider: %s", aiProvider)
		}
	}

	// Read markdown file
	markdown, err := os.ReadFile(markdownFile)
	if err != nil {
		log.Fatalf("Failed to read markdown file: %v", err)
	}

	// Create smith agent
	agent := smith.NewAgent(executor)

	// Generate plan
	vision, err := agent.GenerateVision(string(markdown))
	if err != nil {
		log.Fatalf("Failed to generate vision: %v", err)
	}

	plan := &state.Plan{
		ProjectName:   vision.ProjectName,
		ProductVision: vision.ProductVision,
	}

	for i, epic := range vision.Epics {
		plan.Epics = append(plan.Epics, state.Epic{ID: smith.GenerateID("EPIC", i+1), Name: epic.Name})
	}

	storyCounter := 0
	taskCounter := 0
	for _, epic := range plan.Epics {
		stories, err := agent.GenerateStories(plan.ProductVision, epic.Name, epic.ID)
		if err != nil {
			log.Fatalf("Failed to generate stories for epic %s: %v", epic.Name, err)
		}
		for _, story := range stories.UserStories {
			storyID := smith.GenerateID("STORY", storyCounter+1)
			plan.UserStories = append(plan.UserStories, state.UserStory{ID: storyID, Title: story.Title, Story: story.Story, Priority: story.Priority, EpicID: epic.ID})
			storyCounter++

			tasks, err := agent.GenerateTasks(story.Title, story.Story, storyID)
			if err != nil {
				log.Fatalf("Failed to generate tasks for story %s: %v", story.Title, err)
			}
			for _, task := range tasks.Tasks {
				plan.Tasks = append(plan.Tasks, state.Task{ID: smith.GenerateID("TASK", taskCounter+1), Title: task.Title, Description: task.Description, StoryID: storyID, Dependencies: task.Dependencies, Labels: task.Labels})
				taskCounter++
			}
		}
	}

	// Create Trello client
	key := viper.GetString("trello.key")
	token := viper.GetString("trello.token")
	if key == "" || token == "" {
		log.Fatalf("Trello key or token not found in config")
	}
	trelloClient := trello.NewClient(key, token)

	// Create Trello board
	if boardName == "" {
		boardName = plan.ProjectName
	}

	trelloBoard, err := trelloClient.CreateProjectBoard(boardName)
	if err != nil {
		log.Fatalf("Failed to create Trello board: %v", err)
	}

	trelloPlan := &trello.Plan{
		ProjectName:   plan.ProjectName,
		ProductVision: plan.ProductVision,
		Epics:         []trello.Epic{},
		UserStories:   []trello.UserStory{},
		Tasks:         []trello.Task{},
	}

	for _, epic := range plan.Epics {
		trelloPlan.Epics = append(trelloPlan.Epics, trello.Epic{ID: epic.ID, Name: epic.Name})
	}

	for _, story := range plan.UserStories {
		trelloPlan.UserStories = append(trelloPlan.UserStories, trello.UserStory{ID: story.ID, Title: story.Title, Story: story.Story, Priority: story.Priority, EpicID: story.EpicID})
	}

	for _, task := range plan.Tasks {
		trelloPlan.Tasks = append(trelloPlan.Tasks, trello.Task{ID: task.ID, Title: task.Title, Description: task.Description, StoryID: task.StoryID, Dependencies: task.Dependencies, Labels: task.Labels})
	}

	err = trelloClient.PopulateBoard(trelloBoard.ID, trelloPlan)
	if err != nil {
		log.Fatalf("Failed to populate Trello board: %v", err)
	}

	err = trelloClient.LinkCards(trelloBoard.ID, trelloPlan)
	if err != nil {
		log.Fatalf("Failed to link cards on Trello board: %v", err)
	}

	// Save the plan
	stateManager := state.NewManager()
	err = stateManager.SavePlan(plan)
	if err != nil {
		log.Fatalf("Failed to save plan: %v", err)
	}

	fmt.Printf("Trello board created successfully: %s\n", trelloBoard.URL)
}
