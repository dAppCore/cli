package cli

import (
	"bufio"
	"context"
	core "dappco.re/go"
	"time"
)

func TestUtils_GhAuthenticated_Good(t *core.T) {
	cliFakeCommands(t, map[string]string{"gh": "echo 'Logged in to github.com'\n"})

	core.AssertTrue(t, GhAuthenticated())
	core.AssertNotPanics(t, func() { _ = GhAuthenticated() })
}

func TestUtils_GhAuthenticated_Bad(t *core.T) {
	cliFakeCommands(t, map[string]string{"gh": "echo 'not logged in'\nexit 1\n"})

	core.AssertFalse(t, GhAuthenticated())
	core.AssertNotPanics(t, func() { _ = GhAuthenticated() })
}

func TestUtils_GhAuthenticated_Ugly(t *core.T) {
	cliFakeCommands(t, map[string]string{"gh": "echo 'Logged in as codex'\n"})

	core.AssertTrue(t, GhAuthenticated())
	core.AssertNotPanics(t, func() { _ = GhAuthenticated() })
}

func TestUtils_DefaultYes_Good(t *core.T) {
	cfg := &confirmConfig{}
	DefaultYes()(cfg)

	core.AssertTrue(t, cfg.defaultYes)
	core.AssertFalse(t, cfg.required)
}

func TestUtils_DefaultYes_Bad(t *core.T) {
	var cfg *confirmConfig

	core.AssertPanics(t, func() { DefaultYes()(cfg) })
	core.AssertNil(t, cfg)
}

func TestUtils_DefaultYes_Ugly(t *core.T) {
	cfg := &confirmConfig{defaultYes: false, required: true}
	DefaultYes()(cfg)

	core.AssertTrue(t, cfg.defaultYes)
	core.AssertTrue(t, cfg.required)
}

func TestUtils_Required_Good(t *core.T) {
	cfg := &confirmConfig{}
	Required()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertFalse(t, cfg.defaultYes)
}

func TestUtils_Required_Bad(t *core.T) {
	var cfg *confirmConfig

	core.AssertPanics(t, func() { Required()(cfg) })
	core.AssertNil(t, cfg)
}

func TestUtils_Required_Ugly(t *core.T) {
	cfg := &confirmConfig{defaultYes: true}
	Required()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertTrue(t, cfg.defaultYes)
}

func TestUtils_Timeout_Good(t *core.T) {
	cfg := &confirmConfig{}
	Timeout(time.Second)(cfg)

	core.AssertEqual(t, time.Second, cfg.timeout)
	core.AssertFalse(t, cfg.required)
}

func TestUtils_Timeout_Bad(t *core.T) {
	cfg := &confirmConfig{}
	Timeout(0)(cfg)

	core.AssertEqual(t, time.Duration(0), cfg.timeout)
	core.AssertFalse(t, cfg.defaultYes)
}

func TestUtils_Timeout_Ugly(t *core.T) {
	cfg := &confirmConfig{}
	Timeout(-time.Second)(cfg)

	core.AssertEqual(t, -time.Second, cfg.timeout)
	core.AssertFalse(t, cfg.required)
}

func TestUtils_Confirm_Good(t *core.T) {
	SetStdin(core.NewReader("y\n"))
	defer SetStdin(nil)

	core.AssertTrue(t, Confirm("Continue?"))
}

