package cli

import (
	core "dappco.re/go"
)

func TestOutput_Blank_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Blank() })

	core.AssertContains(t, out, "\n")
	core.AssertNotPanics(t, func() { Blank() })
}

func TestOutput_Blank_Bad(t *core.T) {
	SetStdout(core.Discard)
	defer SetStdout(nil)

	core.AssertNotPanics(t, func() { Blank() })
	core.AssertNotNil(t, stdoutWriter())
}

func TestOutput_Blank_Ugly(t *core.T) {
	out := cliCaptureStdout(t, func() {
		Blank()
		Blank()
	})

	core.AssertContains(t, out, "\n")
	core.AssertTrue(t, core.RuneCount(out) >= 2)
}

func TestOutput_Echo_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Echo("i18n.progress.check") })

	core.AssertContains(t, out, "Checking")
	core.AssertContains(t, out, "...")
}

func TestOutput_Echo_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Echo("") })

	core.AssertContains(t, out, "\n")
	core.AssertNotContains(t, out, "Checking")
}

func TestOutput_Echo_Ugly(t *core.T) {
	out := cliCaptureStdout(t, func() { Echo("i18n.fail.load", "config") })

	core.AssertContains(t, out, "Failed to load config")
	core.AssertNotContains(t, out, "i18n.fail")
}

func TestOutput_Print_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Print("hello %s", "codex") })

	core.AssertEqual(t, "hello codex", out)
	core.AssertContains(t, out, "codex")
}

func TestOutput_Print_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Print("") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestOutput_Print_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Print(":check:") })

	core.AssertEqual(t, "[OK]", out)
	core.AssertNotContains(t, out, ":check:")
}

func TestOutput_Println_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Println("hello %s", "codex") })

	core.AssertContains(t, out, "hello codex")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Println_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Println("") })

	core.AssertContains(t, out, "\n")
	core.AssertEqual(t, 1, core.RuneCount(out))
}

func TestOutput_Println_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Println(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Text_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Text("count:", 2) })

	core.AssertContains(t, out, "count:2")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Text_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Text() })

	core.AssertContains(t, out, "\n")
	core.AssertEqual(t, 1, core.RuneCount(out))
}

func TestOutput_Text_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Text(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestOutput_Success_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Success("done") })

	core.AssertContains(t, out, "done")
	core.AssertContains(t, out, "[OK]")
}

func TestOutput_Success_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Success("") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, "done")
}

func TestOutput_Success_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Success(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestOutput_Successf_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Successf("done %d", 1) })

	core.AssertContains(t, out, "done 1")
	core.AssertContains(t, out, "[OK]")
}

func TestOutput_Successf_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Successf("") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, "done")
}

func TestOutput_Successf_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Successf("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[OK]")
}

func TestOutput_Error_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Error("failed") })

	core.AssertContains(t, out, "failed")
	core.AssertContains(t, out, "[FAIL]")
}

func TestOutput_Error_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Error("") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, "failed")
}

func TestOutput_Error_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Error(":cross:") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, ":cross:")
}

func TestOutput_Errorf_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Errorf("failed %d", 1) })

	core.AssertContains(t, out, "failed 1")
	core.AssertContains(t, out, "[FAIL]")
}

func TestOutput_Errorf_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Errorf("") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, "failed")
}

func TestOutput_Errorf_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Errorf("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[FAIL]")
}

func TestOutput_ErrorWrap_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { ErrorWrap(Err("root"), "wrap") })

	core.AssertContains(t, out, "wrap")
	core.AssertContains(t, out, "root")
}

func TestOutput_ErrorWrap_Bad(t *core.T) {
	out := cliCaptureStderr(t, func() { ErrorWrap(nil, "wrap") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestOutput_ErrorWrap_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { ErrorWrap(Err("root"), "") })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "[FAIL]")
}

func TestOutput_ErrorWrapVerb_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { ErrorWrapVerb(Err("root"), "load", "config") })

	core.AssertContains(t, out, "Failed to load config")
	core.AssertContains(t, out, "root")
}

func TestOutput_ErrorWrapVerb_Bad(t *core.T) {
	out := cliCaptureStderr(t, func() { ErrorWrapVerb(nil, "load", "config") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestOutput_ErrorWrapVerb_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { ErrorWrapVerb(Err("root"), "", "") })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "[FAIL]")
}

func TestOutput_ErrorWrapAction_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { ErrorWrapAction(Err("root"), "connect") })

	core.AssertContains(t, out, "Failed to connect")
	core.AssertContains(t, out, "root")
}

func TestOutput_ErrorWrapAction_Bad(t *core.T) {
	out := cliCaptureStderr(t, func() { ErrorWrapAction(nil, "connect") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestOutput_ErrorWrapAction_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { ErrorWrapAction(Err("root"), "") })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "[FAIL]")
}

func TestOutput_Warn_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Warn("careful") })

	core.AssertContains(t, out, "careful")
	core.AssertContains(t, out, "[WARN]")
}

func TestOutput_Warn_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Warn("") })

	core.AssertContains(t, out, "[WARN]")
	core.AssertNotContains(t, out, "careful")
}

func TestOutput_Warn_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Warn(":warn:") })

	core.AssertContains(t, out, "[WARN]")
	core.AssertNotContains(t, out, ":warn:")
}

