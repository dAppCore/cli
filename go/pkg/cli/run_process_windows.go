// SPDX-License-Identifier: EUPL-1.2

// runProcessOutput over os/exec.
//
// Process creation on Windows is CreateProcess with inheritable handles and a
// single flattened command line — a different model from fork/exec, not the
// same one under other names. os/exec is the maintained implementation of it,
// and CombinedOutput already provides exactly the contract the POSIX path
// builds by hand: one merged stdout+stderr stream, and the tool's own words as
// the error on a non-zero exit. Re-implementing CreateProcess here would ship
// a second process runtime that nothing in this repository can test.
//
// Command resolution goes through exec.LookPath rather than findExecutable
// because PATHEXT is what makes "gh" find gh.exe. Splitting PATH and joining a
// bare name — which is all findExecutable does, correctly, for POSIX — finds
// nothing on Windows.

package cli

import (
	"context"
	"os/exec"

	"dappco.re/go"
)

func runProcessOutput(ctx context.Context, command string, args ...string) core.Result {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return core.Fail(err)
	}
	if command == "" {
		return core.Fail(core.NewError("empty command"))
	}

	commandPath, err := exec.LookPath(command)
	if err != nil {
		return core.Fail(core.E("cli.process", core.Concat("command not found: ", command), err))
	}

	cmd := exec.CommandContext(ctx, commandPath, args...)
	cmd.Env = core.Environ()

	combined, runErr := cmd.CombinedOutput()
	output := string(combined)
	if runErr == nil {
		return core.Ok(output)
	}
	if output != "" {
		return core.Fail(core.NewError(output))
	}

	var exitErr *exec.ExitError
	if core.As(runErr, &exitErr) {
		return core.Fail(core.E("cli.process",
			core.Sprintf("%s exited with status %d", command, exitErr.ExitCode()), nil))
	}
	return core.Fail(runErr)
}
