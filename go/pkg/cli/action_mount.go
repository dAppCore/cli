// Action-to-command bridge.
//
// MountActions projects the Core Action registry — the capability map — onto
// the CLI command tree, the same way the API and MCP surfaces project it into
// REST routes and tool definitions. One map, three surfaces; entitlements gate
// who may invoke what, where.
package cli

import (
	core "dappco.re/go"
)

// HelpGenerator derives a command description from an action's dotted name. It
// lets a binary inject readable help (e.g. via go-i18n grammar) without the
// lean lib taking that dependency — pass one to MountActions; omit it for the
// plain default.
//
//	cli.MountActions(c, func(name string) string { return grammar.Readable(name) })
type HelpGenerator func(name string) string

// MountActions registers every action in c's capability map as a CLI command.
//
// The dotted action name maps to a command path by replacing dots with slashes
// ("git.log" -> "git/log"), so `core git log` resolves to the action. An
// explicit, executable command already bound to a path is left untouched —
// hand-written commands win over projected ones. Invocation routes through
// Action.Run, so disabled or unentitled actions still report their gate at call
// time rather than being silently hidden.
//
// An optional HelpGenerator supplies descriptions for actions that declare
// none; without one, a plain readable form of the name is used.
//
//	cli.MountActions(cli.Core())
//	cli.MountActions(cli.Core(), grammarHelp)
func MountActions(c *core.Core, help ...HelpGenerator) core.Result {
	if c == nil {
		return core.Fail(core.E("cli.MountActions", "nil core", nil))
	}
	var gen HelpGenerator
	if len(help) > 0 {
		gen = help[0]
	}
	for _, name := range c.Actions() {
		path := core.Join("/", core.Split(name, ".")...)

		// Don't clobber an explicit, executable command at this path.
		if existing := c.Command(path); existing.OK {
			if cmd, ok := existing.Value.(*core.Command); ok && cmd.Action != nil {
				continue
			}
		}

		name := name // capture per iteration for the closure
		if r := c.Command(path, core.Command{
			Description: actionHelp(c.Action(name), gen),
			Action: func(opts core.Options) core.Result {
				return c.Action(name).Run(actionContext(), opts)
			},
		}); !r.OK {
			return r
		}
	}
	return core.Ok(nil)
}

// actionContext returns the CLI runtime context when initialised, else a
// background context (e.g. when MountActions runs under test).
func actionContext() core.Context {
	if instance != nil && instance.ctx != nil {
		return instance.ctx
	}
	return core.Background()
}

// actionHelp derives a command description for an action: its own Description
// if set, else an injected generator's output, else a plain readable form of
// the dotted name.
func actionHelp(a *core.Action, gen HelpGenerator) string {
	if a == nil {
		return ""
	}
	if a.Description != "" {
		return a.Description
	}
	if gen != nil {
		if s := gen(a.Name); s != "" {
			return s
		}
	}
	return core.Join(" ", core.Split(a.Name, ".")...)
}
