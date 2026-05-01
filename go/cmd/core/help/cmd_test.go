package help

import (
	. "dappco.re/go"
)

func TestCmd_AddHelpCommands_Good(t *T) {
	c := New()
	AddHelpCommands(c)

	AssertTrue(t, c.Command("help").OK)
	AssertContains(t, c.Command("help").Value.(Command).Description, "help")
}

func TestCmd_AddHelpCommands_Bad(t *T) {
	var c *Core

	AssertPanics(t, func() { AddHelpCommands(c) })
	AssertNil(t, c)
}

func TestCmd_AddHelpCommands_Ugly(t *T) {
	c := New()
	AddHelpCommands(c)
	AddHelpCommands(c)

	AssertTrue(t, c.Command("help").OK)
	AssertNotNil(t, c.Command("help").Value)
}
