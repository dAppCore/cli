package frame

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiDim   = "\033[2m"

	ColourGray500 = "#6b7280"
)

var (
	colorEnabled   = true
	colorEnabledMu sync.RWMutex

	DimStyle  = NewStyle().Dim().Foreground(ColourGray500)
	BoldStyle = NewStyle().Bold()
)

func init() {
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		colorEnabled = false
		return
	}
	if os.Getenv("TERM") == "dumb" {
		colorEnabled = false
	}
}

// ColorEnabled returns true if ANSI color output is enabled.
func ColorEnabled() bool {
	colorEnabledMu.RLock()
	defer colorEnabledMu.RUnlock()
	return colorEnabled
}

// SetColorEnabled enables or disables ANSI color output.
func SetColorEnabled(enabled bool) {
	colorEnabledMu.Lock()
	colorEnabled = enabled
	colorEnabledMu.Unlock()
}

// AnsiStyle represents terminal text styling.
type AnsiStyle struct {
	bold bool
	dim  bool
	fg   string
}

// NewStyle creates a new empty style.
func NewStyle() *AnsiStyle {
	return &AnsiStyle{}
}

// Bold enables bold text.
func (s *AnsiStyle) Bold() *AnsiStyle {
	s.bold = true
	return s
}

// Dim enables dim text.
func (s *AnsiStyle) Dim() *AnsiStyle {
	s.dim = true
	return s
}

// Foreground sets foreground color from hex string.
func (s *AnsiStyle) Foreground(hex string) *AnsiStyle {
	s.fg = fgColorHex(hex)
	return s
}

// Render applies the style to text.
func (s *AnsiStyle) Render(text string) string {
	if s == nil || !ColorEnabled() {
		return text
	}

	var codes []string
	if s.bold {
		codes = append(codes, ansiBold)
	}
	if s.dim {
		codes = append(codes, ansiDim)
	}
	if s.fg != "" {
		codes = append(codes, s.fg)
	}
	if len(codes) == 0 {
		return text
	}

	return strings.Join(codes, "") + text + ansiReset
}

func fgColorHex(hex string) string {
	r, g, b := hexToRGB(hex)
	return "\033[38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return int(r), int(g), int(b)
}

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
	var out strings.Builder
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

func Glyph(code string) string {
	switch code {
	case ":check:":
		return "✓"
	case ":cross:":
		return "✗"
	case ":warn:":
		return "⚠"
	case ":info:":
		return "ℹ"
	case ":dash:":
		return "─"
	default:
		return code
	}
}

func compileGlyphs(x string) string {
	if x == "" {
		return ""
	}
	input := bytes.NewBufferString(x)
	output := bytes.NewBufferString("")

	for {
		r, _, err := input.ReadRune()
		if err != nil {
			break
		}
		if r == ':' {
			output.WriteString(replaceGlyph(input))
		} else {
			output.WriteRune(r)
		}
	}
	return output.String()
}

func replaceGlyph(input *bytes.Buffer) string {
	code := bytes.NewBufferString(":")
	for {
		r, _, err := input.ReadRune()
		if err != nil {
			return code.String()
		}
		if r == ':' && code.Len() == 1 {
			return code.String() + replaceGlyph(input)
		}
		code.WriteRune(r)
		if unicode.IsSpace(r) {
			return code.String()
		}
		if r == ':' {
			return Glyph(code.String())
		}
	}
}
