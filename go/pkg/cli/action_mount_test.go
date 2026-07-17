package cli

import (
	core "dappco.re/go"
)

func TestActionMount_MountActions_Good(t *core.T) {
	c := core.New()
	core.RequireTrue(t, core.CliRegister(c).OK)
	ran := false
	c.Action("demo.ping", func(_ core.Context, _ core.Options) core.Result {
		ran = true
		return core.Ok("pong")
	})

	core.AssertTrue(t, MountActions(c).OK)

	// The action is now reachable as a command at the dotted-name path.
	core.AssertTrue(t, c.Command("demo/ping").OK)

	// Full CLI arg resolution (`demo ping`) routes through Cli.Run to the action.
	core.AssertTrue(t, c.Cli().Run("demo", "ping").OK)
	core.AssertTrue(t, ran)
}

func TestActionMount_MountActions_Bad(t *core.T) {
	// A nil core is reported, not panicked.
	core.AssertFalse(t, MountActions(nil).OK)
}

func TestActionMount_MountActions_HelpGenerator(t *core.T) {
	c := core.New()
	c.Action("widget.spin", func(_ core.Context, _ core.Options) core.Result {
		return core.Ok(nil)
	})

	// An injected generator supplies help for actions without a Description,
	// keeping go-i18n (or any generator) out of the lean lib.
	core.AssertTrue(t, MountActions(c, func(name string) string {
		return "generated: " + name
	}).OK)

	cmd := c.Command("widget/spin").Value.(*core.Command)
	core.AssertEqual(t, "generated: widget.spin", cmd.Description)
}

func TestActionMount_MountActions_HelpGenerator_DescriptionWins(t *core.T) {
	c := core.New()
	d := c.Action("widget.spin", func(_ core.Context, _ core.Options) core.Result {
		return core.Ok(nil)
	})
	d.Description = "explicit description"

	MountActions(c, func(string) string { return "generated" })

	cmd := c.Command("widget/spin").Value.(*core.Command)
	core.AssertEqual(t, "explicit description", cmd.Description) // explicit wins over generator
}

func TestActionMount_MountActions_Ugly(t *core.T) {
	c := core.New()

	// An explicit command must win over an action at the same path.
	explicitRan := false
	c.Command("git/log", core.Command{
		Description: "explicit",
		Action: func(_ core.Options) core.Result {
			explicitRan = true
			return core.Ok(nil)
		},
	})
	c.Action("git.log", func(_ core.Context, _ core.Options) core.Result {
		return core.Ok("from action")
	})

	core.AssertTrue(t, MountActions(c).OK)

	cmd := c.Command("git/log").Value.(*core.Command)
	core.AssertTrue(t, cmd.Run(core.NewOptions()).OK)
	core.AssertTrue(t, explicitRan) // explicit handler, not the action
}
