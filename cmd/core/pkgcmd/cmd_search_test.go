package pkgcmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvePkgSearchPattern_Good(t *testing.T) {
	t.Run("uses flag pattern when set", func(t *testing.T) {
		got := resolvePkgSearchPattern("core-*", []string{"api"})
		assert.Equal(t, "core-*", got)
	})

	t.Run("uses positional pattern when flag is empty", func(t *testing.T) {
		got := resolvePkgSearchPattern("", []string{"api"})
		assert.Equal(t, "api", got)
	})

	t.Run("defaults to wildcard when nothing is provided", func(t *testing.T) {
		got := resolvePkgSearchPattern("", nil)
		assert.Equal(t, "*", got)
	})
}
