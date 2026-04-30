package term

import (
	"dappco.re/go"
	xterm "golang.org/x/term"
)

func IsTerminal(fd int) bool {
	return xterm.IsTerminal(fd)
}

func TerminalSize(fd int) core.Result {
	w, h, err := xterm.GetSize(fd)
	if err != nil {
		return core.Fail(err)
	}
	return core.Ok([]int{w, h})
}
