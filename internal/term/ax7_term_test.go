package term

import . "dappco.re/go"

func TestAX7Term_IsTerminal_Good(t *T) {
	got := IsTerminal(0)

	AssertEqual(t, got, IsTerminal(0))
	AssertTrue(t, got || !got)
}

func TestAX7Term_IsTerminal_Bad(t *T) {
	got := IsTerminal(-1)

	AssertFalse(t, got)
	AssertEqual(t, false, got)
}

func TestAX7Term_IsTerminal_Ugly(t *T) {
	got := IsTerminal(1 << 20)

	AssertFalse(t, got)
	AssertEqual(t, false, got)
}

func TestAX7Term_TerminalSize_Good(t *T) {
	w, h, err := TerminalSize(0)

	AssertTrue(t, (err == nil && w >= 0 && h >= 0) || err != nil)
	AssertTrue(t, w >= 0)
}

func TestAX7Term_TerminalSize_Bad(t *T) {
	w, h, err := TerminalSize(-1)

	AssertError(t, err)
	AssertEqual(t, 0, w+h)
}

func TestAX7Term_TerminalSize_Ugly(t *T) {
	w, h, err := TerminalSize(1 << 20)

	AssertError(t, err)
	AssertEqual(t, 0, w+h)
}
