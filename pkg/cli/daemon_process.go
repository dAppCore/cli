package cli

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// DaemonOptions configures a background process helper.
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
type Daemon struct {
	opts DaemonOptions

	mu       sync.Mutex
	listener net.Listener
	server   *http.Server
	addr     string
	started  bool
}

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
func (d *Daemon) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.started {
		return nil
	}

	if err := d.writePIDFile(); err != nil {
		return err
	}

	if d.opts.HealthAddr != "" {
		if err := d.startHealthServer(ctx); err != nil {
			_ = d.removePIDFile()
			return err
		}
	}

	d.started = true
	return nil
}

// Stop shuts down the health server and removes the PID file.
func (d *Daemon) Stop(ctx context.Context) error {
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

	if err := d.removePIDFile(); err != nil && firstErr == nil {
		firstErr = err
	}

	return firstErr
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

func (d *Daemon) writePIDFile() error {
	if d.opts.PIDFile == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(d.opts.PIDFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(d.opts.PIDFile, []byte(strconv.Itoa(os.Getpid())+"\n"), 0o644)
}

func (d *Daemon) removePIDFile() error {
	if d.opts.PIDFile == "" {
		return nil
	}
	if err := os.Remove(d.opts.PIDFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (d *Daemon) startHealthServer(ctx context.Context) error {
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
		return err
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
			_ = err
		}
	}()

	return nil
}

func writeProbe(w http.ResponseWriter, ok bool) {
	if ok {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok\n")
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	_, _ = io.WriteString(w, "unhealthy\n")
}

func isClosedServerError(err error) bool {
	return err == nil || err == http.ErrServerClosed
}

func isListenerClosedError(err error) bool {
	return err == nil || errors.Is(err, net.ErrClosed)
}
