// Package cli provides the CLI runtime and utilities.
package cli

import (
	"io/fs"
	"iter"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/i18n"
)

// WithCommands returns a CommandSetup that registers a command group.
// The register function receives the Core instance during Main().
//
//	cli.Main(
//	    cli.WithCommands("config", config.AddConfigCommands),
//	    cli.WithCommands("doctor", doctor.AddDoctorCommands),
//	)
func WithCommands(name string, register CommandRegistration, localeFS ...fs.FS) CommandSetup {
	return func(c *core.Core) {
		loadLocaleSources(localeSourcesFromFS(localeFS...)...)
		register(c)
		appendLocales(localeFS...)
	}
}

// CommandRegistration is a function that adds commands to the Core instance.
//
// Example:
//
//	func addCommands(c *core.Core) {
//	    c.Command("ping", core.Command{
//	        Description: "Ping API",
//	        Action: func(opts core.Options) core.Result {
//	            cli.Println("pong")
//	            return core.Ok(nil)
//	        },
//	    })
//	}
type CommandRegistration func(c *core.Core)

var (
	registeredCommands   []CommandRegistration
	registeredCommandsMu core.Mutex
	commandsAttached     bool
	registeredLocales    []fs.FS
)

// RegisterCommands registers a function that adds commands to the CLI.
// Optionally pass a locale fs.FS to provide translations for the commands.
//
//	func init() {
//	    cli.RegisterCommands(AddCommands, locales.FS)
//	}
//
// Example:
//
//	cli.RegisterCommands(func(c *core.Core) {
//	    c.Command("version", core.Command{
//	        Description: "Show version",
//	        Action: func(opts core.Options) core.Result {
//	            cli.Println(cli.SemVer())
//	            return core.Ok(nil)
//	        },
//	    })
//	})
func RegisterCommands(fn CommandRegistration, localeFS ...fs.FS) {
	registeredCommandsMu.Lock()
	registeredCommands = append(registeredCommands, fn)
	attached := commandsAttached && instance != nil && instance.core != nil
	coreInstance := instance
	registeredCommandsMu.Unlock()

	loadLocaleSources(localeSourcesFromFS(localeFS...)...)
	appendLocales(localeFS...)

	// If commands already attached (CLI already running), attach immediately
	if attached {
		fn(coreInstance.core)
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
//
// Example:
//
//	for _, fs := range cli.RegisteredLocales() {
//		_ = fs
//	}
func RegisteredLocales() []fs.FS {
	registeredCommandsMu.Lock()
	defer registeredCommandsMu.Unlock()
	if len(registeredLocales) == 0 {
		return nil
	}
	out := make([]fs.FS, len(registeredLocales))
	copy(out, registeredLocales)
	return out
}

// RegisteredCommands returns an iterator over the registered command functions.
//
// Example:
//
//	for attach := range cli.RegisteredCommands() {
//		_ = attach
//	}
func RegisteredCommands() iter.Seq[CommandRegistration] {
	return func(yield func(CommandRegistration) bool) {
		registeredCommandsMu.Lock()
		snapshot := make([]CommandRegistration, len(registeredCommands))
		copy(snapshot, registeredCommands)
		registeredCommandsMu.Unlock()

		for _, fn := range snapshot {
			if !yield(fn) {
				return
			}
		}
	}
}

// attachRegisteredCommands calls all registered command functions.
// Called by Init() after creating the Core instance.
func attachRegisteredCommands(c *core.Core) {
	registeredCommandsMu.Lock()
	snapshot := make([]CommandRegistration, len(registeredCommands))
	copy(snapshot, registeredCommands)
	commandsAttached = true
	registeredCommandsMu.Unlock()

	for _, fn := range snapshot {
		fn(c)
	}
}
