package cli

import core "dappco.re/go"

func TestLayout_ParseVariant_Good(t *core.T) {
	result := ParseVariant("H[LC]F")
	core.RequireTrue(t, result.OK, result.Error())
	composite := result.Value.(*Composite)

	_, hasHeader := composite.regions[RegionHeader]
	_, hasFooter := composite.regions[RegionFooter]
	headerSlot := composite.regions[RegionHeader]
	_, childHasLeft := headerSlot.child.regions[RegionLeft]
	core.AssertTrue(t, hasHeader)
	core.AssertTrue(t, hasFooter)
	core.AssertTrue(t, childHasLeft)
}

func TestLayout_ParseVariant_Bad(t *core.T) {
	result := ParseVariant("X")
	core.AssertFalse(t, result.OK)

	bracketResult := ParseVariant("H[C")
	core.AssertFalse(t, bracketResult.OK)
	core.AssertContains(t, bracketResult.Error(), "unmatched")
}

func TestLayout_ParseVariant_Ugly(t *core.T) {
	result := ParseVariant("")
	core.RequireTrue(t, result.OK, result.Error())
	composite := result.Value.(*Composite)

	core.AssertEqual(t, "", composite.variant)
	core.AssertLen(t, composite.regions, 0)
}

func TestLayout_Composite_Regions_Good(t *core.T) {
	composite := Layout("HCF")
	var regions []Region
	for region := range composite.Regions() {
		regions = append(regions, region)
	}

	core.AssertLen(t, regions, 3)
	core.AssertTrue(t, composite.regions[RegionContent] != nil)
}

func TestLayout_Composite_Regions_Bad(t *core.T) {
	composite := Layout("Z")
	var regions []Region
	for region := range composite.Regions() {
		regions = append(regions, region)
	}

	core.AssertLen(t, regions, 0)
	core.AssertEmpty(t, composite.regions)
}

func TestLayout_Composite_Regions_Ugly(t *core.T) {
	composite := Layout("")
	count := 0
	for range composite.Regions() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEqual(t, "", composite.variant)
}

func TestLayout_Composite_Slots_Good(t *core.T) {
	composite := Layout("HC")
	slots := map[Region]*Slot{}
	for region, slot := range composite.Slots() {
		slots[region] = slot
	}

	core.AssertLen(t, slots, 2)
	core.AssertEqual(t, RegionHeader, slots[RegionHeader].region)
}

func TestLayout_Composite_Slots_Bad(t *core.T) {
	composite := Layout("Z")
	seen := false
	for range composite.Slots() {
		seen = true
	}

	core.AssertFalse(t, seen)
	core.AssertEmpty(t, composite.regions)
}

func TestLayout_Composite_Slots_Ugly(t *core.T) {
	composite := Layout("H[LC]")
	slot := composite.regions[RegionHeader]
	for range composite.Slots() {
		break
	}

	core.AssertNotNil(t, slot.child)
	core.AssertEqual(t, RegionHeader, slot.region)
}

func TestLayout_StringBlock_Render_Good(t *core.T) {
	got := StringBlock("hello").Render()

	core.AssertEqual(t, "hello", got)
	core.AssertNotEmpty(t, got)
}

