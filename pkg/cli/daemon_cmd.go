package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"forge.lthn.ai/core/go-process"
)

// DaemonCommandConfig configures the generic daemon CLI command group.
type DaemonCommandConfig struct {
	// Name is the command group name (default: "daemon").
	Name string

	// Description is the short description for the command group.
	Description string

	// RunForeground is called when the daemon runs in foreground mode.
	// Receives context (cancelled on SIGINT/SIGTERM) and the started Daemon.
	// If nil, the run command just blocks until signal.
	RunForeground func(ctx context.Context, daemon *process.Daemon) error

	// PIDFile default path.
	PIDFile string

	// HealthAddr default address.
	HealthAddr string

	// ExtraStartArgs returns additional CLI args to pass when re-execing
	// the binary as a background daemon.
	ExtraStartArgs func() []string

	// Flags registers custom persistent flags on the daemon command group.
	Flags func(cmd *Command)
}

// AddDaemonCommand registers start/stop/status/run subcommands on root.
func AddDaemonCommand(root *Command, cfg DaemonCommandConfig) {
	if cfg.Name == "" {
		cfg.Name = "daemon"
	}
	if cfg.Description == "" {
		cfg.Description = "Manage the background daemon"
	}

	daemonCmd := NewGroup(
		cfg.Name,
		cfg.Description,
		fmt.Sprintf("Manage the background daemon process.\n\n"+
			"Subcommands:\n"+
			"  start   - Start the daemon in the background\n"+
			"  stop    - Stop the running daemon\n"+
			"  status  - Show daemon status\n"+
			"  run     - Run in foreground (for development/debugging)"),
	)

	PersistentStringFlag(daemonCmd, &cfg.HealthAddr, "health-addr", "", cfg.HealthAddr,
		"Health check endpoint address (empty to disable)")
	PersistentStringFlag(daemonCmd, &cfg.PIDFile, "pid-file", "", cfg.PIDFile,
		"PID file path (empty to disable)")

	if cfg.Flags != nil {
		cfg.Flags(daemonCmd)
	}

	startCmd := NewCommand("start", "Start the daemon in the background",
		"Re-executes the binary as a background daemon process.\n"+
			"The daemon PID is written to the PID file for later management.",
		func(cmd *Command, args []string) error {
			return daemonRunStart(cfg)
		},
	)

	stopCmd := NewCommand("stop", "Stop the running daemon",
		"Sends SIGTERM to the daemon process identified by the PID file.\n"+
			"Waits for graceful shutdown before returning.",
		func(cmd *Command, args []string) error {
			return daemonRunStop(cfg)
		},
	)

	statusCmd := NewCommand("status", "Show daemon status",
		"Checks if the daemon is running and queries its health endpoint.",
		func(cmd *Command, args []string) error {
			return daemonRunStatus(cfg)
		},
	)

	runCmd := NewCommand("run", "Run the daemon in the foreground",
		"Runs the daemon in the current terminal (blocks until SIGINT/SIGTERM).\n"+
			"Useful for development, debugging, or running under a process manager.",
		func(cmd *Command, args []string) error {
			return daemonRunForeground(cfg)
		},
	)

	daemonCmd.AddCommand(startCmd, stopCmd, statusCmd, runCmd)
	root.AddCommand(daemonCmd)
}

func daemonRunStart(cfg DaemonCommandConfig) error {
	if pid, running := process.ReadPID(cfg.PIDFile); running {
		return fmt.Errorf("daemon already running (PID %d)", pid)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find executable: %w", err)
	}

	args := []string{cfg.Name, "run",
		"--health-addr", cfg.HealthAddr,
		"--pid-file", cfg.PIDFile,
	}

	if cfg.ExtraStartArgs != nil {
		args = append(args, cfg.ExtraStartArgs()...)
	}

	cmd := exec.Command(exePath, args...)
	cmd.Env = append(os.Environ(), "CORE_DAEMON=1")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	if cfg.HealthAddr != "" {
		if process.WaitForHealth(cfg.HealthAddr, 5_000) {
			LogInfo(fmt.Sprintf("Daemon started (PID %d, health %s)", pid, cfg.HealthAddr))
		} else {
			LogInfo(fmt.Sprintf("Daemon started (PID %d, health not yet ready)", pid))
		}
	} else {
		LogInfo(fmt.Sprintf("Daemon started (PID %d)", pid))
	}

	return nil
}

func daemonRunStop(cfg DaemonCommandConfig) error {
	pid, running := process.ReadPID(cfg.PIDFile)
	if !running {
		LogInfo("Daemon is not running")
		return nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	LogInfo(fmt.Sprintf("Stopping daemon (PID %d)", pid))
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to PID %d: %w", pid, err)
	}

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if _, still := process.ReadPID(cfg.PIDFile); !still {
			LogInfo("Daemon stopped")
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}

	LogWarn("Daemon did not stop within 30s, sending SIGKILL")
	_ = proc.Signal(syscall.SIGKILL)
	_ = os.Remove(cfg.PIDFile)
	LogInfo("Daemon killed")
	return nil
}

func daemonRunStatus(cfg DaemonCommandConfig) error {
	pid, running := process.ReadPID(cfg.PIDFile)
	if !running {
		fmt.Println("Daemon is not running")
		return nil
	}

	fmt.Printf("Daemon is running (PID %d)\n", pid)

	if cfg.HealthAddr != "" {
		healthURL := fmt.Sprintf("http://%s/health", cfg.HealthAddr)
		resp, err := http.Get(healthURL)
		if err != nil {
			fmt.Printf("Health: unreachable (%v)\n", err)
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Health: ok")
		} else {
			fmt.Printf("Health: unhealthy (HTTP %d)\n", resp.StatusCode)
		}

		readyURL := fmt.Sprintf("http://%s/ready", cfg.HealthAddr)
		resp2, err := http.Get(readyURL)
		if err == nil {
			defer resp2.Body.Close()
			if resp2.StatusCode == http.StatusOK {
				fmt.Println("Ready:  yes")
			} else {
				fmt.Println("Ready:  no")
			}
		}
	}

	return nil
}

func daemonRunForeground(cfg DaemonCommandConfig) error {
	os.Setenv("CORE_DAEMON", "1")

	daemon := process.NewDaemon(process.DaemonOptions{
		PIDFile:         cfg.PIDFile,
		HealthAddr:      cfg.HealthAddr,
		ShutdownTimeout: 30 * time.Second,
	})

	if err := daemon.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	daemon.SetReady(true)

	ctx := Context()

	if cfg.RunForeground != nil {
		svcErr := make(chan error, 1)
		go func() {
			svcErr <- cfg.RunForeground(ctx, daemon)
		}()

		select {
		case <-ctx.Done():
			LogInfo("Shutting down daemon")
		case err := <-svcErr:
			if err != nil {
				LogError(fmt.Sprintf("Service exited with error: %v", err))
			}
		}
	} else {
		<-ctx.Done()
	}

	return daemon.Stop()
}
