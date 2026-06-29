package cli

import "dappco.re/go"

// ANSI escape codes
const (
	ansiReset     = "\033[0m"
	ansiBold      = "\033[1m"
	ansiDim       = "\033[2m"
	ansiItalic    = "\033[3m"
	ansiUnderline = "\033[4m"
)

var (
	colorEnabled        = true
	colorEnabledMu      core.RWMutex
	asciiDisabledColors bool
)

func init() {
	// NO_COLOR standard: https://no-color.org/
	// If NO_COLOR is set (to any value, including empty), disable colors.
	if _, exists := core.LookupEnv("NO_COLOR"); exists {
		colorEnabled = false
		return
	}

	// TERM=dumb indicates a terminal without color support.
	if core.Getenv("TERM") == "dumb" {
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
// This overrides the NO_COLOR environment variable check.
func SetColorEnabled(enabled bool) {
	colorEnabledMu.Lock()
	colorEnabled = enabled
	if enabled {
		asciiDisabledColors = false
	}
	colorEnabledMu.Unlock()
}

func restoreColorIfASCII() {
	colorEnabledMu.Lock()
	if asciiDisabledColors {
		colorEnabled = true
		asciiDisabledColors = false
	}
	colorEnabledMu.Unlock()
}

// AnsiStyle represents terminal text styling.
// Use NewStyle() to create, chain methods, call Render().
type AnsiStyle struct {
	bold      bool
	dim       bool
	italic    bool
	underline bool
	fg        string
	bg        string
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

// Italic enables italic text.
func (s *AnsiStyle) Italic() *AnsiStyle {
	s.italic = true
	return s
}

// Underline enables underlined text.
func (s *AnsiStyle) Underline() *AnsiStyle {
	s.underline = true
	return s
}

// Foreground sets foreground color from hex string.
func (s *AnsiStyle) Foreground(hex string) *AnsiStyle {
	s.fg = fgColorHex(hex)
	return s
}

// Background sets background color from hex string.
func (s *AnsiStyle) Background(hex string) *AnsiStyle {
	s.bg = bgColorHex(hex)
	return s
}

// Render applies the style to text.
// Returns plain text if NO_COLOR is set or colors are disabled.
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
	if s.italic {
		codes = append(codes, ansiItalic)
	}
	if s.underline {
		codes = append(codes, ansiUnderline)
	}
	if s.fg != "" {
		codes = append(codes, s.fg)
	}
	if s.bg != "" {
		codes = append(codes, s.bg)
	}

	if len(codes) == 0 {
		return text
	}

	return core.Join("", codes...) + text + ansiReset
}

// fgColorHex converts a hex string to an ANSI foreground color code.
func fgColorHex(hex string) string {
	r, g, b := hexToRGB(hex)
	return core.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// bgColorHex converts a hex string to an ANSI background color code.
func bgColorHex(hex string) string {
	r, g, b := hexToRGB(hex)
	return core.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

// stripANSI removes ANSI escape sequences from s, returning only the visible
// text. It is the zero-dependency replacement for charmbracelet/x/ansi.Strip:
// the CLI emits only SGR sequences (see AnsiStyle.Render), but the full CSI and
// OSC forms are handled so display-width and truncation stay correct for any
// styled input. Escape bytes are all ASCII, so scanning by byte is safe over
// UTF-8 content — multibyte runes (>= 0x80) are copied through untouched.
func stripANSI(s string) string {
	if !core.Contains(s, "\033") {
		return s
	}
	out := core.NewBuilder()
	for i := 0; i < len(s); {
		if s[i] != 0x1b { // not ESC — a visible byte
			out.WriteByte(s[i])
			i++
			continue
		}
		if i+1 >= len(s) {
			break // dangling ESC at end of string
		}
		switch s[i+1] {
		case '[': // CSI: ESC [ params… final-byte (0x40–0x7E)
			i += 2
			for i < len(s) && (s[i] < 0x40 || s[i] > 0x7e) {
				i++
			}
			if i < len(s) {
				i++ // consume the final byte
			}
		case ']': // OSC: ESC ] … terminated by BEL or ST (ESC \)
			i += 2
			for i < len(s) {
				if s[i] == 0x07 { // BEL
					i++
					break
				}
				if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '\\' { // ST
					i += 2
					break
				}
				i++
			}
		default: // other two-byte escape (ESC c, ESC M, …)
			i += 2
		}
	}
	return out.String()
}

// hexToRGB converts a hex string to RGB values.
func hexToRGB(hex string) (int, int, int) {
	hex = core.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}
	// Use 8-bit parsing since RGB values are 0-255, avoiding integer overflow on 32-bit systems.
	r := ParseHexByte(hex[0:2])
	g := ParseHexByte(hex[2:4])
	b := ParseHexByte(hex[4:6])
	if !r.OK || !g.OK || !b.OK {
		return 255, 255, 255
	}
	return r.Value.(int), g.Value.(int), b.Value.(int)
}
