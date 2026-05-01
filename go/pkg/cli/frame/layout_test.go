package frame

import (
	core "dappco.re/go"
)

func TestLayout_Composite_Regions_Good(t *core.T) {
	c := Layout("HC")
	var regions []Region
	for r := range c.Regions() {
		regions = append(regions, r)
	}

	core.AssertLen(t, regions, 2)
	core.AssertNotNil(t, c.regions[RegionHeader])
}

func TestLayout_Composite_Regions_Bad(t *core.T) {
	c := Layout("Z")
	var count int
	for range c.Regions() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, c.regions)
}

func TestLayout_Composite_Regions_Ugly(t *core.T) {
	c := Layout("HH")
	var count int
	for range c.Regions() {
		count++
	}

	core.AssertEqual(t, 1, count)
	core.AssertNotNil(t, c.regions[RegionHeader])
}

func TestLayout_Composite_Slots_Good(t *core.T) {
	c := Layout("CF")
	var count int
	for _, slot := range c.Slots() {
		core.AssertNotNil(t, slot)
		count++
	}

	core.AssertEqual(t, 2, count)
}

func TestLayout_Composite_Slots_Bad(t *core.T) {
	c := Layout("Z")
	var count int
	for range c.Slots() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, c.regions)
}

func TestLayout_Composite_Slots_Ugly(t *core.T) {
	c := Layout("C[HF]")
	var child *Composite
	for _, slot := range c.Slots() {
		child = slot.child
	}

	core.AssertNotNil(t, child)
	core.AssertNotNil(t, child.regions[RegionHeader])
}

func TestLayout_StringBlock_Render_Good(t *core.T) {
	got := StringBlock(":check: ready").Render()

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "✓")
}

