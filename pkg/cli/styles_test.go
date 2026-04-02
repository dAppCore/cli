package cli

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestTable_Good(t *testing.T) {
	t.Run("plain table unchanged", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("NAME", "AGE")
		tbl.AddRow("Alice", "30")
		tbl.AddRow("Bob", "25")

		out := tbl.String()
		assert.Contains(t, out, "NAME")
		assert.Contains(t, out, "Alice")
		assert.Contains(t, out, "Bob")
	})

	t.Run("bordered normal", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B").WithBorders(BorderNormal)
		tbl.AddRow("x", "y")

		out := tbl.String()
		assert.True(t, strings.HasPrefix(out, "┌"))
		assert.Contains(t, out, "┐")
		assert.Contains(t, out, "│")
		assert.Contains(t, out, "├")
		assert.Contains(t, out, "┤")
		assert.Contains(t, out, "└")
		assert.Contains(t, out, "┘")
	})

	t.Run("bordered rounded", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("REPO", "STATUS").WithBorders(BorderRounded)
		tbl.AddRow("core", "clean")

		out := tbl.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		assert.True(t, strings.HasPrefix(lines[0], "╭"))
		assert.True(t, strings.HasSuffix(lines[0], "╮"))
		assert.True(t, strings.HasPrefix(lines[len(lines)-1], "╰"))
		assert.True(t, strings.HasSuffix(lines[len(lines)-1], "╯"))
	})

	t.Run("bordered heavy", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("X").WithBorders(BorderHeavy)
		tbl.AddRow("v")

		out := tbl.String()
		assert.Contains(t, out, "┏")
		assert.Contains(t, out, "┓")
		assert.Contains(t, out, "┃")
	})

	t.Run("bordered double", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("X").WithBorders(BorderDouble)
		tbl.AddRow("v")

		out := tbl.String()
		assert.Contains(t, out, "╔")
		assert.Contains(t, out, "╗")
		assert.Contains(t, out, "║")
	})

	t.Run("ASCII theme uses ASCII borders", func(t *testing.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tbl := NewTable("REPO", "STATUS").WithBorders(BorderRounded)
		tbl.AddRow("core", "clean")

		out := tbl.String()
		assert.Contains(t, out, "+")
		assert.Contains(t, out, "-")
		assert.Contains(t, out, "|")
		assert.NotContains(t, out, "╭")
		assert.NotContains(t, out, "╮")
		assert.NotContains(t, out, "│")
	})

	t.Run("bordered structure", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B").WithBorders(BorderRounded)
		tbl.AddRow("x", "y")
		tbl.AddRow("1", "2")

		lines := strings.Split(strings.TrimRight(tbl.String(), "\n"), "\n")
		// Top border, header, separator, 2 data rows, bottom border = 6 lines
		assert.Equal(t, 6, len(lines), "expected 6 lines: border, header, sep, 2 rows, border")
	})

	t.Run("cell style function", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		called := false
		tbl := NewTable("STATUS").
			WithCellStyle(0, func(val string) *AnsiStyle {
				called = true
				if val == "ok" {
					return SuccessStyle
				}
				return ErrorStyle
			})
		tbl.AddRow("ok")
		tbl.AddRow("fail")

		_ = tbl.String()
		assert.True(t, called, "cell style function should be called")
	})

	t.Run("cell style with borders", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("NAME", "STATUS").
			WithBorders(BorderRounded).
			WithCellStyle(1, func(val string) *AnsiStyle {
				return nil // fallback to default
			})
		tbl.AddRow("core", "ok")

		out := tbl.String()
		assert.Contains(t, out, "core")
		assert.Contains(t, out, "ok")
	})

	t.Run("max width truncates", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("LONG_HEADER", "SHORT").WithMaxWidth(25)
		tbl.AddRow("very_long_value_here", "x")

		out := tbl.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		for _, line := range lines {
			w := utf8.RuneCountInString(line)
			assert.LessOrEqual(t, w, 25, "line should not exceed max width: %q", line)
		}
	})

	t.Run("max width with borders", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B").WithBorders(BorderNormal).WithMaxWidth(20)
		tbl.AddRow("hello", "world")

		out := tbl.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		for _, line := range lines {
			w := utf8.RuneCountInString(line)
			assert.LessOrEqual(t, w, 20, "bordered line should not exceed max width: %q", line)
		}
	})

	t.Run("empty table returns empty", func(t *testing.T) {
		tbl := NewTable()
		assert.Equal(t, "", tbl.String())
	})

	t.Run("no headers with borders", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable().WithBorders(BorderNormal)
		tbl.Rows = [][]string{{"a", "b"}, {"c", "d"}}

		out := tbl.String()
		assert.Contains(t, out, "┌")
		// No header separator since no headers
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		// Top border, 2 data rows, bottom border = 4 lines (no header separator)
		assert.Equal(t, 4, len(lines))
	})
}

func TestTable_Bad(t *testing.T) {
	t.Run("short rows padded", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B", "C")
		tbl.AddRow("x") // only 1 cell, 3 columns

		out := tbl.String()
		assert.Contains(t, out, "x")
	})
}

func TestTruncate_Good(t *testing.T) {
	assert.Equal(t, "hel...", Truncate("hello world", 6))
	assert.Equal(t, "hi", Truncate("hi", 6))
	assert.Equal(t, "he", Truncate("hello", 2))
	assert.Equal(t, "東", Truncate("東京", 3))
}

func TestPad_Good(t *testing.T) {
	assert.Equal(t, "hi   ", Pad("hi", 5))
	assert.Equal(t, "hello", Pad("hello", 3))
	assert.Equal(t, "東京  ", Pad("東京", 6))
}
