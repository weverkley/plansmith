package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogLevel represents the verbosity of logging.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of a LogLevel.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel converts a string to a LogLevel.
func ParseLogLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":	
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return INFO, fmt.Errorf("unknown log level: %s", level)
	}
}

type Logger struct {
	logFile    *os.File
	logger     *log.Logger
	currentLevel LogLevel
}

func NewLogger(level LogLevel, logDir string) (*Logger, error) {
	// Create logs directory
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := filepath.Join(logDir, fmt.Sprintf("plansmith_%s.log", timestamp))

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	logger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	return &Logger{
		logFile:    logFile,
		logger:     logger,
		currentLevel: level,
	}, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.currentLevel <= INFO {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.currentLevel <= ERROR {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.currentLevel <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.currentLevel <= WARN {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}

// Global logger instance
var globalLogger *Logger

func InitGlobalLogger(level LogLevel, logDir string) error {
	logger, err := NewLogger(level, logDir)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

func CloseGlobalLogger() {
	if globalLogger != nil {
		globalLogger.Close()
	}
}

func Info(format string, v ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(format, v...)
	}
}

func GlobalLogger() *Logger {
	return globalLogger
}
