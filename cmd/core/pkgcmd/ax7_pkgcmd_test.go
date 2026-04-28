package pkgcmd

import . "dappco.re/go"

func TestAX7Pkg_AddPkgCommands_Good(t *T) {
	c := New()
	AddPkgCommands(c)

	AssertTrue(t, c.Command("pkg/list").OK)
	AssertTrue(t, c.Command("pkg/install").OK)
}

func TestAX7Pkg_AddPkgCommands_Bad(t *T) {
	var c *Core

	AssertPanics(t, func() { AddPkgCommands(c) })
	AssertNil(t, c)
}

func TestAX7Pkg_AddPkgCommands_Ugly(t *T) {
	c := New()
	AddPkgCommands(c)
	AddPkgCommands(c)

	AssertTrue(t, c.Command("pkg/search").OK)
	AssertTrue(t, c.Command("pkg/remove").OK)
}
