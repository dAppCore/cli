package cli

import (
	core "dappco.re/go"
)

func TestPrompt_Prompt_Good(t *core.T) {
	SetStdin(core.NewReader("codex\n"))
	defer SetStdin(nil)
	result := Prompt("Name", "default")
	got, _ := result.Value.(string)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "codex", got)
}

func TestPrompt_Prompt_Bad(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)
	result := Prompt("Name", "default")
	got, _ := result.Value.(string)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "default", got)
}

func TestPrompt_Prompt_Ugly(t *core.T) {
	SetStdin(core.NewReader(""))
	defer SetStdin(nil)
	result := Prompt("Name", "")
	got, _ := result.Value.(string)
	err := cliResultError(result)

	core.AssertError(t, err)
	core.AssertEqual(t, "", got)
}

func TestPrompt_Select_Good(t *core.T) {
	SetStdin(core.NewReader("2\n"))
	defer SetStdin(nil)
	result := Select("Pick", []string{"alpha", "beta"})
	got, _ := result.Value.(string)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "beta", got)
}

func TestPrompt_Select_Bad(t *core.T) {
	SetStdin(core.NewReader("9\n"))
	defer SetStdin(nil)
	result := Select("Pick", []string{"alpha", "beta"})
	got, _ := result.Value.(string)
	err := cliResultError(result)

	core.AssertError(t, err)
	core.AssertEqual(t, "", got)
}

func TestPrompt_Select_Ugly(t *core.T) {
	result := Select("Pick", nil)
	got, _ := result.Value.(string)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "", got)
}

func TestPrompt_MultiSelect_Good(t *core.T) {
	SetStdin(core.NewReader("1 3\n"))
	defer SetStdin(nil)
	result := MultiSelect("Pick", []string{"alpha", "beta", "gamma"})
	got, _ := result.Value.([]string)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertEqual(t, []string{"alpha", "gamma"}, got)
}

func TestPrompt_MultiSelect_Bad(t *core.T) {
	SetStdin(core.NewReader("9\n"))
	defer SetStdin(nil)
	result := MultiSelect("Pick", []string{"alpha", "beta"})
	got, _ := result.Value.([]string)
	err := cliResultError(result)

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestPrompt_MultiSelect_Ugly(t *core.T) {
	result := MultiSelect("Pick", nil)
	got, _ := result.Value.([]string)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertEmpty(t, got)
}
