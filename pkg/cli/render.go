package cli

import (
	"io"
	"strings"

	"dappco.re/go/core"
)

// RenderStyle controls how layouts are rendered.
type RenderStyle int

const (
	RenderFlat RenderStyle = iota
	RenderSimple
	RenderBoxed
)

var currentRenderStyle = RenderFlat

func UseRenderFlat()   { currentRenderStyle = RenderFlat }
func UseRenderSimple() { currentRenderStyle = RenderSimple }
func UseRenderBoxed()  { currentRenderStyle = RenderBoxed }

// Render outputs the layout to terminal.
func (c *Composite) Render() {
	io.WriteString(stdoutWriter(), c.String())
}

// String returns the rendered layout.
func (c *Composite) String() string {
	var sb strings.Builder
	c.renderTo(&sb, 0)
	return sb.String()
}

func (c *Composite) renderTo(sb *strings.Builder, depth int) {
	order := []Region{RegionHeader, RegionLeft, RegionContent, RegionRight, RegionFooter}

	var active []Region
	for _, r := range order {
		if slot, ok := c.regions[r]; ok {
			if len(slot.blocks) > 0 || slot.child != nil {
				active = append(active, r)
			}
		}
	}

	for i, r := range active {
		slot := c.regions[r]
		if i > 0 && currentRenderStyle != RenderFlat {
			c.renderSeparator(sb, depth)
		}
		c.renderSlot(sb, slot, depth)
	}
}

func (c *Composite) renderSeparator(sb *strings.Builder, depth int) {
	indent := strings.Repeat("  ", depth)
	switch currentRenderStyle {
	case RenderBoxed:
		sb.WriteString(indent + Glyph(":tee:") + strings.Repeat(Glyph(":dash:"), 40) + Glyph(":tee:") + "\n")
	case RenderSimple:
		sb.WriteString(indent + strings.Repeat(Glyph(":dash:"), 40) + "\n")
	}
}

func (c *Composite) renderSlot(sb *strings.Builder, slot *Slot, depth int) {
	indent := strings.Repeat("  ", depth)
	for _, block := range slot.blocks {
		for _, line := range core.Split(block.Render(), "\n") {
			if line != "" {
				sb.WriteString(indent + line + "\n")
			}
		}
	}
	if slot.child != nil {
		slot.child.renderTo(sb, depth+1)
	}
}
