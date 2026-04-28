package help

import . "dappco.re/go"

func TestAX7Help_AddHelpCommands_Good(t *T) {
	c := New()
	AddHelpCommands(c)

	AssertTrue(t, c.Command("help").OK)
	AssertContains(t, c.Command("help").Value.(Command).Description, "help")
}

func TestAX7Help_AddHelpCommands_Bad(t *T) {
	var c *Core

	AssertPanics(t, func() { AddHelpCommands(c) })
	AssertNil(t, c)
}

func TestAX7Help_AddHelpCommands_Ugly(t *T) {
	c := New()
	AddHelpCommands(c)
	AddHelpCommands(c)

	AssertTrue(t, c.Command("help").OK)
	AssertNotNil(t, c.Command("help").Value)
}
