// Package cli provides semantic CLI output with zero external dependencies.
package cli

import (
	"time"

	"dappco.re/go/core"
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

func Pad(s string, width int) string {
	if displayWidth(s) >= width {
		return s
	}
	return s + Repeat(" ", width-displayWidth(s))
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
	var width int
	out := core.NewBuilder()
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

func FormatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return core.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return core.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return core.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 30*24*time.Hour:
		return core.Sprintf("%dw ago", int(d.Hours()/(24*7)))
	default:
		return core.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	}
}

type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderNormal
	BorderRounded
	BorderHeavy
	BorderDouble
)

type borderSet struct {
	tl, tr, bl, br string
	h, v           string
	tt, bt, lt, rt string
	x              string
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

type CellStyleFn func(value string) *AnsiStyle

type Table struct {
	Headers      []string
	Rows         [][]string
	Style        TableStyle
	borders      BorderStyle
	cellStyleFns map[int]CellStyleFn
	maxWidth     int
}

type TableStyle struct {
	HeaderStyle *AnsiStyle
	CellStyle   *AnsiStyle
	Separator   string
}

func DefaultTableStyle() TableStyle {
	return TableStyle{HeaderStyle: HeaderStyle, CellStyle: nil, Separator: "  "}
}

func NewTable(headers ...string) *Table {
	return &Table{Headers: headers, Style: DefaultTableStyle()}
}

func (t *Table) AddRow(cells ...string) *Table { t.Rows = append(t.Rows, cells); return t }

func (t *Table) WithBorders(style BorderStyle) *Table { t.borders = style; return t }

func (t *Table) WithCellStyle(col int, fn CellStyleFn) *Table {
	if t.cellStyleFns == nil {
		t.cellStyleFns = make(map[int]CellStyleFn)
	}
	t.cellStyleFns[col] = fn
	return t
}

func (t *Table) WithMaxWidth(w int) *Table { t.maxWidth = w; return t }

func (t *Table) String() string {
	if len(t.Headers) == 0 && len(t.Rows) == 0 {
		return ""
	}
	if t.borders != BorderNone {
		return t.renderBordered()
	}
	return t.renderPlain()
}

func (t *Table) Render() {
	writeString(stdoutWriter(), t.String())
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
		if w := displayWidth(compileGlyphs(h)); w > widths[i] {
			widths[i] = w
		}
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < cols {
				if w := displayWidth(compileGlyphs(cell)); w > widths[i] {
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
		overhead = (cols + 1) + (cols * 2)
	} else {
		overhead = (cols - 1) * len(t.Style.Separator)
	}
	total := overhead
	for _, w := range widths {
		total += w
	}
	if total <= t.maxWidth {
		return
	}
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
	sb := core.NewBuilder()
	sep := t.Style.Separator
	if len(t.Headers) > 0 {
		for i, h := range t.Headers {
			if i > 0 {
				sb.WriteString(sep)
			}
			cell := Pad(Truncate(compileGlyphs(h), widths[i]), widths[i])
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
			cell := Pad(Truncate(compileGlyphs(val), widths[i]), widths[i])
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
	sb := core.NewBuilder()
	sb.WriteString(b.tl)
	for i := range cols {
		sb.WriteString(Repeat(b.h, widths[i]+2))
		if i < cols-1 {
			sb.WriteString(b.tt)
		}
	}
	sb.WriteString(b.tr)
	sb.WriteByte('\n')
	if len(t.Headers) > 0 {
		sb.WriteString(b.v)
		for i := range cols {
			h := ""
			if i < len(t.Headers) {
				h = t.Headers[i]
			}
			cell := Pad(Truncate(compileGlyphs(h), widths[i]), widths[i])
			if t.Style.HeaderStyle != nil {
				cell = t.Style.HeaderStyle.Render(cell)
			}
			sb.WriteByte(' ')
			sb.WriteString(cell)
			sb.WriteByte(' ')
			sb.WriteString(b.v)
		}
		sb.WriteByte('\n')
		sb.WriteString(b.lt)
		for i := range cols {
			sb.WriteString(Repeat(b.h, widths[i]+2))
			if i < cols-1 {
				sb.WriteString(b.x)
			}
		}
		sb.WriteString(b.rt)
		sb.WriteByte('\n')
	}
	for _, row := range t.Rows {
		sb.WriteString(b.v)
		for i := range cols {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			cell := Pad(Truncate(compileGlyphs(val), widths[i]), widths[i])
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
	sb.WriteString(b.bl)
	for i := range cols {
		sb.WriteString(Repeat(b.h, widths[i]+2))
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
