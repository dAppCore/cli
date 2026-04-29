package cli

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/i18n"
)

// ─────────────────────────────────────────────────────────────────────────────
// Error Creation (replace fmt.Errorf)
// ─────────────────────────────────────────────────────────────────────────────

// Err creates a failed Result from a format string.
func Err(format string, args ...any) core.Result {
	return core.Fail(core.E("cli", core.Sprintf(format, args...), nil))
}

// Wrap wraps an error with a message.
// Returns nil if err is nil.
//
//	return cli.Wrap(err, "load config")  // "load config: <original error>"
func Wrap(err error, msg string) core.Result {
	if err == nil {
		return core.Ok(nil)
	}
	return core.Fail(core.E("cli", msg, err))
}

// WrapVerb wraps an error using i18n grammar for "Failed to verb subject".
// Uses the i18n.ActionFailed function for proper grammar composition.
// Returns nil if err is nil.
//
//	return cli.WrapVerb(err, "load", "config")  // "Failed to load config: <original error>"
func WrapVerb(err error, verb, subject string) core.Result {
	if err == nil {
		return core.Ok(nil)
	}
	msg := i18n.ActionFailed(verb, subject)
	return core.Fail(core.E("cli", msg, err))
}

// WrapAction wraps an error using i18n grammar for "Failed to verb".
// Uses the i18n.ActionFailed function for proper grammar composition.
// Returns nil if err is nil.
//
//	return cli.WrapAction(err, "connect")  // "Failed to connect: <original error>"
func WrapAction(err error, verb string) core.Result {
	if err == nil {
		return core.Ok(nil)
	}
	msg := i18n.ActionFailed(verb, "")
	return core.Fail(core.E("cli", msg, err))
}

// ─────────────────────────────────────────────────────────────────────────────
// Error Helpers
// ─────────────────────────────────────────────────────────────────────────────

// Is reports whether any error in err's tree matches target.
// This is a re-export of errors.Is for convenience.
func Is(err, target error) bool {
	return core.Is(err, target)
}

// As finds the first error in err's tree that matches target.
// This is a re-export of errors.As for convenience.
func As(err error, target any) bool {
	return core.As(err, target)
}

// Join returns an error that wraps the given errors.
// This is a re-export of errors.Join for convenience.
func Join(errs ...error) core.Result {
	err := core.ErrorJoin(errs...)
	if err == nil {
		return core.Ok(nil)
	}
	return core.Fail(err)
}

// ExitError represents an error that should cause the CLI to exit with a specific code.
//
//	err := cli.Exit(2, core.NewError("validation failed"))
//	var exitErr *cli.ExitError
//	if cli.As(err, &exitErr) {
//	    cli.Println("exit code:", exitErr.Code)
//	}
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Exit creates a new ExitError with the given code and error.
//
//	return cli.Exit(2, core.NewError("validation failed"))
func Exit(code int, err error) core.Result {
	if err == nil {
		return core.Ok(nil)
	}
	return core.Fail(&ExitError{Code: code, Err: err})
}

// ─────────────────────────────────────────────────────────────────────────────
// Fatal Functions (Deprecated - return error from command instead)
// ─────────────────────────────────────────────────────────────────────────────

// Fatal prints an error message to stderr, logs it, and exits with code 1.
//
// Deprecated: return an error from the command instead.
func Fatal(failure any) {
	var err error
	switch v := failure.(type) {
	case core.Result:
		if v.OK {
			return
		}
		err, _ = v.Value.(error)
		if err == nil {
			err = core.NewError(v.Error())
		}
	case error:
		err = v
	}
	if err != nil {
		LogError("Fatal error", "err", err)
		core.Print(stderrWriter(), "%s", ErrorStyle.Render(Glyph(":cross:")+" "+err.Error()))
		core.Exit(1)
	}
}

// Fatalf prints a formatted error message to stderr, logs it, and exits with code 1.
//
// Deprecated: return an error from the command instead.
func Fatalf(format string, args ...any) {
	msg := core.Sprintf(format, args...)
	LogError("Fatal error", "msg", msg)
	core.Print(stderrWriter(), "%s", ErrorStyle.Render(Glyph(":cross:")+" "+msg))
	core.Exit(1)
}

// FatalWrap prints a wrapped error message to stderr, logs it, and exits with code 1.
// Does nothing if err is nil.
//
// Deprecated: return an error from the command instead.
//
//	cli.FatalWrap(err, "load config")  // Prints "✗ load config: <error>" and exits
func FatalWrap(err error, msg string) {
	if err == nil {
		return
	}
	LogError("Fatal error", "msg", msg, "err", err)
	fullMsg := core.Sprintf("%s: %v", msg, err)
	core.Print(stderrWriter(), "%s", ErrorStyle.Render(Glyph(":cross:")+" "+fullMsg))
	core.Exit(1)
}

// FatalWrapVerb prints a wrapped error using i18n grammar to stderr, logs it, and exits with code 1.
// Does nothing if err is nil.
//
// Deprecated: return an error from the command instead.
//
//	cli.FatalWrapVerb(err, "load", "config")  // Prints "✗ Failed to load config: <error>" and exits
func FatalWrapVerb(err error, verb, subject string) {
	if err == nil {
		return
	}
	msg := i18n.ActionFailed(verb, subject)
	LogError("Fatal error", "msg", msg, "err", err, "verb", verb, "subject", subject)
	fullMsg := core.Sprintf("%s: %v", msg, err)
	core.Print(stderrWriter(), "%s", ErrorStyle.Render(Glyph(":cross:")+" "+fullMsg))
	core.Exit(1)
}
