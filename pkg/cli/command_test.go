package cli

import (
	core "dappco.re/go"
)

func TestCommand_RegisterCommand_Good(t *core.T) {
	c := core.New()
	RegisterCommand(c, "hello", core.Command{Description: "Hello"})

	core.AssertTrue(t, c.Command("hello").OK)
	core.AssertContains(t, c.Command("hello").Value.(*core.Command).Description, "Hello")
}

func TestCommand_RegisterCommand_Bad(t *core.T) {
	var c *core.Core

	core.AssertPanics(t, func() { RegisterCommand(c, "hello", core.Command{}) })
	core.AssertNil(t, c)
}

func TestCommand_RegisterCommand_Ugly(t *core.T) {
	c := core.New()
	RegisterCommand(c, "root", core.Command{Description: "Root"})

	core.AssertTrue(t, c.Command("root").OK)
	core.AssertNotNil(t, c.Command("root").Value)
}

func TestCommand_RequireArgs_Good(t *core.T) {
	opts := core.NewOptions(core.Option{Key: "_arg", Value: "config"})
	got := RequireArgs(opts, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestCommand_RequireArgs_Bad(t *core.T) {
	opts := core.NewOptions()
	got := RequireArgs(opts, 1)

	core.AssertContains(t, got, "requires")
	core.AssertContains(t, got, "1")
}

func TestCommand_RequireArgs_Ugly(t *core.T) {
	opts := core.NewOptions()
	got := RequireArgs(opts, 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestCommand_RequireExactArgs_Good(t *core.T) {
	opts := core.NewOptions(core.Option{Key: "_arg", Value: "config"})
	got := RequireExactArgs(opts, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestCommand_RequireExactArgs_Bad(t *core.T) {
	opts := core.NewOptions(core.Option{Key: "_arg", Value: "extra"})
	got := RequireExactArgs(opts, 0)

	core.AssertContains(t, got, "accepts no arguments")
	core.AssertNotEmpty(t, got)
}

func TestCommand_RequireExactArgs_Ugly(t *core.T) {
	opts := core.NewOptions()
	got := RequireExactArgs(opts, 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}
