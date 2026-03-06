package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"forge.lthn.ai/core/go-crypt/crypt/openpgp"
	"forge.lthn.ai/core/go/pkg/framework"
	"forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go-io/workspace"
	"github.com/spf13/cobra"
)

// AppName is the default CLI application name.
// Override with WithAppName before calling Main.
var AppName = "core"

// Build-time variables set via ldflags (SemVer 2.0.0):
//
//	go build -ldflags="-X forge.lthn.ai/core/cli/pkg/cli.AppVersion=1.2.0 \
//	  -X forge.lthn.ai/core/cli/pkg/cli.BuildCommit=df94c24 \
//	  -X forge.lthn.ai/core/cli/pkg/cli.BuildDate=2026-02-06 \
//	  -X forge.lthn.ai/core/cli/pkg/cli.BuildPreRelease=dev.8"
var (
	AppVersion      = "0.0.0"
	BuildCommit     = "unknown"
	BuildDate       = "unknown"
	BuildPreRelease = ""
)

// SemVer returns the full SemVer 2.0.0 version string.
//   - Release:  1.2.0
//   - Pre-release: 1.2.0-dev.8
//   - Full:     1.2.0-dev.8+df94c24.20260206
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

// WithAppName sets the application name used in help text and shell completion.
// Call before Main for variant binaries (e.g. "lem", "devops").
//
//	cli.WithAppName("lem")
//	cli.Main()
func WithAppName(name string) {
	AppName = name
}

// Main initialises and runs the CLI application.
// Pass command services via WithCommands to register CLI commands
// through the Core framework lifecycle.
//
//	cli.Main(
//	    cli.WithCommands("config", config.AddConfigCommands),
//	    cli.WithCommands("doctor", doctor.AddDoctorCommands),
//	)
//
// Exits with code 1 on error or panic.
func Main(commands ...framework.Option) {
	// Recovery from panics
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovered from panic", "error", r, "stack", string(debug.Stack()))
			Shutdown()
			Fatal(fmt.Errorf("panic: %v", r))
		}
	}()

	// Core services load first, then command services
	services := []framework.Option{
		framework.WithName("i18n", NewI18nService(I18nOptions{})),
		framework.WithName("log", NewLogService(log.Options{
			Level: log.LevelInfo,
		})),
		framework.WithName("crypt", openpgp.New),
		framework.WithName("workspace", workspace.New),
	}
	services = append(services, commands...)

	// Initialise CLI runtime with services
	if err := Init(Options{
		AppName:  AppName,
		Version:  SemVer(),
		Services: services,
	}); err != nil {
		Error(err.Error())
		os.Exit(1)
	}
	defer Shutdown()

	// Add completion command to the CLI's root
	RootCmd().AddCommand(newCompletionCmd())

	if err := Execute(); err != nil {
		code := 1
		var exitErr *ExitError
		if As(err, &exitErr) {
			code = exitErr.Code
		}
		Error(err.Error())
		os.Exit(code)
	}
}

// newCompletionCmd creates the shell completion command using the current AppName.
func newCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: fmt.Sprintf(`Generate shell completion script for the specified shell.

To load completions:

Bash:
  $ source <(%s completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ %s completion bash > /etc/bash_completion.d/%s
  # macOS:
  $ %s completion bash > $(brew --prefix)/etc/bash_completion.d/%s

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ %s completion zsh > "${fpath[1]}/_%s"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ %s completion fish | source

  # To load completions for each session, execute once:
  $ %s completion fish > ~/.config/fish/completions/%s.fish

PowerShell:
  PS> %s completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> %s completion powershell > %s.ps1
  # and source this file from your PowerShell profile.
`, AppName, AppName, AppName, AppName, AppName,
			AppName, AppName, AppName, AppName, AppName,
			AppName, AppName, AppName),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
}
