package cli

import (
	core "dappco.re/go"
)

func TestRender_UseRenderFlat_Good(t *core.T) {
	UseRenderSimple()
	UseRenderFlat()

	core.AssertEqual(t, RenderFlat, currentRenderStyle)
	core.AssertNotEqual(t, RenderSimple, currentRenderStyle)
}

func TestRender_UseRenderFlat_Bad(t *core.T) {
	currentRenderStyle = RenderBoxed
	UseRenderFlat()

	core.AssertEqual(t, RenderFlat, currentRenderStyle)
	core.AssertNotEqual(t, RenderBoxed, currentRenderStyle)
}

func TestRender_UseRenderFlat_Ugly(t *core.T) {
	UseRenderFlat()
	UseRenderFlat()

	core.AssertEqual(t, RenderFlat, currentRenderStyle)
	core.AssertNotPanics(t, UseRenderFlat)
}

func TestRender_UseRenderSimple_Good(t *core.T) {
	UseRenderSimple()

	core.AssertEqual(t, RenderSimple, currentRenderStyle)
	core.AssertNotEqual(t, RenderFlat, currentRenderStyle)
}

func TestRender_UseRenderSimple_Bad(t *core.T) {
	currentRenderStyle = RenderBoxed
	UseRenderSimple()

	core.AssertEqual(t, RenderSimple, currentRenderStyle)
	core.AssertNotEqual(t, RenderBoxed, currentRenderStyle)
}

func TestRender_UseRenderSimple_Ugly(t *core.T) {
	UseRenderSimple()
	UseRenderSimple()

	core.AssertEqual(t, RenderSimple, currentRenderStyle)
	core.AssertNotPanics(t, UseRenderSimple)
}

func TestRender_UseRenderBoxed_Good(t *core.T) {
	UseRenderBoxed()

	core.AssertEqual(t, RenderBoxed, currentRenderStyle)
	core.AssertNotEqual(t, RenderFlat, currentRenderStyle)
}

func TestRender_UseRenderBoxed_Bad(t *core.T) {
	currentRenderStyle = RenderSimple
	UseRenderBoxed()

	core.AssertEqual(t, RenderBoxed, currentRenderStyle)
	core.AssertNotEqual(t, RenderSimple, currentRenderStyle)
}

func TestRender_UseRenderBoxed_Ugly(t *core.T) {
	UseRenderBoxed()
	UseRenderBoxed()

	core.AssertEqual(t, RenderBoxed, currentRenderStyle)
	core.AssertNotPanics(t, UseRenderBoxed)
}

func TestRender_Composite_Render_Good(t *core.T) {
	out := cliCaptureStdout(t, func() { Layout("C").C("content").Render() })

	core.AssertContains(t, out, "content")
	core.AssertContains(t, out, "\n")
}

func TestRender_Composite_Render_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { Layout("Z").Render() })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestRender_Composite_Render_Ugly(t *core.T) {
	UseRenderSimple()
	out := cliCaptureStdout(t, func() { Layout("HC").H("h").C("c").Render() })

	core.AssertContains(t, out, "h")
	core.AssertContains(t, out, "c")
}

func TestRender_Composite_String_Good(t *core.T) {
	got := Layout("C").C("content").String()

	core.AssertContains(t, got, "content")
	core.AssertContains(t, got, "\n")
}

func TestRender_Composite_String_Bad(t *core.T) {
	got := Layout("Z").String()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestRender_Composite_String_Ugly(t *core.T) {
	UseRenderBoxed()
	got := Layout("HC").H("h").C("c").String()

	core.AssertContains(t, got, "h")
	core.AssertContains(t, got, "c")
}
