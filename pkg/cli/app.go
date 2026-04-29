package cli

import (
	"embed"
	"io/fs"
	"runtime/debug"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/i18n"
)

//go:embed locales/*.json
var cliLocaleFS embed.FS

// AppName is the default CLI application name.
// Override with WithAppName before calling Main.
var AppName = "core"

// Build-time variables set via ldflags (SemVer 2.0.0):
//
//	go build -ldflags="-X dappco.re/go/cli/pkg/cli.AppVersion=1.2.0 \
//	  -X dappco.re/go/cli/pkg/cli.BuildCommit=df94c24 \
//	  -X dappco.re/go/cli/pkg/cli.BuildDate=2026-02-06 \
//	  -X dappco.re/go/cli/pkg/cli.BuildPreRelease=dev.8"
var (
	AppVersion      = "0.0.0"
	BuildCommit     = "unknown"
	BuildDate       = "unknown"
	BuildPreRelease = ""
)

// SemVer returns the full SemVer 2.0.0 version string.
//
// Examples:
//
//	// Release only:
//	// AppVersion=1.2.0 -> 1.2.0
//	cli.AppVersion = "1.2.0"
//	fmt.Println(cli.SemVer())
//
//	// Pre-release + commit + date:
//	// AppVersion=1.2.0, BuildPreRelease=dev.8, BuildCommit=df94c24, BuildDate=20260206
//	// -> 1.2.0-dev.8+df94c24.20260206
func SemVer() string {
	v := AppVersion
	if BuildPreRelease != "" {
		v += "-" + BuildPreRelease
	}
	if BuildCommit != "unknown" {
		v += "+" + BuildCommit
		if BuildDate != "unknown" {
			v += "." + BuildDate
		}
	}
	return v
}

// WithAppName sets the application name used in help text.
// Call before Main for variant binaries (e.g. "lem", "devops").
//
//	cli.WithAppName("lem")
//	cli.Main()
func WithAppName(name string) {
	AppName = name
}

// LocaleSource pairs a filesystem with a directory for loading translations.
type LocaleSource = i18n.FSSource

// WithLocales returns a locale source for use with MainWithLocales.
//
// Example:
//
//	fs := embed.FS{}
//	locales := cli.WithLocales(fs, "locales")
//	cli.MainWithLocales([]cli.LocaleSource{locales})
func WithLocales(fsys fs.FS, dir string) LocaleSource {
	return LocaleSource{FS: fsys, Dir: dir}
}

// CommandSetup is a function that registers commands on the CLI after init.
//
// Example:
//
//	cli.Main(
//	    cli.WithCommands("doctor", doctor.AddDoctorCommands),
//	)
type CommandSetup func(c *core.Core)

// Main initialises and runs the CLI with the framework's built-in translations.
//
// Example:
//
//	cli.WithAppName("core")
//	cli.Main(config.AddConfigCommands)
func Main(commands ...CommandSetup) {
	MainWithLocales(nil, commands...)
}

// MainWithLocales initialises and runs the CLI with additional translation sources.
//
// Example:
//
//	locales := []cli.LocaleSource{cli.WithLocales(embeddedLocales, "locales")}
//	cli.MainWithLocales(locales, doctor.AddDoctorCommands)
func MainWithLocales(locales []LocaleSource, commands ...CommandSetup) {
	// Recovery from panics
	defer func() {
		if r := recover(); r != nil {
			core.Error("recovered from panic", "error", r, "stack", string(debug.Stack()))
			Shutdown()
			Fatal(core.E("Main", core.Sprintf("panic: %v", r), nil))
		}
	}()

	// Build locale sources: framework built-in + caller's extras + registered packages
	extraFS := []i18n.FSSource{
		{FS: cliLocaleFS, Dir: "locales"},
	}
	extraFS = append(extraFS, locales...)
	for _, lfs := range RegisteredLocales() {
		extraFS = append(extraFS, i18n.FSSource{FS: lfs, Dir: "."})
	}

	// Initialise CLI runtime
	if r := Init(Options{
		AppName:     AppName,
		Version:     SemVer(),
		I18nSources: extraFS,
	}); !r.OK {
		Error(r.Error())
		core.Exit(1)
	}
	defer Shutdown()

	c := Core()

	// Set banner on the CLI
	cl := c.Cli()
	if cl != nil {
		cl.SetBanner(func(_ *core.Cli) string {
			return core.Concat(AppName, " ", SemVer())
		})
	}

	// Run command setup functions
	for _, setup := range commands {
		setup(c)
	}

	if r := Execute(); !r.OK {
		code := 1
		var exitErr *ExitError
		if err, ok := r.Value.(error); ok && As(err, &exitErr) {
			code = exitErr.Code
		}
		Error(r.Error())
		c.Exit(code)
	}
}
