// Package cli provides semantic CLI output with zero external dependencies.
package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

// Tailwind colour palette (hex strings)
const (
	ColourBlue50     = "#eff6ff"
	ColourBlue100    = "#dbeafe"
	ColourBlue200    = "#bfdbfe"
	ColourBlue300    = "#93c5fd"
	ColourBlue400    = "#60a5fa"
	ColourBlue500    = "#3b82f6"
	ColourBlue600    = "#2563eb"
	ColourBlue700    = "#1d4ed8"
	ColourGreen400   = "#4ade80"
	ColourGreen500   = "#22c55e"
	ColourGreen600   = "#16a34a"
	ColourRed400     = "#f87171"
	ColourRed500     = "#ef4444"
	ColourRed600     = "#dc2626"
	ColourAmber400   = "#fbbf24"
	ColourAmber500   = "#f59e0b"
	ColourAmber600   = "#d97706"
	ColourOrange500  = "#f97316"
	ColourYellow500  = "#eab308"
	ColourEmerald500 = "#10b981"
	ColourPurple500  = "#a855f7"
	ColourViolet400  = "#a78bfa"
	ColourViolet500  = "#8b5cf6"
	ColourIndigo500  = "#6366f1"
	ColourCyan500    = "#06b6d4"
	ColourGray50     = "#f9fafb"
	ColourGray100    = "#f3f4f6"
	ColourGray200    = "#e5e7eb"
	ColourGray300    = "#d1d5db"
	ColourGray400    = "#9ca3af"
	ColourGray500    = "#6b7280"
	ColourGray600    = "#4b5563"
	ColourGray700    = "#374151"
	ColourGray800    = "#1f2937"
	ColourGray900    = "#111827"
)

// Core styles
var (
	SuccessStyle  = NewStyle().Bold().Foreground(ColourGreen500)
	ErrorStyle    = NewStyle().Bold().Foreground(ColourRed500)
	WarningStyle  = NewStyle().Bold().Foreground(ColourAmber500)
	InfoStyle     = NewStyle().Foreground(ColourBlue400)
	SecurityStyle = NewStyle().Bold().Foreground(ColourPurple500)
	DimStyle      = NewStyle().Dim().Foreground(ColourGray500)
	MutedStyle    = NewStyle().Foreground(ColourGray600)
	BoldStyle     = NewStyle().Bold()
	KeyStyle      = NewStyle().Foreground(ColourGray400)
	ValueStyle    = NewStyle().Foreground(ColourGray200)
	AccentStyle   = NewStyle().Foreground(ColourCyan500)
	LinkStyle     = NewStyle().Foreground(ColourBlue500).Underline()
	HeaderStyle   = NewStyle().Bold().Foreground(ColourGray200)
	TitleStyle    = NewStyle().Bold().Foreground(ColourBlue500)
	CodeStyle     = NewStyle().Foreground(ColourGray300)
	NumberStyle   = NewStyle().Foreground(ColourBlue300)
	RepoStyle     = NewStyle().Bold().Foreground(ColourBlue500)
)

// Truncate shortens a string to max length with ellipsis.
func Truncate(s string, max int) string {
	if max <= 0 || s == "" {
		return ""
	}
	if displayWidth(s) <= max {
		return s
	}
	if max <= 3 {
		return truncateByWidth(s, max)
	}
	return truncateByWidth(s, max-3) + "..."
}

// Pad right-pads a string to width.
func Pad(s string, width int) string {
	if displayWidth(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-displayWidth(s))
}

func displayWidth(s string) int {
	return runewidth.StringWidth(ansi.Strip(s))
}

func truncateByWidth(s string, max int) string {
	if max <= 0 || s == "" {
		return ""
	}

	plain := ansi.Strip(s)
	if displayWidth(plain) <= max {
		return plain
	}

	var (
		width int
		out   strings.Builder
	)
	for _, r := range plain {
		rw := runewidth.RuneWidth(r)
		if width+rw > max {
			break
		}
		out.WriteRune(r)
		width += rw
	}
	return out.String()
}

// FormatAge formats a time as human-readable age (e.g., "2h ago", "3d ago").
func FormatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dw ago", int(d.Hours()/(24*7)))
	default:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Border Styles
// ─────────────────────────────────────────────────────────────────────────────

// BorderStyle selects the box-drawing character set for table borders.
type BorderStyle int

const (
	// BorderNone disables borders (default).
	BorderNone BorderStyle = iota
	// BorderNormal uses standard box-drawing: ┌─┬┐ │ ├─┼┤ └─┴┘
	BorderNormal
	// BorderRounded uses rounded corners: ╭─┬╮ │ ├─┼┤ ╰─┴╯
	BorderRounded
	// BorderHeavy uses heavy box-drawing: ┏━┳┓ ┃ ┣━╋┫ ┗━┻┛
	BorderHeavy
	// BorderDouble uses double-line box-drawing: ╔═╦╗ ║ ╠═╬╣ ╚═╩╝
	BorderDouble
)

