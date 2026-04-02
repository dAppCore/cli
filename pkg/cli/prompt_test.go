package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	oldErr := os.Stderr
	r, w, err := os.Pipe()
	if !assert.NoError(t, err) {
		return ""
	}
	os.Stderr = w

	defer func() {
		os.Stderr = oldErr
	}()

	fn()

	if !assert.NoError(t, w.Close()) {
		return ""
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if !assert.NoError(t, err) {
		return ""
	}
	return buf.String()
}

func captureStdoutStderr(t *testing.T, fn func()) (string, string) {
	t.Helper()

	oldOut := os.Stdout
	oldErr := os.Stderr
	rOut, wOut, err := os.Pipe()
	if !assert.NoError(t, err) {
		return "", ""
	}
	rErr, wErr, err := os.Pipe()
	if !assert.NoError(t, err) {
		return "", ""
	}
	os.Stdout = wOut
	os.Stderr = wErr

	defer func() {
		os.Stdout = oldOut
		os.Stderr = oldErr
	}()

	fn()

	if !assert.NoError(t, wOut.Close()) {
		return "", ""
	}
	if !assert.NoError(t, wErr.Close()) {
		return "", ""
	}

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	_, err = io.Copy(&outBuf, rOut)
	if !assert.NoError(t, err) {
		return "", ""
	}
	_, err = io.Copy(&errBuf, rErr)
	if !assert.NoError(t, err) {
		return "", ""
	}
	return outBuf.String(), errBuf.String()
}

func TestPrompt_Good(t *testing.T) {
	SetStdin(strings.NewReader("hello\n"))
	defer SetStdin(nil) // reset

	val, err := Prompt("Name", "")
	assert.NoError(t, err)
	assert.Equal(t, "hello", val)
}

func TestPrompt_Good_Default(t *testing.T) {
	SetStdin(strings.NewReader("\n"))
	defer SetStdin(nil)

	val, err := Prompt("Name", "world")
	assert.NoError(t, err)
	assert.Equal(t, "world", val)
}

func TestPrompt_Bad_EOFUsesDefault(t *testing.T) {
	SetStdin(strings.NewReader(""))
	defer SetStdin(nil)

	val, err := Prompt("Name", "world")
	assert.NoError(t, err)
	assert.Equal(t, "world", val)
}

func TestPrompt_Bad_EOFWithoutDefaultReturnsError(t *testing.T) {
	SetStdin(strings.NewReader(""))
	defer SetStdin(nil)

	val, err := Prompt("Name", "")
	assert.ErrorIs(t, err, io.EOF)
	assert.Empty(t, val)
}

