package cli

import (
	"context"
	core "dappco.re/go"
	"time"
)

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
