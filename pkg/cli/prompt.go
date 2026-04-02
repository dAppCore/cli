package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var stdin io.Reader = os.Stdin

// SetStdin overrides the default stdin reader for testing.
// Pass nil to restore the real os.Stdin reader.
func SetStdin(r io.Reader) {
	if r == nil {
		stdin = os.Stdin
		return
	}
	stdin = r
}

// newReader wraps stdin in a bufio.Reader if it isn't one already.
func newReader() *bufio.Reader {
	if br, ok := stdin.(*bufio.Reader); ok {
		return br
	}
	return bufio.NewReader(stdin)
}

// Prompt asks for text input with a default value.
func Prompt(label, defaultVal string) (string, error) {
	label = compileGlyphs(label)
	defaultVal = compileGlyphs(defaultVal)
	if defaultVal != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, defaultVal)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}

	r := newReader()
	input, err := r.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil {
		if !errors.Is(err, io.EOF) {
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

	fmt.Fprintln(os.Stderr, compileGlyphs(label))
	for i, opt := range options {
		fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, compileGlyphs(opt))
	}
	fmt.Fprintf(os.Stderr, "Choose [1-%d]: ", len(options))

	r := newReader()
	input, err := r.ReadString('\n')
	if err != nil && strings.TrimSpace(input) == "" {
		promptHint("No input received. Selection cancelled.")
		return "", Wrap(err, "selection cancelled")
	}

	trimmed := strings.TrimSpace(input)
	n, err := strconv.Atoi(trimmed)
	if err != nil || n < 1 || n > len(options) {
		promptHint(fmt.Sprintf("Please enter a number between 1 and %d.", len(options)))
		return "", Err("invalid selection %q: choose a number between 1 and %d", trimmed, len(options))
	}
	return options[n-1], nil
}

// MultiSelect presents checkboxes (space-separated numbers).
func MultiSelect(label string, options []string) ([]string, error) {
	if len(options) == 0 {
		return []string{}, nil
	}

	fmt.Fprintln(os.Stderr, compileGlyphs(label))
	for i, opt := range options {
		fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, compileGlyphs(opt))
	}
	fmt.Fprintf(os.Stderr, "Choose (space-separated) [1-%d]: ", len(options))

	r := newReader()
	input, err := r.ReadString('\n')
	trimmed := strings.TrimSpace(input)
	if err != nil && trimmed == "" {
		return []string{}, nil
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	selected, parseErr := parseMultiSelection(trimmed, len(options))
	if parseErr != nil {
		return nil, Wrap(parseErr, fmt.Sprintf("invalid selection %q", trimmed))
	}

	selectedOptions := make([]string, 0, len(selected))
	for _, idx := range selected {
		selectedOptions = append(selectedOptions, options[idx])
	}
	return selectedOptions, nil
}
