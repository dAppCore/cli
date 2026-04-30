package cli

import (
	"bufio"
	"io"

	"dappco.re/go"
)

func newReader() *bufio.Reader {
	if br, ok := stdinReader().(*bufio.Reader); ok {
		return br
	}
	return bufio.NewReader(stdinReader())
}

// Prompt asks for text input with a default value.
func Prompt(label, defaultVal string) core.Result {
	label = compileGlyphs(label)
	defaultVal = compileGlyphs(defaultVal)
	if defaultVal != "" {
		if r := core.WriteString(stderrWriter(), core.Sprintf("%s [%s]: ", label, defaultVal)); !r.OK {
			return r
		}
	} else {
		if r := core.WriteString(stderrWriter(), core.Sprintf("%s: ", label)); !r.OK {
			return r
		}
	}

	r := newReader()
	input, err := r.ReadString('\n')
	input = core.Trim(input)
	if err != nil {
		if !core.Is(err, io.EOF) {
			return core.Fail(err)
		}
		if input == "" {
			if defaultVal != "" {
				return core.Ok(defaultVal)
			}
			return core.Fail(err)
		}
	}
	if input == "" {
		return core.Ok(defaultVal)
	}
	return core.Ok(input)
}

// Select presents numbered options and returns the selected value.
func Select(label string, options []string) core.Result {
	if len(options) == 0 {
		return core.Ok("")
	}

	core.Print(stderrWriter(), "%s", compileGlyphs(label))
	for i, opt := range options {
		core.Print(stderrWriter(), "  %d. %s", i+1, compileGlyphs(opt))
	}
	if r := core.WriteString(stderrWriter(), core.Sprintf("Choose [1-%d]: ", len(options))); !r.OK {
		return r
	}

	r := newReader()
	input, err := r.ReadString('\n')
	if err != nil && core.Trim(input) == "" {
		promptHint("No input received. Selection cancelled.")
		return Wrap(err, "selection cancelled")
	}

	trimmed := core.Trim(input)
	parsed := Atoi(trimmed)
	if !parsed.OK {
		promptHint(core.Sprintf("Please enter a number between 1 and %d.", len(options)))
		return Err("invalid selection %q: choose a number between 1 and %d", trimmed, len(options))
	}
	n := parsed.Value.(int)
	if n < 1 || n > len(options) {
		promptHint(core.Sprintf("Please enter a number between 1 and %d.", len(options)))
		return Err("invalid selection %q: choose a number between 1 and %d", trimmed, len(options))
	}
	return core.Ok(options[n-1])
}

// MultiSelect presents checkboxes (space-separated numbers).
func MultiSelect(label string, options []string) core.Result {
	if len(options) == 0 {
		return core.Ok([]string{})
	}

	core.Print(stderrWriter(), "%s", compileGlyphs(label))
	for i, opt := range options {
		core.Print(stderrWriter(), "  %d. %s", i+1, compileGlyphs(opt))
	}
	if r := core.WriteString(stderrWriter(), core.Sprintf("Choose (space-separated) [1-%d]: ", len(options))); !r.OK {
		return r
	}

	r := newReader()
	input, err := r.ReadString('\n')
	trimmed := core.Trim(input)
	if err != nil && trimmed == "" {
		return core.Ok([]string{})
	}
	if err != nil && !core.Is(err, io.EOF) {
		return core.Fail(err)
	}

	selectedResult := parseMultiSelection(trimmed, len(options))
	if !selectedResult.OK {
		return Wrap(selectedResult.Value.(error), core.Sprintf("invalid selection %q", trimmed))
	}
	selected := selectedResult.Value.([]int)

	selectedOptions := make([]string, 0, len(selected))
	for _, idx := range selected {
		selectedOptions = append(selectedOptions, options[idx])
	}
	return core.Ok(selectedOptions)
}
