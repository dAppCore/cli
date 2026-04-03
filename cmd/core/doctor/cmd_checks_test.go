package doctor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredChecksIncludesGo(t *testing.T) {
	checks := requiredChecks()

	var found bool
	for _, c := range checks {
		if c.command == "go" {
			found = true
			assert.Equal(t, "version", c.versionFlag)
			break
		}
	}

	assert.True(t, found, "required checks should include the Go compiler")
}