func TestOutput_Warnf_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Warnf("careful %d", 1) })

	core.AssertContains(t, out, "careful 1")
	core.AssertContains(t, out, "[WARN]")
}

func TestOutput_Warnf_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Warnf("") })

	core.AssertContains(t, out, "[WARN]")
	core.AssertNotContains(t, out, "careful")
}

func TestOutput_Warnf_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Warnf("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[WARN]")
}

func TestOutput_Info_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Info("ready") })

	core.AssertContains(t, out, "ready")
	core.AssertContains(t, out, "[INFO]")
}

func TestOutput_Info_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Info("") })

	core.AssertContains(t, out, "[INFO]")
	core.AssertNotContains(t, out, "ready")
}

func TestOutput_Info_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Info(":info:") })

	core.AssertContains(t, out, "[INFO]")
	core.AssertNotContains(t, out, ":info:")
}

func TestOutput_Infof_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Infof("ready %d", 1) })

	core.AssertContains(t, out, "ready 1")
	core.AssertContains(t, out, "[INFO]")
}

func TestOutput_Infof_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Infof("") })

	core.AssertContains(t, out, "[INFO]")
	core.AssertNotContains(t, out, "ready")
}

func TestOutput_Infof_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Infof("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[INFO]")
}

func TestOutput_Dim_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Dim("quiet") })

	core.AssertContains(t, out, "quiet")
	core.AssertNotContains(t, out, "\033")
}

func TestOutput_Dim_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Dim("") })

	core.AssertContains(t, out, "\n")
	core.AssertNotContains(t, out, "quiet")
}

func TestOutput_Dim_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Dim(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestOutput_ProgressDone_Good(t *core.T) {
	out := cliCaptureStderr(t, func() { ProgressDone() })

	core.AssertContains(t, out, "\033[2K")
	core.AssertContains(t, out, "\r")
}

func TestOutput_ProgressDone_Bad(t *core.T) {
	SetStderr(core.Discard)
	defer SetStderr(nil)

	core.AssertNotPanics(t, func() { ProgressDone() })
	core.AssertNotNil(t, stderrWriter())
}

func TestOutput_ProgressDone_Ugly(t *core.T) {
	out := cliCaptureStderr(t, func() {
		ProgressDone()
		ProgressDone()
	})

	core.AssertContains(t, out, "\r")
	core.AssertTrue(t, core.Contains(out, "\033[2K"))
}

func TestOutput_Progress_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Progress("check", 1, 2, "repo") })

	core.AssertContains(t, out, "Checking")
	core.AssertContains(t, out, "1/2")
}

func TestOutput_Progress_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Progress("", 0, 0) })

	core.AssertContains(t, out, "0/0")
	core.AssertContains(t, out, "\r")
}

func TestOutput_Progress_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Progress("tie", -1, 3, "") })

	core.AssertContains(t, out, "-1/3")
	core.AssertContains(t, out, "Tying")
}

func TestOutput_Label_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Label("workspace", "/tmp") })

	core.AssertContains(t, out, "Workspace:")
	core.AssertContains(t, out, "/tmp")
}

func TestOutput_Label_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Label("", "") })

	core.AssertNotContains(t, out, ":")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Label_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Label(":check:", ":warn:") })

	core.AssertContains(t, out, "[OK]:")
	core.AssertContains(t, out, "[WARN]")
}

func TestOutput_Task_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Task("go", "Running") })

	core.AssertContains(t, out, "[go]")
	core.AssertContains(t, out, "Running")
}

func TestOutput_Task_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Task("", "") })

	core.AssertContains(t, out, "[]")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Task_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Task(":check:", ":warn:") })

	core.AssertContains(t, out, "[[OK]]")
	core.AssertContains(t, out, "[WARN]")
}

func TestOutput_Section_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Section("audit") })

	core.AssertContains(t, out, "AUDIT")
	core.AssertContains(t, out, "─")
}

func TestOutput_Section_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Section("") })

	core.AssertContains(t, out, "─")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Section_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Section(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestOutput_Hint_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Hint("fix", "run tests") })

	core.AssertContains(t, out, "fix:")
	core.AssertContains(t, out, "run tests")
}

func TestOutput_Hint_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Hint("", "") })

	core.AssertContains(t, out, ":")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Hint_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Hint(":check:", ":warn:") })

	core.AssertContains(t, out, "[OK]:")
	core.AssertContains(t, out, "[WARN]")
}

func TestOutput_Severity_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Severity("critical", "sql injection") })

	core.AssertContains(t, out, "[critical]")
	core.AssertContains(t, out, "sql injection")
}

func TestOutput_Severity_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Severity("unknown", "message") })

	core.AssertContains(t, out, "[unknown]")
	core.AssertContains(t, out, "message")
}

func TestOutput_Severity_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Severity("", "") })

	core.AssertContains(t, out, "[]")
	core.AssertContains(t, out, "\n")
}

func TestOutput_Result_Good(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStdout(t, func() { Result(true, "passed") })

	core.AssertContains(t, out, "passed")
	core.AssertContains(t, out, "[OK]")
}

func TestOutput_Result_Bad(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Result(false, "failed") })

	core.AssertContains(t, out, "failed")
	core.AssertContains(t, out, "[FAIL]")
}

func TestOutput_Result_Ugly(t *core.T) {
	cliPlainCLI(t)
	out := cliCaptureStderr(t, func() { Result(false, "") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, "failed")
}
