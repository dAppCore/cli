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
