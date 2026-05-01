package cli

import (
	"context"
	"io" // Note: AX-6 - io.Reader/io.Writer stream IO contract for prompt writes and interception.
	"syscall"
	"time"
	"unicode" // Note: AX-6 - unicode.IsSpace/IsDigit classify interactive selection tokens.

	"dappco.re/go"
)

func GhAuthenticated() bool {
	result := runProcessOutput(context.Background(), "gh", "auth", "status")
	output := processOutput(result.Value)
	authenticated := core.Contains(output, "Logged in")
	if authenticated {
		LogWarn("GitHub CLI authenticated", "user", core.Username())
	} else {
		LogWarn("GitHub CLI not authenticated", "user", core.Username())
	}
	return authenticated
}

func processCore() *core.Core {
	if instance != nil && instance.core != nil {
		return instance.core
	}
	return core.New()
}

func runProcessOutput(ctx context.Context, command string, args ...string) core.Result {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return core.Fail(err)
	}
	commandResult := findExecutable(command)
	if !commandResult.OK {
		return commandResult
	}
	commandPath := commandResult.Value.(string)

	var pipe [2]int
	if err := syscall.Pipe(pipe[:]); err != nil {
		return core.Fail(err)
	}
	readFD, writeFD := pipe[0], pipe[1]
	defer syscall.Close(readFD)

	argv := append([]string{commandPath}, args...)
	pid, err := syscall.ForkExec(commandPath, argv, &syscall.ProcAttr{
		Env:   core.Environ(),
		Files: []uintptr{0, uintptr(writeFD), uintptr(writeFD)},
	})
	syscall.Close(writeFD)
	if err != nil {
		return core.Fail(err)
	}

	out := core.NewBuilder()
	buf := make([]byte, 4096)
	for {
		n, readErr := syscall.Read(readFD, buf)
		if n > 0 {
			out.WriteString(string(buf[:n]))
		}
		if readErr != nil {
			if readErr == syscall.EINTR {
				continue
			}
			break
		}
		if n == 0 {
			break
		}
	}

	var status syscall.WaitStatus
	if _, err := syscall.Wait4(pid, &status, 0, nil); err != nil {
		return core.Fail(err)
	}
	output := out.String()
	if status.ExitStatus() == 0 {
		return core.Ok(output)
	}
	if output != "" {
		return core.Fail(core.NewError(output))
	}
	return core.Fail(core.E("cli.process", core.Sprintf("%s exited with status %d", command, status.ExitStatus()), nil))
}

func findExecutable(command string) core.Result {
	if command == "" {
		return core.Fail(core.NewError("empty command"))
	}
	if core.Contains(command, string(core.PathSeparator)) {
		if r := core.Stat(command); r.OK {
			return core.Ok(command)
		}
		return core.Fail(core.E("cli.process", core.Concat("command not found: ", command), nil))
	}
	for _, dir := range core.Split(core.Getenv("PATH"), string(core.PathListSeparator)) {
		if dir == "" {
			continue
		}
		candidate := core.PathJoin(dir, command)
		if r := core.Stat(candidate); r.OK {
			return core.Ok(candidate)
		}
	}
	return core.Fail(core.E("cli.process", core.Concat("command not found: ", command), nil))
}

func processOutput(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case nil:
		return ""
	default:
		return core.Sprint(v)
	}
}

type ConfirmOption func(*confirmConfig)

type confirmConfig struct {
	defaultYes bool
	required   bool
	timeout    time.Duration
}

func promptHint(msg string) {
	core.Print(stderrWriter(), "%s", DimStyle.Render(compileGlyphs(msg)))
}

func promptWarning(msg string) {
	core.Print(stderrWriter(), "%s", WarningStyle.Render(compileGlyphs(msg)))
}

func DefaultYes() ConfirmOption {
	return func(c *confirmConfig) { c.defaultYes = true }
}

func Required() ConfirmOption {
	return func(c *confirmConfig) { c.required = true }
}

func Timeout(d time.Duration) ConfirmOption {
	return func(c *confirmConfig) { c.timeout = d }
}

