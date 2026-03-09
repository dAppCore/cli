package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectMode(t *testing.T) {
	t.Run("daemon mode from env", func(t *testing.T) {
		t.Setenv("CORE_DAEMON", "1")
		assert.Equal(t, ModeDaemon, DetectMode())
	})

	t.Run("mode string", func(t *testing.T) {
		assert.Equal(t, "interactive", ModeInteractive.String())
		assert.Equal(t, "pipe", ModePipe.String())
		assert.Equal(t, "daemon", ModeDaemon.String())
		assert.Equal(t, "unknown", Mode(99).String())
	})
}
