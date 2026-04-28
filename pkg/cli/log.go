package cli

import (
	"dappco.re/go"
)

// LogLevel aliases for convenience.
type LogLevel = core.Level

const (
	LogLevelQuiet = core.LevelQuiet
	LogLevelError = core.LevelError
	LogLevelWarn  = core.LevelWarn
	LogLevelInfo  = core.LevelInfo
	LogLevelDebug = core.LevelDebug
)

// LogDebug logs a debug message if the default logger is available.
//
//	cli.LogDebug("cache miss", "key", cacheKey)
func LogDebug(msg string, keyvals ...any) { core.Debug(msg, keyvals...) }

// LogInfo logs an info message.
//
//	cli.LogInfo("configuration reloaded", "path", configPath)
func LogInfo(msg string, keyvals ...any) { core.Info(msg, keyvals...) }

// LogWarn logs a warning message.
//
//	cli.LogWarn("GitHub CLI not authenticated", "user", username)
func LogWarn(msg string, keyvals ...any) { core.Warn(msg, keyvals...) }

// LogError logs an error message.
//
//	cli.LogError("Fatal error", "err", err)
func LogError(msg string, keyvals ...any) { core.Error(msg, keyvals...) }

// LogSecurity logs a security-sensitive message.
//
//	cli.LogSecurity("login attempt", "user", "admin")
func LogSecurity(msg string, keyvals ...any) { core.Security(msg, keyvals...) }

// LogSecurityf logs a formatted security-sensitive message.
//
//	cli.LogSecurityf("login attempt from %s", username)
func LogSecurityf(format string, args ...any) {
	core.Security(core.Sprintf(format, args...))
}