type borderSet struct {
	tl, tr, bl, br string // corners
	h, v           string // horizontal, vertical
	tt, bt, lt, rt string // tees (top, bottom, left, right)
	x              string // cross
}

var borderSets = map[BorderStyle]borderSet{
	BorderNormal:  {"┌", "┐", "└", "┘", "─", "│", "┬", "┴", "├", "┤", "┼"},
	BorderRounded: {"╭", "╮", "╰", "╯", "─", "│", "┬", "┴", "├", "┤", "┼"},
	BorderHeavy:   {"┏", "┓", "┗", "┛", "━", "┃", "┳", "┻", "┣", "┫", "╋"},
	BorderDouble:  {"╔", "╗", "╚", "╝", "═", "║", "╦", "╩", "╠", "╣", "╬"},
}

var borderSetsASCII = map[BorderStyle]borderSet{
	BorderNormal:  {"+", "+", "+", "+", "-", "|", "+", "+", "+", "+", "+"},
	BorderRounded: {"+", "+", "+", "+", "-", "|", "+", "+", "+", "+", "+"},
	BorderHeavy:   {"+", "+", "+", "+", "=", "|", "+", "+", "+", "+", "+"},
	BorderDouble:  {"+", "+", "+", "+", "=", "|", "+", "+", "+", "+", "+"},
}

// CellStyleFn returns a style based on the cell's raw value.
// Return nil to use the table's default CellStyle.
type CellStyleFn func(value string) *AnsiStyle

// ─────────────────────────────────────────────────────────────────────────────
// Table
// ─────────────────────────────────────────────────────────────────────────────

// Table renders tabular data with aligned columns.
// Supports optional box-drawing borders and per-column cell styling.
//
//	t := cli.NewTable("REPO", "STATUS", "BRANCH").
//	    WithBorders(cli.BorderRounded).
//	    WithCellStyle(1, func(val string) *cli.AnsiStyle {
//	        if val == "clean" { return cli.SuccessStyle }
//	        return cli.WarningStyle
//	    })
//	t.AddRow("core-php", "clean", "main")
//	t.Render()
type Table struct {
	Headers      []string
	Rows         [][]string
	Style        TableStyle
	borders      BorderStyle
	cellStyleFns map[int]CellStyleFn
	maxWidth     int
}

// TableStyle configures the appearance of table output.
type TableStyle struct {
	HeaderStyle *AnsiStyle
	CellStyle   *AnsiStyle
	Separator   string
}

// DefaultTableStyle returns sensible defaults.
func DefaultTableStyle() TableStyle {
	return TableStyle{
		HeaderStyle: HeaderStyle,
		CellStyle:   nil,
		Separator:   "  ",
	}
}

// NewTable creates a table with headers.
func NewTable(headers ...string) *Table {
	return &Table{
		Headers: headers,
		Style:   DefaultTableStyle(),
	}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(cells ...string) *Table {
	t.Rows = append(t.Rows, cells)
	return t
}

// WithBorders enables box-drawing borders on the table.
func (t *Table) WithBorders(style BorderStyle) *Table {
	t.borders = style
	return t
}

// WithCellStyle sets a per-column style function.
// The function receives the raw cell value and returns a style.
func (t *Table) WithCellStyle(col int, fn CellStyleFn) *Table {
	if t.cellStyleFns == nil {
		t.cellStyleFns = make(map[int]CellStyleFn)
	}
	t.cellStyleFns[col] = fn
	return t
}

// WithMaxWidth sets the maximum table width, truncating columns to fit.
func (t *Table) WithMaxWidth(w int) *Table {
	t.maxWidth = w
	return t
}

// String renders the table.
func (t *Table) String() string {
	if len(t.Headers) == 0 && len(t.Rows) == 0 {
		return ""
	}

	if t.borders != BorderNone {
		return t.renderBordered()
	}
	return t.renderPlain()
}

// Render prints the table to stdout.
func (t *Table) Render() {
	fmt.Print(t.String())
}

func (t *Table) colCount() int {
	cols := len(t.Headers)
	if cols == 0 && len(t.Rows) > 0 {
		cols = len(t.Rows[0])
	}
	return cols
}

func (t *Table) columnWidths() []int {
	cols := t.colCount()
	widths := make([]int, cols)

	for i, h := range t.Headers {
		if w := displayWidth(h); w > widths[i] {
			widths[i] = w
		}
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < cols {
				if w := displayWidth(cell); w > widths[i] {
					widths[i] = w
				}
			}
		}
	}

	if t.maxWidth > 0 {
		t.constrainWidths(widths)
	}
	return widths
}

