package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plansmith/pkg/logging"
	"plansmith/pkg/tui"
	"strings" // Added import
	"time"    // Added import

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "plansmith",
	Short: "PlanSmith is an AI-powered project planning tool.",
	Long: `PlanSmith is an interactive, chat-like terminal application that \'crafts\'
raw project ideas (from Markdown) into fully-formed, actionable Kanban boards (in Trello),
with human review at every step.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize logger after viper has read the config
		parsedLogLevel, err := logging.ParseLogLevel(logLevel)
		if err != nil {
			log.Fatalf("Invalid log level: %v", err)
		}
		err = logging.InitGlobalLogger(parsedLogLevel)
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
		defer logging.CloseGlobalLogger()

		logging.Info("Starting PlanSmith application")
		StartTUI()
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
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".plansmith" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".plansmith"))
		viper.AddConfigPath(".") // optionally look for config in the working directory
		viper.SetConfigName("config")
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

func StartTUI() {
	logging.Info("Initializing TUI")

	model := tui.NewModel()

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
			saveChatHistory(m.GetChatMessages())
		}
	}

	logging.Info("TUI program finished")
}

func saveChatHistory(messages []tui.ChatMessage) {
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
