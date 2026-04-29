package i18n

import (
	"io/fs"
	"sort"
	"sync"
	"text/template"
	"unicode"

	"dappco.re/go"
)

// FSSource pairs a filesystem with a directory containing locale JSON files.
type FSSource struct {
	FS  fs.FS
	Dir string
}

// FSLoader loads flattened JSON translations from a filesystem.
type FSLoader struct {
	fsys fs.FS
	dir  string
}

// NewFSLoader creates a filesystem-backed locale loader.
func NewFSLoader(fsys fs.FS, dir string) *FSLoader {
	return &FSLoader{fsys: fsys, dir: dir}
}

// Service stores the active CLI translations.
type Service struct {
	mu       sync.RWMutex
	messages map[string]string
	lang     string
}

var defaultService = &Service{
	messages: make(map[string]string),
	lang:     "en",
}

// Default returns the process-wide CLI translation service.
func Default() *Service {
	return defaultService
}

// AddLoader merges translations from loader into the service.
func (s *Service) AddLoader(loader *FSLoader) core.Result {
	if s == nil {
		return core.Fail(core.NewError("i18n: nil service"))
	}
	if loader == nil {
		return core.Fail(core.NewError("i18n: nil loader"))
	}
	messagesResult := loader.Load(s.lang)
	if !messagesResult.OK {
		return messagesResult
	}
	messages := messagesResult.Value.(map[string]string)
	s.mu.Lock()
	for key, value := range messages {
		s.messages[key] = value
	}
	s.mu.Unlock()
	return core.Ok(nil)
}

// Load reads the best matching locale file for lang.
func (l *FSLoader) Load(lang string) core.Result {
	if l == nil || l.fsys == nil {
		return core.Fail(core.NewError("i18n: nil filesystem"))
	}
	dir := l.dir
	if dir == "" {
		dir = "."
	}

	candidates := localeCandidates(lang)
	var firstErr error
	for _, candidate := range candidates {
		messages := l.loadFile(core.Join("/", dir, candidate+".json"))
		if messages.OK {
			return messages
		}
		err, _ := messages.Value.(error)
		if firstErr == nil && !core.Is(err, fs.ErrNotExist) {
			firstErr = err
		}
	}

	languagesResult := l.Languages()
	if !languagesResult.OK {
		return languagesResult
	}
	languages := languagesResult.Value.([]string)
	if len(languages) == 0 {
		if firstErr != nil {
			return core.Fail(firstErr)
		}
		return core.Fail(core.NewError("i18n: no locale files found"))
	}
	return l.loadFile(core.Join("/", dir, languages[0]+".json"))
}

// Languages returns the locale tags available in the loader directory.
func (l *FSLoader) Languages() core.Result {
	if l == nil || l.fsys == nil {
		return core.Fail(core.NewError("i18n: nil filesystem"))
	}
	dir := l.dir
	if dir == "" {
		dir = "."
	}
	entriesResult := core.ReadDir(l.fsys, dir)
	if !entriesResult.OK {
		return entriesResult
	}
	entries := entriesResult.Value.([]core.FsDirEntry)
	languages := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !core.HasSuffix(entry.Name(), ".json") {
			continue
		}
		languages = append(languages, core.TrimSuffix(entry.Name(), ".json"))
	}
	sort.Strings(languages)
	return core.Ok(languages)
}

func (l *FSLoader) loadFile(name string) core.Result {
	dataResult := core.ReadFSFile(l.fsys, name)
	if !dataResult.OK {
		return dataResult
	}
	data := dataResult.Value.([]byte)
	var raw map[string]any
	if r := core.JSONUnmarshal(data, &raw); !r.OK {
		return r
	}
	out := make(map[string]string)
	flatten("", raw, out)
	return core.Ok(out)
}

func localeCandidates(lang string) []string {
	lang = core.Trim(lang)
	if lang == "" {
		lang = "en"
	}
	candidates := []string{lang}
	if normalized := core.Replace(lang, "_", "-"); normalized != lang {
		candidates = append(candidates, normalized)
	}
	if idx := firstIndexAny(lang, "-_"); idx > 0 {
		candidates = append(candidates, lang[:idx])
	}
	candidates = append(candidates, "en")
	return uniqueStrings(candidates)
}