func TestUtils_Confirm_Bad(t *core.T) {
	SetStdin(core.NewReader("n\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, Confirm("Continue?", DefaultYes()))
}

func TestUtils_Confirm_Ugly(t *core.T) {
	SetStdin(core.NewReader("maybe\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, Confirm("Continue?"))
}

func TestUtils_ConfirmAction_Good(t *core.T) {
	SetStdin(core.NewReader("yes\n"))
	defer SetStdin(nil)

	core.AssertTrue(t, ConfirmAction("install", "package"))
}

func TestUtils_ConfirmAction_Bad(t *core.T) {
	SetStdin(core.NewReader("no\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, ConfirmAction("remove", "package", DefaultYes()))
}

func TestUtils_ConfirmAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertTrue(t, ConfirmAction("deploy", "", DefaultYes()))
}

func TestUtils_ConfirmDangerousAction_Good(t *core.T) {
	SetStdin(bufio.NewReader(core.NewReader("y\ny\n")))
	defer SetStdin(nil)

	core.AssertTrue(t, ConfirmDangerousAction("remove", "package"))
}

func TestUtils_ConfirmDangerousAction_Bad(t *core.T) {
	SetStdin(core.NewReader("n\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, ConfirmDangerousAction("remove", "package"))
}

func TestUtils_ConfirmDangerousAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("y\nn\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, ConfirmDangerousAction("remove", "package"))
}

func TestUtils_WithDefault_Good(t *core.T) {
	cfg := &questionConfig{}
	WithDefault("codex")(cfg)

	core.AssertEqual(t, "codex", cfg.defaultValue)
	core.AssertFalse(t, cfg.required)
}

func TestUtils_WithDefault_Bad(t *core.T) {
	var cfg *questionConfig

	core.AssertPanics(t, func() { WithDefault("codex")(cfg) })
	core.AssertNil(t, cfg)
}

func TestUtils_WithDefault_Ugly(t *core.T) {
	cfg := &questionConfig{defaultValue: "old"}
	WithDefault("")(cfg)

	core.AssertEqual(t, "", cfg.defaultValue)
	core.AssertFalse(t, cfg.required)
}

func TestUtils_WithValidator_Good(t *core.T) {
	cfg := &questionConfig{}
	WithValidator(func(string) error { return nil })(cfg)

	core.AssertNotNil(t, cfg.validator)
	core.AssertNoError(t, cfg.validator("ok"))
}

func TestUtils_WithValidator_Bad(t *core.T) {
	cfg := &questionConfig{}
	WithValidator(nil)(cfg)

	core.AssertNil(t, cfg.validator)
	core.AssertFalse(t, cfg.required)
}

func TestUtils_WithValidator_Ugly(t *core.T) {
	cfg := &questionConfig{}
	WithValidator(func(string) error { return Err("invalid") })(cfg)

	core.AssertError(t, cfg.validator("bad"))
	core.AssertNotNil(t, cfg.validator)
}

func TestUtils_RequiredInput_Good(t *core.T) {
	cfg := &questionConfig{}
	RequiredInput()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertEqual(t, "", cfg.defaultValue)
}

func TestUtils_RequiredInput_Bad(t *core.T) {
	var cfg *questionConfig

	core.AssertPanics(t, func() { RequiredInput()(cfg) })
	core.AssertNil(t, cfg)
}

func TestUtils_RequiredInput_Ugly(t *core.T) {
	cfg := &questionConfig{defaultValue: "fallback"}
	RequiredInput()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertEqual(t, "fallback", cfg.defaultValue)
}

func TestUtils_Question_Good(t *core.T) {
	SetStdin(core.NewReader("codex\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "codex", Question("Name?"))
}

func TestUtils_Question_Bad(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "default", Question("Name?", WithDefault("default")))
}

func TestUtils_Question_Ugly(t *core.T) {
	SetStdin(core.NewReader("bad\ngood\n"))
	defer SetStdin(nil)
	got := Question("Name?", WithValidator(func(v string) error {
		if v == "bad" {
			return Err("bad")
		}
		return nil
	}))

	core.AssertEqual(t, "good", got)
}

func TestUtils_QuestionAction_Good(t *core.T) {
	SetStdin(core.NewReader("codex\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "codex", QuestionAction("name", "agent"))
}

func TestUtils_QuestionAction_Bad(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "fallback", QuestionAction("name", "agent", WithDefault("fallback")))
}

func TestUtils_QuestionAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("value\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "value", QuestionAction("", ""))
}

func TestUtils_WithDisplay_Good(t *core.T) {
	cfg := &chooseConfig[int]{}
	WithDisplay(func(v int) string { return Sprintf("item-%d", v) })(cfg)

	core.AssertNotNil(t, cfg.displayFn)
	core.AssertEqual(t, "item-2", cfg.displayFn(2))
}

func TestUtils_WithDisplay_Bad(t *core.T) {
	cfg := &chooseConfig[int]{}
	WithDisplay[int](nil)(cfg)

	core.AssertNil(t, cfg.displayFn)
	core.AssertFalse(t, cfg.filter)
}

func TestUtils_WithDisplay_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDisplay(func(v string) string { return core.Upper(v) })(cfg)

	core.AssertEqual(t, "CODEX", cfg.displayFn("codex"))
	core.AssertFalse(t, cfg.multi)
}

func TestUtils_WithDefaultIndex_Good(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDefaultIndex[string](1)(cfg)

	core.AssertEqual(t, 1, cfg.defaultN)
	core.AssertFalse(t, cfg.filter)
}

func TestUtils_WithDefaultIndex_Bad(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDefaultIndex[string](-1)(cfg)

	core.AssertEqual(t, -1, cfg.defaultN)
	core.AssertFalse(t, cfg.multi)
}

func TestUtils_WithDefaultIndex_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDefaultIndex[string](99)(cfg)

	core.AssertEqual(t, 99, cfg.defaultN)
	core.AssertFalse(t, cfg.filter)
}

func TestUtils_Filter_Good(t *core.T) {
	cfg := &chooseConfig[string]{}
	Filter[string]()(cfg)

	core.AssertTrue(t, cfg.filter)
	core.AssertFalse(t, cfg.multi)
}

func TestUtils_Filter_Bad(t *core.T) {
	var cfg *chooseConfig[string]

	core.AssertPanics(t, func() { Filter[string]()(cfg) })
	core.AssertNil(t, cfg)
}

func TestUtils_Filter_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{multi: true}
	Filter[string]()(cfg)

	core.AssertTrue(t, cfg.filter)
	core.AssertTrue(t, cfg.multi)
}

func TestUtils_Multi_Good(t *core.T) {
	cfg := &chooseConfig[string]{}
	Multi[string]()(cfg)

	core.AssertTrue(t, cfg.multi)
	core.AssertFalse(t, cfg.filter)
}

func TestUtils_Multi_Bad(t *core.T) {
	var cfg *chooseConfig[string]

	core.AssertPanics(t, func() { Multi[string]()(cfg) })
	core.AssertNil(t, cfg)
}

func TestUtils_Multi_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{filter: true}
	Multi[string]()(cfg)

	core.AssertTrue(t, cfg.multi)
	core.AssertTrue(t, cfg.filter)
}

func TestUtils_Display_Good(t *core.T) {
	cfg := &chooseConfig[int]{}
	Display(func(v int) string { return Sprintf("n=%d", v) })(cfg)

	core.AssertNotNil(t, cfg.displayFn)
	core.AssertEqual(t, "n=3", cfg.displayFn(3))
}

func TestUtils_Display_Bad(t *core.T) {
	cfg := &chooseConfig[int]{}
	Display[int](nil)(cfg)

	core.AssertNil(t, cfg.displayFn)
	core.AssertFalse(t, cfg.filter)
}

func TestUtils_Display_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{}
	Display(func(v string) string { return v + v })(cfg)

	core.AssertEqual(t, "aa", cfg.displayFn("a"))
	core.AssertFalse(t, cfg.multi)
}

func TestUtils_Choose_Good(t *core.T) {
	SetStdin(core.NewReader("2\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "beta", Choose("Pick", []string{"alpha", "beta"}))
}

func TestUtils_Choose_Bad(t *core.T) {
	got := Choose("Pick", []string{})

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestUtils_Choose_Ugly(t *core.T) {
	SetStdin(core.NewReader("alp\n1\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "alpha", Choose("Pick", []string{"alpha", "beta"}, Filter[string]()))
}

func TestUtils_ChooseAction_Good(t *core.T) {
	SetStdin(core.NewReader("1\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "alpha", ChooseAction("pick", "agent", []string{"alpha", "beta"}))
}

func TestUtils_ChooseAction_Bad(t *core.T) {
	got := ChooseAction("pick", "agent", []string{})

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestUtils_ChooseAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "beta", ChooseAction("pick", "", []string{"alpha", "beta"}, WithDefaultIndex[string](1)))
}

func TestUtils_ChooseMulti_Good(t *core.T) {
	SetStdin(core.NewReader("1 3\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, []string{"alpha", "gamma"}, ChooseMulti("Pick", []string{"alpha", "beta", "gamma"}))
}

func TestUtils_ChooseMulti_Bad(t *core.T) {
	got := ChooseMulti("Pick", []string{})

	core.AssertNil(t, got)
	core.AssertEmpty(t, got)
}

func TestUtils_ChooseMulti_Ugly(t *core.T) {
	SetStdin(core.NewReader("gam\n1\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, []string{"gamma"}, ChooseMulti("Pick", []string{"alpha", "beta", "gamma"}, Filter[string]()))
}

func TestUtils_ChooseMultiAction_Good(t *core.T) {
	SetStdin(core.NewReader("2\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, []string{"beta"}, ChooseMultiAction("pick", "agent", []string{"alpha", "beta"}))
}

func TestUtils_ChooseMultiAction_Bad(t *core.T) {
	got := ChooseMultiAction("pick", "agent", []string{})

	core.AssertNil(t, got)
	core.AssertEmpty(t, got)
}

func TestUtils_ChooseMultiAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertNil(t, ChooseMultiAction("pick", "", []string{"alpha"}))
}

func TestUtils_GitClone_Good(t *core.T) {
	cliFakeCommands(t, map[string]string{"gh": "echo 'Logged in'\nexit 0\n"})

	err := cliResultError(GitClone(context.Background(), "org", "repo", "target"))
	core.AssertNoError(t, err)
}

func TestUtils_GitClone_Bad(t *core.T) {
	cliFakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'clone failed'\nexit 2\n",
	})

	err := cliResultError(GitClone(context.Background(), "org", "repo", "target"))
	core.AssertError(t, err)
}

func TestUtils_GitClone_Ugly(t *core.T) {
	cliFakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'ok'\nexit 0\n",
	})

	err := cliResultError(GitClone(nil, "org", "repo", "target"))
	core.AssertNoError(t, err)
}

func TestUtils_GitCloneRef_Good(t *core.T) {
	cliFakeCommands(t, map[string]string{"gh": "echo 'Logged in'\nexit 0\n"})

	err := cliResultError(GitCloneRef(context.Background(), "org", "repo", "target", "main"))
	core.AssertNoError(t, err)
}

func TestUtils_GitCloneRef_Bad(t *core.T) {
	cliFakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'already exists'\nexit 2\n",
	})

	err := cliResultError(GitCloneRef(context.Background(), "org", "repo", "target", "main"))
	core.AssertError(t, err)
}

func TestUtils_GitCloneRef_Ugly(t *core.T) {
	cliFakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'ok'\nexit 0\n",
	})

	err := cliResultError(GitCloneRef(nil, "org", "repo", "target", ""))
	core.AssertNoError(t, err)
}
