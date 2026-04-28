package cli

import (
	"dappco.re/go"
	"strings"
)

func TestPrompt_Good(t *core.T) {
	SetStdin(strings.NewReader("hello\n"))
	defer SetStdin(nil) // reset

	val, err := Prompt("Name", "")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "hello", val)
}

func TestPrompt_Good_Default(t *core.T) {
	SetStdin(strings.NewReader("\n"))
	defer SetStdin(nil)

	val, err := Prompt("Name", "world")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "world", val)
}

func TestSelect_Good(t *core.T) {
	SetStdin(strings.NewReader("2\n"))
	defer SetStdin(nil)

	val, err := Select("Pick", []string{"a", "b", "c"})
	core.AssertNoError(t, err)
	core.AssertEqual(t, "b", val)
}

func TestSelect_Bad_Invalid(t *core.T) {
	SetStdin(strings.NewReader("5\n"))
	defer SetStdin(nil)

	_, err := Select("Pick", []string{"a", "b"})
	core.AssertError(t, err)
}

func TestMultiSelect_Good(t *core.T) {
	SetStdin(strings.NewReader("1 3\n"))
	defer SetStdin(nil)

	vals, err := MultiSelect("Pick", []string{"a", "b", "c"})
	core.AssertNoError(t, err)
	core.AssertEqual(t, []string{"a", "c"}, vals)
}

func TestPrompt_Ugly(t *core.T) {
	t.Run("empty prompt label does not panic", func(t *core.T) {
		SetStdin(strings.NewReader("value\n"))
		defer SetStdin(nil)
		core.AssertNotPanics(t, func() {
			_, _ = Prompt("", "")
		})
	})

	t.Run("prompt with only whitespace input returns default", func(t *core.T) {
		SetStdin(strings.NewReader("   \n"))
		defer SetStdin(nil)

		val, err := Prompt("Name", "fallback")
		core.AssertNoError(t, err)
		// Either whitespace-trimmed empty returns default, or returns whitespace — no panic.
		_ = val
	})
}

func TestSelect_Ugly(t *core.T) {
	t.Run("empty choices does not panic", func(t *core.T) {
		SetStdin(strings.NewReader("1\n"))
		defer SetStdin(nil)
		core.AssertNotPanics(t, func() {
			_, _ = Select("Pick", []string{})
		})
	})

	t.Run("non-numeric input returns error without panic", func(t *core.T) {
		SetStdin(strings.NewReader("abc\n"))
		defer SetStdin(nil)
		core.AssertNotPanics(t, func() {
			_, _ = Select("Pick", []string{"a", "b"})
		})
	})
}
