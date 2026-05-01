package cli

import (
	core "dappco.re/go"
)

func TestIo_SetStdin_Good(t *core.T) {
	reader := core.NewReader("input")
	SetStdin(reader)
	defer SetStdin(nil)
	core.AssertEqual(t, reader, stdinReader())
}

func TestIo_SetStdin_Bad(t *core.T) {
	SetStdin(nil)
	core.AssertNotNil(t, stdinReader())
	core.AssertNotPanics(t, func() { SetStdin(nil) })
}

func TestIo_SetStdin_Ugly(t *core.T) {
	first := core.NewReader("")
	SetStdin(first)
	SetStdin(nil)
	core.AssertNotEqual(t, first, stdinReader())
}

func TestIo_SetStdout_Good(t *core.T) {
	out := core.NewBuilder()
	SetStdout(out)
	defer SetStdout(nil)
	core.AssertEqual(t, out, stdoutWriter())
}

func TestIo_SetStdout_Bad(t *core.T) {
	SetStdout(nil)
	core.AssertNotNil(t, stdoutWriter())
	core.AssertNotPanics(t, func() { SetStdout(nil) })
}

func TestIo_SetStdout_Ugly(t *core.T) {
	out := core.NewBuilder()
	SetStdout(out)
	writeString(stdoutWriter(), "x")
	core.AssertEqual(t, "x", out.String())
}

func TestIo_SetStderr_Good(t *core.T) {
	errOut := core.NewBuilder()
	SetStderr(errOut)
	defer SetStderr(nil)
	core.AssertEqual(t, errOut, stderrWriter())
}

func TestIo_SetStderr_Bad(t *core.T) {
	SetStderr(nil)
	core.AssertNotNil(t, stderrWriter())
	core.AssertNotPanics(t, func() { SetStderr(nil) })
}

func TestIo_SetStderr_Ugly(t *core.T) {
	errOut := core.NewBuilder()
	SetStderr(errOut)
	writeString(stderrWriter(), "x")
	core.AssertEqual(t, "x", errOut.String())
}
