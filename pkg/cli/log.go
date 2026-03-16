package cli

import (
	"forge.lthn.ai/core/go-log"
)

// LogLevel aliases for convenience.
type LogLevel = log.Level

const (
	LogLevelQuiet = log.LevelQuiet
	LogLevelError = log.LevelError
	LogLevelWarn  = log.LevelWarn
	LogLevelInfo  = log.LevelInfo
	LogLevelDebug = log.LevelDebug
)

// LogDebug logs a debug message if the default logger is available.
func LogDebug(msg string, keyvals ...any) { log.Debug(msg, keyvals...) }

// LogInfo logs an info message.
func LogInfo(msg string, keyvals ...any) { log.Info(msg, keyvals...) }

// LogWarn logs a warning message.
func LogWarn(msg string, keyvals ...any) { log.Warn(msg, keyvals...) }

// LogError logs an error message.
func LogError(msg string, keyvals ...any) { log.Error(msg, keyvals...) }
