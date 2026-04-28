package doctor

import . "dappco.re/go"

func TestAX7Doctor_AddDoctorCommands_Good(t *T) {
	c := New()
	AddDoctorCommands(c)

	AssertTrue(t, c.Command("doctor").OK)
	AssertContains(t, c.Command("doctor").Value.(Command).Description, "environment")
}

func TestAX7Doctor_AddDoctorCommands_Bad(t *T) {
	var c *Core

	AssertPanics(t, func() { AddDoctorCommands(c) })
	AssertNil(t, c)
}

func TestAX7Doctor_AddDoctorCommands_Ugly(t *T) {
	c := New()
	AddDoctorCommands(c)
	AddDoctorCommands(c)

	AssertTrue(t, c.Command("doctor").OK)
	AssertNotNil(t, c.Command("doctor").Value)
}
