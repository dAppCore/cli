// SPDX-Licence-Identifier: EUPL-1.2

// pkgRunProcess over os/exec.
//
// Process creation on Windows is CreateProcess with inheritable handles and a
// single flattened command line — a different model from fork/exec, not the
// same one under other names. os/exec is the maintained implementation of it,
// and CombinedOutput already gives exactly the contract the POSIX path builds
// by hand: one merged stdout+stderr stream, and the tool's own words as the
// error on a non-zero exit.
//
// Resolution goes through exec.LookPath rather than pkgFindExecutable, because
// PATHEXT is what makes "git" find git.exe. Splitting PATH and joining a bare
// name — which is what pkgFindExecutable does, correctly, for POSIX — finds
// nothing here.

package pkgcmd

import (
	"os/exec"

	core "dappco.re/go"
)

func pkgRunProcess(dir, command string, args ...string) core.Result {
	if command == "" {
		return core.Fail(core.NewError("empty command"))
	}

	commandPath, err := exec.LookPath(command)
	if err != nil {
		return core.Fail(core.E("pkg.process", core.Concat("command not found: ", command), err))
	}

	cmd := exec.Command(commandPath, args...)
	cmd.Dir = dir
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
		return core.Fail(core.E("pkg.process",
			core.Sprintf("%s exited with status %d", command, exitErr.ExitCode()), nil))
	}
	return core.Fail(runErr)
}
