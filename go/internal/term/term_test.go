package term

import (
	. "dappco.re/go"
)

func TestTerm_IsTerminal_Good(t *T) {
	got := IsTerminal(0)
	AssertEqual(t, got, IsTerminal(0))
	AssertNotPanics(t, func() { _ = IsTerminal(1) })
}

func TestTerm_IsTerminal_Bad(t *T) {
	got := IsTerminal(-1)
	AssertFalse(t, got)
	AssertEqual(t, false, got)
}

func TestTerm_IsTerminal_Ugly(t *T) {
	got := IsTerminal(1 << 20)
	AssertFalse(t, got)
	AssertNotPanics(t, func() { _ = IsTerminal(0) })
}

func TestTerm_TerminalSize_Good(t *T) {
	r := TerminalSize(0)
	AssertNotPanics(t, func() { _ = TerminalSize(0) })
	AssertNotNil(t, r)
}

func TestTerm_TerminalSize_Bad(t *T) {
	r := TerminalSize(-1)
	AssertFalse(t, r.OK)
	AssertNotEmpty(t, r.Error())
}

func TestTerm_TerminalSize_Ugly(t *T) {
	r := TerminalSize(1 << 20)
	AssertFalse(t, r.OK)
	AssertNotEmpty(t, r.Error())
}
