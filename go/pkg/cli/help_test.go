package cli

import (
	core "dappco.re/go"
)

func TestHelp_commandCatalog_Good(t *core.T) {
	c := core.New()
	c.Command("alpha", core.Command{
		Description: "the alpha command",
		Action:      func(core.Options) core.Result { return core.Ok(nil) },
	})
	c.Command("group/leaf", core.Command{
		Description: "a leaf",
		Action:      func(core.Options) core.Result { return core.Ok(nil) },
	})

	byPath := map[string]string{}
	for _, e := range commandCatalog(c) {
		byPath[e.Path] = e.Desc
	}

	core.AssertEqual(t, "the alpha command", byPath["alpha"])
	core.AssertEqual(t, "a leaf", byPath["group leaf"]) // slash path rendered space-separated
}

func TestHelp_commandCatalog_Bad(t *core.T) {
	core.AssertEqual(t, 0, len(commandCatalog(nil)))
}

func TestHelp_commandCatalog_Ugly(t *core.T) {
	c := core.New()
	c.Command("group/leaf", core.Command{
		Description: "a leaf",
		Action:      func(core.Options) core.Result { return core.Ok(nil) },
	})

	for _, e := range commandCatalog(c) {
		// The auto-created "group" parent placeholder is not runnable and must
		// not appear in the catalog.
		core.AssertNotEqual(t, "group", e.Path)
	}
}
