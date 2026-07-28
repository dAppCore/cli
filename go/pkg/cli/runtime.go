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
	"time"

	"dappco.re/go"
)

var (
	instance *runtime
	initLock = core.New().Lock("cli.runtime.init")
	once     any // Legacy test reset hook; Init gates on initLock + instance.
)

// runtime is the CLI's internal Core runtime.
type runtime struct {
	core   *core.Core
	ctx    context.Context
	cancel context.CancelFunc
}

// Options configures the CLI runtime.
//
// Example:
//
//	opts := cli.Options{
//		AppName: "core",
//		Version: "1.0.0",
//	}
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
//
// Example:
//
//	r := cli.Init(cli.Options{AppName: "core"})
//	if !r.OK { panic(r.Error()) }
//	defer cli.Shutdown()
func Init(opts Options) core.Result {
	initLock.Mutex.Lock()
	defer initLock.Mutex.Unlock()
	if instance != nil {
		return core.Ok(nil)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create Core instance, then stand up the CLI via the package's own
	// service (registers the primitive once) — never core.WithCli() here, which
	// would force the default cli back in and double-register.
	c := core.New(
		core.WithOption("name", opts.AppName),
	)
	if r := Register(c); !r.OK {
		cancel()
		return r
	}
	c.App().Name = opts.AppName
	c.App().Version = opts.Version

	// Register signal service
	signalSvc := &signalService{
		cancel:   cancel,
		sigChan:  make(chan string, 1),
		stopLock: c.Lock("cli.signal.stop"),
	}
	if opts.OnReload != nil {
		signalSvc.onReload = opts.OnReload
	}
	if r := c.Service("signal", core.Service{
		OnStart: func() core.Result {
			return signalSvc.start(ctx)
		},
		OnStop: func() core.Result {
			return signalSvc.stop()
		},
	}); !r.OK {
		cancel()
		return r
	}

	// Register additional services
	for _, svc := range opts.Services {
		if svc.Name != "" {
			if r := c.Service(svc.Name, svc); !r.OK {
				cancel()
				return r
			}
		}
	}

	instance = &runtime{
		core:   c,
		ctx:    ctx,
		cancel: cancel,
	}

	r := c.ServiceStartup(ctx, nil)
	if !r.OK {
		return r
	}

	loadLocaleSources(opts.I18nSources...)

	// Attach registered commands AFTER Core startup so i18n is available
	if r := attachRegisteredCommands(c); !r.OK {
		return r
	}
	return core.Ok(nil)
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

// Execute runs the CLI via core.Cli().Run().
// Returns a Result when the command fails.
//
// Example:
//
//	if r := cli.Execute(); !r.OK {
//		cli.Warn("command failed:", "err", r.Error())
//	}
func Execute() core.Result {
	mustInit()
	cl := instance.core.Cli()
	if cl == nil {
		return core.Fail(core.E("cli.Execute", "CLI service not available", nil))
	}

	// CLI presentation belongs to pkg/cli, not surface-agnostic core/go: when no
	// command is given, render our own catalog (space-separated paths,
	// descriptions resolved via T) instead of core/go's slash-path PrintHelp.
	//
	// A leading help flag means the same thing. Without this, core/go's Cli
	// reads --help as a command name, fails to find one, prints its own
	// slash-path catalog and returns an error — so `core --help` showed a
	// different catalog from bare `core` and exited 1, which is the one
	// invocation every user and every script tries first.
	args := core.FilterArgs(core.Args()[1:])
	if len(args) == 0 || isCatalogRequest(args[0]) {
		PrintHelp()
		return core.Ok(nil) // no command → show help → success (not an error)
	}

	result := cl.Run()
	if !result.OK {
		return result
	}
	return core.Ok(nil)
}

// isCatalogRequest reports whether an argument asks for the command catalog.
//
// Only checked in first position: `core --help` wants the catalog, while
// `core build --help` is asking about `build` and belongs to that command.
//
//	isCatalogRequest("--help") // true
//	isCatalogRequest("build")  // false
func isCatalogRequest(arg string) bool {
	switch arg {
	case "--help", "-help", "-h", "--h":
		return true
	default:
		return false
	}
}

// Run executes the CLI and watches an external context for cancellation.
// If the context is cancelled first, the runtime is shut down and the
// command error is returned if execution failed during shutdown.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//	if r := cli.Run(ctx); !r.OK {
//		cli.Error(r.Error())
//	}
func Run(ctx context.Context) core.Result {
	mustInit()
	if ctx == nil {
		ctx = context.Background()
	}

	resultCh := make(chan core.Result, 1)
	go func() {
		resultCh <- Execute()
	}()

	select {
	case r := <-resultCh:
		return r
	case <-ctx.Done():
		Shutdown()
		if r := <-resultCh; !r.OK {
			return r
		}
		return core.Fail(ctx.Err())
	}
}

// RunWithTimeout returns a shutdown helper that waits for the runtime to stop
// for up to timeout before giving up. It is intended for deferred cleanup.
//
// Example:
//
//	stop := cli.RunWithTimeout(5 * time.Second)
//	defer stop()
func RunWithTimeout(timeout time.Duration) func() {
	return func() {
		if timeout <= 0 {
			Shutdown()
			return
		}

		done := make(chan struct{})
		go func() {
			Shutdown()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(timeout):
			// Give up waiting, but let the shutdown goroutine finish in the background.
		}
	}
}

// Context returns the CLI's root context.
// Cancelled on SIGINT/SIGTERM.
//
// Example:
//
//	if ctx := cli.Context(); ctx != nil {
//		_ = ctx
//	}
func Context() context.Context {
	mustInit()
	return instance.ctx
}

// Shutdown gracefully shuts down the CLI.
//
// Example:
//
//	cli.Shutdown()
func Shutdown() {
	if instance == nil {
		return
	}
	instance.cancel()
	if r := instance.core.ServiceShutdown(context.WithoutCancel(instance.ctx)); !r.OK {
		LogWarn("CLI service shutdown failed", "err", r.Error())
	}
}

// --- Signal Srv (internal) ---

type signalService struct {
	cancel context.CancelFunc
	// sigChan carries signal *names* ("SIGINT", "SIGTERM", "SIGHUP") — the
	// same vocabulary core.Signal's "signal.received" action publishes, and
	// portable in a way golang.org/x/sys/unix is not.
	sigChan  chan string
	onReload func() error
	stopLock *core.Lock
	stopped  bool
}

func (s *signalService) start(ctx context.Context) core.Result {
	go func() {
		for {
			select {
			case sig, ok := <-s.sigChan:
				if !ok {
					return
				}
				switch sig {
				case "SIGHUP":
					if s.onReload != nil {
						if err := s.onReload(); err != nil {
							LogError("reload failed", "err", err)
						} else {
							LogInfo("configuration reloaded")
						}
					}
				case "SIGINT", "SIGTERM":
					s.cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return core.Ok(nil)
}

func (s *signalService) stop() core.Result {
	s.stopLock.Mutex.Lock()
	defer s.stopLock.Mutex.Unlock()
	if s.stopped {
		return core.Ok(nil)
	}
	s.stopped = true
	close(s.sigChan)
	return core.Ok(nil)
}
