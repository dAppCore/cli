package i18n

import (
	core "dappco.re/go"
	"testing/fstest"
)

func i18nResultError(r core.Result) error {
	if r.OK {
		return nil
	}
	if err, ok := r.Value.(error); ok {
		return err
	}
	return core.NewError(r.Error())
}

func i18nLocaleFS() fstest.MapFS {
	return fstest.MapFS{
		"locales/en.json": {Data: []byte(`{"greeting":"Hello {{.Name}}","nested":{"value":"ready"}}`)},
		"locales/fr.json": {Data: []byte(`{"greeting":"Bonjour {{.Name}}"}`)},
	}
}

func TestI18n_NewFSLoader_Good(t *core.T) {
	loader := NewFSLoader(i18nLocaleFS(), "locales")

	core.AssertNotNil(t, loader)
	core.AssertEqual(t, "locales", loader.dir)
}

func TestI18n_NewFSLoader_Bad(t *core.T) {
	loader := NewFSLoader(nil, "locales")

	messagesResult := loader.Load("en")
	err := i18nResultError(messagesResult)
	core.AssertError(t, err)
}

func TestI18n_NewFSLoader_Ugly(t *core.T) {
	loader := NewFSLoader(i18nLocaleFS(), "")

	languagesResult := loader.Languages()
	languages, _ := languagesResult.Value.([]string)
	err := i18nResultError(languagesResult)
	core.AssertNoError(t, err)
	core.AssertEmpty(t, languages)
}

func TestI18n_Default_Good(t *core.T) {
	svc := Default()

	core.AssertNotNil(t, svc)
	core.AssertEqual(t, "en", svc.lang)
}

func TestI18n_Default_Bad(t *core.T) {
	first := Default()
	second := Default()

	core.AssertEqual(t, first, second)
	core.AssertNotNil(t, second)
}

func TestI18n_Default_Ugly(t *core.T) {
	got := Default().T("i18n.progress.check")

	core.AssertEqual(t, "Checking...", got)
	core.AssertNotNil(t, Default())
}

func TestI18n_Service_AddLoader_Good(t *core.T) {
	svc := &Service{messages: make(map[string]string), lang: "en"}
	err := i18nResultError(svc.AddLoader(NewFSLoader(i18nLocaleFS(), "locales")))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "Hello Codex", svc.T("greeting", map[string]any{"Name": "Codex"}))
}

func TestI18n_Service_AddLoader_Bad(t *core.T) {
	svc := &Service{messages: make(map[string]string), lang: "en"}
	err := i18nResultError(svc.AddLoader(nil))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "nil loader")
}

func TestI18n_Service_AddLoader_Ugly(t *core.T) {
	var svc *Service
	err := i18nResultError(svc.AddLoader(NewFSLoader(i18nLocaleFS(), "locales")))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "nil service")
}

func TestI18n_FSLoader_Load_Good(t *core.T) {
	messagesResult := NewFSLoader(i18nLocaleFS(), "locales").Load("fr")
	messages, _ := messagesResult.Value.(map[string]string)
	err := i18nResultError(messagesResult)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "Bonjour {{.Name}}", messages["greeting"])
}

func TestI18n_FSLoader_Load_Bad(t *core.T) {
	messagesResult := NewFSLoader(fstest.MapFS{}, "locales").Load("en")
	err := i18nResultError(messagesResult)

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "file does not exist")
}

func TestI18n_FSLoader_Load_Ugly(t *core.T) {
	messagesResult := NewFSLoader(i18nLocaleFS(), "locales").Load("en-GB")
	messages, _ := messagesResult.Value.(map[string]string)
	err := i18nResultError(messagesResult)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "ready", messages["nested.value"])
}

// Load falls back to the first available locale when the requested tag is absent.
func TestI18n_FSLoader_Load_Fallback(t *core.T) {
	fsys := fstest.MapFS{
		"locales/de.json": {Data: []byte(`{"greeting":"Hallo"}`)},
	}
	messagesResult := NewFSLoader(fsys, "locales").Load("zz")
	messages, _ := messagesResult.Value.(map[string]string)
	err := i18nResultError(messagesResult)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "Hallo", messages["greeting"])
}

// Load surfaces a malformed-JSON error rather than masking it as not-found.
func TestI18n_FSLoader_Load_MalformedSurfaces(t *core.T) {
	fsys := fstest.MapFS{
		"locales/en.json": {Data: []byte(`{not json`)},
	}
	messagesResult := NewFSLoader(fsys, "locales").Load("en")
	err := i18nResultError(messagesResult)

	core.AssertError(t, err)
}