func Confirm(prompt string, opts ...ConfirmOption) bool {
	cfg := &confirmConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	prompt = compileGlyphs(prompt)

	var suffix string
	if cfg.required {
		suffix = "[y/n] "
	} else if cfg.defaultYes {
		suffix = "[Y/n] "
	} else {
		suffix = "[y/N] "
	}

	if cfg.timeout > 0 {
		suffix = core.Sprintf("%s(auto in %s) ", suffix, cfg.timeout.Round(time.Second))
	}

	reader := newReader()

	for {
		io.WriteString(stderrWriter(), core.Sprintf("%s %s", prompt, suffix))

		var response string
		var readErr error

		if cfg.timeout > 0 {
			resultChan := make(chan string, 1)
			errChan := make(chan error, 1)
			go func() {
				line, err := reader.ReadString('\n')
				resultChan <- line
				errChan <- err
			}()

			select {
			case response = <-resultChan:
				readErr = <-errChan
				response = core.Lower(core.Trim(response))
			case <-time.After(cfg.timeout):
				core.Print(stderrWriter(), "")
				return cfg.defaultYes
			}
		} else {
			line, err := reader.ReadString('\n')
			readErr = err
			if err != nil && line == "" {
				return cfg.defaultYes
			}
			response = line
			response = core.Lower(core.Trim(response))
		}

		if response == "" {
			if readErr == nil && cfg.required {
				promptHint("Please enter y or n, then press Enter.")
				continue
			}
			if cfg.required {
				return cfg.defaultYes
			}
			return cfg.defaultYes
		}

		if response == "y" || response == "yes" {
			return true
		}
		if response == "n" || response == "no" {
			return false
		}

		if cfg.required {
			promptHint("Please enter y or n, then press Enter.")
			continue
		}
		return cfg.defaultYes
	}
}

func ConfirmAction(verb, subject string, opts ...ConfirmOption) bool {
	question := title(verb) + " " + subject + "?"
	return Confirm(question, opts...)
}

func ConfirmDangerousAction(verb, subject string) bool {
	question := title(verb) + " " + subject + "?"
	if !Confirm(question, Required()) {
		return false
	}
	confirm := "Really " + verb + " " + subject + "?"
	return Confirm(confirm, Required())
}

type QuestionOption func(*questionConfig)

type questionConfig struct {
	defaultValue string
	required     bool
	validator    func(string) error
}

func WithDefault(value string) QuestionOption {
	return func(c *questionConfig) { c.defaultValue = value }
}

func WithValidator(fn func(string) error) QuestionOption {
	return func(c *questionConfig) { c.validator = fn }
}

func RequiredInput() QuestionOption {
	return func(c *questionConfig) { c.required = true }
}

func Question(prompt string, opts ...QuestionOption) string {
	cfg := &questionConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	prompt = compileGlyphs(prompt)
	reader := newReader()
	for {
		if cfg.defaultValue != "" {
			io.WriteString(stderrWriter(), core.Sprintf("%s [%s] ", prompt, compileGlyphs(cfg.defaultValue)))
		} else {
			io.WriteString(stderrWriter(), core.Sprintf("%s ", prompt))
		}
		response, err := reader.ReadString('\n')
		response = core.Trim(response)
		if err != nil && response == "" {
			return cfg.defaultValue
		}
		if response == "" {
			if cfg.required {
				promptHint("Please enter a value, then press Enter.")
				continue
			}
			response = cfg.defaultValue
		}
		if cfg.validator != nil {
			if err := cfg.validator(response); err != nil {
				promptWarning(core.Sprintf("Invalid: %v", err))
				continue
			}
		}
		return response
	}
}

func QuestionAction(verb, subject string, opts ...QuestionOption) string {
	question := title(verb) + " " + subject + "?"
	return Question(question, opts...)
}

type ChooseOption[T any] func(*chooseConfig[T])

type chooseConfig[T any] struct {
	displayFn func(T) string
	defaultN  int
	filter    bool
	multi     bool
}

func WithDisplay[T any](fn func(T) string) ChooseOption[T] {
	return func(c *chooseConfig[T]) { c.displayFn = fn }
}

func WithDefaultIndex[T any](idx int) ChooseOption[T] {
	return func(c *chooseConfig[T]) { c.defaultN = idx }
}

func Filter[T any]() ChooseOption[T] {
	return func(c *chooseConfig[T]) { c.filter = true }
}

func Multi[T any]() ChooseOption[T] {
	return func(c *chooseConfig[T]) { c.multi = true }
}

func Display[T any](fn func(T) string) ChooseOption[T] {
	return WithDisplay[T](fn)
}

