// help_grammar.go provides cmd/core's go-i18n-backed help generator, injected
// into cli.MountActions. This keeps the grammar dependency (dappco.re/go/i18n)
// in the binary, not the lean cli lib — pkg/cli only knows the HelpGenerator
// func type.
package main

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	grammar "dappco.re/go/i18n"
)

// grammarHelp composes a readable command description from a dotted action name
// using go-i18n's grammar engine: the last segment is the imperative action
// (title-cased), earlier segments the definite-phrase subject. A heuristic —
// "close to readable" without per-action grammatical hints — but a real
// exercise of the engine (article elision, phonetics, irregular forms).
//
//	grammarHelp("demo.echo")    // "Echo the demo"
//	grammarHelp("brain.recall") // "Recall the brain"
func grammarHelp(name string) string {
	parts := core.Split(name, ".")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	action := grammar.Title(parts[len(parts)-1])
	if len(parts) == 1 {
		return action
	}
	subject := core.Join(" ", parts[:len(parts)-1]...)
	return action + " " + grammar.DefinitePhrase(subject)
}

// mountActionsWithGrammar projects the action registry onto the CLI with
// go-i18n-generated help. Passed to cli.Main as the final setup so explicit
// commands win path collisions.
func mountActionsWithGrammar(c *core.Core) core.Result {
	return cli.MountActions(c, grammarHelp)
}
