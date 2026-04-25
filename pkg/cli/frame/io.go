package frame

import (
	"io"
	"os"

	xterm "golang.org/x/term"
)

type Writer = io.Writer

func stdoutWriter() Writer {
	return os.Stdout
}

func stderrWriter() Writer {
	return os.Stderr
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

func terminalSize(fd int) (w, h int, err error) {
	return xterm.GetSize(fd)
}