// Load normalises underscore locale tags onto their hyphenated candidate.
func TestI18n_FSLoader_Load_Underscore(t *core.T) {
	fsys := fstest.MapFS{
		"locales/pt-BR.json": {Data: []byte(`{"greeting":"Ola"}`)},
	}
	messagesResult := NewFSLoader(fsys, "locales").Load("pt_BR")
	messages, _ := messagesResult.Value.(map[string]string)
	err := i18nResultError(messagesResult)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "Ola", messages["greeting"])
}

func TestI18n_FSLoader_Languages_Good(t *core.T) {
	languagesResult := NewFSLoader(i18nLocaleFS(), "locales").Languages()
	languages, _ := languagesResult.Value.([]string)
	err := i18nResultError(languagesResult)

	core.AssertNoError(t, err)
	core.AssertEqual(t, []string{"en", "fr"}, languages)
}

func TestI18n_FSLoader_Languages_Bad(t *core.T) {
	languagesResult := NewFSLoader(nil, "locales").Languages()
	err := i18nResultError(languagesResult)

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "nil filesystem")
}

func TestI18n_FSLoader_Languages_Ugly(t *core.T) {
	languagesResult := NewFSLoader(fstest.MapFS{"locales/readme.txt": {Data: []byte("x")}}, "locales").Languages()
	languages, _ := languagesResult.Value.([]string)
	err := i18nResultError(languagesResult)

	core.AssertNoError(t, err)
	core.AssertEmpty(t, languages)
}

func TestI18n_T_Good(t *core.T) {
	got := T("i18n.fail.load", "config")

	core.AssertEqual(t, "Failed to load config", got)
	core.AssertContains(t, got, "config")
}

