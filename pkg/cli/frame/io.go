package frame

import (
	"dappco.re/go"
	xterm "golang.org/x/term"
)

type Writer = core.Writer

// TODO(mantis-558): Replace the direct x/term calls below with
// dappco.re/go terminal primitives once core/go exposes them.

func stdoutWriter() Writer {
	return core.Stdout()
}

func stderrWriter() Writer {
	return core.Stderr()
}

func writerFileDescriptor(w Writer) (int, bool) {
	file, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return 0, false
	}
	return int(file.Fd()), true
}

func isTerminal(fd int) bool {
	return xterm.IsTerminal(fd)
}

func terminalSize(fd int) core.Result {
	w, h, err := xterm.GetSize(fd)
	if err != nil {
		return core.Fail(err)
	}
	return core.Ok([]int{w, h})
}
