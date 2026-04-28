package cli

import (
	"strings"

	"dappco.re/go"
	"unicode/utf8"
)

func TestTable_Good(t *core.T) {
	t.Run("plain table unchanged", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("NAME", "AGE")
		tbl.AddRow("Alice", "30")
		tbl.AddRow("Bob", "25")

		out := tbl.String()
		core.AssertContains(t, out, "NAME")
		core.AssertContains(t, out, "Alice")
		core.AssertContains(t, out, "Bob")
	})

	t.Run("bordered normal", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B").WithBorders(BorderNormal)
		tbl.AddRow("x", "y")

		out := tbl.String()
		core.AssertTrue(t, strings.HasPrefix(out, "┌"))
		core.AssertContains(t, out, "┐")
		core.AssertContains(t, out, "│")
		core.AssertContains(t, out, "├")
		core.AssertContains(t, out, "┤")
		core.AssertContains(t, out, "└")
		core.AssertContains(t, out, "┘")
	})

	t.Run("bordered rounded", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("REPO", "STATUS").WithBorders(BorderRounded)
		tbl.AddRow("core", "clean")

		out := tbl.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		core.AssertTrue(t, strings.HasPrefix(lines[0], "╭"))
		core.AssertTrue(t, strings.HasSuffix(lines[0], "╮"))
		core.AssertTrue(t, strings.HasPrefix(lines[len(lines)-1], "╰"))
		core.AssertTrue(t, strings.HasSuffix(lines[len(lines)-1], "╯"))
	})

	t.Run("bordered heavy", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("X").WithBorders(BorderHeavy)
		tbl.AddRow("v")

		out := tbl.String()
		core.AssertContains(t, out, "┏")
		core.AssertContains(t, out, "┓")
		core.AssertContains(t, out, "┃")
	})

	t.Run("bordered double", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("X").WithBorders(BorderDouble)
		tbl.AddRow("v")

		out := tbl.String()
		core.AssertContains(t, out, "╔")
		core.AssertContains(t, out, "╗")
		core.AssertContains(t, out, "║")
	})

	t.Run("ASCII theme uses ASCII borders", func(t *core.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tbl := NewTable("REPO", "STATUS").WithBorders(BorderRounded)
		tbl.AddRow("core", "clean")

		out := tbl.String()
		core.AssertContains(t, out, "+")
		core.AssertContains(t, out, "-")
		core.AssertContains(t, out, "|")
		core.AssertNotContains(t, out, "╭")
		core.AssertNotContains(t, out, "╮")
		core.AssertNotContains(t, out, "│")
	})

	t.Run("bordered structure", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B").WithBorders(BorderRounded)
		tbl.AddRow("x", "y")
		tbl.AddRow("1", "2")

		lines := strings.Split(strings.TrimRight(tbl.String(), "\n"), "\n")
		core.
			// Top border, header, separator, 2 data rows, bottom border = 6 lines
			AssertEqual(t, 6, len(lines), "expected 6 lines: border, header, sep, 2 rows, border")
	})

	t.Run("cell style function", func(t *core.T) {
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
		core.AssertTrue(t, called, "cell style function should be called")
	})

	t.Run("cell style with borders", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("NAME", "STATUS").
			WithBorders(BorderRounded).
			WithCellStyle(1, func(val string) *AnsiStyle {
				return nil // fallback to default
			})
		tbl.AddRow("core", "ok")

		out := tbl.String()
		core.AssertContains(t, out, "core")
		core.AssertContains(t, out, "ok")
	})

	t.Run("glyph shortcodes render in headers and cells", func(t *core.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tbl := NewTable(":check: NAME", "STATUS").
			WithBorders(BorderRounded)
		tbl.AddRow("core", ":warn:")

		out := tbl.String()
		core.AssertContains(t, out, "[OK] NAME")
		core.AssertContains(t, out, "[WARN]")
	})

	t.Run("max width truncates", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("LONG_HEADER", "SHORT").WithMaxWidth(25)
		tbl.AddRow("very_long_value_here", "x")

		out := tbl.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		for _, line := range lines {
			w := utf8.RuneCountInString(line)
			core.AssertLessOrEqual(t, w, 25, core.Sprintf("line should not exceed max width: %q", line))
		}
	})

	t.Run("max width with borders", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B").WithBorders(BorderNormal).WithMaxWidth(20)
		tbl.AddRow("hello", "world")

		out := tbl.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		for _, line := range lines {
			w := utf8.RuneCountInString(line)
			core.AssertLessOrEqual(t, w, 20, core.Sprintf("bordered line should not exceed max width: %q", line))
		}
	})

	t.Run("empty table returns empty", func(t *core.T) {
		tbl := NewTable()
		core.AssertEqual(t, "", tbl.String())
	})

	t.Run("no headers with borders", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable().WithBorders(BorderNormal)
		tbl.Rows = [][]string{{"a", "b"}, {"c", "d"}}

		out := tbl.String()
		core.AssertContains(t, out, "┌")
		// No header separator since no headers
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		core.
			// Top border, 2 data rows, bottom border = 4 lines (no header separator)
			AssertEqual(t, 4, len(lines))
	})
}

func TestTable_Bad(t *core.T) {
	t.Run("short rows padded", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A", "B", "C")
		tbl.AddRow("x") // only 1 cell, 3 columns

		out := tbl.String()
		core.AssertContains(t, out, "x")
	})
}

func TestTable_Ugly(t *core.T) {
	t.Run("no columns no panic", func(t *core.T) {
		core.AssertNotPanics(t, func() {
			tbl := NewTable()
			tbl.AddRow()
			_ = tbl.String()
		})
	})

	t.Run("cell style function returning nil does not panic", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("A").WithCellStyle(0, func(_ string) *AnsiStyle {
			return nil
		})
		tbl.AddRow("value")
		core.AssertNotPanics(t, func() {
			_ = tbl.String()
		})
	})

	t.Run("max width of 1 does not panic", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tbl := NewTable("HEADER").WithMaxWidth(1)
		tbl.AddRow("data")
		core.AssertNotPanics(t, func() {
			_ = tbl.String()
		})
	})
}

func TestTruncate_Good(t *core.T) {
	core.AssertEqual(t, "hel...", Truncate("hello world", 6))
	core.AssertEqual(t, "hi", Truncate("hi", 6))
	core.AssertEqual(t, "he", Truncate("hello", 2))
	core.AssertEqual(t, "東", Truncate("東京", 3))
}

func TestTruncate_Ugly(t *core.T) {
	t.Run("zero max does not panic", func(t *core.T) {
		core.AssertNotPanics(t, func() {
			_ = Truncate("hello", 0)
		})
	})
}

func TestPad_Good(t *core.T) {
	core.AssertEqual(t, "hi   ", Pad("hi", 5))
	core.AssertEqual(t, "hello", Pad("hello", 3))
	core.AssertEqual(t, "東京  ", Pad("東京", 6))
}

func TestStyled_Good_NilStyle(t *core.T) {
	restoreThemeAndColors(t)
	UseASCII()
	core.AssertEqual(t, "hello [OK]", Styled(nil, "hello :check:"))
}

func TestStyledf_Good_NilStyle(t *core.T) {
	restoreThemeAndColors(t)
	UseASCII()
	core.AssertEqual(t, "value: [WARN]", Styledf(nil, "value: %s", ":warn:"))
}

func TestPad_Ugly(t *core.T) {
	t.Run("zero width does not panic", func(t *core.T) {
		core.AssertNotPanics(t, func() {
			_ = Pad("hello", 0)
		})
	})
}
