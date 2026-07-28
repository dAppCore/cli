// SPDX-License-Identifier: EUPL-1.2

//go:build unix

// Liveness and signalling for StopPIDFile, over POSIX signals.
//
// syscall.Kill with signal 0 is the standard "does this process exist, and may
// I signal it" probe: it performs the permission check and delivers nothing.
// EPERM therefore means the process is alive and owned by somebody else, which
// still answers the question being asked.

package cli

import (
	"syscall"

	"dappco.re/go"
)

var processAlive = func(pid int) bool {
	err := syscall.Kill(pid, syscall.Signal(0))
	return err == nil || core.Is(err, syscall.EPERM)
}

var processSignal = func(pid int, sig syscall.Signal) core.Result {
	if err := syscall.Kill(pid, sig); err != nil {
		return core.Fail(err)
	}
	return core.Ok(nil)
}
