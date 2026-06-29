package doctor

import (
	. "dappco.re/go"
)

func TestCmdCommands_AddDoctorCommands_Good(t *T) {
	c := New()
	AddDoctorCommands(c)

	AssertTrue(t, c.Command("doctor").OK)
	AssertContains(t, c.Command("doctor").Value.(*Command).Description, "environment")
	AssertTrue(t, c.Command("doctor/environment").OK)
	AssertTrue(t, c.Command("doctor/commands").OK)
	AssertTrue(t, c.Command("doctor/checks").OK)
	AssertTrue(t, c.Command("doctor/install").OK)
}

func TestCmdCommands_AddDoctorCommands_Bad(t *T) {
	var c *Core

	AssertPanics(t, func() { AddDoctorCommands(c) })
	AssertNil(t, c)
}

func TestCmdCommands_AddDoctorCommands_Ugly(t *T) {
	c := New()
	AddDoctorCommands(c)
	AddDoctorCommands(c)

	AssertTrue(t, c.Command("doctor").OK)
	AssertNotNil(t, c.Command("doctor").Value)
}