func TestLayout_StringBlock_Render_Bad(t *core.T) {
	got := StringBlock("").Render()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestLayout_StringBlock_Render_Ugly(t *core.T) {
	cliPlainCLI(t)
	got := StringBlock(":check:").Render()

	core.AssertEqual(t, "[OK]", got)
	core.AssertNotContains(t, got, ":check:")
}

func TestLayout_Layout_Good(t *core.T) {
	composite := Layout("HC")
	composite.H("header").C("body")

	core.AssertContains(t, composite.String(), "header")
	core.AssertContains(t, composite.String(), "body")
}

func TestLayout_Layout_Bad(t *core.T) {
	composite := Layout("Z")
	composite.C("body")

	core.AssertEqual(t, "Z", composite.variant)
	core.AssertEqual(t, "", composite.String())
}

func TestLayout_Layout_Ugly(t *core.T) {
	composite := Layout("H[LC]").H("outer")

	core.AssertNotNil(t, composite.regions[RegionHeader].child)
	core.AssertContains(t, composite.String(), "outer")
}

func TestLayout_Composite_H_Good(t *core.T) {
	composite := Layout("H")
	composite.H("header")

	core.AssertContains(t, composite.String(), "header")
	core.AssertLen(t, composite.regions[RegionHeader].blocks, 1)
}

func TestLayout_Composite_H_Bad(t *core.T) {
	composite := Layout("C")
	composite.H("header")

	core.AssertNotContains(t, composite.String(), "header")
	core.AssertNil(t, composite.regions[RegionHeader])
}

func TestLayout_Composite_H_Ugly(t *core.T) {
	composite := Layout("H")
	composite.H("one", "two")

	core.AssertContains(t, composite.String(), "one")
	core.AssertContains(t, composite.String(), "two")
}

func TestLayout_Composite_L_Good(t *core.T) {
	composite := Layout("L")
	composite.L("left")

	core.AssertContains(t, composite.String(), "left")
	core.AssertLen(t, composite.regions[RegionLeft].blocks, 1)
}

func TestLayout_Composite_L_Bad(t *core.T) {
	composite := Layout("C")
	composite.L("left")

	core.AssertNotContains(t, composite.String(), "left")
	core.AssertNil(t, composite.regions[RegionLeft])
}

func TestLayout_Composite_L_Ugly(t *core.T) {
	composite := Layout("L")
	composite.L("one", "two")

	core.AssertContains(t, composite.String(), "one")
	core.AssertContains(t, composite.String(), "two")
}

func TestLayout_Composite_C_Good(t *core.T) {
	composite := Layout("C")
	composite.C("content")

	core.AssertContains(t, composite.String(), "content")
	core.AssertLen(t, composite.regions[RegionContent].blocks, 1)
}

func TestLayout_Composite_C_Bad(t *core.T) {
	composite := Layout("H")
	composite.C("content")

	core.AssertNotContains(t, composite.String(), "content")
	core.AssertNil(t, composite.regions[RegionContent])
}

func TestLayout_Composite_C_Ugly(t *core.T) {
	composite := Layout("C")
	composite.C("one", "two")

	core.AssertContains(t, composite.String(), "one")
	core.AssertContains(t, composite.String(), "two")
}

func TestLayout_Composite_R_Good(t *core.T) {
	composite := Layout("R")
	composite.R("right")

	core.AssertContains(t, composite.String(), "right")
	core.AssertLen(t, composite.regions[RegionRight].blocks, 1)
}

func TestLayout_Composite_R_Bad(t *core.T) {
	composite := Layout("C")
	composite.R("right")

	core.AssertNotContains(t, composite.String(), "right")
	core.AssertNil(t, composite.regions[RegionRight])
}

func TestLayout_Composite_R_Ugly(t *core.T) {
	composite := Layout("R")
	composite.R("one", "two")

	core.AssertContains(t, composite.String(), "one")
	core.AssertContains(t, composite.String(), "two")
}

func TestLayout_Composite_F_Good(t *core.T) {
	composite := Layout("F")
	composite.F("footer")

	core.AssertContains(t, composite.String(), "footer")
	core.AssertLen(t, composite.regions[RegionFooter].blocks, 1)
}

func TestLayout_Composite_F_Bad(t *core.T) {
	composite := Layout("C")
	composite.F("footer")

	core.AssertNotContains(t, composite.String(), "footer")
	core.AssertNil(t, composite.regions[RegionFooter])
}

func TestLayout_Composite_F_Ugly(t *core.T) {
	composite := Layout("F")
	composite.F("one", "two")

	core.AssertContains(t, composite.String(), "one")
	core.AssertContains(t, composite.String(), "two")
}