func TestSelect_Good(t *testing.T) {
	SetStdin(strings.NewReader("2\n"))
	defer SetStdin(nil)

	val, err := Select("Pick", []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Equal(t, "b", val)
}

func TestSelect_Bad_Invalid(t *testing.T) {
	SetStdin(strings.NewReader("5\n"))
	defer SetStdin(nil)

	var err error
	stderr := captureStderr(t, func() {
		_, err = Select("Pick", []string{"a", "b"})
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "choose a number between 1 and 2")
	assert.Contains(t, stderr, "Please enter a number between 1 and 2.")
}

func TestSelect_Bad_EOF(t *testing.T) {
	SetStdin(strings.NewReader(""))
	defer SetStdin(nil)

	_, err := Select("Pick", []string{"a", "b"})
	assert.ErrorIs(t, err, io.EOF)
}

func TestSelect_Good_EmptyOptions(t *testing.T) {
	val, err := Select("Pick", nil)
	assert.NoError(t, err)
	assert.Empty(t, val)
}

func TestMultiSelect_Good(t *testing.T) {
	SetStdin(strings.NewReader("1 3\n"))
	defer SetStdin(nil)

	vals, err := MultiSelect("Pick", []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "c"}, vals)
}

func TestMultiSelect_Good_CommasAndRanges(t *testing.T) {
	SetStdin(strings.NewReader("1-2,4\n"))
	defer SetStdin(nil)

	vals, err := MultiSelect("Pick", []string{"a", "b", "c", "d"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "d"}, vals)
}

func TestMultiSelect_Bad_EOFReturnsEmptySelection(t *testing.T) {
	SetStdin(strings.NewReader(""))
	defer SetStdin(nil)

	vals, err := MultiSelect("Pick", []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Empty(t, vals)
}

func TestMultiSelect_Good_EOFWithInput(t *testing.T) {
	SetStdin(strings.NewReader("1 3"))
	defer SetStdin(nil)

	vals, err := MultiSelect("Pick", []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "c"}, vals)
}

func TestMultiSelect_Good_DedupesSelections(t *testing.T) {
	SetStdin(strings.NewReader("1 1 2-3 2\n"))
	defer SetStdin(nil)

	vals, err := MultiSelect("Pick", []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, vals)
}

func TestMultiSelect_Bad_InvalidInput(t *testing.T) {
	SetStdin(strings.NewReader("1 foo\n"))
	defer SetStdin(nil)

	_, err := MultiSelect("Pick", []string{"a", "b", "c"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid selection")
}

func TestMultiSelect_Good_EmptyOptions(t *testing.T) {
	vals, err := MultiSelect("Pick", nil)
	assert.NoError(t, err)
	assert.Empty(t, vals)
}

func TestConfirm_Good(t *testing.T) {
	SetStdin(strings.NewReader("y\n"))
	defer SetStdin(nil)

	assert.True(t, Confirm("Proceed?"))
}

func TestConfirm_Bad_EOFUsesDefault(t *testing.T) {
	SetStdin(strings.NewReader(""))
	defer SetStdin(nil)

	assert.False(t, Confirm("Proceed?", Required()))
	assert.True(t, Confirm("Proceed?", DefaultYes(), Required()))
}

func TestConfirm_Good_RequiredReprompts(t *testing.T) {
	SetStdin(strings.NewReader("\ny\n"))
	defer SetStdin(nil)

	assert.True(t, Confirm("Proceed?", Required()))
}

func TestQuestion_Good(t *testing.T) {
	SetStdin(strings.NewReader("alice\n"))
	defer SetStdin(nil)

	val := Question("Name:")
	assert.Equal(t, "alice", val)
}

func TestQuestion_Bad_EOFReturnsDefault(t *testing.T) {
	SetStdin(strings.NewReader(""))
	defer SetStdin(nil)

	assert.Equal(t, "anonymous", Question("Name:", WithDefault("anonymous")))
	assert.Equal(t, "", Question("Name:", RequiredInput()))
}

func TestQuestion_Good_RequiredReprompts(t *testing.T) {
	SetStdin(strings.NewReader("\nalice\n"))
	defer SetStdin(nil)

	val := Question("Name:", RequiredInput())
	assert.Equal(t, "alice", val)
}

func TestChoose_Good_DefaultIndex(t *testing.T) {
	SetStdin(strings.NewReader("\n"))
	defer SetStdin(nil)

	val := Choose("Pick", []string{"a", "b", "c"}, WithDefaultIndex[string](1))
	assert.Equal(t, "b", val)
}

func TestChoose_Good_EmptyRepromptsWithoutDefault(t *testing.T) {
	SetStdin(strings.NewReader("\n2\n"))
	defer SetStdin(nil)

	val := Choose("Pick", []string{"a", "b", "c"})
	assert.Equal(t, "b", val)
}

func TestChoose_Bad_EOFWithoutDefaultReturnsZeroValue(t *testing.T) {
	SetStdin(strings.NewReader(""))
	defer SetStdin(nil)

	val := Choose("Pick", []string{"a", "b", "c"})
	assert.Empty(t, val)
}

func TestChooseMulti_Good_EmptyWithoutDefaultReturnsNone(t *testing.T) {
	SetStdin(strings.NewReader("\n"))
	defer SetStdin(nil)

	vals := ChooseMulti("Pick", []string{"a", "b", "c"})
	assert.Empty(t, vals)
}

func TestChoose_Good_Filter(t *testing.T) {
	SetStdin(strings.NewReader("ap\n2\n"))
	defer SetStdin(nil)

	val := Choose("Pick", []string{"apple", "apricot", "banana"}, Filter[string]())
	assert.Equal(t, "apricot", val)
}

func TestChoose_Bad_FilteredDefaultDoesNotFallBackToFirstVisible(t *testing.T) {
	SetStdin(strings.NewReader("ap\n\n2\n"))
	defer SetStdin(nil)

	val := Choose("Pick", []string{"apple", "banana", "apricot"}, WithDefaultIndex[string](1), Filter[string]())
	assert.Equal(t, "apricot", val)
}

func TestChoose_Bad_InvalidNumberUsesStderrHint(t *testing.T) {
	SetStdin(strings.NewReader("5\n2\n"))
	defer SetStdin(nil)

	var val string
	stderr := captureStderr(t, func() {
		val = Choose("Pick", []string{"a", "b"})
	})

	assert.Equal(t, "b", val)
	assert.Contains(t, stderr, "Please enter a number between 1 and 2.")
}

func TestChooseMulti_Good_Filter(t *testing.T) {
	SetStdin(strings.NewReader("ap\n1 2\n"))
	defer SetStdin(nil)

	vals := ChooseMulti("Pick", []string{"apple", "apricot", "banana"}, Filter[string]())
	assert.Equal(t, []string{"apple", "apricot"}, vals)
}

func TestChooseMulti_Bad_FilteredDefaultDoesNotFallBackToFirstVisible(t *testing.T) {
	SetStdin(strings.NewReader("ap\n\n2\n"))
	defer SetStdin(nil)

	vals := ChooseMulti("Pick", []string{"apple", "banana", "apricot"}, WithDefaultIndex[string](1), Filter[string]())
	assert.Equal(t, []string{"apricot"}, vals)
}

func TestChooseMulti_Good_Commas(t *testing.T) {
	SetStdin(strings.NewReader("1,3\n"))
	defer SetStdin(nil)

	vals := ChooseMulti("Pick", []string{"a", "b", "c"})
	assert.Equal(t, []string{"a", "c"}, vals)
}

func TestChooseMulti_Good_CommasAndRanges(t *testing.T) {
	SetStdin(strings.NewReader("1-2,4\n"))
	defer SetStdin(nil)

	vals := ChooseMulti("Pick", []string{"a", "b", "c", "d"})
	assert.Equal(t, []string{"a", "b", "d"}, vals)
}

func TestChooseMulti_Good_DefaultIndex(t *testing.T) {
	SetStdin(strings.NewReader("\n"))
	defer SetStdin(nil)

	vals := ChooseMulti("Pick", []string{"a", "b", "c"}, WithDefaultIndex[string](1))
	assert.Equal(t, []string{"b"}, vals)
}

func TestSetStdin_Good_ResetNil(t *testing.T) {
	original := stdin
	t.Cleanup(func() { stdin = original })

	SetStdin(strings.NewReader("hello\n"))
	assert.NotSame(t, os.Stdin, stdin)

	SetStdin(nil)
	assert.Same(t, os.Stdin, stdin)
}

func TestPromptHints_Good_UseStderr(t *testing.T) {
	oldOut := os.Stdout
	oldErr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	promptHint("try again")
	promptWarning("invalid")

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	_, _ = io.Copy(&stdout, rOut)
	_, _ = io.Copy(&stderr, rErr)

	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "try again")
	assert.Contains(t, stderr.String(), "invalid")
}

func TestPrompt_Good_WritesToStderr(t *testing.T) {
	SetStdin(strings.NewReader("hello\n"))
	defer SetStdin(nil)

	stdout, stderr := captureStdoutStderr(t, func() {
		val, err := Prompt("Name", "")
		assert.NoError(t, err)
		assert.Equal(t, "hello", val)
	})

	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "Name:")
}
