package server

import (
	"fmt"
	"time"
)

// LogLevel represents log severity
type LogLevel int

const (
	// LogDebug is debug level
	LogDebug LogLevel = iota
	// LogInfo is info level
	LogInfo
	// LogWarn is warning level
	LogWarn
	// LogError is error level
	LogError
)

// Logger is a simple logger
type Logger struct {
	level  LogLevel
	prefix string
}

// NewLogger creates a new Logger
func NewLogger(level LogLevel, prefix string) *Logger {
	return &Logger{level: level, prefix: prefix}
}

// Log logs a message at the given level
func (l *Logger) Log(level LogLevel, msg string) {
	if level >= l.level {
		fmt.Printf("[%s] %s %s: %s\n", l.prefix, time.Now().Format(time.RFC3339), levelName(level), msg)
	}
}

// levelName converts level to string - internal helper
func levelName(l LogLevel) string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// SetLevel changes the log level - dead code (never called)
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// formatJSON formats a log entry as JSON - dead code
func formatJSON(prefix, level, msg string) string {
	return fmt.Sprintf(`{"prefix":%q,"level":%q,"msg":%q}`, prefix, level, msg)
}
