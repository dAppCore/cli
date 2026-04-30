package config

import (
	. "dappco.re/go"
)

func TestCmd_AddConfigCommands_Good(t *T) {
	c := New()
	AddConfigCommands(c)

	AssertTrue(t, c.Command("config/get").OK)
	AssertTrue(t, c.Command("config/set").OK)
}

func TestCmd_AddConfigCommands_Bad(t *T) {
	var c *Core

	AssertPanics(t, func() { AddConfigCommands(c) })
	AssertNil(t, c)
}

func TestCmd_AddConfigCommands_Ugly(t *T) {
	c := New()
	AddConfigCommands(c)
	AddConfigCommands(c)

	AssertTrue(t, c.Command("config/list").OK)
	AssertTrue(t, c.Command("config/path").OK)
}