func (t *Table) constrainWidths(widths []int) {
	cols := len(widths)
	overhead := 0
	if t.borders != BorderNone {
		// │ cell │ cell │ = (cols+1) verticals + 2*cols padding spaces
		overhead = (cols + 1) + (cols * 2)
	} else {
		// separator between columns
		overhead = (cols - 1) * len(t.Style.Separator)
	}

	total := overhead
	for _, w := range widths {
		total += w
	}

	if total <= t.maxWidth {
		return
	}

	// Shrink widest columns first until we fit.
	budget := max(t.maxWidth-overhead, cols)
	for total-overhead > budget {
		maxIdx, maxW := 0, 0
		for i, w := range widths {
			if w > maxW {
				maxIdx, maxW = i, w
			}
		}
		widths[maxIdx]--
		total--
	}
}

func (t *Table) resolveStyle(col int, value string) *AnsiStyle {
	if t.cellStyleFns != nil {
		if fn, ok := t.cellStyleFns[col]; ok {
			if s := fn(value); s != nil {
				return s
			}
		}
	}
	return t.Style.CellStyle
}

func (t *Table) renderPlain() string {
	widths := t.columnWidths()

	var sb strings.Builder
	sep := t.Style.Separator

	if len(t.Headers) > 0 {
		for i, h := range t.Headers {
			if i > 0 {
				sb.WriteString(sep)
			}
			cell := Pad(Truncate(h, widths[i]), widths[i])
			if t.Style.HeaderStyle != nil {
				cell = t.Style.HeaderStyle.Render(cell)
			}
			sb.WriteString(cell)
		}
		sb.WriteByte('\n')
	}

	for _, row := range t.Rows {
		for i := range t.colCount() {
			if i > 0 {
				sb.WriteString(sep)
			}
			val := ""
			if i < len(row) {
				val = row[i]
			}
			cell := Pad(Truncate(val, widths[i]), widths[i])
			if style := t.resolveStyle(i, val); style != nil {
				cell = style.Render(cell)
			}
			sb.WriteString(cell)
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func (t *Table) renderBordered() string {
	b := tableBorderSet(t.borders)
	widths := t.columnWidths()
	cols := t.colCount()

	var sb strings.Builder

	// Top border: ╭──────┬──────╮
	sb.WriteString(b.tl)
	for i := range cols {
		sb.WriteString(strings.Repeat(b.h, widths[i]+2))
		if i < cols-1 {
			sb.WriteString(b.tt)
		}
	}
	sb.WriteString(b.tr)
	sb.WriteByte('\n')

	// Header row
	if len(t.Headers) > 0 {
		sb.WriteString(b.v)
		for i := range cols {
			h := ""
			if i < len(t.Headers) {
				h = t.Headers[i]
			}
			cell := Pad(Truncate(h, widths[i]), widths[i])
			if t.Style.HeaderStyle != nil {
				cell = t.Style.HeaderStyle.Render(cell)
			}
			sb.WriteByte(' ')
			sb.WriteString(cell)
			sb.WriteByte(' ')
			sb.WriteString(b.v)
		}
		sb.WriteByte('\n')

		// Header separator: ├──────┼──────┤
		sb.WriteString(b.lt)
		for i := range cols {
			sb.WriteString(strings.Repeat(b.h, widths[i]+2))
			if i < cols-1 {
				sb.WriteString(b.x)
			}
		}
		sb.WriteString(b.rt)
		sb.WriteByte('\n')
	}

	// Data rows
	for _, row := range t.Rows {
		sb.WriteString(b.v)
		for i := range cols {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			cell := Pad(Truncate(val, widths[i]), widths[i])
			if style := t.resolveStyle(i, val); style != nil {
				cell = style.Render(cell)
			}
			sb.WriteByte(' ')
			sb.WriteString(cell)
			sb.WriteByte(' ')
			sb.WriteString(b.v)
		}
		sb.WriteByte('\n')
	}

	// Bottom border: ╰──────┴──────╯
	sb.WriteString(b.bl)
	for i := range cols {
		sb.WriteString(strings.Repeat(b.h, widths[i]+2))
		if i < cols-1 {
			sb.WriteString(b.bt)
		}
	}
	sb.WriteString(b.br)
	sb.WriteByte('\n')

	return sb.String()
}

func tableBorderSet(style BorderStyle) borderSet {
	if currentTheme == ThemeASCII {
		if b, ok := borderSetsASCII[style]; ok {
			return b
		}
	}
	if b, ok := borderSets[style]; ok {
		return b
	}
	return borderSet{}
}
