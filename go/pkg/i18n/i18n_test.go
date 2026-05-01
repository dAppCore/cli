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
