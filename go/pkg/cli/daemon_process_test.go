package cli

import (
	"context"
	"net"
	"net/http"
	"syscall"
	"time"

	core "dappco.re/go"
)

// withProcessMocks swaps the package signal/alive/clock indirection vars
// for the duration of fn and restores them afterwards, so StopPIDFile can
// be driven deterministically without real OS processes.
func withProcessMocks(t *core.T, alive func(int) bool, signal func(int, syscall.Signal) core.Result, fn func()) {
	t.Helper()
	origAlive, origSignal, origSleep, origNow := processAlive, processSignal, processSleep, processNow
	processAlive = alive
	processSignal = signal
	processSleep = func(time.Duration) {}
	now := time.Now()
	processNow = func() time.Time { now = now.Add(time.Second); return now }
	t.Cleanup(func() {
		processAlive, processSignal, processSleep, processNow = origAlive, origSignal, origSleep, origNow
	})
	fn()
}

func TestDaemonProcess_NewDaemon_Good(t *core.T) {
	d := NewDaemon(DaemonOptions{PIDFile: core.Path(t.TempDir(), "daemon.pid")})

	core.AssertNotNil(t, d)
	core.AssertEqual(t, "/health", d.opts.HealthPath)
}

func TestDaemonProcess_NewDaemon_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{})

	core.AssertEqual(t, "", d.opts.PIDFile)
	core.AssertEqual(t, "/ready", d.opts.ReadyPath)
}

func TestDaemonProcess_NewDaemon_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthPath: "/h", ReadyPath: "/r"})

	core.AssertEqual(t, "/h", d.opts.HealthPath)
	core.AssertEqual(t, "/r", d.opts.ReadyPath)
}

func TestDaemonProcess_Daemon_Start_Good(t *core.T) {
	pid := core.Path(t.TempDir(), "daemon.pid")
	d := NewDaemon(DaemonOptions{PIDFile: pid})

	core.AssertNoError(t, cliResultError(d.Start(context.Background())))
	core.AssertNoError(t, cliResultError(d.Stop(context.Background())))
}

func TestDaemonProcess_Daemon_Start_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{PIDFile: core.Path(t.TempDir(), "missing", "daemon.pid")})

	core.AssertNoError(t, cliResultError(d.Start(nil)))
	core.AssertNoError(t, cliResultError(d.Stop(nil)))
}

func TestDaemonProcess_Daemon_Start_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0"})

	core.AssertNoError(t, cliResultError(d.Start(context.Background())))
	core.AssertNoError(t, cliResultError(d.Stop(context.Background())))
}

func TestDaemonProcess_Daemon_Stop_Good(t *core.T) {
	d := NewDaemon(DaemonOptions{PIDFile: core.Path(t.TempDir(), "daemon.pid")})
	core.RequireNoError(t, cliResultError(d.Start(context.Background())))

	core.AssertNoError(t, cliResultError(d.Stop(context.Background())))
	core.AssertFalse(t, d.started)
}

func TestDaemonProcess_Daemon_Stop_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{})

	core.AssertNoError(t, cliResultError(d.Stop(nil)))
	core.AssertFalse(t, d.started)
}

func TestDaemonProcess_Daemon_Stop_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0"})
	core.RequireNoError(t, cliResultError(d.Start(context.Background())))

	core.AssertNoError(t, cliResultError(d.Stop(nil)))
	core.AssertEqual(t, "", d.addr)
}

func TestDaemonProcess_Daemon_HealthAddr_Good(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0"})
	core.RequireNoError(t, cliResultError(d.Start(context.Background())))
	defer d.Stop(context.Background())

	core.AssertNotEmpty(t, d.HealthAddr())
}

func TestDaemonProcess_Daemon_HealthAddr_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{})

	core.AssertEqual(t, "", d.HealthAddr())
	core.AssertFalse(t, d.started)
}

func TestDaemonProcess_Daemon_HealthAddr_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:9999"})

	core.AssertEqual(t, "127.0.0.1:9999", d.HealthAddr())
	core.AssertFalse(t, d.started)
}

