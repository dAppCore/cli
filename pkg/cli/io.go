package cli

import (
	"io" // Note: AX-6 — io.Reader/io.Writer is the public stdin/stdout/stderr interception contract.
	"os" // Note: AX-6 — os.Stdin/os.Stdout are the structural defaults intercepted by SetStdin/SetStdout.

	"dappco.re/go/core"
)

type Reader = io.Reader
type Writer = io.Writer

var (
	stdin Reader = os.Stdin

	stdoutOverride Writer
	stderrOverride Writer

	ioMu core.RWMutex
)

// SetStdin overrides the default stdin reader for testing.
// Pass nil to restore the real os.Stdin reader.
func SetStdin(r Reader) {
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
func SetStdout(w Writer) {
	ioMu.Lock()
	defer ioMu.Unlock()
	stdoutOverride = w
}

// SetStderr overrides the default stderr writer.
// Pass nil to restore writes to os.Stderr.
func SetStderr(w Writer) {
	ioMu.Lock()
	defer ioMu.Unlock()
	stderrOverride = w
}

func stdinReader() Reader {
	ioMu.RLock()
	defer ioMu.RUnlock()
	return stdin
}

func stdoutWriter() Writer {
	ioMu.RLock()
	defer ioMu.RUnlock()
	if stdoutOverride != nil {
		return stdoutOverride
	}
	return os.Stdout
}

func stderrWriter() Writer {
	ioMu.RLock()
	defer ioMu.RUnlock()
	if stderrOverride != nil {
		return stderrOverride
	}
	return os.Stderr
}

func writeString(w Writer, s string) {
	if w == nil {
		return
	}
	_, _ = w.Write([]byte(s))
}

func isEOF(err error) bool {
	return core.Is(err, io.EOF)
}

func writerFileDescriptor(w Writer) (int, bool) {
	file, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return 0, false
	}
	return int(file.Fd()), true
}
