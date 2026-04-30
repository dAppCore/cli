package cli

import (
	core "dappco.re/go"
	"testing/fstest"
)

func TestCommands_WithCommands_Good(t *core.T) {
	c := core.New()
	setup := WithCommands("x", func(c *core.Core) { c.Command("x", core.Command{}) })

	setup(c)
	core.AssertTrue(t, c.Command("x").OK)
}

func TestCommands_WithCommands_Bad(t *core.T) {
	c := core.New()
	setup := WithCommands("", func(c *core.Core) { c.Command("empty", core.Command{}) }, nil)

	setup(c)
	core.AssertTrue(t, c.Command("empty").OK)
}

func TestCommands_WithCommands_Ugly(t *core.T) {
	c := core.New()
	called := 0
	setup := WithCommands("x", func(_ *core.Core) { called++ })

	setup(c)
	core.AssertEqual(t, 1, called)
}

func TestCommands_RegisterCommands_Good(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(c *core.Core) { c.Command("registered", core.Command{}) })

	var count int
	for range RegisteredCommands() {
		count++
	}
	core.AssertEqual(t, 1, count)
}

func TestCommands_RegisterCommands_Bad(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(*core.Core) {}, nil)

	core.AssertEmpty(t, RegisteredLocales())
	core.AssertNotNil(t, RegisteredCommands())
}

func TestCommands_RegisterCommands_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "registered"})))
	RegisterCommands(func(c *core.Core) { c.Command("late", core.Command{}) })

	core.AssertTrue(t, Core().Command("late").OK)
}

func TestCommands_RegisteredLocales_Good(t *core.T) {
	resetGlobals(t)
	fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
	RegisterCommands(func(*core.Core) {}, fs)

	core.AssertLen(t, RegisteredLocales(), 1)
	core.AssertNotNil(t, RegisteredLocales()[0])
}

func TestCommands_RegisteredLocales_Bad(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(*core.Core) {}, nil)

	core.AssertNil(t, RegisteredLocales())
	core.AssertEmpty(t, RegisteredLocales())
}

func TestCommands_RegisteredLocales_Ugly(t *core.T) {
	resetGlobals(t)
	fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
	RegisterCommands(func(*core.Core) {}, fs, nil)

	core.AssertLen(t, RegisteredLocales(), 1)
	core.AssertNotNil(t, RegisteredLocales()[0])
}

func TestCommands_RegisteredCommands_Good(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(c *core.Core) { c.Command("one", core.Command{}) })

	var count int
	for fn := range RegisteredCommands() {
		core.AssertNotNil(t, fn)
		count++
	}
	core.AssertEqual(t, 1, count)
}

func TestCommands_RegisteredCommands_Bad(t *core.T) {
	resetGlobals(t)
	var count int
	for range RegisteredCommands() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, RegisteredLocales())
}

func TestCommands_RegisteredCommands_Ugly(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(*core.Core) {})
	RegisterCommands(func(*core.Core) {})

	var count int
	for range RegisteredCommands() {
		count++
	}
	core.AssertEqual(t, 2, count)
}
