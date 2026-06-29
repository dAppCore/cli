// actions_demo.go demonstrates the action→CLI bridge (cli.MountActions).
//
// A capability registered as a named Action becomes a CLI endpoint with no
// hand-written command — the same capability map the API and MCP surfaces
// project. This is a harness demonstration: the wired ecosystem modules still
// expose their functionality as explicit commands (they predate the action
// registry), so there is nothing real to project yet. Delete this file once
// those modules register their capabilities as actions.
package main

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// AddDemoActions registers a sample action so `core demo echo <text>` exercises
// the capability-map → CLI projection end to end.
func AddDemoActions(c *core.Core) core.Result {
	c.Action("demo.echo", func(_ core.Context, opts core.Options) core.Result {
		cli.Println("demo.echo (via action→CLI projection) arg=%q", opts.String("_arg"))
		return core.Ok(nil)
	})
	return core.Ok(nil)
}
