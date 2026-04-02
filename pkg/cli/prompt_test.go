package cli

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	_, err := Select("Pick", []string{"a", "b"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "choose a number between 1 and 2")
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

func TestChooseMulti_Good_Filter(t *testing.T) {
	SetStdin(strings.NewReader("ap\n1 2\n"))
	defer SetStdin(nil)

	vals := ChooseMulti("Pick", []string{"apple", "apricot", "banana"}, Filter[string]())
	assert.Equal(t, []string{"apple", "apricot"}, vals)
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
