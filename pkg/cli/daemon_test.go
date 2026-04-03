package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectMode_Good(t *testing.T) {
	t.Setenv("CORE_DAEMON", "1")
	assert.Equal(t, ModeDaemon, DetectMode())
}

func TestDetectMode_Bad(t *testing.T) {
	t.Setenv("CORE_DAEMON", "0")
	mode := DetectMode()
	assert.NotEqual(t, ModeDaemon, mode)
}

func TestDetectMode_Ugly(t *testing.T) {
	// Mode.String() covers all branches including the default unknown case.
	assert.Equal(t, "interactive", ModeInteractive.String())
	assert.Equal(t, "pipe", ModePipe.String())
	assert.Equal(t, "daemon", ModeDaemon.String())
	assert.Equal(t, "unknown", Mode(99).String())
}
