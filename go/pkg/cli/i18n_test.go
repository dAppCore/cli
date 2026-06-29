package cli

import core "dappco.re/go"

func TestI18n_T_Good(t *core.T) {
	result := T("some.key")

	core.AssertEqual(t, "some.key", result)
	core.AssertNotEmpty(t, result)
}

func TestI18n_T_Bad(t *core.T) {
	result := T("cmd.doctor.issues", map[string]any{"Count": 0})

	core.AssertNotPanics(t, func() { _ = T("cmd.doctor.issues") })
	core.AssertNotEmpty(t, result)
}

func TestI18n_T_Ugly(t *core.T) {
	result := T("")

	core.AssertEqual(t, "", result)
	core.AssertEmpty(t, result)
}

// T routes through the live Core i18n service once the runtime is
// initialised, resolving keys from the bundled en.json locale.
func TestI18n_T_LiveInstance(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "i18ntest"})))

	core.AssertNotNil(t, instance)
	core.AssertNotEmpty(t, T("cmd.doctor.short"))
}

// T falls back to the magic-key translator for keys absent from the
// Core service even when the runtime is initialised.
func TestI18n_T_LiveInstanceFallback(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "i18ntest"})))

	core.AssertEqual(t, "Workspace:", T("i18n.label.workspace"))
}