func TestDaemonProcess_StopPIDFile_Good(t *core.T) {
	err := cliResultError(StopPIDFile(core.Path(t.TempDir(), "missing.pid"), time.Millisecond))

	core.AssertNoError(t, err)
	core.AssertNil(t, err)
}

func TestDaemonProcess_StopPIDFile_Bad(t *core.T) {
	err := cliResultError(StopPIDFile("", time.Millisecond))

	core.AssertNoError(t, err)
	core.AssertNil(t, err)
}

func TestDaemonProcess_StopPIDFile_Ugly(t *core.T) {
	path := core.Path(t.TempDir(), "bad.pid")
	core.RequireTrue(t, core.WriteFile(path, []byte("not-a-pid"), 0o644).OK)

	err := cliResultError(StopPIDFile(path, time.Millisecond))
	core.AssertError(t, err)
}

// StopPIDFile sends SIGTERM, observes the process exit, then removes the file.
func TestDaemonProcess_StopPIDFile_TermExits(t *core.T) {
	path := core.Path(t.TempDir(), "live.pid")
	core.RequireTrue(t, core.WriteFile(path, []byte("4242\n"), 0o644).OK)

	var sentTerm bool
	withProcessMocks(t,
		func(int) bool { return false },
		func(_ int, sig syscall.Signal) core.Result {
			if sig == syscall.SIGTERM {
				sentTerm = true
			}
			return core.Ok(nil)
		},
		func() {
			core.AssertNoError(t, cliResultError(StopPIDFile(path, time.Second)))
		},
	)

	core.AssertTrue(t, sentTerm)
	core.AssertTrue(t, pidFileRemoved(path))
}

// StopPIDFile escalates to SIGKILL when the process survives SIGTERM.
func TestDaemonProcess_StopPIDFile_KillEscalates(t *core.T) {
	path := core.Path(t.TempDir(), "stubborn.pid")
	core.RequireTrue(t, core.WriteFile(path, []byte("4243\n"), 0o644).OK)

	calls := 0
	var sentKill bool
	withProcessMocks(t,
		// Alive through the SIGTERM wait loop, gone after SIGKILL wait begins.
		func(int) bool { calls++; return calls <= 3 },
		func(_ int, sig syscall.Signal) core.Result {
			if sig == syscall.SIGKILL {
				sentKill = true
			}
			return core.Ok(nil)
		},
		func() {
			core.AssertNoError(t, cliResultError(StopPIDFile(path, time.Second)))
		},
	)

	core.AssertTrue(t, sentKill)
}

// StopPIDFile treats an ESRCH signal error as "already gone" and proceeds.
func TestDaemonProcess_StopPIDFile_AlreadyGone(t *core.T) {
	path := core.Path(t.TempDir(), "gone.pid")
	core.RequireTrue(t, core.WriteFile(path, []byte("4244\n"), 0o644).OK)

	withProcessMocks(t,
		func(int) bool { return false },
		func(int, syscall.Signal) core.Result { return core.Fail(syscall.ESRCH) },
		func() {
			core.AssertNoError(t, cliResultError(StopPIDFile(path, time.Second)))
		},
	)

	core.AssertTrue(t, pidFileRemoved(path))
}

// Stop tolerates a PID file that was already removed out from under it,
// exercising removePIDFile's IsNotExist branch.
func TestDaemonProcess_Daemon_Stop_PIDAlreadyGone(t *core.T) {
	pid := core.Path(t.TempDir(), "vanish.pid")
	d := NewDaemon(DaemonOptions{PIDFile: pid})
	core.RequireNoError(t, cliResultError(d.Start(context.Background())))
	core.RequireTrue(t, core.Remove(pid).OK)

	core.AssertNoError(t, cliResultError(d.Stop(context.Background())))
}

// pidFileRemoved reports whether path no longer exists on disk.
func pidFileRemoved(path string) bool {
	r := core.Stat(path)
	if r.OK {
		return false
	}
	err, _ := r.Value.(error)
	return core.IsNotExist(err)
}

