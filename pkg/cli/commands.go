// Package cli provides the CLI runtime and utilities.
package cli

import (
	"context"
	"io/fs"
	"iter"
	"sync"

	"forge.lthn.ai/core/go/pkg/core"
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
// WithCommands creates a framework Option that registers a command group.
// Optionally pass a locale fs.FS as the third argument to provide translations.
//
//	cli.WithCommands("dev", dev.AddDevCommands, locales.FS)
func WithCommands(name string, register func(root *Command), localeFS ...fs.FS) core.Option {
	return core.WithName("cmd."+name, func(c *core.Core) (any, error) {
		svc := &commandService{core: c, name: name, register: register}
		if len(localeFS) > 0 {
			svc.localeFS = localeFS[0]
		}
		return svc, nil
	})
}

type commandService struct {
	core     *core.Core
	name     string
	register func(root *Command)
	localeFS fs.FS
}

func (s *commandService) OnStartup(_ context.Context) error {
	if root, ok := s.core.App.(*cobra.Command); ok {
		s.register(root)
		// Auto-set Short/Long from i18n keys derived from command name.
		// The Conclave's i18n service has already loaded all translations
		// from sibling services' LocaleProvider before commands attach.
		s.applyI18n(root)
	}
	return nil
}

// applyI18n walks commands added by this service and sets Short/Long
// from derived i18n keys if they're empty or still raw keys.
func (s *commandService) applyI18n(root *cobra.Command) {
	for _, cmd := range root.Commands() {
		key := "cmd." + cmd.Name()
		// Only set if Short is empty or looks like a raw key (contains dots)
		if cmd.Short == "" || cmd.Short == key+".short" {
			if translated := T(key + ".short"); translated != key+".short" {
				cmd.Short = translated
			}
		}
		if cmd.Long == "" || cmd.Long == key+".long" {
			if translated := T(key + ".long"); translated != key+".long" {
				cmd.Long = translated
			}
		}
	}
}

// Locales implements core.LocaleProvider.
func (s *commandService) Locales() fs.FS {
	return s.localeFS
}

// CommandRegistration is a function that adds commands to the root.
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

