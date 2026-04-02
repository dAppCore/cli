// Package cli provides the CLI runtime and utilities.
package cli

import (
	"io/fs"
	"iter"
	"sync"

	"dappco.re/go/core"
	"forge.lthn.ai/core/go-i18n"
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
		loadLocaleSources(localeSourcesFromFS(localeFS...)...)
		if root, ok := c.App().Runtime.(*cobra.Command); ok {
			register(root)
		}
		appendLocales(localeFS...)
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
	registeredCommands = append(registeredCommands, fn)
	attached := commandsAttached && instance != nil && instance.root != nil
	root := instance
	registeredCommandsMu.Unlock()

	loadLocaleSources(localeSourcesFromFS(localeFS...)...)
	appendLocales(localeFS...)

	// If commands already attached (CLI already running), attach immediately
	if attached {
		fn(root.root)
	}
}

// appendLocales appends non-nil locale filesystems to the registry.
func appendLocales(localeFS ...fs.FS) {
	var nonempty []fs.FS
	for _, lfs := range localeFS {
		if lfs != nil {
			nonempty = append(nonempty, lfs)
		}
	}
	if len(nonempty) == 0 {
		return
	}
	registeredCommandsMu.Lock()
	registeredLocales = append(registeredLocales, nonempty...)
	registeredCommandsMu.Unlock()
}

func localeSourcesFromFS(localeFS ...fs.FS) []LocaleSource {
	sources := make([]LocaleSource, 0, len(localeFS))
	for _, lfs := range localeFS {
		if lfs != nil {
			sources = append(sources, LocaleSource{FS: lfs, Dir: "."})
		}
	}
	return sources
}

func loadLocaleSources(sources ...LocaleSource) {
	svc := i18n.Default()
	if svc == nil {
		return
	}
	for _, src := range sources {
		if src.FS == nil {
			continue
		}
		if err := svc.AddLoader(i18n.NewFSLoader(src.FS, src.Dir)); err != nil {
			LogDebug("failed to load locale source", "dir", src.Dir, "err", err)
		}
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