// StopPIDFile fails when the process survives even SIGKILL.
func TestDaemonProcess_StopPIDFile_KillFails(t *core.T) {
	path := core.Path(t.TempDir(), "immortal.pid")
	core.RequireTrue(t, core.WriteFile(path, []byte("4245\n"), 0o644).OK)

	withProcessMocks(t,
		func(int) bool { return true }, // never dies
		func(int, syscall.Signal) core.Result { return core.Ok(nil) },
		func() {
			err := cliResultError(StopPIDFile(path, time.Second))
			core.AssertError(t, err)
			core.AssertContains(t, err.Error(), "did not exit")
		},
	)
}

// StopPIDFile surfaces a non-ESRCH signal error from the initial SIGTERM.
func TestDaemonProcess_StopPIDFile_SignalError(t *core.T) {
	path := core.Path(t.TempDir(), "perm.pid")
	core.RequireTrue(t, core.WriteFile(path, []byte("4246\n"), 0o644).OK)

	withProcessMocks(t,
		func(int) bool { return true },
		func(int, syscall.Signal) core.Result { return core.Fail(syscall.EPERM) },
		func() {
			err := cliResultError(StopPIDFile(path, time.Second))
			core.AssertError(t, err)
		},
	)
}

// parsePID rejects empty, non-numeric, and non-positive PIDs.
func TestDaemonProcess_parsePID_Good(t *core.T) {
	r := parsePID("123")

	core.AssertTrue(t, r.OK)
	core.AssertEqual(t, 123, r.Value.(int))
}

func TestDaemonProcess_parsePID_Bad(t *core.T) {
	r := parsePID("")

	core.AssertFalse(t, r.OK)
}

func TestDaemonProcess_parsePID_Ugly(t *core.T) {
	r := parsePID("0")

	core.AssertFalse(t, r.OK)
}

// isProcessGone recognises ESRCH and rejects unrelated errors.
func TestDaemonProcess_isProcessGone_Good(t *core.T) {
	core.AssertTrue(t, isProcessGone(syscall.ESRCH))
}

func TestDaemonProcess_isProcessGone_Bad(t *core.T) {
	core.AssertFalse(t, isProcessGone(syscall.EPERM))
}

func TestDaemonProcess_isProcessGone_Ugly(t *core.T) {
	core.AssertFalse(t, isProcessGone(nil))
}

// isListenerClosedError recognises nil and net.ErrClosed as benign.
func TestDaemonProcess_isListenerClosedError_Good(t *core.T) {
	core.AssertTrue(t, isListenerClosedError(net.ErrClosed))
}

func TestDaemonProcess_isListenerClosedError_Bad(t *core.T) {
	core.AssertFalse(t, isListenerClosedError(core.NewError("boom")))
}

func TestDaemonProcess_isListenerClosedError_Ugly(t *core.T) {
	core.AssertTrue(t, isListenerClosedError(nil))
}

// writeProbe is exercised end-to-end against the live health server.
func TestDaemonProcess_writeProbe_Healthy(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0", HealthCheck: func() bool { return true }})
	core.RequireNoError(t, cliResultError(d.Start(context.Background())))
	defer d.Stop(context.Background())

	resp, err := http.Get("http://" + d.HealthAddr() + "/health")
	core.RequireNoError(t, err)
	defer resp.Body.Close()

	core.AssertEqual(t, http.StatusOK, resp.StatusCode)
}

func TestDaemonProcess_writeProbe_Unhealthy(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0", HealthCheck: func() bool { return false }})
	core.RequireNoError(t, cliResultError(d.Start(context.Background())))
	defer d.Stop(context.Background())

	resp, err := http.Get("http://" + d.HealthAddr() + "/ready")
	core.RequireNoError(t, err)
	defer resp.Body.Close()

	core.AssertEqual(t, http.StatusServiceUnavailable, resp.StatusCode)
}
