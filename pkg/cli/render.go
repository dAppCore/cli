package cli

import (
	"io"

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
	sb := core.NewBuilder()
	c.renderTo(sb, 0)
	return sb.String()
}

func (c *Composite) renderTo(sb io.StringWriter, depth int) {
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

func (c *Composite) renderSeparator(sb io.StringWriter, depth int) {
	indent := Repeat("  ", depth)
	switch currentRenderStyle {
	case RenderBoxed:
		_, _ = sb.WriteString(indent + Glyph(":tee:") + Repeat(Glyph(":dash:"), 40) + Glyph(":tee:") + "\n")
	case RenderSimple:
		_, _ = sb.WriteString(indent + Repeat(Glyph(":dash:"), 40) + "\n")
	}
}

func (c *Composite) renderSlot(sb io.StringWriter, slot *Slot, depth int) {
	indent := Repeat("  ", depth)
	for _, block := range slot.blocks {
		for _, line := range core.Split(block.Render(), "\n") {
			if line != "" {
				_, _ = sb.WriteString(indent + line + "\n")
			}
		}
	}
	if slot.child != nil {
		slot.child.renderTo(sb, depth+1)
	}
}