func TestI18n_T_Bad(t *core.T) {
	got := T("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestI18n_T_Ugly(t *core.T) {
	got := T("unregistered.message.id")

	core.AssertEqual(t, "unregistered.message.id", got)
	core.AssertContains(t, got, "message")
}

func TestI18n_Service_T_Good(t *core.T) {
	svc := &Service{messages: map[string]string{"hello": "Hello {{.Name}}"}, lang: "en"}
	got := svc.T("hello", map[string]any{"Name": "Agent"})

	core.AssertEqual(t, "Hello Agent", got)
	core.AssertContains(t, got, "Agent")
}

func TestI18n_Service_T_Bad(t *core.T) {
	var svc *Service
	got := svc.T("missing.key")

	core.AssertEqual(t, "missing.key", got)
	core.AssertContains(t, got, "missing")
}

func TestI18n_Service_T_Ugly(t *core.T) {
	svc := &Service{messages: map[string]string{"bad": "{{"}, lang: "en"}
	got := svc.T("bad", "value")

	core.AssertEqual(t, "{{", got)
	core.AssertNotEmpty(t, got)
}

func TestI18n_Title_Good(t *core.T) {
	got := Title("load config")

	core.AssertEqual(t, "Load Config", got)
	core.AssertContains(t, got, "Config")
}

func TestI18n_Title_Bad(t *core.T) {
	got := Title("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestI18n_Title_Ugly(t *core.T) {
	got := Title("re-load config")

	core.AssertEqual(t, "Re-Load Config", got)
	core.AssertContains(t, got, "-Load")
}

func TestI18n_Progress_Good(t *core.T) {
	got := Progress("check")

	core.AssertEqual(t, "Checking...", got)
	core.AssertContains(t, got, "...")
}

func TestI18n_Progress_Bad(t *core.T) {
	got := Progress("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestI18n_Progress_Ugly(t *core.T) {
	got := Progress("tie")

	core.AssertEqual(t, "Tying...", got)
	core.AssertContains(t, got, "Tying")
}

func TestI18n_ActionFailed_Good(t *core.T) {
	got := ActionFailed("load", "config")

	core.AssertEqual(t, "Failed to load config", got)
	core.AssertContains(t, got, "load")
}

func TestI18n_ActionFailed_Bad(t *core.T) {
	got := ActionFailed("", "config")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestI18n_ActionFailed_Ugly(t *core.T) {
	got := ActionFailed("  LOAD  ", "")

	core.AssertEqual(t, "Failed to load", got)
	core.AssertContains(t, got, "load")
}

func TestI18n_Label_Good(t *core.T) {
	got := Label("workspace")

	core.AssertEqual(t, "Workspace:", got)
	core.AssertContains(t, got, ":")
}

func TestI18n_Label_Bad(t *core.T) {
	got := Label("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestI18n_Label_Ugly(t *core.T) {
	got := Label("  git status  ")

	core.AssertEqual(t, "Git Status:", got)
	core.AssertContains(t, got, "Status")
}

// actionResult is reached via the "i18n.done." magic-key branch of T.
func TestI18n_ActionResult_Good(t *core.T) {
	got := T("i18n.done.deploy", "service")

	core.AssertEqual(t, "Service deployed", got)
	core.AssertContains(t, got, "deployed")
}

func TestI18n_ActionResult_Bad(t *core.T) {
	got := T("i18n.done.", "service")

	core.AssertEqual(t, "i18n.done.", got)
	core.AssertContains(t, got, "done")
}

func TestI18n_ActionResult_Ugly(t *core.T) {
	got := T("i18n.done.sync")

	core.AssertEqual(t, "Synced", got)
	core.AssertNotEmpty(t, got)
}

// pastTense + isVowel drive the verb morphology in the "i18n.done." branch.
func TestI18n_PastTense_Good(t *core.T) {
	got := T("i18n.done.create")

	core.AssertEqual(t, "Created", got)
	core.AssertContains(t, got, "Created")
}

func TestI18n_PastTense_Bad(t *core.T) {
	got := T("i18n.done.copy")

	core.AssertEqual(t, "Copied", got)
	core.AssertContains(t, got, "ied")
}

func TestI18n_PastTense_Ugly(t *core.T) {
	got := T("i18n.done.deploy")

	core.AssertEqual(t, "Deployed", got)
	core.AssertContains(t, got, "yed")
}

// subjectArg resolves the first string-bearing argument across map shapes.
func TestI18n_SubjectArg_Good(t *core.T) {
	got := T("i18n.fail.load", map[string]any{"Subject": "config"})

	core.AssertEqual(t, "Failed to load config", got)
	core.AssertContains(t, got, "config")
}

func TestI18n_SubjectArg_Bad(t *core.T) {
	got := T("i18n.fail.load", map[string]int{"Count": 7})

	core.AssertEqual(t, "Failed to load 7", got)
	core.AssertContains(t, got, "7")
}

func TestI18n_SubjectArg_Ugly(t *core.T) {
	got := T("i18n.fail.load", map[string]string{"Name": "agent"})

	core.AssertEqual(t, "Failed to load agent", got)
	core.AssertContains(t, got, "agent")
}

// firstMapString falls through every candidate key and returns "" when none match.
func TestI18n_FirstMapString_Good(t *core.T) {
	got := T("i18n.fail.start", map[string]any{"Item": "daemon"})

	core.AssertEqual(t, "Failed to start daemon", got)
	core.AssertContains(t, got, "daemon")
}

func TestI18n_FirstMapString_Bad(t *core.T) {
	got := T("i18n.fail.start", map[string]any{"Unknown": "x"})

	core.AssertEqual(t, "Failed to start", got)
	core.AssertNotContains(t, got, "x")
}

func TestI18n_FirstMapString_Ugly(t *core.T) {
	got := T("i18n.fail.start", map[string]string{"Value": "broker"})

	core.AssertEqual(t, "Failed to start broker", got)
	core.AssertContains(t, got, "broker")
}

// templateData maps positional args onto Arg1..ArgN plus the canonical aliases.
func TestI18n_TemplateData_Good(t *core.T) {
	svc := &Service{messages: map[string]string{"k": "{{.Arg1}}-{{.Arg2}}"}, lang: "en"}
	got := svc.T("k", "alpha", "beta")

	core.AssertEqual(t, "alpha-beta", got)
	core.AssertContains(t, got, "alpha")
}

func TestI18n_TemplateData_Bad(t *core.T) {
	svc := &Service{messages: map[string]string{"k": "{{.Item}}"}, lang: "en"}
	got := svc.T("k")

	core.AssertEqual(t, "<no value>", got)
	core.AssertNotEmpty(t, got)
}

func TestI18n_TemplateData_Ugly(t *core.T) {
	svc := &Service{messages: map[string]string{"k": "{{.Count}}"}, lang: "en"}
	got := svc.T("k", map[string]int{"Count": 3})

	core.AssertEqual(t, "3", got)
	core.AssertContains(t, got, "3")
}

// renderMagicKey routes label/progress prefixes and returns "" for unknown shapes.
func TestI18n_RenderMagicKey_Good(t *core.T) {
	got := T("i18n.label.workspace")

	core.AssertEqual(t, "Workspace:", got)
	core.AssertContains(t, got, ":")
}

func TestI18n_RenderMagicKey_Bad(t *core.T) {
	got := T("i18n.unknown.thing")

	core.AssertEqual(t, "i18n.unknown.thing", got)
	core.AssertContains(t, got, "unknown")
}

func TestI18n_RenderMagicKey_Ugly(t *core.T) {
	got := T("i18n.progress.deploy")

	core.AssertEqual(t, "Deploying...", got)
	core.AssertContains(t, got, "...")
}
