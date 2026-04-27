package cli

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestIO_SetStdin_Good(t *testing.T) {
	defer SetStdin(nil)

	r := strings.NewReader("hello\n")
	SetStdin(r)

	if got := stdinReader(); got != r {
		t.Errorf("stdinReader: expected override, got %T", got)
	}
}

func TestIO_SetStdin_Bad(t *testing.T) {
	defer SetStdin(nil)

	SetStdin(strings.NewReader("data"))
	// Passing nil should restore the real os.Stdin.
	SetStdin(nil)

	if got := stdinReader(); got != os.Stdin {
		t.Errorf("stdinReader: expected os.Stdin after SetStdin(nil), got %T", got)
	}
}

func TestIO_SetStdin_Ugly(t *testing.T) {
	// Concurrent SetStdin + stdinReader must not race.
	// Uses RWMutex internally — this test exercises the lock path.
	defer SetStdin(nil)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			SetStdin(strings.NewReader("concurrent"))
		}()
		go func() {
			defer wg.Done()
			_ = stdinReader()
		}()
	}
	wg.Wait()
}

func TestIO_SetStdout_Good(t *testing.T) {
	defer SetStdout(nil)

	buf := &bytes.Buffer{}
	SetStdout(buf)

	if got := stdoutWriter(); got != buf {
		t.Errorf("stdoutWriter: expected override, got %T", got)
	}
}

func TestIO_SetStdout_Bad(t *testing.T) {
	defer SetStdout(nil)

	SetStdout(&bytes.Buffer{})
	// Passing nil should clear the override — writes return to os.Stdout.
	SetStdout(nil)

	if got := stdoutWriter(); got != os.Stdout {
		t.Errorf("stdoutWriter: expected os.Stdout after SetStdout(nil), got %T", got)
	}
}

func TestIO_SetStdout_Ugly(t *testing.T) {
	// Override + writes must be observed through the injected writer.
	defer SetStdout(nil)

	buf := &bytes.Buffer{}
	SetStdout(buf)

	Println("hello %s", "world")

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("stdout override: expected 'hello world' in buffer, got %q", buf.String())
	}
}

func TestIO_SetStderr_Good(t *testing.T) {
	defer SetStderr(nil)

	buf := &bytes.Buffer{}
	SetStderr(buf)

	if got := stderrWriter(); got != buf {
		t.Errorf("stderrWriter: expected override, got %T", got)
	}
}

func TestIO_SetStderr_Bad(t *testing.T) {
	defer SetStderr(nil)

	SetStderr(&bytes.Buffer{})
	// Passing nil should clear the override — writes return to os.Stderr.
	SetStderr(nil)

	if got := stderrWriter(); got != os.Stderr {
		t.Errorf("stderrWriter: expected os.Stderr after SetStderr(nil), got %T", got)
	}
}

func TestIO_SetStderr_Ugly(t *testing.T) {
	// Concurrent readers and writers across all three streams must not race.
	defer func() {
		SetStdin(nil)
		SetStdout(nil)
		SetStderr(nil)
	}()

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(6)
		go func() { defer wg.Done(); SetStdout(&bytes.Buffer{}) }()
		go func() { defer wg.Done(); SetStderr(&bytes.Buffer{}) }()
		go func() { defer wg.Done(); SetStdin(strings.NewReader("x")) }()
		go func() { defer wg.Done(); _ = stdoutWriter() }()
		go func() { defer wg.Done(); _ = stderrWriter() }()
		go func() { defer wg.Done(); _ = stdinReader() }()
	}
	wg.Wait()
}
