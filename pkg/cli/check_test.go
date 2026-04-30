package cli

import (
	core "dappco.re/go"
)

func TestCheck_Check_Good(t *core.T) {
	check := Check("audit")

	core.AssertNotNil(t, check)
	core.AssertContains(t, check.String(), "audit")
}

func TestCheck_Check_Bad(t *core.T) {
	check := Check("")

	core.AssertNotNil(t, check)
	core.AssertNotContains(t, check.String(), "\n")
}

func TestCheck_Check_Ugly(t *core.T) {
	check := Check(":check:")

	core.AssertContains(t, check.String(), "✓")
	core.AssertNotContains(t, check.String(), ":check:")
}

func TestCheck_CheckBuilder_Pass_Good(t *core.T) {
	check := Check("audit").Pass()

	core.AssertEqual(t, "passed", check.status)
	core.AssertNotNil(t, check.style)
}

func TestCheck_CheckBuilder_Pass_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Pass() })
	core.AssertNil(t, check)
}

func TestCheck_CheckBuilder_Pass_Ugly(t *core.T) {
	check := Check("audit").Fail().Pass()

	core.AssertEqual(t, "passed", check.status)
	core.AssertEqual(t, Glyph(":check:"), check.icon)
}

func TestCheck_CheckBuilder_Fail_Good(t *core.T) {
	check := Check("audit").Fail()

	core.AssertEqual(t, "failed", check.status)
	core.AssertNotNil(t, check.style)
}

func TestCheck_CheckBuilder_Fail_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Fail() })
	core.AssertNil(t, check)
}

func TestCheck_CheckBuilder_Fail_Ugly(t *core.T) {
	check := Check("audit").Pass().Fail()

	core.AssertEqual(t, "failed", check.status)
	core.AssertEqual(t, Glyph(":cross:"), check.icon)
}

func TestCheck_CheckBuilder_Skip_Good(t *core.T) {
	check := Check("audit").Skip()

	core.AssertEqual(t, "skipped", check.status)
	core.AssertNotNil(t, check.style)
}

func TestCheck_CheckBuilder_Skip_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Skip() })
	core.AssertNil(t, check)
}

func TestCheck_CheckBuilder_Skip_Ugly(t *core.T) {
	check := Check("audit").Fail().Skip()

	core.AssertEqual(t, "skipped", check.status)
	core.AssertEqual(t, Glyph(":skip:"), check.icon)
}

func TestCheck_CheckBuilder_Warn_Good(t *core.T) {
	check := Check("audit").Warn()

	core.AssertEqual(t, "warning", check.status)
	core.AssertNotNil(t, check.style)
}

func TestCheck_CheckBuilder_Warn_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Warn() })
	core.AssertNil(t, check)
}

func TestCheck_CheckBuilder_Warn_Ugly(t *core.T) {
	check := Check("audit").Pass().Warn()

	core.AssertEqual(t, "warning", check.status)
	core.AssertEqual(t, Glyph(":warn:"), check.icon)
}

func TestCheck_CheckBuilder_Duration_Good(t *core.T) {
	check := Check("audit").Duration("1s")

	core.AssertEqual(t, "1s", check.duration)
	core.AssertContains(t, check.String(), "1s")
}

func TestCheck_CheckBuilder_Duration_Bad(t *core.T) {
	check := Check("audit").Duration("")

	core.AssertEqual(t, "", check.duration)
	core.AssertNotContains(t, check.String(), "1s")
}

func TestCheck_CheckBuilder_Duration_Ugly(t *core.T) {
	check := Check("audit").Duration("∞")

	core.AssertEqual(t, "∞", check.duration)
	core.AssertContains(t, check.String(), "∞")
}

func TestCheck_CheckBuilder_Message_Good(t *core.T) {
	check := Check("audit").Message("ready")

	core.AssertEqual(t, "ready", check.status)
	core.AssertContains(t, check.String(), "ready")
}

func TestCheck_CheckBuilder_Message_Bad(t *core.T) {
	check := Check("audit").Message("")

	core.AssertEqual(t, "", check.status)
	core.AssertNotContains(t, check.String(), "ready")
}

func TestCheck_CheckBuilder_Message_Ugly(t *core.T) {
	check := Check("audit").Message(":check:")

	core.AssertEqual(t, ":check:", check.status)
	core.AssertContains(t, check.String(), "✓")
}

func TestCheck_CheckBuilder_String_Good(t *core.T) {
	got := Check("audit").Pass().String()

	core.AssertContains(t, got, "audit")
	core.AssertContains(t, got, "passed")
}

func TestCheck_CheckBuilder_String_Bad(t *core.T) {
	got := Check("").String()

	core.AssertNotContains(t, got, "\n")
	core.AssertNotContains(t, got, "passed")
}

func TestCheck_CheckBuilder_String_Ugly(t *core.T) {
	got := Check(":check:").Warn().String()

	core.AssertContains(t, got, "✓")
	core.AssertContains(t, got, "warning")
}

func TestCheck_CheckBuilder_Print_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Check("audit").Pass().Print() })

	core.AssertContains(t, out, "audit")
	core.AssertContains(t, out, "passed")
}

func TestCheck_CheckBuilder_Print_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Check("").Print() })

	core.AssertContains(t, out, "\n")
	core.AssertNotContains(t, out, "passed")
}

func TestCheck_CheckBuilder_Print_Ugly(t *core.T) {
	out := cliCaptureStdout(t, func() { Check(":check:").Warn().Print() })

	core.AssertContains(t, out, "✓")
	core.AssertContains(t, out, "warning")
}
