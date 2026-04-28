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
	"os"        // Note: signal handling exception; released core lacks OnSignal.
	"os/signal" // Note: signal handling exception; paired with os.Signal.
	"time"

	"dappco.re/go"
	"golang.org/x/sys/unix"
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
//	err := cli.Init(cli.Options{AppName: "core"})
//	if err != nil { panic(err) }
//	defer cli.Shutdown()
func Init(opts Options) error {
	initLock.Mutex.Lock()
	defer initLock.Mutex.Unlock()
	if instance != nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create Core instance with CLI service (registered automatically by core.New)
	c := core.New(
		core.WithOption("name", opts.AppName),
	)
	c.App().Name = opts.AppName
	c.App().Version = opts.Version

	// Register signal service
	signalSvc := &signalService{
		cancel:   cancel,
		sigChan:  make(chan os.Signal, 1),
		stopLock: c.Lock("cli.signal.stop"),
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
		ctx:    ctx,
		cancel: cancel,
	}

	r := c.ServiceStartup(ctx, nil)
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return err
		}
		return nil
	}

	loadLocaleSources(opts.I18nSources...)

	// Attach registered commands AFTER Core startup so i18n is available
	attachRegisteredCommands(c)
	return nil
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
// Returns an error if the command fails.
//
// Example:
//
//	if err := cli.Execute(); err != nil {
//		cli.Warn("command failed:", "err", err)
//	}
func Execute() error {
	mustInit()
	cl := instance.core.Cli()
	if cl == nil {
		return core.E("cli.Execute", "CLI service not available", nil)
	}
	result := cl.Run()
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
	}
	return nil
}

// Run executes the CLI and watches an external context for cancellation.
// If the context is cancelled first, the runtime is shut down and the
// command error is returned if execution failed during shutdown.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//	if err := cli.Run(ctx); err != nil {
//		cli.Error(err.Error())
//	}
func Run(ctx context.Context) error {
	mustInit()
	if ctx == nil {
		ctx = context.Background()
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- Execute()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		Shutdown()
		if err := <-errCh; err != nil {
			return err
		}
		return ctx.Err()
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
	_ = instance.core.ServiceShutdown(context.WithoutCancel(instance.ctx))
}

// --- Signal Srv (internal) ---

type signalService struct {
	cancel   context.CancelFunc
	sigChan  chan os.Signal
	onReload func() error
	stopLock *core.Lock
	stopped  bool
}

func (s *signalService) start(ctx context.Context) core.Result {
	signals := []os.Signal{unix.SIGINT, unix.SIGTERM}
	if s.onReload != nil {
		signals = append(signals, unix.SIGHUP)
	}
	signal.Notify(s.sigChan, signals...)

	go func() {
		for {
			select {
			case sig, ok := <-s.sigChan:
				if !ok {
					return
				}
				switch sig {
				case unix.SIGHUP:
					if s.onReload != nil {
						if err := s.onReload(); err != nil {
							LogError("reload failed", "err", err)
						} else {
							LogInfo("configuration reloaded")
						}
					}
				case unix.SIGINT, unix.SIGTERM:
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
	signal.Stop(s.sigChan)
	close(s.sigChan)
	return core.Ok(nil)
}
