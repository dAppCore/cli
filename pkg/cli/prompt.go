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
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}

	r := newReader()
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal, nil
	}
	return input, nil
}

// Select presents numbered options and returns the selected value.
func Select(label string, options []string) (string, error) {
	fmt.Println(label)
	for i, opt := range options {
		fmt.Printf("  %d. %s\n", i+1, opt)
	}
	fmt.Printf("Choose [1-%d]: ", len(options))

	r := newReader()
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	n, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || n < 1 || n > len(options) {
		return "", errors.New("invalid selection")
	}
	return options[n-1], nil
}

// MultiSelect presents checkboxes (space-separated numbers).
func MultiSelect(label string, options []string) ([]string, error) {
	fmt.Println(label)
	for i, opt := range options {
		fmt.Printf("  %d. %s\n", i+1, opt)
	}
	fmt.Printf("Choose (space-separated) [1-%d]: ", len(options))

	r := newReader()
	input, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	var selected []string
	normalized := strings.NewReplacer(",", " ").Replace(input)
	for _, s := range strings.Fields(normalized) {
		if strings.Contains(s, "-") {
			parts := strings.SplitN(s, "-", 2)
			if len(parts) != 2 {
				continue
			}
			start, startErr := strconv.Atoi(parts[0])
			end, endErr := strconv.Atoi(parts[1])
			if startErr != nil || endErr != nil || start < 1 || end > len(options) || start > end {
				continue
			}
			for n := start; n <= end; n++ {
				selected = append(selected, options[n-1])
			}
			continue
		}

		n, convErr := strconv.Atoi(s)
		if convErr != nil || n < 1 || n > len(options) {
			continue
		}
		selected = append(selected, options[n-1])
	}
	return selected, nil
}
