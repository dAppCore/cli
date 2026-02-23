// Package cli provides the CLI runtime and utilities.
package cli

import (
	"context"
	"iter"
	"sync"

	"forge.lthn.ai/core/go/pkg/framework"
	"github.com/spf13/cobra"
)

// WithCommands creates a framework Option that registers a command group.
// The register function receives the root command during service startup,
// allowing commands to participate in the Core lifecycle.
//
//	cli.Main(
//	    cli.WithCommands("config", config.AddConfigCommands),
//	    cli.WithCommands("doctor", doctor.AddDoctorCommands),
//	)
func WithCommands(name string, register func(root *Command)) framework.Option {
	return framework.WithName("cmd."+name, func(c *framework.Core) (any, error) {
		return &commandService{core: c, register: register}, nil
	})
}

type commandService struct {
	core     *framework.Core
	register func(root *Command)
}

func (s *commandService) OnStartup(_ context.Context) error {
	if root, ok := s.core.App.(*cobra.Command); ok {
		s.register(root)
	}
	return nil
}

// CommandRegistration is a function that adds commands to the root.
type CommandRegistration func(root *cobra.Command)

var (
	registeredCommands   []CommandRegistration
	registeredCommandsMu sync.Mutex
	commandsAttached     bool
)

// RegisterCommands registers a function that adds commands to the CLI.
// Call this in your package's init() to register commands.
//
//	func init() {
//	    cli.RegisterCommands(AddCommands)
//	}
//
//	func AddCommands(root *cobra.Command) {
//	    root.AddCommand(myCmd)
//	}
func RegisterCommands(fn CommandRegistration) {
	registeredCommandsMu.Lock()
	defer registeredCommandsMu.Unlock()
	registeredCommands = append(registeredCommands, fn)

	// If commands already attached (CLI already running), attach immediately
	if commandsAttached && instance != nil && instance.root != nil {
		fn(instance.root)
	}
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

