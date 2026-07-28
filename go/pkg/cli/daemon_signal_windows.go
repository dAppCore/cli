// SPDX-License-Identifier: EUPL-1.2

// Liveness and signalling for StopPIDFile, over the Win32 process API.
//
// Windows has no signals, so the two things the POSIX path gets from Kill are
// obtained separately: liveness from waiting on the process handle for zero
// milliseconds, and termination from TerminateProcess.
//
// Two honest differences from the POSIX behaviour, rather than a pretence of
// parity:
//
//   - SIGTERM is not graceful here. Asking a process to shut down cleanly on
//     Windows means a console control event, which only reaches processes that
//     share the caller's console — not the detached daemon this file exists to
//     stop. So both SIGTERM and SIGKILL terminate. StopPIDFile's escalation
//     from one to the other still runs; it simply has nothing gentler to try.
//   - Liveness is decided by WaitForSingleObject, not GetExitCodeProcess. The
//     latter reports STILL_ACTIVE (259) for a running process, which is
//     indistinguishable from a process that genuinely exited with code 259.
//
// A pid that cannot be opened is reported as ESRCH so that isProcessGone keeps
// working unchanged: the caller already speaks that vocabulary, and mapping to
// it here is cheaper than teaching every call site a second one.

package cli

import (
	"syscall"

	"dappco.re/go"
)

// processAccess is what both operations need: SYNCHRONIZE to wait on the
// handle, PROCESS_TERMINATE to end it, PROCESS_QUERY_INFORMATION so a denied
// pid is distinguishable from an absent one.
const processAccess = syscall.SYNCHRONIZE | syscall.PROCESS_TERMINATE | syscall.PROCESS_QUERY_INFORMATION

var processAlive = func(pid int) bool {
	if pid <= 0 {
		return false
	}
	handle, err := syscall.OpenProcess(processAccess, false, uint32(pid))
	if err != nil {
		// ERROR_ACCESS_DENIED means the process exists and belongs to somebody
		// else — the same answer EPERM gives on POSIX.
		return core.Is(err, syscall.ERROR_ACCESS_DENIED)
	}
	defer func() { _ = syscall.CloseHandle(handle) }()

	state, err := syscall.WaitForSingleObject(handle, 0)
	if err != nil {
		return false
	}
	return state == uint32(syscall.WAIT_TIMEOUT)
}

var processSignal = func(pid int, sig syscall.Signal) core.Result {
	if pid <= 0 {
		return core.Fail(syscall.ESRCH)
	}
	if sig != syscall.SIGTERM && sig != syscall.SIGKILL {
		return core.Fail(core.E("cli.processSignal",
			core.Sprintf("Windows has no signals; only SIGTERM and SIGKILL map to terminating a process, got %v", sig), nil))
	}

	handle, err := syscall.OpenProcess(processAccess, false, uint32(pid))
	if err != nil {
		if core.Is(err, syscall.ERROR_ACCESS_DENIED) {
			return core.Fail(err)
		}
		// Anything else at open time means there is no such process to signal.
		return core.Fail(syscall.ESRCH)
	}
	defer func() { _ = syscall.CloseHandle(handle) }()

	if err := syscall.TerminateProcess(handle, 1); err != nil {
		return core.Fail(err)
	}
	return core.Ok(nil)
}
