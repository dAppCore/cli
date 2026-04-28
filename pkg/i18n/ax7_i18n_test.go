package i18n

import (
	"testing/fstest"

	core "dappco.re/go"
)

func ax7LocaleFS() fstest.MapFS {
	return fstest.MapFS{
		"locales/en.json": {Data: []byte(`{"greeting":"Hello {{.Name}}","nested":{"value":"ready"}}`)},
		"locales/fr.json": {Data: []byte(`{"greeting":"Bonjour {{.Name}}"}`)},
	}
}

func TestAX7I18N_NewFSLoader_Good(t *core.T) {
	loader := NewFSLoader(ax7LocaleFS(), "locales")

	core.AssertNotNil(t, loader)
	core.AssertEqual(t, "locales", loader.dir)
}

func TestAX7I18N_NewFSLoader_Bad(t *core.T) {
	loader := NewFSLoader(nil, "locales")

	_, err := loader.Load("en")
	core.AssertError(t, err)
}

func TestAX7I18N_NewFSLoader_Ugly(t *core.T) {
	loader := NewFSLoader(ax7LocaleFS(), "")

	languages, err := loader.Languages()
	core.AssertNoError(t, err)
	core.AssertEmpty(t, languages)
}

func TestAX7I18N_Default_Good(t *core.T) {
	svc := Default()

	core.AssertNotNil(t, svc)
	core.AssertEqual(t, "en", svc.lang)
}

func TestAX7I18N_Default_Bad(t *core.T) {
	first := Default()
	second := Default()

	core.AssertEqual(t, first, second)
	core.AssertNotNil(t, second)
}

func TestAX7I18N_Default_Ugly(t *core.T) {
	got := Default().T("i18n.progress.check")

	core.AssertEqual(t, "Checking...", got)
	core.AssertNotNil(t, Default())
}

func TestAX7I18N_Service_AddLoader_Good(t *core.T) {
	svc := &Service{messages: make(map[string]string), lang: "en"}
	err := svc.AddLoader(NewFSLoader(ax7LocaleFS(), "locales"))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "Hello Codex", svc.T("greeting", map[string]any{"Name": "Codex"}))
}

func TestAX7I18N_Service_AddLoader_Bad(t *core.T) {
	svc := &Service{messages: make(map[string]string), lang: "en"}
	err := svc.AddLoader(nil)

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "nil loader")
}

func TestAX7I18N_Service_AddLoader_Ugly(t *core.T) {
	var svc *Service
	err := svc.AddLoader(NewFSLoader(ax7LocaleFS(), "locales"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "nil service")
}

func TestAX7I18N_FSLoader_Load_Good(t *core.T) {
	messages, err := NewFSLoader(ax7LocaleFS(), "locales").Load("fr")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "Bonjour {{.Name}}", messages["greeting"])
}

func TestAX7I18N_FSLoader_Load_Bad(t *core.T) {
	_, err := NewFSLoader(fstest.MapFS{}, "locales").Load("en")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "file does not exist")
}

func TestAX7I18N_FSLoader_Load_Ugly(t *core.T) {
	messages, err := NewFSLoader(ax7LocaleFS(), "locales").Load("en-GB")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "ready", messages["nested.value"])
}

func TestAX7I18N_FSLoader_Languages_Good(t *core.T) {
	languages, err := NewFSLoader(ax7LocaleFS(), "locales").Languages()

	core.AssertNoError(t, err)
	core.AssertEqual(t, []string{"en", "fr"}, languages)
}

func TestAX7I18N_FSLoader_Languages_Bad(t *core.T) {
	_, err := NewFSLoader(nil, "locales").Languages()

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "nil filesystem")
}

func TestAX7I18N_FSLoader_Languages_Ugly(t *core.T) {
	languages, err := NewFSLoader(fstest.MapFS{"locales/readme.txt": {Data: []byte("x")}}, "locales").Languages()

	core.AssertNoError(t, err)
	core.AssertEmpty(t, languages)
}

func TestAX7I18N_T_Good(t *core.T) {
	got := T("i18n.fail.load", "config")

	core.AssertEqual(t, "Failed to load config", got)
	core.AssertContains(t, got, "config")
}

func TestAX7I18N_T_Bad(t *core.T) {
	got := T("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7I18N_T_Ugly(t *core.T) {
	got := T("unregistered.message.id")

	core.AssertEqual(t, "unregistered.message.id", got)
	core.AssertContains(t, got, "message")
}

func TestAX7I18N_Service_T_Good(t *core.T) {
	svc := &Service{messages: map[string]string{"hello": "Hello {{.Name}}"}, lang: "en"}
	got := svc.T("hello", map[string]any{"Name": "Agent"})

	core.AssertEqual(t, "Hello Agent", got)
	core.AssertContains(t, got, "Agent")
}

func TestAX7I18N_Service_T_Bad(t *core.T) {
	var svc *Service
	got := svc.T("missing.key")

	core.AssertEqual(t, "missing.key", got)
	core.AssertContains(t, got, "missing")
}

func TestAX7I18N_Service_T_Ugly(t *core.T) {
	svc := &Service{messages: map[string]string{"bad": "{{"}, lang: "en"}
	got := svc.T("bad", "value")

	core.AssertEqual(t, "{{", got)
	core.AssertNotEmpty(t, got)
}

func TestAX7I18N_Title_Good(t *core.T) {
	got := Title("load config")

	core.AssertEqual(t, "Load Config", got)
	core.AssertContains(t, got, "Config")
}

func TestAX7I18N_Title_Bad(t *core.T) {
	got := Title("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7I18N_Title_Ugly(t *core.T) {
	got := Title("re-load config")

	core.AssertEqual(t, "Re-Load Config", got)
	core.AssertContains(t, got, "-Load")
}

func TestAX7I18N_Progress_Good(t *core.T) {
	got := Progress("check")

	core.AssertEqual(t, "Checking...", got)
	core.AssertContains(t, got, "...")
}

func TestAX7I18N_Progress_Bad(t *core.T) {
	got := Progress("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7I18N_Progress_Ugly(t *core.T) {
	got := Progress("tie")

	core.AssertEqual(t, "Tying...", got)
	core.AssertContains(t, got, "Tying")
}

func TestAX7I18N_ActionFailed_Good(t *core.T) {
	got := ActionFailed("load", "config")

	core.AssertEqual(t, "Failed to load config", got)
	core.AssertContains(t, got, "load")
}

func TestAX7I18N_ActionFailed_Bad(t *core.T) {
	got := ActionFailed("", "config")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7I18N_ActionFailed_Ugly(t *core.T) {
	got := ActionFailed("  LOAD  ", "")

	core.AssertEqual(t, "Failed to load", got)
	core.AssertContains(t, got, "load")
}

func TestAX7I18N_Label_Good(t *core.T) {
	got := Label("workspace")

	core.AssertEqual(t, "Workspace:", got)
	core.AssertContains(t, got, ":")
}

func TestAX7I18N_Label_Bad(t *core.T) {
	got := Label("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7I18N_Label_Ugly(t *core.T) {
	got := Label("  git status  ")

	core.AssertEqual(t, "Git Status:", got)
	core.AssertContains(t, got, "Status")
}