func Choose[T any](prompt string, items []T, opts ...ChooseOption[T]) T {
	var zero T
	if len(items) == 0 {
		return zero
	}
	cfg := &chooseConfig[T]{
		displayFn: func(item T) string { return core.Sprint(item) },
		defaultN:  -1,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	prompt = compileGlyphs(prompt)
	reader := newReader()
	visible := make([]int, len(items))
	for i := range items {
		visible[i] = i
	}
	allVisible := append([]int(nil), visible...)

	for {
		renderChoices(prompt, items, visible, cfg.displayFn, cfg.defaultN, cfg.filter)
		if cfg.filter {
			io.WriteString(stderrWriter(), core.Sprintf("Enter number [1-%d] or filter: ", len(visible)))
		} else {
			io.WriteString(stderrWriter(), core.Sprintf("Enter number [1-%d]: ", len(visible)))
		}
		response, err := reader.ReadString('\n')
		response = core.Trim(response)
		if err != nil && response == "" {
			if idx, ok := defaultVisibleIndex(visible, cfg.defaultN); ok {
				return items[idx]
			}
			var zero T
			return zero
		}
		if response == "" {
			if cfg.filter && len(visible) != len(allVisible) {
				visible = append([]int(nil), allVisible...)
				promptHint("Filter cleared.")
				continue
			}
			if idx, ok := defaultVisibleIndex(visible, cfg.defaultN); ok {
				return items[idx]
			}
			if cfg.defaultN >= 0 {
				promptHint("Default selection is not available in the current list. Narrow the list or choose another number.")
				continue
			}
			promptHint(core.Sprintf("Please enter a number between 1 and %d.", len(visible)))
			continue
		}
		if parsed := Atoi(response); parsed.OK {
			n := parsed.Value.(int)
			if n >= 1 && n <= len(visible) {
				return items[visible[n-1]]
			}
			promptHint(core.Sprintf("Please enter a number between 1 and %d.", len(visible)))
			continue
		}
		if cfg.filter {
			nextVisible := filterVisible(items, visible, response, cfg.displayFn)
			if len(nextVisible) == 0 {
				promptHint(core.Sprintf("No matches for %q. Try a shorter search term or clear the filter.", response))
				continue
			}
			visible = nextVisible
			continue
		}
		promptHint(core.Sprintf("Please enter a number between 1 and %d.", len(visible)))
	}
}

func ChooseAction[T any](verb, subject string, items []T, opts ...ChooseOption[T]) T {
	question := title(verb) + " " + subject + ":"
	return Choose(question, items, opts...)
}

func ChooseMulti[T any](prompt string, items []T, opts ...ChooseOption[T]) []T {
	if len(items) == 0 {
		return nil
	}
	cfg := &chooseConfig[T]{
		displayFn: func(item T) string { return core.Sprint(item) },
		defaultN:  -1,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	prompt = compileGlyphs(prompt)
	reader := newReader()
	visible := make([]int, len(items))
	for i := range items {
		visible[i] = i
	}
	for {
		renderChoices(prompt, items, visible, cfg.displayFn, -1, cfg.filter)
		if cfg.filter {
			io.WriteString(stderrWriter(), "Enter numbers (e.g., 1 3 5 or 1-3), or filter text, or empty for none: ")
		} else {
			io.WriteString(stderrWriter(), "Enter numbers (e.g., 1 3 5 or 1-3) or empty for none: ")
		}
		response, _ := reader.ReadString('\n')
		response = core.Trim(response)
		if response == "" {
			return nil
		}
		selectedResult := parseMultiSelection(response, len(visible))
		if !selectedResult.OK {
			if cfg.filter && !looksLikeMultiSelectionInput(response) {
				nextVisible := filterVisible(items, visible, response, cfg.displayFn)
				if len(nextVisible) == 0 {
					promptHint(core.Sprintf("No matches for %q. Try a shorter search term or clear the filter.", response))
					continue
				}
				visible = nextVisible
				continue
			}
			promptWarning(core.Sprintf("Invalid selection %q: enter numbers like 1 3 or 1-3.", response))
			continue
		}
		selected := selectedResult.Value.([]int)
		result := make([]T, 0, len(selected))
		for _, idx := range selected {
			result = append(result, items[visible[idx]])
		}
		return result
	}
}

func renderChoices[T any](prompt string, items []T, visible []int, displayFn func(T) string, defaultN int, filter bool) {
	core.Print(stderrWriter(), "%s", prompt)
	for i, idx := range visible {
		marker := " "
		if defaultN >= 0 && idx == defaultN {
			marker = "*"
		}
		core.Print(stderrWriter(), "  %s%d. %s", marker, i+1, compileGlyphs(displayFn(items[idx])))
	}
	if filter {
		core.Print(stderrWriter(), "  (type to filter the list)")
	}
}

func defaultVisibleIndex(visible []int, defaultN int) (int, bool) {
	if defaultN < 0 {
		return 0, false
	}
	for _, idx := range visible {
		if idx == defaultN {
			return idx, true
		}
	}
	return 0, false
}

func filterVisible[T any](items []T, visible []int, query string, displayFn func(T) string) []int {
	q := core.Lower(core.Trim(query))
	if q == "" {
		return visible
	}
	filtered := make([]int, 0, len(visible))
	for _, idx := range visible {
		if core.Contains(core.Lower(displayFn(items[idx])), q) {
			filtered = append(filtered, idx)
		}
	}
	return filtered
}

func looksLikeMultiSelectionInput(input string) bool {
	hasDigit := false
	for _, r := range input {
		switch {
		case unicode.IsSpace(r), r == '-' || r == ',':
			continue
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			return false
		}
	}
	return hasDigit
}

func parseMultiSelection(input string, maxItems int) core.Result {
	selected := make(map[int]bool)
	normalized := core.Replace(input, ",", " ")
	for _, part := range fields(normalized) {
		if core.Contains(part, "-") {
			rangeParts := core.Split(part, "-")
			if len(rangeParts) != 2 {
				return Err("invalid range: %s", part)
			}
			startResult := Atoi(rangeParts[0])
			if !startResult.OK {
				return Err("invalid range start: %s", rangeParts[0])
			}
			start := startResult.Value.(int)
			endResult := Atoi(rangeParts[1])
			if !endResult.OK {
				return Err("invalid range end: %s", rangeParts[1])
			}
			end := endResult.Value.(int)
			if start < 1 || start > maxItems || end < 1 || end > maxItems || start > end {
				return Err("range out of bounds: %s", part)
			}
			for i := start; i <= end; i++ {
				selected[i-1] = true
			}
		} else {
			nResult := Atoi(part)
			if !nResult.OK {
				return Err("invalid number: %s", part)
			}
			n := nResult.Value.(int)
			if n < 1 || n > maxItems {
				return Err("number out of range: %d", n)
			}
			selected[n-1] = true
		}
	}
	result := make([]int, 0, len(selected))
	for i := range maxItems {
		if selected[i] {
			result = append(result, i)
		}
	}
	return core.Ok(result)
}

// fields splits a string on whitespace runs, returning non-empty tokens.
// Equivalent to strings.Fields without importing the stdlib package directly.
func fields(s string) []string {
	var parts []string
	start := -1
	for i, r := range s {
		if unicode.IsSpace(r) {
			if start >= 0 {
				parts = append(parts, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		parts = append(parts, s[start:])
	}
	return parts
}

func ChooseMultiAction[T any](verb, subject string, items []T, opts ...ChooseOption[T]) []T {
	question := title(verb) + " " + subject + ":"
	return ChooseMulti(question, items, opts...)
}

func GitClone(ctx context.Context, org, repo, path string) core.Result {
	return GitCloneRef(ctx, org, repo, path, "")
}

func GitCloneRef(ctx context.Context, org, repo, path, ref string) core.Result {
	if GhAuthenticated() {
		httpsURL := core.Sprintf("https://github.com/%s/%s.git", org, repo)
		args := ghRepoCloneArgs(httpsURL, path, ref)
		result := runProcessOutput(ctx, "gh", args...)
		output := processOutput(result.Value)
		if result.OK {
			return core.Ok(nil)
		}
		errStr := core.Trim(output)
		if core.Contains(errStr, "already exists") {
			return core.Fail(core.NewError(errStr))
		}
	}
	args := gitCloneArgs(core.Sprintf("git@github.com:%s/%s.git", org, repo), path, ref)
	result := runProcessOutput(ctx, "git", args...)
	if !result.OK {
		errStr := core.Trim(processOutput(result.Value))
		if errStr == "" {
			return result
		}
		return core.Fail(core.NewError(errStr))
	}
	return core.Ok(nil)
}

func ghRepoCloneArgs(repoURL, path, ref string) []string {
	args := []string{"repo", "clone", "--", repoURL, path}
	if ref != "" {
		args = append(args, "--", "--branch", ref, "--single-branch")
	}
	return args
}

func gitCloneArgs(repoURL, path, ref string) []string {
	args := []string{"clone"}
	if ref != "" {
		args = append(args, "--branch", ref, "--single-branch")
	}
	args = append(args, "--", repoURL, path)
	return args
}
