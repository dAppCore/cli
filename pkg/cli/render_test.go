package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompositeRender_GlyphTheme(t *testing.T) {
	prevStyle := currentRenderStyle
	t.Cleanup(func() {
		currentRenderStyle = prevStyle
	})

	restoreThemeAndColors(t)
	UseASCII()

	c := Layout("HCF")
	c.H("header").C("content").F("footer")

	UseRenderSimple()
	out := c.String()
	assert.Contains(t, out, strings.Repeat("-", 40))

	UseRenderBoxed()
	out = c.String()
	assert.Contains(t, out, "+")
	assert.Contains(t, out, strings.Repeat("-", 40))
}
