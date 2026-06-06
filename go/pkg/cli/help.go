// help.go is the CLI's own command-catalog renderer.
//
// core/go's Cli.PrintHelp prints the raw internal command path (slash form) and
// treats Description as an i18n key, suppressing anything that doesn't resolve.
// CLI presentation is pkg/cli's concern — not surface-agnostic core/go — so the
// catalog is rendered here: paths space-separated (the user-facing form) and
// descriptions resolved through T (which translates keys and passes literal or
// generated text through). Wired in via Execute for the no-command case.
package cli

import (
	core "dappco.re/go"
)

// helpEntry is one row of the command catalog: a space-separated command path
// and its resolved description.
type helpEntry struct {
	Path string
	Desc string
}

// commandCatalog returns the runnable commands registered on c, each with its
// path rendered space-separated and its description resolved via T. Non-runnable
// group placeholders (no Action) are omitted.
func commandCatalog(c *core.Core) []helpEntry {
	if c == nil {
		return nil
	}
	var entries []helpEntry
	for _, path := range c.Commands() {
		r := c.Command(path)
		if !r.OK {
			continue
		}
		cmd, ok := r.Value.(*core.Command)
		if !ok || cmd.Action == nil {
			continue
		}
		entries = append(entries, helpEntry{
			Path: core.Join(" ", core.Split(path, "/")...),
			Desc: T(cmd.Description),
		})
	}
	return entries
}

// PrintHelp renders the command catalog to stdout: a header line followed by
// each runnable command (space-separated path) and its description.
//
//	cli.PrintHelp()
func PrintHelp() {
	c := Core()
	if c == nil {
		return
	}
	if AppName != "" {
		Println("%s %s", BoldStyle.Render(AppName), DimStyle.Render(SemVer()))
	}
	Blank()
	Println("%s", HeaderStyle.Render("Commands:"))
	for _, e := range commandCatalog(c) {
		Println("  %s  %s", RepoStyle.Render(e.Path), DimStyle.Render(e.Desc))
	}
}