func TestLayout_StringBlock_Render_Bad(t *core.T) {
	got := StringBlock("").Render()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestLayout_StringBlock_Render_Ugly(t *core.T) {
	got := StringBlock(":missing:").Render()

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestLayout_Layout_Good(t *core.T) {
	c := Layout("HCF")

	core.AssertNotNil(t, c.regions[RegionHeader])
	core.AssertNotNil(t, c.regions[RegionFooter])
}

func TestLayout_Layout_Bad(t *core.T) {
	c := Layout("Z")

	core.AssertEqual(t, "Z", c.variant)
	core.AssertEmpty(t, c.regions)
}

func TestLayout_Layout_Ugly(t *core.T) {
	c := Layout("C[HF]")

	core.AssertNotNil(t, c.regions[RegionContent])
	core.AssertNotNil(t, c.regions[RegionContent].child)
}

func TestLayout_ParseVariant_Good(t *core.T) {
	result := ParseVariant("HCF")
	c, _ := result.Value.(*Composite)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestLayout_ParseVariant_Bad(t *core.T) {
	result := ParseVariant("Z")
	c, _ := result.Value.(*Composite)
	err := cliResultError(result)

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestLayout_ParseVariant_Ugly(t *core.T) {
	result := ParseVariant("C[HF")
	c, _ := result.Value.(*Composite)
	err := cliResultError(result)

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestLayout_Composite_H_Good(t *core.T) {
	c := Layout("H").H("header")

	core.AssertEqual(t, c, c.H("more"))
	core.AssertLen(t, c.regions[RegionHeader].blocks, 2)
}

func TestLayout_Composite_H_Bad(t *core.T) {
	c := Layout("C").H("header")

	core.AssertNil(t, c.regions[RegionHeader])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestLayout_Composite_H_Ugly(t *core.T) {
	c := Layout("H").H(123)

	core.AssertEqual(t, "123", c.regions[RegionHeader].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionHeader].blocks, 1)
}

func TestLayout_Composite_L_Good(t *core.T) {
	c := Layout("L").L("left")

	core.AssertEqual(t, c, c.L("more"))
	core.AssertLen(t, c.regions[RegionLeft].blocks, 2)
}

func TestLayout_Composite_L_Bad(t *core.T) {
	c := Layout("C").L("left")

	core.AssertNil(t, c.regions[RegionLeft])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestLayout_Composite_L_Ugly(t *core.T) {
	c := Layout("L").L(StringBlock(":check:"))

	core.AssertEqual(t, "✓", c.regions[RegionLeft].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionLeft].blocks, 1)
}

func TestLayout_Composite_C_Good(t *core.T) {
	c := Layout("C").C("content")

	core.AssertEqual(t, c, c.C("more"))
	core.AssertLen(t, c.regions[RegionContent].blocks, 2)
}

func TestLayout_Composite_C_Bad(t *core.T) {
	c := Layout("H").C("content")

	core.AssertNil(t, c.regions[RegionContent])
	core.AssertNotNil(t, c.regions[RegionHeader])
}

func TestLayout_Composite_C_Ugly(t *core.T) {
	c := Layout("C").C("")

	core.AssertEqual(t, "", c.regions[RegionContent].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionContent].blocks, 1)
}

func TestLayout_Composite_R_Good(t *core.T) {
	c := Layout("R").R("right")

	core.AssertEqual(t, c, c.R("more"))
	core.AssertLen(t, c.regions[RegionRight].blocks, 2)
}

func TestLayout_Composite_R_Bad(t *core.T) {
	c := Layout("C").R("right")

	core.AssertNil(t, c.regions[RegionRight])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestLayout_Composite_R_Ugly(t *core.T) {
	c := Layout("R").R(RegionRight)

	core.AssertEqual(t, "82", c.regions[RegionRight].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionRight].blocks, 1)
}

func TestLayout_Composite_F_Good(t *core.T) {
	c := Layout("F").F("footer")

	core.AssertEqual(t, c, c.F("more"))
	core.AssertLen(t, c.regions[RegionFooter].blocks, 2)
}

func TestLayout_Composite_F_Bad(t *core.T) {
	c := Layout("C").F("footer")

	core.AssertNil(t, c.regions[RegionFooter])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestLayout_Composite_F_Ugly(t *core.T) {
	c := Layout("F").F(nil)

	core.AssertEqual(t, "<nil>", c.regions[RegionFooter].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionFooter].blocks, 1)
}

func TestLayout_Composite_Regions_Good(t *core.T) {
	c := Layout("HC")
	var regions []Region
	for r := range c.Regions() {
		regions = append(regions, r)
	}

	core.AssertLen(t, regions, 2)
	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
}

func TestLayout_Composite_Regions_Bad(t *core.T) {
	c := Layout("Z")
	var regions []Region
	for r := range c.Regions() {
		regions = append(regions, r)
	}

	core.AssertEmpty(t, regions)
	core.AssertFalse(t, layoutHasRegion(c, RegionContent))
}

func TestLayout_Composite_Regions_Ugly(t *core.T) {
	c := Layout("HH")
	var count int
	for range c.Regions() {
		count++
	}

	core.AssertEqual(t, 1, count)
	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
}

func TestLayout_Composite_Slots_Good(t *core.T) {
	c := Layout("CF")
	var count int
	for _, slot := range c.Slots() {
		core.AssertNotNil(t, slot)
		count++
	}

	core.AssertEqual(t, 2, count)
	core.AssertTrue(t, layoutHasRegion(c, RegionFooter))
}

func TestLayout_Composite_Slots_Bad(t *core.T) {
	c := Layout("Z")
	var count int
	for range c.Slots() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, c.regions)
}

func TestLayout_Composite_Slots_Ugly(t *core.T) {
	c := Layout("C[HF]")
	var child *Composite
	for _, slot := range c.Slots() {
		child = slot.child
	}

	core.AssertNotNil(t, child)
	core.AssertTrue(t, layoutHasRegion(child, RegionHeader))
}

func TestLayout_StringBlock_Render_Good(t *core.T) {
	got := StringBlock(":check: ready").Render()

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "✓")
}

func TestLayout_StringBlock_Render_Bad(t *core.T) {
	got := StringBlock("").Render()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestLayout_StringBlock_Render_Ugly(t *core.T) {
	got := StringBlock(":unknown:").Render()

	core.AssertEqual(t, ":unknown:", got)
	core.AssertContains(t, got, "unknown")
}

func TestLayout_Layout_Good(t *core.T) {
	c := Layout("HCF")

	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
	core.AssertTrue(t, layoutHasRegion(c, RegionFooter))
}

func TestLayout_Layout_Bad(t *core.T) {
	c := Layout("Z")

	core.AssertEqual(t, "Z", c.variant)
	core.AssertEmpty(t, c.regions)
}

func TestLayout_Layout_Ugly(t *core.T) {
	c := Layout("C[HF]")

	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
	core.AssertNotNil(t, c.regions[RegionContent].child)
}

func TestLayout_ParseVariant_Good(t *core.T) {
	result := ParseVariant("HCF")
	c, _ := result.Value.(*Composite)
	err := cliResultError(result)

	core.AssertNoError(t, err)
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestLayout_ParseVariant_Bad(t *core.T) {
	result := ParseVariant("Z")
	c, _ := result.Value.(*Composite)
	err := cliResultError(result)

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestLayout_ParseVariant_Ugly(t *core.T) {
	result := ParseVariant("C[HF")
	c, _ := result.Value.(*Composite)
	err := cliResultError(result)

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestLayout_Composite_H_Good(t *core.T) {
	c := Layout("H").H("header")

	core.AssertEqual(t, c, c.H("more"))
	core.AssertLen(t, c.regions[RegionHeader].blocks, 2)
}

func TestLayout_Composite_H_Bad(t *core.T) {
	c := Layout("C").H("header")

	core.AssertFalse(t, layoutHasRegion(c, RegionHeader))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestLayout_Composite_H_Ugly(t *core.T) {
	c := Layout("H").H(123)

	core.AssertEqual(t, "123", c.regions[RegionHeader].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionHeader].blocks, 1)
}

func TestLayout_Composite_L_Good(t *core.T) {
	c := Layout("L").L("left")

	core.AssertEqual(t, c, c.L("more"))
	core.AssertLen(t, c.regions[RegionLeft].blocks, 2)
}

func TestLayout_Composite_L_Bad(t *core.T) {
	c := Layout("C").L("left")

	core.AssertFalse(t, layoutHasRegion(c, RegionLeft))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestLayout_Composite_L_Ugly(t *core.T) {
	c := Layout("L").L(StringBlock(":check:"))

	core.AssertEqual(t, "✓", c.regions[RegionLeft].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionLeft].blocks, 1)
}

func TestLayout_Composite_C_Good(t *core.T) {
	c := Layout("C").C("content")

	core.AssertEqual(t, c, c.C("more"))
	core.AssertLen(t, c.regions[RegionContent].blocks, 2)
}

func TestLayout_Composite_C_Bad(t *core.T) {
	c := Layout("H").C("content")

	core.AssertFalse(t, layoutHasRegion(c, RegionContent))
	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
}

func TestLayout_Composite_C_Ugly(t *core.T) {
	c := Layout("C").C("")

	core.AssertEqual(t, "", c.regions[RegionContent].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionContent].blocks, 1)
}

func TestLayout_Composite_R_Good(t *core.T) {
	c := Layout("R").R("right")

	core.AssertEqual(t, c, c.R("more"))
	core.AssertLen(t, c.regions[RegionRight].blocks, 2)
}

func TestLayout_Composite_R_Bad(t *core.T) {
	c := Layout("C").R("right")

	core.AssertFalse(t, layoutHasRegion(c, RegionRight))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestLayout_Composite_R_Ugly(t *core.T) {
	c := Layout("R").R(RegionRight)

	core.AssertEqual(t, "82", c.regions[RegionRight].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionRight].blocks, 1)
}

func TestLayout_Composite_F_Good(t *core.T) {
	c := Layout("F").F("footer")

	core.AssertEqual(t, c, c.F("more"))
	core.AssertLen(t, c.regions[RegionFooter].blocks, 2)
}

func TestLayout_Composite_F_Bad(t *core.T) {
	c := Layout("C").F("footer")

	core.AssertFalse(t, layoutHasRegion(c, RegionFooter))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestLayout_Composite_F_Ugly(t *core.T) {
	c := Layout("F").F(nil)

	core.AssertEqual(t, "<nil>", c.regions[RegionFooter].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionFooter].blocks, 1)
}
