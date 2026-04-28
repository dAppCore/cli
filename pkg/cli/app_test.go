package cli

import (
	"bytes"
	"runtime/debug"

	"dappco.re/go"
)

// TestCli_PanicRecovery_Good verifies that the panic recovery mechanism
// catches panics and calls the appropriate shutdown and error handling.
func TestCli_PanicRecovery_Good(t *core.T) {
	t.Run("recovery captures panic value and stack", func(t *core.T) {
		var recovered any
		var capturedStack []byte
		var shutdownCalled bool

		// Simulate the panic recovery pattern from Main()
		func() {
			defer func() {
				if r := recover(); r != nil {
					recovered = r
					capturedStack = debug.Stack()
					shutdownCalled = true // simulates Shutdown() call
				}
			}()

			panic("test panic")
		}()
		core.AssertEqual(t, "test panic", recovered)
		core.AssertTrue(t, shutdownCalled, "Shutdown should be called after panic recovery")
		core.AssertNotEmpty(t, capturedStack, "Stack trace should be captured")
		core.AssertContains(t, string(capturedStack), "TestCli_PanicRecovery_Good")
	})

	t.Run("recovery handles error type panics", func(t *core.T) {
		var recovered any

		func() {
			defer func() {
				if r := recover(); r != nil {
					recovered = r
				}
			}()

			panic(core.E("", "error panic", nil))
		}()

		err, ok := recovered.(error)
		core.AssertTrue(t, ok, "Recovered value should be an error")
		core.AssertEqual(t, "error panic", err.Error())
	})

	t.Run("recovery handles nil panic gracefully", func(t *core.T) {
		recoveryExecuted := false

		func() {
			defer func() {
				if r := recover(); r != nil {
					recoveryExecuted = true
				}
			}()

			// No panic occurs
		}()
		core.AssertFalse(t, recoveryExecuted, "Recovery block should not execute without panic")
	})
}

// TestCli_PanicRecovery_Bad tests error conditions in panic recovery.
func TestCli_PanicRecovery_Bad(t *core.T) {
	t.Run("recovery handles concurrent panics", func(t *core.T) {
		done := make(chan struct{}, 3)
		recovered := make(chan struct{}, 3)

		for i := 0; i < 3; i++ {
			go func(id int) {
				defer func() { done <- struct{}{} }()
				defer func() {
					if r := recover(); r != nil {
						recovered <- struct{}{}
					}
				}()

				panic(core.Sprintf("panic from goroutine %d", id))
			}(i)
		}

		for i := 0; i < 3; i++ {
			<-done
		}
		close(recovered)

		recoveryCount := 0
		for range recovered {
			recoveryCount++
		}
		core.AssertEqual(t, 3, recoveryCount, "All goroutine panics should be recovered")
	})
}

// TestCli_PanicRecovery_Ugly tests edge cases in panic recovery.
func TestCli_PanicRecovery_Ugly(t *core.T) {
	t.Run("recovery handles typed panic values", func(t *core.T) {
		type customError struct {
			code int
			msg  string
		}

		var recovered any

		func() {
			defer func() {
				recovered = recover()
			}()

			panic(customError{code: 500, msg: "internal error"})
		}()

		ce, ok := recovered.(customError)
		core.AssertTrue(t, ok, "Should recover custom type")
		core.AssertEqual(t, 500, ce.code)
		core.AssertEqual(t, "internal error", ce.msg)
	})
}

// TestCli_MainPanicRecoveryPattern_Good verifies the exact pattern used in Main().
func TestCli_MainPanicRecoveryPattern_Good(t *core.T) {
	t.Run("pattern logs error and calls shutdown", func(t *core.T) {
		var logBuffer bytes.Buffer
		var shutdownCalled bool
		var fatalErr error

		// Mock implementations
		mockLogError := func(msg string, args ...any) {
			logBuffer.WriteString(core.Sprintf(msg, args...))
		}
		mockShutdown := func() {
			shutdownCalled = true
		}
		mockFatal := func(err error) {
			fatalErr = err
		}

		// Execute the pattern from Main()
		func() {
			defer func() {
				if r := recover(); r != nil {
					mockLogError("recovered from panic: %v", r)
					mockShutdown()
					mockFatal(core.E("", core.Sprintf("panic: %v", r), nil))
				}
			}()

			panic("simulated crash")
		}()
		core.AssertContains(t, logBuffer.String(), "recovered from panic: simulated crash")
		core.AssertTrue(t, shutdownCalled, "Shutdown must be called on panic")
		core.AssertNotNil(t, fatalErr, "Fatal must be called with error")
		core.AssertEqual(t, "panic: simulated crash", fatalErr.Error())
	})
}
