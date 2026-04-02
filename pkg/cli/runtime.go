// Package cli provides the CLI runtime and utilities.
//
// The CLI uses the Core framework for its own runtime. Usage is simple:
//
//	cli.Init(cli.Options{AppName: "core"})
//	defer cli.Shutdown()
//
//	cli.Success("Done!")
//	cli.Error("Failed")
//	if cli.Confirm("Proceed?") { ... }
//
//	// When you need the Core instance
//	c := cli.Core()
package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"dappco.re/go/core"
	"github.com/spf13/cobra"
)

var (
	instance *runtime
	once     sync.Once
)

// runtime is the CLI's internal Core runtime.
type runtime struct {
	core   *core.Core
	root   *cobra.Command
	ctx    context.Context
	cancel context.CancelFunc
}

// Options configures the CLI runtime.
type Options struct {
	AppName     string
	Version     string
	Services    []core.Service // Additional services to register
	I18nSources []LocaleSource // Additional i18n translation sources

	// OnReload is called when SIGHUP is received (daemon mode).
	// Use for configuration reloading. Leave nil to ignore SIGHUP.
	OnReload func() error
}

// Init initialises the global CLI runtime.
// Call this once at startup (typically in main.go or cmd.Execute).
func Init(opts Options) error {
	var initErr error
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())

		// Create root command
		rootCmd := &cobra.Command{
			Use:           opts.AppName,
			Version:       opts.Version,
			SilenceErrors: true,
			SilenceUsage:  true,
		}

		// Create Core with app identity
		c := core.New(core.Options{
			{Key: "name", Value: opts.AppName},
		})
		c.App().Version = opts.Version
		c.App().Runtime = rootCmd

		// Register signal service
		signalSvc := &signalService{
			cancel:  cancel,
			sigChan: make(chan os.Signal, 1),
		}
		if opts.OnReload != nil {
			signalSvc.onReload = opts.OnReload
		}
		c.Service("signal", core.Service{
			OnStart: func() core.Result {
				return signalSvc.start(ctx)
			},
			OnStop: func() core.Result {
				return signalSvc.stop()
			},
		})

		// Register additional services
		for _, svc := range opts.Services {
			if svc.Name != "" {
				c.Service(svc.Name, svc)
			}
		}

		instance = &runtime{
			core:   c,
			root:   rootCmd,
			ctx:    ctx,
			cancel: cancel,
		}

		r := c.ServiceStartup(ctx, nil)
		if !r.OK {
			if err, ok := r.Value.(error); ok {
				initErr = err
			}
			return
		}

		loadLocaleSources(opts.I18nSources...)

		// Attach registered commands AFTER Core startup so i18n is available
		attachRegisteredCommands(rootCmd)
	})
	return initErr
}

func mustInit() {
	if instance == nil {
		panic("cli not initialised - call cli.Init() first")
	}
}

// --- Core Access ---

// Core returns the CLI's framework Core instance.
func Core() *core.Core {
	mustInit()
	return instance.core
}

// RootCmd returns the CLI's root cobra command.
func RootCmd() *cobra.Command {
	mustInit()
	return instance.root
}

// Execute runs the CLI root command.
// Returns an error if the command fails.
func Execute() error {
	mustInit()
	return instance.root.Execute()
}

// Context returns the CLI's root context.
// Cancelled on SIGINT/SIGTERM.
func Context() context.Context {
	mustInit()
	return instance.ctx
}

// Shutdown gracefully shuts down the CLI.
func Shutdown() {
	if instance == nil {
		return
	}
	instance.cancel()
	_ = instance.core.ServiceShutdown(instance.ctx)
}

// --- Signal Srv (internal) ---

type signalService struct {
	cancel       context.CancelFunc
	sigChan      chan os.Signal
	onReload     func() error
	shutdownOnce sync.Once
}

func (s *signalService) start(ctx context.Context) core.Result {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	if s.onReload != nil {
		signals = append(signals, syscall.SIGHUP)
	}
	signal.Notify(s.sigChan, signals...)

	go func() {
		for {
			select {
			case sig := <-s.sigChan:
				switch sig {
				case syscall.SIGHUP:
					if s.onReload != nil {
						if err := s.onReload(); err != nil {
							LogError("reload failed", "err", err)
						} else {
							LogInfo("configuration reloaded")
						}
					}
				case syscall.SIGINT, syscall.SIGTERM:
					s.cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return core.Result{OK: true}
}

func (s *signalService) stop() core.Result {
	s.shutdownOnce.Do(func() {
		signal.Stop(s.sigChan)
		close(s.sigChan)
	})
	return core.Result{OK: true}
}