func firstIndexAny(s, chars string) int {
	for i, r := range s {
		for _, c := range chars {
			if r == c {
				return i
			}
		}
	}
	return -1
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func flatten(prefix string, value any, out map[string]string) {
	switch v := value.(type) {
	case string:
		if prefix != "" {
			out[prefix] = v
		}
	case map[string]any:
		for key, child := range v {
			next := key
			if prefix != "" {
				next = prefix + "." + key
			}
			flatten(next, child, out)
		}
	}
}

// T translates a message ID using the default service.
func T(messageID string, args ...any) string {
	return Default().T(messageID, args...)
}

// T translates a message ID using the service.
func (s *Service) T(messageID string, args ...any) string {
	if messageID == "" {
		return ""
	}
	if s == nil {
		return messageID
	}

	s.mu.RLock()
	text, ok := s.messages[messageID]
	s.mu.RUnlock()
	if ok {
		return renderTemplate(text, templateData(args...))
	}
	if msg := renderMagicKey(messageID, args...); msg != "" {
		return msg
	}
	return messageID
}

func renderTemplate(text string, data any) string {
	if !core.Contains(text, "{{") {
		return text
	}
	tmpl, err := template.New("i18n").Option("missingkey=zero").Funcs(template.FuncMap{
		"title":        Title,
		"label":        Label,
		"progress":     Progress,
		"actionFailed": ActionFailed,
	}).Parse(text)
	if err != nil {
		return text
	}
	buf := core.NewBufferString("")
	if err := tmpl.Execute(buf, data); err != nil {
		return text
	}
	return buf.String()
}

func templateData(args ...any) any {
	if len(args) == 0 {
		return nil
	}
	if len(args) == 1 {
		switch v := args[0].(type) {
		case map[string]any, map[string]string, map[string]int:
			return v
		default:
			return map[string]any{
				"Item":  v,
				"Value": v,
				"Name":  v,
				"Count": v,
			}
		}
	}
	data := make(map[string]any, len(args)+4)
	for i, arg := range args {
		data[core.Sprintf("Arg%d", i+1)] = arg
	}
	data["Item"] = args[0]
	data["Value"] = args[0]
	data["Name"] = args[0]
	data["Count"] = args[0]
	return data
}

func renderMagicKey(messageID string, args ...any) string {
	switch {
	case core.HasPrefix(messageID, "i18n.fail."):
		return ActionFailed(core.TrimPrefix(messageID, "i18n.fail."), subjectArg(args...))
	case core.HasPrefix(messageID, "i18n.done."):
		return actionResult(core.TrimPrefix(messageID, "i18n.done."), subjectArg(args...))
	case core.HasPrefix(messageID, "i18n.label."):
		return Label(core.TrimPrefix(messageID, "i18n.label."))
	case core.HasPrefix(messageID, "i18n.progress."):
		return Progress(core.TrimPrefix(messageID, "i18n.progress."))
	default:
		return ""
	}
}

func subjectArg(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	switch v := args[0].(type) {
	case string:
		return v
	case map[string]any:
		return firstMapString(v, "Subject", "Item", "Name", "Value")
	case map[string]string:
		return firstMapString(v, "Subject", "Item", "Name", "Value")
	case map[string]int:
		for _, key := range []string{"Count", "Total"} {
			if value, ok := v[key]; ok {
				return core.Sprint(value)
			}
		}
	}
	return core.Sprint(args[0])
}

func firstMapString[M ~map[string]V, V any](m M, keys ...string) string {
	for _, key := range keys {
		if value, ok := m[key]; ok {
			return core.Sprint(value)
		}
	}
	return ""
}

// Title capitalises the first rune after whitespace or hyphen separators.
func Title(s string) string {
	b := core.NewBuilder()
	b.Grow(len(s))
	capNext := true
	for _, r := range s {
		if unicode.IsLetter(r) && capNext {
			r = unicode.ToUpper(r)
		}
		b.WriteRune(r)
		capNext = unicode.IsSpace(r) || r == '-'
	}
	return b.String()
}

// Progress returns a simple gerund progress phrase.
func Progress(verb string) string {
	verb = core.Trim(verb)
	if verb == "" {
		return ""
	}
	return Title(gerund(verb)) + "..."
}

// ActionFailed returns a failure phrase.
func ActionFailed(verb, subject string) string {
	verb = core.Trim(core.Lower(verb))
	if verb == "" {
		return ""
	}
	subject = core.Trim(subject)
	if subject == "" {
		return "Failed to " + verb
	}
	return "Failed to " + verb + " " + subject
}

// Label returns a title-cased label with a colon suffix.
func Label(word string) string {
	word = core.Trim(word)
	if word == "" {
		return ""
	}
	return Title(word) + ":"
}

func actionResult(verb, subject string) string {
	verb = core.Trim(core.Lower(verb))
	if verb == "" {
		return ""
	}
	result := pastTense(verb)
	subject = core.Trim(subject)
	if subject == "" {
		return Title(result)
	}
	return Title(subject) + " " + result
}

func gerund(verb string) string {
	verb = core.Lower(core.Trim(verb))
	switch {
	case core.HasSuffix(verb, "ie"):
		return core.TrimSuffix(verb, "ie") + "ying"
	case core.HasSuffix(verb, "e") && !core.HasSuffix(verb, "ee"):
		return core.TrimSuffix(verb, "e") + "ing"
	default:
		return verb + "ing"
	}
}

func pastTense(verb string) string {
	verb = core.Lower(core.Trim(verb))
	switch {
	case core.HasSuffix(verb, "e"):
		return verb + "d"
	case core.HasSuffix(verb, "y") && len(verb) > 1 && !isVowel(rune(verb[len(verb)-2])):
		return core.TrimSuffix(verb, "y") + "ied"
	default:
		return verb + "ed"
	}
}

func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}
