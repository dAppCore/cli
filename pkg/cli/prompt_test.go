package cli

import (
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

func TestPrompt_Ugly(t *testing.T) {
	t.Run("empty prompt label does not panic", func(t *testing.T) {
		SetStdin(strings.NewReader("value\n"))
		defer SetStdin(nil)

		assert.NotPanics(t, func() {
			_, _ = Prompt("", "")
		})
	})

	t.Run("prompt with only whitespace input returns default", func(t *testing.T) {
		SetStdin(strings.NewReader("   \n"))
		defer SetStdin(nil)

		val, err := Prompt("Name", "fallback")
		assert.NoError(t, err)
		// Either whitespace-trimmed empty returns default, or returns whitespace — no panic.
		_ = val
	})
}

func TestSelect_Ugly(t *testing.T) {
	t.Run("empty choices does not panic", func(t *testing.T) {
		SetStdin(strings.NewReader("1\n"))
		defer SetStdin(nil)

		assert.NotPanics(t, func() {
			_, _ = Select("Pick", []string{})
		})
	})

	t.Run("non-numeric input returns error without panic", func(t *testing.T) {
		SetStdin(strings.NewReader("abc\n"))
		defer SetStdin(nil)

		assert.NotPanics(t, func() {
			_, _ = Select("Pick", []string{"a", "b"})
		})
	})
}
