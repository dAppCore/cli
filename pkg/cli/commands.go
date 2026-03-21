// Package cli provides the CLI runtime and utilities.
package cli

import (
	"io/fs"
	"iter"
	"sync"

	"dappco.re/go/core"
	"github.com/spf13/cobra"
)

// WithCommands returns a CommandSetup that registers a command group.
// The register function receives the root cobra command during Main().
//
//	cli.Main(
//	    cli.WithCommands("config", config.AddConfigCommands),
//	    cli.WithCommands("doctor", doctor.AddDoctorCommands),
//	)
func WithCommands(name string, register func(root *Command), localeFS ...fs.FS) CommandSetup {
	return func(c *core.Core) {
		if root, ok := c.App().Runtime.(*cobra.Command); ok {
			register(root)
		}
		// Register locale FS if provided
		if len(localeFS) > 0 && localeFS[0] != nil {
			registeredCommandsMu.Lock()
			registeredLocales = append(registeredLocales, localeFS[0])
			registeredCommandsMu.Unlock()
		}
	}
}

// CommandRegistration is a function that adds commands to the CLI root.
type CommandRegistration func(root *cobra.Command)

var (
	registeredCommands   []CommandRegistration
	registeredCommandsMu sync.Mutex
	commandsAttached     bool
	registeredLocales    []fs.FS
)

// RegisterCommands registers a function that adds commands to the CLI.
// Optionally pass a locale fs.FS to provide translations for the commands.
//
//	func init() {
//	    cli.RegisterCommands(AddCommands, locales.FS)
//	}
func RegisterCommands(fn CommandRegistration, localeFS ...fs.FS) {
	registeredCommandsMu.Lock()
	defer registeredCommandsMu.Unlock()
	registeredCommands = append(registeredCommands, fn)
	for _, lfs := range localeFS {
		if lfs != nil {
			registeredLocales = append(registeredLocales, lfs)
		}
	}

	// If commands already attached (CLI already running), attach immediately
	if commandsAttached && instance != nil && instance.root != nil {
		fn(instance.root)
	}
}

// RegisteredLocales returns all locale filesystems registered by command packages.
func RegisteredLocales() []fs.FS {
	registeredCommandsMu.Lock()
	defer registeredCommandsMu.Unlock()
	return registeredLocales
}

// RegisteredCommands returns an iterator over the registered command functions.
func RegisteredCommands() iter.Seq[CommandRegistration] {
	return func(yield func(CommandRegistration) bool) {
		registeredCommandsMu.Lock()
		defer registeredCommandsMu.Unlock()
		for _, fn := range registeredCommands {
			if !yield(fn) {
				return
			}
		}
	}
}

// attachRegisteredCommands calls all registered command functions.
// Called by Init() after creating the root command.
func attachRegisteredCommands(root *cobra.Command) {
	registeredCommandsMu.Lock()
	defer registeredCommandsMu.Unlock()

	for _, fn := range registeredCommands {
		fn(root)
	}
	commandsAttached = true
}
