package cli

import (
	core "dappco.re/go"
)

func TestStrings_Sprintf_Good(t *core.T) {
	got := Sprintf("hello %s", "codex")

	core.AssertEqual(t, "hello codex", got)
	core.AssertContains(t, got, "codex")
}

func TestStrings_Sprintf_Bad(t *core.T) {
	got := Sprintf("%s", "bad")

	core.AssertEqual(t, "bad", got)
	core.AssertContains(t, got, "bad")
}

func TestStrings_Sprintf_Ugly(t *core.T) {
	got := Sprintf("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStrings_Sprint_Good(t *core.T) {
	got := Sprint("count:", 2)

	core.AssertEqual(t, "count:2", got)
	core.AssertContains(t, got, "2")
}

func TestStrings_Sprint_Bad(t *core.T) {
	got := Sprint()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStrings_Sprint_Ugly(t *core.T) {
	got := Sprint(nil, "x")

	core.AssertEqual(t, "<nil>x", got)
	core.AssertContains(t, got, "nil")
}

func TestStrings_Styled_Good(t *core.T) {
	cliPlainCLI(t)
	got := Styled(NewStyle().Bold(), ":check: ready")

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "[OK]")
}

func TestStrings_Styled_Bad(t *core.T) {
	cliPlainCLI(t)
	got := Styled(nil, ":missing:")

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestStrings_Styled_Ugly(t *core.T) {
	cliPlainCLI(t)
	got := Styled(NewStyle(), "")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStrings_Styledf_Good(t *core.T) {
	cliPlainCLI(t)
	got := Styledf(NewStyle().Bold(), "%s", ":check:")

	core.AssertEqual(t, "[OK]", got)
	core.AssertContains(t, got, "[OK]")
}

func TestStrings_Styledf_Bad(t *core.T) {
	got := Styledf(nil, "")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStrings_Styledf_Ugly(t *core.T) {
	got := Styledf(nil, "%s", "bad")

	core.AssertEqual(t, "bad", got)
	core.AssertContains(t, got, "bad")
}

func TestStrings_SuccessStr_Good(t *core.T) {
	cliPlainCLI(t)
	got := SuccessStr("done")

	core.AssertContains(t, got, "done")
	core.AssertContains(t, got, "[OK]")
}

func TestStrings_SuccessStr_Bad(t *core.T) {
	cliPlainCLI(t)
	got := SuccessStr("")

	core.AssertContains(t, got, "[OK]")
	core.AssertNotContains(t, got, "done")
}

func TestStrings_SuccessStr_Ugly(t *core.T) {
	cliPlainCLI(t)
	got := SuccessStr(":check:")

	core.AssertContains(t, got, "[OK]")
	core.AssertNotContains(t, got, ":check:")
}

func TestStrings_ErrorStr_Good(t *core.T) {
	cliPlainCLI(t)
	got := ErrorStr("failed")

	core.AssertContains(t, got, "failed")
	core.AssertContains(t, got, "[FAIL]")
}

func TestStrings_ErrorStr_Bad(t *core.T) {
	cliPlainCLI(t)
	got := ErrorStr("")

	core.AssertContains(t, got, "[FAIL]")
	core.AssertNotContains(t, got, "failed")
}

func TestStrings_ErrorStr_Ugly(t *core.T) {
	cliPlainCLI(t)
	got := ErrorStr(":cross:")

	core.AssertContains(t, got, "[FAIL]")
	core.AssertNotContains(t, got, ":cross:")
}

func TestStrings_WarnStr_Good(t *core.T) {
	cliPlainCLI(t)
	got := WarnStr("careful")

	core.AssertContains(t, got, "careful")
	core.AssertContains(t, got, "[WARN]")
}

func TestStrings_WarnStr_Bad(t *core.T) {
	cliPlainCLI(t)
	got := WarnStr("")

	core.AssertContains(t, got, "[WARN]")
	core.AssertNotContains(t, got, "careful")
}

func TestStrings_WarnStr_Ugly(t *core.T) {
	cliPlainCLI(t)
	got := WarnStr(":warn:")

	core.AssertContains(t, got, "[WARN]")
	core.AssertNotContains(t, got, ":warn:")
}

func TestStrings_InfoStr_Good(t *core.T) {
	cliPlainCLI(t)
	got := InfoStr("ready")

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "[INFO]")
}

func TestStrings_InfoStr_Bad(t *core.T) {
	cliPlainCLI(t)
	got := InfoStr("")

	core.AssertContains(t, got, "[INFO]")
	core.AssertNotContains(t, got, "ready")
}

func TestStrings_InfoStr_Ugly(t *core.T) {
	cliPlainCLI(t)
	got := InfoStr(":info:")

	core.AssertContains(t, got, "[INFO]")
	core.AssertNotContains(t, got, ":info:")
}

func TestStrings_DimStr_Good(t *core.T) {
	cliPlainCLI(t)
	got := DimStr("quiet")

	core.AssertEqual(t, "quiet", got)
	core.AssertContains(t, got, "quiet")
}

func TestStrings_DimStr_Bad(t *core.T) {
	cliPlainCLI(t)
	got := DimStr("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStrings_DimStr_Ugly(t *core.T) {
	cliPlainCLI(t)
	got := DimStr(":check:")

	core.AssertEqual(t, "[OK]", got)
	core.AssertNotContains(t, got, ":check:")
}

func TestStrings_Repeat_Good(t *core.T) {
	got := Repeat("ab", 3)
	core.AssertEqual(t, "ababab", got)
	core.AssertContains(t, got, "ab")
}

func TestStrings_Repeat_Bad(t *core.T) {
	got := Repeat("ab", -1)
	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStrings_Repeat_Ugly(t *core.T) {
	got := Repeat("", 5)
	core.AssertEqual(t, "", got)
	core.AssertNotPanics(t, func() { _ = Repeat("x", 0) })
}

func TestStrings_LastIndex_Good(t *core.T) {
	got := LastIndex("a/b/c", "/")
	core.AssertEqual(t, 3, got)
	core.AssertGreater(t, got, 0)
}

func TestStrings_LastIndex_Bad(t *core.T) {
	got := LastIndex("abc", "/")
	core.AssertEqual(t, -1, got)
	core.AssertLess(t, got, 0)
}

func TestStrings_LastIndex_Ugly(t *core.T) {
	got := LastIndex("abc", "")
	core.AssertEqual(t, 3, got)
	core.AssertNotPanics(t, func() { _ = LastIndex("", "") })
}

func TestStrings_Atoi_Good(t *core.T) {
	r := Atoi("42")
	core.AssertTrue(t, r.OK)
	core.AssertEqual(t, 42, r.Value)
}

func TestStrings_Atoi_Bad(t *core.T) {
	r := Atoi("nope")
	core.AssertFalse(t, r.OK)
	core.AssertContains(t, r.Error(), "invalid")
}

func TestStrings_Atoi_Ugly(t *core.T) {
	r := Atoi("")
	core.AssertFalse(t, r.OK)
	core.AssertNotEmpty(t, r.Error())
}

func TestStrings_ParseHexByte_Good(t *core.T) {
	r := ParseHexByte("ff")
	core.AssertTrue(t, r.OK)
	core.AssertEqual(t, 255, r.Value)
}

func TestStrings_ParseHexByte_Bad(t *core.T) {
	r := ParseHexByte("zz")
	core.AssertFalse(t, r.OK)
	core.AssertNotEmpty(t, r.Error())
}

func TestStrings_ParseHexByte_Ugly(t *core.T) {
	r := ParseHexByte("fff")
	core.AssertFalse(t, r.OK)
	core.AssertContains(t, r.Error(), "range")
}
