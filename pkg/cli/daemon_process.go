package cli

import (
	"context"
	"net"      // Note: AX-6 — daemon HTTP server boundary.
	"net/http" // Note: AX-6 — daemon HTTP server boundary.
	"syscall"  // Note: AX-6 — SIGTERM/SIGINT signal handling.
	"time"

	"dappco.re/go"
)

// DaemonOptions configures a background process helper.
//
//	daemon := cli.NewDaemon(cli.DaemonOptions{
//	    PIDFile:    "/tmp/core.pid",
//	    HealthAddr: "127.0.0.1:8080",
//	})
type DaemonOptions struct {
	// PIDFile stores the current process ID on Start and removes it on Stop.
	PIDFile string

	// HealthAddr binds the HTTP health server.
	// Pass an empty string to disable the server.
	HealthAddr string

	// HealthPath serves the liveness probe endpoint.
	HealthPath string

	// ReadyPath serves the readiness probe endpoint.
	ReadyPath string

	// HealthCheck reports whether the process is healthy.
	// Defaults to true when nil.
	HealthCheck func() bool

	// ReadyCheck reports whether the process is ready to serve traffic.
	// Defaults to HealthCheck when nil, or true when both are nil.
	ReadyCheck func() bool
}

// Daemon manages a PID file and optional HTTP health endpoints.
//
//	daemon := cli.NewDaemon(cli.DaemonOptions{PIDFile: "/tmp/core.pid"})
//	_ = daemon.Start(context.Background())
type Daemon struct {
	opts DaemonOptions

	mu       core.Mutex
	listener net.Listener
	server   *http.Server
	addr     string
	started  bool
}

var (
	processNow   = time.Now
	processSleep = time.Sleep
	processAlive = func(pid int) bool {
		err := syscall.Kill(pid, syscall.Signal(0))
		return err == nil || core.Is(err, syscall.EPERM)
	}
	processSignal = func(pid int, sig syscall.Signal) core.Result {
		if err := syscall.Kill(pid, sig); err != nil {
			return core.Fail(err)
		}
		return core.Ok(nil)
	}
	processPollInterval = 100 * time.Millisecond
	processShutdownWait = 30 * time.Second
)

// NewDaemon creates a daemon helper with sensible defaults.
func NewDaemon(opts DaemonOptions) *Daemon {
	if opts.HealthPath == "" {
		opts.HealthPath = "/health"
	}
	if opts.ReadyPath == "" {
		opts.ReadyPath = "/ready"
	}
	return &Daemon{opts: opts}
}

// Start writes the PID file and starts the health server, if configured.
func (d *Daemon) Start(ctx context.Context) core.Result {
	if ctx == nil {
		ctx = context.Background()
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.started {
		return core.Ok(nil)
	}

	if r := d.writePIDFile(); !r.OK {
		return r
	}

	if d.opts.HealthAddr != "" {
		if r := d.startHealthServer(ctx); !r.OK {
			if cleanup := d.removePIDFile(); !cleanup.OK {
				LogWarn("failed to remove PID file after daemon startup error", "err", cleanup.Error())
			}
			return r
		}
	}

	d.started = true
	return core.Ok(nil)
}

// Stop shuts down the health server and removes the PID file.
func (d *Daemon) Stop(ctx context.Context) core.Result {
	if ctx == nil {
		ctx = context.Background()
	}

	d.mu.Lock()
	server := d.server
	listener := d.listener
	d.server = nil
	d.listener = nil
	d.addr = ""
	d.started = false
	d.mu.Unlock()

	var firstErr error

	if server != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil && !isClosedServerError(err) {
			firstErr = err
		}
	}

	if listener != nil {
		if err := listener.Close(); err != nil && !isListenerClosedError(err) && firstErr == nil {
			firstErr = err
		}
	}

	if r := d.removePIDFile(); !r.OK && firstErr == nil {
		firstErr, _ = r.Value.(error)
	}

	if firstErr != nil {
		return core.Fail(firstErr)
	}
	return core.Ok(nil)
}

// HealthAddr returns the bound health server address, if running.
func (d *Daemon) HealthAddr() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.addr != "" {
		return d.addr
	}
	return d.opts.HealthAddr
}

