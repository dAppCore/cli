// SPDX-License-Identifier: EUPL-1.2

//go:build unix

// runProcessOutput over POSIX fork/exec.
//
// The child's stdout and stderr share one pipe on purpose: callers read the
// combined stream, and on a non-zero exit they surface that text as the error
// itself rather than a generic "exit status N". Keep that contract — the
// Windows implementation reproduces it, and gh/git callers depend on the
// message being the tool's own words.

package cli

import (
	"context"
	"syscall"

	"dappco.re/go"
)

func runProcessOutput(ctx context.Context, command string, args ...string) core.Result {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return core.Fail(err)
	}
	commandResult := findExecutable(command)
	if !commandResult.OK {
		return commandResult
	}
	commandPath := commandResult.Value.(string)

	var pipe [2]int
	if err := syscall.Pipe(pipe[:]); err != nil {
		return core.Fail(err)
	}
	readFD, writeFD := pipe[0], pipe[1]
	defer syscall.Close(readFD)

	argv := append([]string{commandPath}, args...)
	pid, err := syscall.ForkExec(commandPath, argv, &syscall.ProcAttr{
		Env:   core.Environ(),
		Files: []uintptr{0, uintptr(writeFD), uintptr(writeFD)},
	})
	syscall.Close(writeFD)
	if err != nil {
		return core.Fail(err)
	}

	out := core.NewBuilder()
	buf := make([]byte, 4096)
	for {
		n, readErr := syscall.Read(readFD, buf)
		if n > 0 {
			out.WriteString(string(buf[:n]))
		}
		if readErr != nil {
			if readErr == syscall.EINTR {
				continue
			}
			break
		}
		if n == 0 {
			break
		}
	}

	var status syscall.WaitStatus
	if _, err := syscall.Wait4(pid, &status, 0, nil); err != nil {
		return core.Fail(err)
	}
	output := out.String()
	if status.ExitStatus() == 0 {
		return core.Ok(output)
	}
	if output != "" {
		return core.Fail(core.NewError(output))
	}
	return core.Fail(core.E("cli.process", core.Sprintf("%s exited with status %d", command, status.ExitStatus()), nil))
}
