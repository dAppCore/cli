package cli

import (
	core "dappco.re/go"
	"time"
)

func TestStyles_Pad_Good(t *core.T) {
	got := Pad("go", 4)

	core.AssertEqual(t, "go  ", got)
	core.AssertLen(t, got, 4)
}

func TestStyles_Pad_Bad(t *core.T) {
	got := Pad("long", 2)

	core.AssertEqual(t, "long", got)
	core.AssertLen(t, got, 4)
}

func TestStyles_Pad_Ugly(t *core.T) {
	got := Pad("", 3)

	core.AssertEqual(t, "   ", got)
	core.AssertLen(t, got, 3)
}

func TestStyles_FormatAge_Good(t *core.T) {
	got := FormatAge(time.Now().Add(-2 * time.Minute))

	core.AssertContains(t, got, "m ago")
	core.AssertNotEqual(t, "just now", got)
}

func TestStyles_FormatAge_Bad(t *core.T) {
	got := FormatAge(time.Now().Add(time.Minute))

	core.AssertEqual(t, "just now", got)
	core.AssertNotEmpty(t, got)
}

func TestStyles_FormatAge_Ugly(t *core.T) {
	got := FormatAge(time.Now().Add(-45 * 24 * time.Hour))

	core.AssertContains(t, got, "mo ago")
	core.AssertNotEmpty(t, got)
}

func TestStyles_DefaultTableStyle_Good(t *core.T) {
	style := DefaultTableStyle()

	core.AssertNotNil(t, style.HeaderStyle)
	core.AssertEqual(t, "  ", style.Separator)
}

func TestStyles_DefaultTableStyle_Bad(t *core.T) {
	style := DefaultTableStyle()

	core.AssertNil(t, style.CellStyle)
	core.AssertNotNil(t, style.HeaderStyle)
}

func TestStyles_DefaultTableStyle_Ugly(t *core.T) {
	style := DefaultTableStyle()
	style.Separator = "|"

	core.AssertEqual(t, "|", style.Separator)
	core.AssertEqual(t, "  ", DefaultTableStyle().Separator)
}

func TestStyles_NewTable_Good(t *core.T) {
	table := NewTable("Name", "Status")

	core.AssertEqual(t, []string{"Name", "Status"}, table.Headers)
	core.AssertNotNil(t, table.Style.HeaderStyle)
}

func TestStyles_NewTable_Bad(t *core.T) {
	table := NewTable()

	core.AssertEmpty(t, table.Headers)
	core.AssertEqual(t, "", table.String())
}

func TestStyles_NewTable_Ugly(t *core.T) {
	table := NewTable(":check:")

	core.AssertEqual(t, []string{":check:"}, table.Headers)
	core.AssertContains(t, table.String(), "✓")
}

func TestStyles_Table_AddRow_Good(t *core.T) {
	table := NewTable("Name").AddRow("codex")

	core.AssertLen(t, table.Rows, 1)
	core.AssertEqual(t, []string{"codex"}, table.Rows[0])
}

func TestStyles_Table_AddRow_Bad(t *core.T) {
	table := NewTable("Name").AddRow()

	core.AssertLen(t, table.Rows, 1)
	core.AssertEmpty(t, table.Rows[0])
}

func TestStyles_Table_AddRow_Ugly(t *core.T) {
	table := NewTable().AddRow("orphan")

	core.AssertContains(t, table.String(), "orphan")
	core.AssertLen(t, table.Rows, 1)
}

func TestStyles_Table_WithBorders_Good(t *core.T) {
	table := NewTable("Name").WithBorders(BorderRounded)

	core.AssertEqual(t, BorderRounded, table.borders)
	core.AssertContains(t, table.String(), "╭")
}

func TestStyles_Table_WithBorders_Bad(t *core.T) {
	table := NewTable("Name").WithBorders(BorderNone)

	core.AssertEqual(t, BorderNone, table.borders)
	core.AssertNotContains(t, table.String(), "╭")
}

func TestStyles_Table_WithBorders_Ugly(t *core.T) {
	cliPlainCLI(t)
	table := NewTable("Name").WithBorders(BorderHeavy)

	core.AssertEqual(t, BorderHeavy, table.borders)
	core.AssertContains(t, table.String(), "+")
}

func TestStyles_Table_WithCellStyle_Good(t *core.T) {
	table := NewTable("Name").WithCellStyle(0, func(string) *AnsiStyle { return NewStyle().Bold() })

	core.AssertNotNil(t, table.cellStyleFns[0])
	core.AssertEqual(t, table, table.WithCellStyle(1, nil))
}

func TestStyles_Table_WithCellStyle_Bad(t *core.T) {
	table := NewTable("Name").WithCellStyle(-1, nil)

	core.AssertNotNil(t, table.cellStyleFns)
	core.AssertNil(t, table.cellStyleFns[-1])
}

func TestStyles_Table_WithCellStyle_Ugly(t *core.T) {
	table := NewTable("Name").WithCellStyle(0, func(value string) *AnsiStyle {
		if value == "hot" {
			return NewStyle().Bold()
		}
		return nil
	})

	core.AssertNotNil(t, table.cellStyleFns[0]("hot"))
}

func TestStyles_Table_WithMaxWidth_Good(t *core.T) {
	table := NewTable("Name").WithMaxWidth(10)

	core.AssertEqual(t, 10, table.maxWidth)
	core.AssertEqual(t, table, table.WithMaxWidth(20))
}

func TestStyles_Table_WithMaxWidth_Bad(t *core.T) {
	table := NewTable("Name").WithMaxWidth(0)

	core.AssertEqual(t, 0, table.maxWidth)
	core.AssertContains(t, table.String(), "Name")
}

func TestStyles_Table_WithMaxWidth_Ugly(t *core.T) {
	table := NewTable("Name").AddRow("abcdef").WithMaxWidth(5)

	core.AssertContains(t, table.String(), "...")
	core.AssertEqual(t, 5, table.maxWidth)
}

func TestStyles_Table_String_Good(t *core.T) {
	got := NewTable("Name").AddRow("codex").String()

	core.AssertContains(t, got, "Name")
	core.AssertContains(t, got, "codex")
}

func TestStyles_Table_String_Bad(t *core.T) {
	got := NewTable().String()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStyles_Table_String_Ugly(t *core.T) {
	got := NewTable("Name").WithBorders(BorderDouble).String()

	core.AssertContains(t, got, "Name")
	core.AssertContains(t, got, "╔")
}

func TestStyles_Table_Render_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { NewTable("Name").AddRow("codex").Render() })

	core.AssertContains(t, out, "Name")
	core.AssertContains(t, out, "codex")
}

func TestStyles_Table_Render_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { NewTable().Render() })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestStyles_Table_Render_Ugly(t *core.T) {
	out := cliCaptureStdout(t, func() { NewTable("Name").WithBorders(BorderNormal).Render() })

	core.AssertContains(t, out, "Name")
	core.AssertContains(t, out, "┌")
}

func TestStyles_Truncate_Good(t *core.T) {
	got := Truncate("abcdef", 4)

	core.AssertEqual(t, "a...", got)
	core.AssertLen(t, got, 4)
}

func TestStyles_Truncate_Bad(t *core.T) {
	got := Truncate("abcdef", 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStyles_Truncate_Ugly(t *core.T) {
	got := Truncate("go", 10)

	core.AssertEqual(t, "go", got)
	core.AssertNotContains(t, got, "...")
}
