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
func Prompt(label, defaultVal string) (string, error) {
	label = compileGlyphs(label)
	defaultVal = compileGlyphs(defaultVal)
	if defaultVal != "" {
		io.WriteString(stderrWriter(), core.Sprintf("%s [%s]: ", label, defaultVal))
	} else {
		io.WriteString(stderrWriter(), core.Sprintf("%s: ", label))
	}

	r := newReader()
	input, err := r.ReadString('\n')
	input = core.Trim(input)
	if err != nil {
		if !core.Is(err, io.EOF) {
			return "", err
		}
		if input == "" {
			if defaultVal != "" {
				return defaultVal, nil
			}
			return "", err
		}
	}
	if input == "" {
		return defaultVal, nil
	}
	return input, nil
}

// Select presents numbered options and returns the selected value.
func Select(label string, options []string) (string, error) {
	if len(options) == 0 {
		return "", nil
	}

	core.Print(stderrWriter(), "%s", compileGlyphs(label))
	for i, opt := range options {
		core.Print(stderrWriter(), "  %d. %s", i+1, compileGlyphs(opt))
	}
	io.WriteString(stderrWriter(), core.Sprintf("Choose [1-%d]: ", len(options)))

	r := newReader()
	input, err := r.ReadString('\n')
	if err != nil && core.Trim(input) == "" {
		promptHint("No input received. Selection cancelled.")
		return "", Wrap(err, "selection cancelled")
	}

	trimmed := core.Trim(input)
	n, err := Atoi(trimmed)
	if err != nil || n < 1 || n > len(options) {
		promptHint(core.Sprintf("Please enter a number between 1 and %d.", len(options)))
		return "", Err("invalid selection %q: choose a number between 1 and %d", trimmed, len(options))
	}
	return options[n-1], nil
}

// MultiSelect presents checkboxes (space-separated numbers).
func MultiSelect(label string, options []string) ([]string, error) {
	if len(options) == 0 {
		return []string{}, nil
	}

	core.Print(stderrWriter(), "%s", compileGlyphs(label))
	for i, opt := range options {
		core.Print(stderrWriter(), "  %d. %s", i+1, compileGlyphs(opt))
	}
	io.WriteString(stderrWriter(), core.Sprintf("Choose (space-separated) [1-%d]: ", len(options)))

	r := newReader()
	input, err := r.ReadString('\n')
	trimmed := core.Trim(input)
	if err != nil && trimmed == "" {
		return []string{}, nil
	}
	if err != nil && !core.Is(err, io.EOF) {
		return nil, err
	}

	selected, parseErr := parseMultiSelection(trimmed, len(options))
	if parseErr != nil {
		return nil, Wrap(parseErr, core.Sprintf("invalid selection %q", trimmed))
	}

	selectedOptions := make([]string, 0, len(selected))
	for _, idx := range selected {
		selectedOptions = append(selectedOptions, options[idx])
	}
	return selectedOptions, nil
}
