package frame

import (
	"iter"

	"dappco.re/go"
)

// Region represents one of the 5 HLCRF regions.
type Region rune

const (
	// RegionHeader is the top region of the layout.
	RegionHeader Region = 'H'
	// RegionLeft is the left sidebar region.
	RegionLeft Region = 'L'
	// RegionContent is the main content region.
	RegionContent Region = 'C'
	// RegionRight is the right sidebar region.
	RegionRight Region = 'R'
	// RegionFooter is the bottom region of the layout.
	RegionFooter Region = 'F'
)

// Composite represents an HLCRF layout node.
type Composite struct {
	variant string
	path    string
	regions map[Region]*Slot
	parent  *Composite
}

// Regions returns an iterator over the regions in the composite.
func (c *Composite) Regions() iter.Seq[Region] {
	return func(yield func(Region) bool) {
		for r := range c.regions {
			if !yield(r) {
				return
			}
		}
	}
}

// Slots returns an iterator over the slots in the composite.
func (c *Composite) Slots() iter.Seq2[Region, *Slot] {
	return func(yield func(Region, *Slot) bool) {
		for r, s := range c.regions {
			if !yield(r, s) {
				return
			}
		}
	}
}

// Slot holds content for a region.
type Slot struct {
	region Region
	path   string
	blocks []Renderable
	child  *Composite
}

// Renderable is anything that can be rendered to terminal.
type Renderable interface {
	Render() string
}

// StringBlock is a simple string that implements Renderable.
type StringBlock string

// Render returns the string content.
func (s StringBlock) Render() string { return compileGlyphs(string(s)) }

// Layout creates a new layout from a variant string.
func Layout(variant string) *Composite {
	r := ParseVariant(variant)
	if !r.OK {
		return &Composite{variant: variant, regions: make(map[Region]*Slot)}
	}
	return r.Value.(*Composite)
}

// ParseVariant parses a variant string like "H[LC]C[HCF]F".
func ParseVariant(variant string) core.Result {
	c := &Composite{
		variant: variant,
		path:    "",
		regions: make(map[Region]*Slot),
	}

	i := 0
	for i < len(variant) {
		r := Region(variant[i])
		if !isValidRegion(r) {
			return core.Fail(core.Errorf("ParseVariant: invalid region: %c", r))
		}

		slot := &Slot{region: r, path: string(r)}
		c.regions[r] = slot
		i++

		if i < len(variant) && variant[i] == '[' {
			end := findMatchingBracket(variant, i)
			if end == -1 {
				return core.Fail(core.Errorf("ParseVariant: unmatched bracket at %d", i))
			}
			nestedResult := ParseVariant(variant[i+1 : end])
			if !nestedResult.OK {
				return nestedResult
			}
			nested := nestedResult.Value.(*Composite)
			nested.path = string(r) + "-"
			nested.parent = c
			slot.child = nested
			i = end + 1
		}
	}
	return core.Ok(c)
}

func isValidRegion(r Region) bool {
	return r == 'H' || r == 'L' || r == 'C' || r == 'R' || r == 'F'
}

func findMatchingBracket(s string, start int) int {
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// H adds content to Header region.
func (c *Composite) H(items ...any) *Composite { c.addToRegion(RegionHeader, items...); return c }

// L adds content to Left region.
func (c *Composite) L(items ...any) *Composite { c.addToRegion(RegionLeft, items...); return c }

// C adds content to Content region.
func (c *Composite) C(items ...any) *Composite { c.addToRegion(RegionContent, items...); return c }

// R adds content to Right region.
func (c *Composite) R(items ...any) *Composite { c.addToRegion(RegionRight, items...); return c }

// F adds content to Footer region.
func (c *Composite) F(items ...any) *Composite { c.addToRegion(RegionFooter, items...); return c }

func (c *Composite) addToRegion(r Region, items ...any) {
	slot, ok := c.regions[r]
	if !ok {
		return
	}
	for _, item := range items {
		slot.blocks = append(slot.blocks, toRenderable(item))
	}
}

func toRenderable(item any) Renderable {
	switch v := item.(type) {
	case Renderable:
		return v
	case string:
		return StringBlock(v)
	default:
		return StringBlock(core.Sprint(v))
	}
}

func layoutHasRegion(c *Composite, target Region) bool {
	if c == nil {
		return false
	}
	_, ok := c.regions[target]
	return ok
}
