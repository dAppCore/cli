// Package cli provides the CLI runtime and utilities.
package cli

import (
	"os"

	"dappco.re/go/core/cli/internal/term"
)

// Mode represents the CLI execution mode.
//
//	mode := cli.DetectMode()
//	if mode == cli.ModeDaemon {
//	    cli.LogInfo("running headless")
//	}
type Mode int

const (
	// ModeInteractive indicates TTY attached with coloured output.
	ModeInteractive Mode = iota
	// ModePipe indicates stdout is piped, colours disabled.
	ModePipe
	// ModeDaemon indicates headless execution, log-only output.
	ModeDaemon
)

// String returns the string representation of the Mode.
func (m Mode) String() string {
	switch m {
	case ModeInteractive:
		return "interactive"
	case ModePipe:
		return "pipe"
	case ModeDaemon:
		return "daemon"
	default:
		return "unknown"
	}
}

// DetectMode determines the execution mode based on environment.
//
//	mode := cli.DetectMode()
//	// cli.ModeDaemon when CORE_DAEMON=1
//	// cli.ModePipe when stdout is not a terminal
//	// cli.ModeInteractive otherwise
func DetectMode() Mode {
	if os.Getenv("CORE_DAEMON") == "1" {
		return ModeDaemon
	}
	if !IsTTY() {
		return ModePipe
	}
	return ModeInteractive
}

// IsTTY returns true if stdout is a terminal.
//
//	if cli.IsTTY() {
//	    cli.Success("interactive output enabled")
//	}
func IsTTY() bool {
	if f, ok := stdoutWriter().(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

// IsStdinTTY returns true if stdin is a terminal.
//
//	if !cli.IsStdinTTY() {
//	    cli.Warn("input is piped")
//	}
func IsStdinTTY() bool {
	if f, ok := stdinReader().(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

// IsStderrTTY returns true if stderr is a terminal.
//
//	if cli.IsStderrTTY() {
//	    cli.Progress("load", 1, 3, "config")
//	}
func IsStderrTTY() bool {
	if f, ok := stderrWriter().(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}
