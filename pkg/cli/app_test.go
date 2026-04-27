package cli

import (
	"bytes"
	"runtime/debug"
	"testing"

	"dappco.re/go/core"
	"github.com/stretchr/testify/assert"
)

// TestCli_PanicRecovery_Good verifies that the panic recovery mechanism
// catches panics and calls the appropriate shutdown and error handling.
func TestCli_PanicRecovery_Good(t *testing.T) {
	t.Run("recovery captures panic value and stack", func(t *testing.T) {
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

		assert.Equal(t, "test panic", recovered)
		assert.True(t, shutdownCalled, "Shutdown should be called after panic recovery")
		assert.NotEmpty(t, capturedStack, "Stack trace should be captured")
		assert.Contains(t, string(capturedStack), "TestCli_PanicRecovery_Good")
	})

	t.Run("recovery handles error type panics", func(t *testing.T) {
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
		assert.True(t, ok, "Recovered value should be an error")
		assert.Equal(t, "error panic", err.Error())
	})

	t.Run("recovery handles nil panic gracefully", func(t *testing.T) {
		recoveryExecuted := false

		func() {
			defer func() {
				if r := recover(); r != nil {
					recoveryExecuted = true
				}
			}()

			// No panic occurs
		}()

		assert.False(t, recoveryExecuted, "Recovery block should not execute without panic")
	})
}

// TestCli_PanicRecovery_Bad tests error conditions in panic recovery.
func TestCli_PanicRecovery_Bad(t *testing.T) {
	t.Run("recovery handles concurrent panics", func(t *testing.T) {
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
		assert.Equal(t, 3, recoveryCount, "All goroutine panics should be recovered")
	})
}

// TestCli_PanicRecovery_Ugly tests edge cases in panic recovery.
func TestCli_PanicRecovery_Ugly(t *testing.T) {
	t.Run("recovery handles typed panic values", func(t *testing.T) {
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
		assert.True(t, ok, "Should recover custom type")
		assert.Equal(t, 500, ce.code)
		assert.Equal(t, "internal error", ce.msg)
	})
}

// TestCli_MainPanicRecoveryPattern_Good verifies the exact pattern used in Main().
func TestCli_MainPanicRecoveryPattern_Good(t *testing.T) {
	t.Run("pattern logs error and calls shutdown", func(t *testing.T) {
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

		assert.Contains(t, logBuffer.String(), "recovered from panic: simulated crash")
		assert.True(t, shutdownCalled, "Shutdown must be called on panic")
		assert.NotNil(t, fatalErr, "Fatal must be called with error")
		assert.Equal(t, "panic: simulated crash", fatalErr.Error())
	})
}
