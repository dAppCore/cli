package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go-process"
	"forge.lthn.ai/core/go-scm/manifest"
)

// AddServiceCommands registers core start/stop/list/restart as top-level commands.
func AddServiceCommands(root *cli.Command) {
	startCmd := cli.NewCommand("start", "Start a project daemon",
		"Reads .core/manifest.yaml and starts the named daemon (or the default).\n"+
			"The daemon runs detached in the background.",
		func(cmd *cli.Command, args []string) error {
			return runStart(args)
		},
	)

	stopCmd := cli.NewCommand("stop", "Stop a project daemon",
		"Stops the named daemon for the current project, or all daemons if no name given.",
		func(cmd *cli.Command, args []string) error {
			return runStop(args)
		},
	)

	listCmd := cli.NewCommand("list", "List running daemons",
		"Shows all running daemons tracked in ~/.core/daemons/.",
		func(cmd *cli.Command, args []string) error {
			return runList()
		},
	)

	restartCmd := cli.NewCommand("restart", "Restart a project daemon",
		"Stops then starts the named daemon.",
		func(cmd *cli.Command, args []string) error {
			if err := runStop(args); err != nil {
				return err
			}
			return runStart(args)
		},
	)

	root.AddCommand(startCmd, stopCmd, listCmd, restartCmd)
}

func runStart(args []string) error {
	m, projectDir, err := findManifest()
	if err != nil {
		return err
	}

	daemonName, spec, err := resolveDaemon(m, args)
	if err != nil {
		return err
	}

	reg := process.DefaultRegistry()

	// Check if already running.
	if _, ok := reg.Get(m.Code, daemonName); ok {
		return fmt.Errorf("%s/%s is already running", m.Code, daemonName)
	}

	// Resolve binary.
	binary := spec.Binary
	if binary == "" {
		return fmt.Errorf("daemon %q has no binary specified", daemonName)
	}

	binPath, err := exec.LookPath(binary)
	if err != nil {
		return fmt.Errorf("binary %q not found in PATH: %w", binary, err)
	}

	// Launch detached.
	cmd := exec.Command(binPath, spec.Args...)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "CORE_DAEMON=1")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", daemonName, err)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	// Wait for health if configured.
	health := spec.Health
	if health != "" && health != "127.0.0.1:0" {
		if process.WaitForHealth(health, 5000) {
			cli.LogInfo(fmt.Sprintf("Started %s/%s (PID %d, health %s)", m.Code, daemonName, pid, health))
		} else {
			cli.LogInfo(fmt.Sprintf("Started %s/%s (PID %d, health not yet ready)", m.Code, daemonName, pid))
		}
	} else {
		cli.LogInfo(fmt.Sprintf("Started %s/%s (PID %d)", m.Code, daemonName, pid))
	}

	// Register in the daemon registry.
	if err := reg.Register(process.DaemonEntry{
		Code:    m.Code,
		Daemon:  daemonName,
		PID:     pid,
		Health:  health,
		Project: projectDir,
		Binary:  binPath,
	}); err != nil {
		cli.LogWarn(fmt.Sprintf("Daemon started but registry failed: %v", err))
	}

	return nil
}

func runStop(args []string) error {
	reg := process.DefaultRegistry()

	m, _, err := findManifest()
	if err != nil {
		return err
	}

	// If a specific daemon name was given, stop only that one.
	if len(args) > 0 {
		return stopDaemon(reg, m.Code, args[0])
	}

	// No args: stop all daemons for this project.
	entries, err := reg.List()
	if err != nil {
		return err
	}

	stopped := 0
	for _, e := range entries {
		if e.Code == m.Code {
			if err := stopDaemon(reg, e.Code, e.Daemon); err != nil {
				cli.LogError(fmt.Sprintf("Failed to stop %s/%s: %v", e.Code, e.Daemon, err))
			} else {
				stopped++
			}
		}
	}

	if stopped == 0 {
		cli.LogInfo("No running daemons for " + m.Code)
	}

	return nil
}

func stopDaemon(reg *process.Registry, code, daemon string) error {
	entry, ok := reg.Get(code, daemon)
	if !ok {
		return fmt.Errorf("%s/%s is not running", code, daemon)
	}

	proc, err := os.FindProcess(entry.PID)
	if err != nil {
		return fmt.Errorf("process %d not found: %w", entry.PID, err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to signal PID %d: %w", entry.PID, err)
	}

	// Wait for process to exit, escalate to SIGKILL after 30s.
	// Poll the process directly via Signal(0) rather than relying on
	// the daemon to self-unregister, which avoids PID reuse issues.
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			// Process is gone.
			_ = reg.Unregister(code, daemon)
			cli.LogInfo(fmt.Sprintf("Stopped %s/%s (PID %d)", code, daemon, entry.PID))
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}

	cli.LogWarn(fmt.Sprintf("%s/%s did not stop within 30s, sending SIGKILL", code, daemon))
	_ = proc.Signal(syscall.SIGKILL)
	_ = reg.Unregister(code, daemon)
	cli.LogInfo(fmt.Sprintf("Killed %s/%s (PID %d)", code, daemon, entry.PID))
	return nil
}

func runList() error {
	reg := process.DefaultRegistry()
	entries, err := reg.List()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No running daemons")
		return nil
	}

	fmt.Printf("%-20s %-12s %-8s %-24s %s\n", "CODE", "DAEMON", "PID", "HEALTH", "PROJECT")
	for _, e := range entries {
		project := e.Project
		if project == "" {
			project = "-"
		}
		fmt.Printf("%-20s %-12s %-8d %-24s %s\n", e.Code, e.Daemon, e.PID, e.Health, project)
	}

	return nil
}

// findManifest walks from cwd up to / looking for .core/manifest.yaml.
func findManifest() (*manifest.Manifest, string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}

	for {
		path := filepath.Join(dir, ".core", "manifest.yaml")
		data, err := os.ReadFile(path)
		if err == nil {
			m, err := manifest.Parse(data)
			if err != nil {
				return nil, "", fmt.Errorf("invalid manifest at %s: %w", path, err)
			}
			return m, dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return nil, "", fmt.Errorf("no .core/manifest.yaml found (checked cwd and parent directories)")
}

// resolveDaemon finds the daemon entry by name or returns the default.
func resolveDaemon(m *manifest.Manifest, args []string) (string, manifest.DaemonSpec, error) {
	if len(args) > 0 {
		name := args[0]
		spec, ok := m.Daemons[name]
		if !ok {
			return "", manifest.DaemonSpec{}, fmt.Errorf("daemon %q not found in manifest (available: %v)", name, daemonNames(m))
		}
		return name, spec, nil
	}

	name, spec, ok := m.DefaultDaemon()
	if !ok {
		return "", manifest.DaemonSpec{}, fmt.Errorf("no default daemon in manifest (use: core start <name>)")
	}
	return name, spec, nil
}

func daemonNames(m *manifest.Manifest) []string {
	var names []string
	for name := range m.Daemons {
		names = append(names, name)
	}
	return names
}
