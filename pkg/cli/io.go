package cli

import (
	"io" // Note: AX-6 — io.Reader/io.Writer is the public stdin/stdout/stderr interception contract.
	"os" // Note: AX-6 — os.Stdin/os.Stdout are the structural defaults intercepted by SetStdin/SetStdout.

	"dappco.re/go/core"
)

var (
	stdin io.Reader = os.Stdin

	stdoutOverride io.Writer
	stderrOverride io.Writer

	ioMu core.RWMutex
)

// SetStdin overrides the default stdin reader for testing.
// Pass nil to restore the real os.Stdin reader.
func SetStdin(r io.Reader) {
	ioMu.Lock()
	defer ioMu.Unlock()
	if r == nil {
		stdin = os.Stdin
		return
	}
	stdin = r
}

// SetStdout overrides the default stdout writer.
// Pass nil to restore writes to os.Stdout.
func SetStdout(w io.Writer) {
	ioMu.Lock()
	defer ioMu.Unlock()
	stdoutOverride = w
}

// SetStderr overrides the default stderr writer.
// Pass nil to restore writes to os.Stderr.
func SetStderr(w io.Writer) {
	ioMu.Lock()
	defer ioMu.Unlock()
	stderrOverride = w
}

func stdinReader() io.Reader {
	ioMu.RLock()
	defer ioMu.RUnlock()
	return stdin
}

func stdoutWriter() io.Writer {
	ioMu.RLock()
	defer ioMu.RUnlock()
	if stdoutOverride != nil {
		return stdoutOverride
	}
	return os.Stdout
}

func stderrWriter() io.Writer {
	ioMu.RLock()
	defer ioMu.RUnlock()
	if stderrOverride != nil {
		return stderrOverride
	}
	return os.Stderr
}
