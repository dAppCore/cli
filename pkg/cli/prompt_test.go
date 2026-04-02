package cli

import (
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
}

func TestMultiSelect_Good(t *testing.T) {
	SetStdin(strings.NewReader("1 3\n"))
	defer SetStdin(nil)

	vals, err := MultiSelect("Pick", []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "c"}, vals)
}

func TestConfirm_Good(t *testing.T) {
	SetStdin(strings.NewReader("y\n"))
	defer SetStdin(nil)

	assert.True(t, Confirm("Proceed?"))
}

func TestQuestion_Good(t *testing.T) {
	SetStdin(strings.NewReader("alice\n"))
	defer SetStdin(nil)

	val := Question("Name:")
	assert.Equal(t, "alice", val)
}

func TestChoose_Good_DefaultIndex(t *testing.T) {
	SetStdin(strings.NewReader("\n"))
	defer SetStdin(nil)

	val := Choose("Pick", []string{"a", "b", "c"}, WithDefaultIndex[string](1))
	assert.Equal(t, "b", val)
}

func TestSetStdin_Good_ResetNil(t *testing.T) {
	original := stdin
	t.Cleanup(func() { stdin = original })

	SetStdin(strings.NewReader("hello\n"))
	assert.NotSame(t, os.Stdin, stdin)

	SetStdin(nil)
	assert.Same(t, os.Stdin, stdin)
}
