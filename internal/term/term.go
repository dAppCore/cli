package term

import xterm "golang.org/x/term"

func IsTerminal(fd int) bool {
	return xterm.IsTerminal(fd)
}

func TerminalSize(fd int) (w, h int, err error) {
	return xterm.GetSize(fd)
}