// StopPIDFile sends SIGTERM to the process identified by pidFile, waits for it
// to exit, escalates to SIGKILL after the timeout, and then removes the file.
//
// If the PID file does not exist, StopPIDFile returns nil.
func StopPIDFile(pidFile string, timeout time.Duration) core.Result {
	if pidFile == "" {
		return core.Ok(nil)
	}
	if timeout <= 0 {
		timeout = processShutdownWait
	}

	rawPIDResult := core.ReadFile(pidFile)
	if !rawPIDResult.OK {
		err, _ := rawPIDResult.Value.(error)
		if core.IsNotExist(err) {
			return core.Ok(nil)
		}
		return rawPIDResult
	}
	rawPID := rawPIDResult.Value.([]byte)

	pidResult := parsePID(core.Trim(string(rawPID)))
	if !pidResult.OK {
		return core.Fail(core.E("StopPIDFile", core.Sprintf("parse pid file %q", pidFile), pidResult.Value.(error)))
	}
	pid := pidResult.Value.(int)

	if r := processSignal(pid, syscall.SIGTERM); !r.OK && !isProcessGone(r.Value.(error)) {
		return r
	}

	deadline := processNow().Add(timeout)
	for processAlive(pid) && processNow().Before(deadline) {
		processSleep(processPollInterval)
	}

	if processAlive(pid) {
		if r := processSignal(pid, syscall.SIGKILL); !r.OK && !isProcessGone(r.Value.(error)) {
			return r
		}

		deadline = processNow().Add(processShutdownWait)
		for processAlive(pid) && processNow().Before(deadline) {
			processSleep(processPollInterval)
		}

		if processAlive(pid) {
			return core.Fail(core.E("StopPIDFile", core.Sprintf("process %d did not exit after SIGKILL", pid), nil))
		}
	}

	return core.Remove(pidFile)
}

func parsePID(raw string) core.Result {
	if raw == "" {
		return core.Fail(core.NewError("empty pid"))
	}
	pidResult := Atoi(raw)
	if !pidResult.OK {
		return pidResult
	}
	pid := pidResult.Value.(int)
	if pid <= 0 {
		return core.Fail(core.E("parsePID", core.Sprintf("invalid pid %d", pid), nil))
	}
	return core.Ok(pid)
}

func isProcessGone(err error) bool {
	return core.Is(err, syscall.ESRCH)
}

func (d *Daemon) writePIDFile() core.Result {
	if d.opts.PIDFile == "" {
		return core.Ok(nil)
	}

	if r := core.MkdirAll(core.PathDir(d.opts.PIDFile), 0o755); !r.OK {
		return r
	}
	return core.WriteFile(d.opts.PIDFile, []byte(core.Sprintf("%d", core.Getpid())+"\n"), 0o644)
}

func (d *Daemon) removePIDFile() core.Result {
	if d.opts.PIDFile == "" {
		return core.Ok(nil)
	}
	r := core.Remove(d.opts.PIDFile)
	if !r.OK {
		err, _ := r.Value.(error)
		if !core.IsNotExist(err) {
			return r
		}
	}
	return core.Ok(nil)
}

func (d *Daemon) startHealthServer(ctx context.Context) core.Result {
	mux := http.NewServeMux()
	healthCheck := d.opts.HealthCheck
	if healthCheck == nil {
		healthCheck = func() bool { return true }
	}
	readyCheck := d.opts.ReadyCheck
	if readyCheck == nil {
		readyCheck = healthCheck
	}

	mux.HandleFunc(d.opts.HealthPath, func(w http.ResponseWriter, r *http.Request) {
		writeProbe(w, healthCheck())
	})
	mux.HandleFunc(d.opts.ReadyPath, func(w http.ResponseWriter, r *http.Request) {
		writeProbe(w, readyCheck())
	})

	listener, err := net.Listen("tcp", d.opts.HealthAddr)
	if err != nil {
		return core.Fail(err)
	}

	server := &http.Server{
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	d.listener = listener
	d.server = server
	d.addr = listener.Addr().String()

	go func() {
		err := server.Serve(listener)
		if err != nil && !isClosedServerError(err) {
			LogWarn("daemon health server stopped unexpectedly", "err", err)
		}
	}()

	return core.Ok(nil)
}

func writeProbe(w http.ResponseWriter, ok bool) {
	if ok {
		w.WriteHeader(http.StatusOK)
		if r := core.WriteString(w, "ok\n"); !r.OK {
			LogWarn("failed to write health response", "err", r.Error())
		}
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	if r := core.WriteString(w, "unhealthy\n"); !r.OK {
		LogWarn("failed to write health response", "err", r.Error())
	}
}

func isClosedServerError(err error) bool {
	return err == nil || err == http.ErrServerClosed
}

func isListenerClosedError(err error) bool {
	return err == nil || core.Is(err, net.ErrClosed)
}
