package cli

import (
	"unicode"

	"dappco.re/go"
	clii18n "dappco.re/go/cli/pkg/i18n"
)

// T translates a key using the CLI's Core i18n service.
// Falls back to the CLI translation service if CLI is not initialised.
//
//	label := cli.T("cmd.doctor.required")
//	msg := cli.T("cmd.doctor.issues", map[string]any{"Count": 3})
func T(key string, args ...any) string {
	if instance != nil && instance.core != nil {
		i18n := instance.core.I18n()
		if i18n != nil && i18n.Translator().OK {
			if r := i18n.Translate(key, args...); r.OK {
				if translated, ok := r.Value.(string); ok {
					return translated
				}
			}
		}
	}
	return clii18n.Default().T(key, args...)
}

func wordLabel(word string) string {
	return T("i18n.label." + word)
}

func progressMessage(verb string) string {
	return T("i18n.progress." + verb)
}

func actionFailed(verb, subject string) string {
	return T("i18n.fail."+verb, subject)
}

func title(s string) string {
	b := core.NewBuilder()
	b.Grow(len(s))
	capNext := true
	for _, r := range s {
		if unicode.IsLetter(r) && capNext {
			r = unicode.ToUpper(r)
		}
		b.WriteRune(r)
		capNext = unicode.IsSpace(r) || r == '-'
	}
	return b.String()
}
