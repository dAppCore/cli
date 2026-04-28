package doctor

import . "dappco.re/go"

func TestRequiredChecksIncludesGo(t *T) {
	checks := requiredChecks()

	var found bool
	for _, c := range checks {
		if c.command == "go" {
			found = true
			AssertEqual(t, "version", c.versionFlag)
			break
		}
	}
	AssertTrue(t, found, "required checks should include the Go compiler")
}
