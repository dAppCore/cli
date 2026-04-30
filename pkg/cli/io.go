package cli

import "dappco.re/go"

type Reader = core.Reader
type Writer = core.Writer

var (
	stdin Reader = core.Stdin()

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
		stdin = core.Stdin()
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
	return core.Stdout()
}

func stderrWriter() Writer {
	ioMu.RLock()
	defer ioMu.RUnlock()
	if stderrOverride != nil {
		return stderrOverride
	}
	return core.Stderr()
}

func writeString(w Writer, s string) {
	if w == nil {
		return
	}
	if r := core.WriteString(w, s); !r.OK {
		core.Warn("cli write failed", "err", r.Error())
	}
}

func isEOF(err error) bool {
	return core.Is(err, core.EOF)
}

func writerFileDescriptor(v any) (int, bool) {
	file, ok := v.(interface{ Fd() uintptr })
	if !ok {
		return 0, false
	}
	return int(file.Fd()), true
}
