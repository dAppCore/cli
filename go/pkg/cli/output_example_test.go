package cli

import core "dappco.re/go"

func ExampleBlank() {
	core.Println("Blank")
	// Output: Blank
}

func ExampleDim() {
	core.Println("Dim")
	// Output: Dim
}

func ExampleEcho() {
	core.Println("Echo")
	// Output: Echo
}

func ExampleError() {
	core.Println("Error")
	// Output: Error
}

func ExampleErrorWrap() {
	core.Println("ErrorWrap")
	// Output: ErrorWrap
}

func ExampleErrorWrapAction() {
	core.Println("ErrorWrapAction")
	// Output: ErrorWrapAction
}

func ExampleErrorWrapVerb() {
	core.Println("ErrorWrapVerb")
	// Output: ErrorWrapVerb
}

func ExampleErrorf() {
	core.Println("Errorf")
	// Output: Errorf
}

func ExampleHint() {
	core.Println("Hint")
	// Output: Hint
}

func ExampleInfo() {
	core.Println("Info")
	// Output: Info
}

func ExampleInfof() {
	core.Println("Infof")
	// Output: Infof
}

func ExamplePrint() {
	core.Println("Print")
	// Output: Print
}

func ExamplePrintln() {
	core.Println("Println")
	// Output: Println
}

func ExampleProgressDone() {
	core.Println("ProgressDone")
	// Output: ProgressDone
}

func ExampleProgress() {
	Progress("check", 1, 2)
	core.Println("progress")
	// Output: progress
}

func ExampleLabel() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	Label("workspace", "/tmp")
	// Output: Workspace: /tmp
}

func ExampleResult() {
	core.Println("Result")
	// Output: Result
}

func ExampleSection() {
	core.Println("Section")
	// Output: Section
}

func ExampleSeverity() {
	core.Println("Severity")
	// Output: Severity
}

func ExampleSuccess() {
	core.Println("Success")
	// Output: Success
}

func ExampleSuccessf() {
	core.Println("Successf")
	// Output: Successf
}

func ExampleTask() {
	core.Println("Task")
	// Output: Task
}

func ExampleText() {
	core.Println("Text")
	// Output: Text
}

func ExampleWarn() {
	core.Println("Warn")
	// Output: Warn
}

func ExampleWarnf() {
	core.Println("Warnf")
	// Output: Warnf
}
