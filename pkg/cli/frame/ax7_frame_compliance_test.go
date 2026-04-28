package frame

import (
	"bytes"
	"strings"
	"time"

	core "dappco.re/go"
	tea "github.com/charmbracelet/bubbletea"
)

type ax7FrameModel string

func (m ax7FrameModel) View(_, _ int) string { return string(m) }

type ax7FrameInteractive struct {
	view    string
	updates int
}

func (m *ax7FrameInteractive) View(_, _ int) string { return m.view }
func (m *ax7FrameInteractive) Init() tea.Cmd        { return nil }
func (m *ax7FrameInteractive) Update(_ tea.Msg) (FrameModel, tea.Cmd) {
	m.updates++
	return m, nil
}

func ax7FrameWithBuffer() (*Frame, *bytes.Buffer) {
	var buf bytes.Buffer
	f := NewFrame("HLCF").WithOutput(&buf)
	return f, &buf
}

func TestAX7Frame_ModelFunc_View_Good(t *core.T) {
	m := ModelFunc(func(width, height int) string { return core.Sprintf("%dx%d", width, height) })

	core.AssertEqual(t, "4x2", m.View(4, 2))
	core.AssertContains(t, m.View(3, 5), "3x5")
}

func TestAX7Frame_ModelFunc_View_Bad(t *core.T) {
	m := ModelFunc(func(_, _ int) string { return "" })

	core.AssertEqual(t, "", m.View(0, 0))
	core.AssertEmpty(t, m.View(-1, -1))
}

func TestAX7Frame_ModelFunc_View_Ugly(t *core.T) {
	m := ModelFunc(func(width, _ int) string { return strings.Repeat("x", width) })

	core.AssertEqual(t, "xxx", m.View(3, 1))
	core.AssertNotPanics(t, func() { _ = m.View(0, 0) })
}

func TestAX7Frame_NewFrame_Good(t *core.T) {
	f := NewFrame("HCF")

	core.AssertEqual(t, "HCF", f.variant)
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestAX7Frame_NewFrame_Bad(t *core.T) {
	f := NewFrame("Z")

	core.AssertEqual(t, "Z", f.variant)
	core.AssertEmpty(t, f.String())
}

func TestAX7Frame_NewFrame_Ugly(t *core.T) {
	f := NewFrame("")

	core.AssertEqual(t, "", f.variant)
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestAX7Frame_Frame_WithOutput_Good(t *core.T) {
	f, buf := ax7FrameWithBuffer()

	core.AssertEqual(t, f, f.WithOutput(buf))
	core.AssertEqual(t, buf, f.out)
}

func TestAX7Frame_Frame_WithOutput_Bad(t *core.T) {
	f := NewFrame("C")
	original := f.out

	core.AssertEqual(t, f, f.WithOutput(nil))
	core.AssertEqual(t, original, f.out)
}

func TestAX7Frame_Frame_WithOutput_Ugly(t *core.T) {
	f := NewFrame("C")
	var buf bytes.Buffer

	f.WithOutput(&buf).Content(StaticModel("content")).Run()
	core.AssertContains(t, buf.String(), "content")
}

func TestAX7Frame_Frame_Header_Good(t *core.T) {
	f := NewFrame("HC").Header(StaticModel("header"))

	core.AssertContains(t, f.String(), "header")
	core.AssertNotNil(t, f.models[RegionHeader])
}

func TestAX7Frame_Frame_Header_Bad(t *core.T) {
	f := NewFrame("C").Header(StaticModel("ignored"))

	core.AssertNotContains(t, f.String(), "ignored")
	core.AssertNotNil(t, f.models[RegionHeader])
}

func TestAX7Frame_Frame_Header_Ugly(t *core.T) {
	f := NewFrame("HC").Header(StaticModel(":check: header"))

	core.AssertContains(t, f.String(), "header")
	core.AssertContains(t, f.String(), "✓")
}

func TestAX7Frame_Frame_Left_Good(t *core.T) {
	f := NewFrame("LC").Left(StaticModel("left")).Content(StaticModel("content"))

	core.AssertContains(t, f.String(), "left")
	core.AssertContains(t, f.String(), "content")
}

func TestAX7Frame_Frame_Left_Bad(t *core.T) {
	f := NewFrame("C").Left(StaticModel("left"))

	core.AssertNotContains(t, f.String(), "left")
	core.AssertNotNil(t, f.models[RegionLeft])
}

func TestAX7Frame_Frame_Left_Ugly(t *core.T) {
	f := NewFrame("LC").Left(StaticModel(""))

	core.AssertNotPanics(t, func() { _ = f.String() })
	core.AssertEqual(t, "", f.models[RegionLeft].View(1, 1))
}

func TestAX7Frame_Frame_Content_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertContains(t, f.String(), "content")
	core.AssertNotNil(t, f.models[RegionContent])
}

func TestAX7Frame_Frame_Content_Bad(t *core.T) {
	f := NewFrame("H").Content(StaticModel("content"))

	core.AssertContains(t, f.String(), "content")
	core.AssertNotNil(t, f.models[RegionContent])
}

func TestAX7Frame_Frame_Content_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel(":warn:"))

	core.AssertContains(t, f.String(), "⚠")
	core.AssertNotEmpty(t, f.String())
}

func TestAX7Frame_Frame_Right_Good(t *core.T) {
	f := NewFrame("CR").Content(StaticModel("content")).Right(StaticModel("right"))

	core.AssertContains(t, f.String(), "right")
	core.AssertContains(t, f.String(), "content")
}

func TestAX7Frame_Frame_Right_Bad(t *core.T) {
	f := NewFrame("C").Right(StaticModel("right"))

	core.AssertNotContains(t, f.String(), "right")
	core.AssertNotNil(t, f.models[RegionRight])
}

func TestAX7Frame_Frame_Right_Ugly(t *core.T) {
	f := NewFrame("CR").Right(StaticModel(""))

	core.AssertNotPanics(t, func() { _ = f.String() })
	core.AssertEqual(t, "", f.models[RegionRight].View(1, 1))
}

func TestAX7Frame_Frame_Footer_Good(t *core.T) {
	f := NewFrame("CF").Content(StaticModel("content")).Footer(StaticModel("footer"))

	core.AssertContains(t, f.String(), "footer")
	core.AssertContains(t, f.String(), "content")
}

func TestAX7Frame_Frame_Footer_Bad(t *core.T) {
	f := NewFrame("C").Footer(StaticModel("footer"))

	core.AssertNotContains(t, f.String(), "footer")
	core.AssertNotNil(t, f.models[RegionFooter])
}

func TestAX7Frame_Frame_Footer_Ugly(t *core.T) {
	f := NewFrame("CF").Footer(StaticModel(":info: footer"))

	core.AssertContains(t, f.String(), "footer")
	core.AssertContains(t, f.String(), "ℹ")
}

func TestAX7Frame_Frame_Navigate_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel("second"))

	core.AssertContains(t, f.String(), "second")
	core.AssertLen(t, f.history, 1)
}

func TestAX7Frame_Frame_Navigate_Bad(t *core.T) {
	f := NewFrame("C")
	f.Navigate(StaticModel("second"))

	core.AssertContains(t, f.String(), "second")
	core.AssertEmpty(t, f.history)
}

func TestAX7Frame_Frame_Navigate_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel(""))

	core.AssertNotContains(t, f.String(), "first")
	core.AssertLen(t, f.history, 1)
}

func TestAX7Frame_Frame_Back_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel("second"))

	core.AssertTrue(t, f.Back())
	core.AssertContains(t, f.String(), "first")
}

func TestAX7Frame_Frame_Back_Bad(t *core.T) {
	f := NewFrame("C")

	core.AssertFalse(t, f.Back())
	core.AssertEmpty(t, f.history)
}

func TestAX7Frame_Frame_Back_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel("second"))
	f.Back()

	core.AssertFalse(t, f.Back())
	core.AssertContains(t, f.String(), "first")
}

func TestAX7Frame_Frame_Stop_Good(t *core.T) {
	f := NewFrame("C")
	f.Stop()

	_, closed := <-f.done
	core.AssertFalse(t, closed)
}

func TestAX7Frame_Frame_Stop_Bad(t *core.T) {
	f := NewFrame("C")
	f.Stop()

	core.AssertNotPanics(t, func() { f.Stop() })
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestAX7Frame_Frame_Stop_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertNotPanics(t, func() { f.RunFor(time.Nanosecond) })
	core.AssertContains(t, f.String(), "content")
}

func TestAX7Frame_Frame_Send_Good(t *core.T) {
	f := NewFrame("C")

	core.AssertNotPanics(t, func() { f.Send(struct{}{}) })
	core.AssertNil(t, f.program)
}

func TestAX7Frame_Frame_Send_Bad(t *core.T) {
	var f *Frame

	core.AssertPanics(t, func() { f.Send(struct{}{}) })
	core.AssertNil(t, f)
}

func TestAX7Frame_Frame_Send_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertNotPanics(t, func() { f.Send(nil) })
	core.AssertContains(t, f.String(), "content")
}

func TestAX7Frame_Frame_WithKeyMap_Good(t *core.T) {
	f := NewFrame("C")
	km := DefaultKeyMap()

	core.AssertEqual(t, f, f.WithKeyMap(km))
	core.AssertEqual(t, km, f.keyMap)
}

func TestAX7Frame_Frame_WithKeyMap_Bad(t *core.T) {
	f := NewFrame("C")
	f.WithKeyMap(KeyMap{})

	core.AssertEqual(t, tea.KeyType(0), f.keyMap.Quit)
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestAX7Frame_Frame_WithKeyMap_Ugly(t *core.T) {
	f := NewFrame("HC")
	f.WithKeyMap(DefaultKeyMap())

	core.AssertNotPanics(t, func() { _, _ = f.Update(tea.KeyMsg{Type: f.keyMap.FocusNext}) })
	core.AssertEqual(t, RegionHeader, f.Focused())
}

func TestAX7Frame_Frame_Focused_Good(t *core.T) {
	f := NewFrame("HC")
	f.Focus(RegionHeader)

	core.AssertEqual(t, RegionHeader, f.Focused())
	core.AssertNotEqual(t, RegionContent, f.Focused())
}

func TestAX7Frame_Frame_Focused_Bad(t *core.T) {
	f := NewFrame("C")
	f.Focus(RegionHeader)

	core.AssertEqual(t, RegionContent, f.Focused())
	core.AssertNotEqual(t, RegionHeader, f.Focused())
}

func TestAX7Frame_Frame_Focused_Ugly(t *core.T) {
	f := NewFrame("")
	f.Focus(RegionContent)

	core.AssertEqual(t, RegionContent, f.Focused())
	core.AssertEmpty(t, f.layout.regions)
}

func TestAX7Frame_Frame_Focus_Good(t *core.T) {
	f := NewFrame("HCF")
	f.Focus(RegionFooter)

	core.AssertEqual(t, RegionFooter, f.Focused())
	core.AssertNotEqual(t, RegionContent, f.Focused())
}

func TestAX7Frame_Frame_Focus_Bad(t *core.T) {
	f := NewFrame("C")
	f.Focus(RegionLeft)

	core.AssertEqual(t, RegionContent, f.Focused())
	core.AssertFalse(t, layoutHasRegion(f.layout, RegionLeft))
}

func TestAX7Frame_Frame_Focus_Ugly(t *core.T) {
	f := NewFrame("HLCF")
	f.Focus(RegionLeft)

	core.AssertEqual(t, RegionLeft, f.Focused())
	core.AssertTrue(t, layoutHasRegion(f.layout, RegionLeft))
}

func TestAX7Frame_Frame_Init_Good(t *core.T) {
	f := NewFrame("C").Content(&ax7FrameInteractive{view: "content"})
	cmd := f.Init()

	core.AssertNil(t, cmd)
	core.AssertNotPanics(t, func() { _ = f.Init() })
}

func TestAX7Frame_Frame_Init_Bad(t *core.T) {
	f := NewFrame("C")
	cmd := f.Init()

	core.AssertNil(t, cmd)
	core.AssertNotPanics(t, func() { _ = f.Init() })
}

func TestAX7Frame_Frame_Init_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("plain"))
	cmd := f.Init()

	core.AssertNil(t, cmd)
	core.AssertNotPanics(t, func() { _ = f.Init() })
}

func TestAX7Frame_Frame_Update_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))
	model, cmd := f.Update(tea.WindowSizeMsg{Width: 40, Height: 10})

	core.AssertEqual(t, f, model)
	core.AssertNil(t, cmd)
}

func TestAX7Frame_Frame_Update_Bad(t *core.T) {
	f := NewFrame("C")
	model, cmd := f.Update(tea.KeyMsg{Type: DefaultKeyMap().Back})

	core.AssertEqual(t, f, model)
	core.AssertNil(t, cmd)
}

func TestAX7Frame_Frame_Update_Ugly(t *core.T) {
	m := &ax7FrameInteractive{view: "content"}
	f := NewFrame("C").Content(m)
	_, _ = f.Update(struct{}{})

	core.AssertEqual(t, 1, m.updates)
	core.AssertContains(t, f.String(), "content")
}

func TestAX7Frame_Frame_View_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertContains(t, f.View(), "content")
	core.AssertNotEmpty(t, f.View())
}

func TestAX7Frame_Frame_View_Bad(t *core.T) {
	f := NewFrame("C")

	core.AssertEqual(t, "", f.View())
	core.AssertEmpty(t, f.View())
}

func TestAX7Frame_Frame_View_Ugly(t *core.T) {
	f := NewFrame("HCF").Header(StaticModel("h")).Content(StaticModel("c")).Footer(StaticModel("f"))

	core.AssertContains(t, f.View(), "h")
	core.AssertContains(t, f.View(), "f")
}

func TestAX7Frame_Frame_Run_Good(t *core.T) {
	f, buf := ax7FrameWithBuffer()
	f.Content(StaticModel("content"))

	f.Run()
	core.AssertContains(t, buf.String(), "content")
}

func TestAX7Frame_Frame_Run_Bad(t *core.T) {
	f, buf := ax7FrameWithBuffer()

	f.Run()
	core.AssertEqual(t, "", buf.String())
}

func TestAX7Frame_Frame_Run_Ugly(t *core.T) {
	f, buf := ax7FrameWithBuffer()
	f.Content(StaticModel(":check:"))

	f.Run()
	core.AssertContains(t, buf.String(), "✓")
}

func TestAX7Frame_Frame_RunFor_Good(t *core.T) {
	f, buf := ax7FrameWithBuffer()
	f.Content(StaticModel("content"))

	f.RunFor(time.Nanosecond)
	core.AssertContains(t, buf.String(), "content")
}

func TestAX7Frame_Frame_RunFor_Bad(t *core.T) {
	f, buf := ax7FrameWithBuffer()

	f.RunFor(0)
	core.AssertEqual(t, "", buf.String())
}

func TestAX7Frame_Frame_RunFor_Ugly(t *core.T) {
	f, buf := ax7FrameWithBuffer()
	f.Content(StaticModel(""))

	f.RunFor(time.Nanosecond)
	core.AssertEqual(t, "", buf.String())
}

func TestAX7Frame_Frame_String_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertEqual(t, "content\n", f.String())
	core.AssertContains(t, f.String(), "\n")
}

func TestAX7Frame_Frame_String_Bad(t *core.T) {
	f := NewFrame("C")

	core.AssertEqual(t, "", f.String())
	core.AssertEmpty(t, f.String())
}

func TestAX7Frame_Frame_String_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("\033[31mred\033[0m"))

	core.AssertEqual(t, "red\n", f.String())
	core.AssertNotContains(t, f.String(), "\033")
}

func TestAX7Frame_Composite_Regions_Good(t *core.T) {
	c := Layout("HC")
	var regions []Region
	for r := range c.Regions() {
		regions = append(regions, r)
	}

	core.AssertLen(t, regions, 2)
	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
}

func TestAX7Frame_Composite_Regions_Bad(t *core.T) {
	c := Layout("Z")
	var regions []Region
	for r := range c.Regions() {
		regions = append(regions, r)
	}

	core.AssertEmpty(t, regions)
	core.AssertFalse(t, layoutHasRegion(c, RegionContent))
}

func TestAX7Frame_Composite_Regions_Ugly(t *core.T) {
	c := Layout("HH")
	var count int
	for range c.Regions() {
		count++
	}

	core.AssertEqual(t, 1, count)
	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
}

func TestAX7Frame_Composite_Slots_Good(t *core.T) {
	c := Layout("CF")
	var count int
	for _, slot := range c.Slots() {
		core.AssertNotNil(t, slot)
		count++
	}

	core.AssertEqual(t, 2, count)
	core.AssertTrue(t, layoutHasRegion(c, RegionFooter))
}

func TestAX7Frame_Composite_Slots_Bad(t *core.T) {
	c := Layout("Z")
	var count int
	for range c.Slots() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, c.regions)
}

func TestAX7Frame_Composite_Slots_Ugly(t *core.T) {
	c := Layout("C[HF]")
	var child *Composite
	for _, slot := range c.Slots() {
		child = slot.child
	}

	core.AssertNotNil(t, child)
	core.AssertTrue(t, layoutHasRegion(child, RegionHeader))
}

func TestAX7Frame_StringBlock_Render_Good(t *core.T) {
	got := StringBlock(":check: ready").Render()

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "✓")
}

func TestAX7Frame_StringBlock_Render_Bad(t *core.T) {
	got := StringBlock("").Render()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7Frame_StringBlock_Render_Ugly(t *core.T) {
	got := StringBlock(":unknown:").Render()

	core.AssertEqual(t, ":unknown:", got)
	core.AssertContains(t, got, "unknown")
}

func TestAX7Frame_Layout_Good(t *core.T) {
	c := Layout("HCF")

	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
	core.AssertTrue(t, layoutHasRegion(c, RegionFooter))
}

func TestAX7Frame_Layout_Bad(t *core.T) {
	c := Layout("Z")

	core.AssertEqual(t, "Z", c.variant)
	core.AssertEmpty(t, c.regions)
}

func TestAX7Frame_Layout_Ugly(t *core.T) {
	c := Layout("C[HF]")

	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
	core.AssertNotNil(t, c.regions[RegionContent].child)
}

func TestAX7Frame_ParseVariant_Good(t *core.T) {
	c, err := ParseVariant("HCF")

	core.AssertNoError(t, err)
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestAX7Frame_ParseVariant_Bad(t *core.T) {
	c, err := ParseVariant("Z")

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestAX7Frame_ParseVariant_Ugly(t *core.T) {
	c, err := ParseVariant("C[HF")

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestAX7Frame_Composite_H_Good(t *core.T) {
	c := Layout("H").H("header")

	core.AssertEqual(t, c, c.H("more"))
	core.AssertLen(t, c.regions[RegionHeader].blocks, 2)
}

func TestAX7Frame_Composite_H_Bad(t *core.T) {
	c := Layout("C").H("header")

	core.AssertFalse(t, layoutHasRegion(c, RegionHeader))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestAX7Frame_Composite_H_Ugly(t *core.T) {
	c := Layout("H").H(123)

	core.AssertEqual(t, "123", c.regions[RegionHeader].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionHeader].blocks, 1)
}

func TestAX7Frame_Composite_L_Good(t *core.T) {
	c := Layout("L").L("left")

	core.AssertEqual(t, c, c.L("more"))
	core.AssertLen(t, c.regions[RegionLeft].blocks, 2)
}

func TestAX7Frame_Composite_L_Bad(t *core.T) {
	c := Layout("C").L("left")

	core.AssertFalse(t, layoutHasRegion(c, RegionLeft))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestAX7Frame_Composite_L_Ugly(t *core.T) {
	c := Layout("L").L(StringBlock(":check:"))

	core.AssertEqual(t, "✓", c.regions[RegionLeft].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionLeft].blocks, 1)
}

func TestAX7Frame_Composite_C_Good(t *core.T) {
	c := Layout("C").C("content")

	core.AssertEqual(t, c, c.C("more"))
	core.AssertLen(t, c.regions[RegionContent].blocks, 2)
}

func TestAX7Frame_Composite_C_Bad(t *core.T) {
	c := Layout("H").C("content")

	core.AssertFalse(t, layoutHasRegion(c, RegionContent))
	core.AssertTrue(t, layoutHasRegion(c, RegionHeader))
}

func TestAX7Frame_Composite_C_Ugly(t *core.T) {
	c := Layout("C").C("")

	core.AssertEqual(t, "", c.regions[RegionContent].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionContent].blocks, 1)
}

func TestAX7Frame_Composite_R_Good(t *core.T) {
	c := Layout("R").R("right")

	core.AssertEqual(t, c, c.R("more"))
	core.AssertLen(t, c.regions[RegionRight].blocks, 2)
}

func TestAX7Frame_Composite_R_Bad(t *core.T) {
	c := Layout("C").R("right")

	core.AssertFalse(t, layoutHasRegion(c, RegionRight))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestAX7Frame_Composite_R_Ugly(t *core.T) {
	c := Layout("R").R(RegionRight)

	core.AssertEqual(t, "82", c.regions[RegionRight].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionRight].blocks, 1)
}

func TestAX7Frame_Composite_F_Good(t *core.T) {
	c := Layout("F").F("footer")

	core.AssertEqual(t, c, c.F("more"))
	core.AssertLen(t, c.regions[RegionFooter].blocks, 2)
}

func TestAX7Frame_Composite_F_Bad(t *core.T) {
	c := Layout("C").F("footer")

	core.AssertFalse(t, layoutHasRegion(c, RegionFooter))
	core.AssertTrue(t, layoutHasRegion(c, RegionContent))
}

func TestAX7Frame_Composite_F_Ugly(t *core.T) {
	c := Layout("F").F(nil)

	core.AssertEqual(t, "<nil>", c.regions[RegionFooter].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionFooter].blocks, 1)
}

func TestAX7Frame_StatusLine_Good(t *core.T) {
	got := StatusLine("core", "ready").View(80, 1)

	core.AssertContains(t, got, "core")
	core.AssertContains(t, got, "ready")
}

func TestAX7Frame_StatusLine_Bad(t *core.T) {
	got := StatusLine("").View(80, 1)

	core.AssertNotContains(t, got, "ready")
	core.AssertEqual(t, "", got)
}

func TestAX7Frame_StatusLine_Ugly(t *core.T) {
	got := StatusLine("abcdef").View(3, 1)

	core.AssertEqual(t, "abc", got)
	core.AssertLen(t, got, 3)
}

func TestAX7Frame_LineModel_View_Good(t *core.T) {
	m := &statusLineModel{title: "core", pairs: []string{"ready"}}

	core.AssertContains(t, m.View(80, 1), "core")
	core.AssertContains(t, m.View(80, 1), "ready")
}

func TestAX7Frame_LineModel_View_Bad(t *core.T) {
	m := &statusLineModel{}

	core.AssertEqual(t, "", m.View(80, 1))
	core.AssertEmpty(t, m.View(0, 1))
}

func TestAX7Frame_LineModel_View_Ugly(t *core.T) {
	m := &statusLineModel{title: ":check:"}

	core.AssertEqual(t, "✓", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "✓")
}

func TestAX7Frame_KeyHints_Good(t *core.T) {
	got := KeyHints("q quit", "enter open").View(80, 1)

	core.AssertContains(t, got, "q quit")
	core.AssertContains(t, got, "enter open")
}

func TestAX7Frame_KeyHints_Bad(t *core.T) {
	got := KeyHints().View(80, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7Frame_KeyHints_Ugly(t *core.T) {
	got := KeyHints("abcdef").View(2, 1)

	core.AssertEqual(t, "ab", got)
	core.AssertLen(t, got, 2)
}

func TestAX7Frame_HintsModel_View_Good(t *core.T) {
	m := &keyHintsModel{hints: []string{"tab focus"}}

	core.AssertContains(t, m.View(80, 1), "tab")
	core.AssertContains(t, m.View(80, 1), "focus")
}

func TestAX7Frame_HintsModel_View_Bad(t *core.T) {
	m := &keyHintsModel{}

	core.AssertEqual(t, "", m.View(80, 1))
	core.AssertEmpty(t, m.View(0, 1))
}

func TestAX7Frame_HintsModel_View_Ugly(t *core.T) {
	m := &keyHintsModel{hints: []string{":cross:"}}

	core.AssertEqual(t, "✗", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "✗")
}

func TestAX7Frame_Breadcrumb_Good(t *core.T) {
	got := Breadcrumb("core", "dev").View(80, 1)

	core.AssertContains(t, got, "core")
	core.AssertContains(t, got, "dev")
}

func TestAX7Frame_Breadcrumb_Bad(t *core.T) {
	got := Breadcrumb().View(80, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7Frame_Breadcrumb_Ugly(t *core.T) {
	got := Breadcrumb("abcdef").View(3, 1)

	core.AssertEqual(t, "abc", got)
	core.AssertLen(t, got, 3)
}

func TestAX7Frame_StaticModel_Good(t *core.T) {
	m := StaticModel("content")
	got := m.View(80, 1)

	core.AssertEqual(t, "content", got)
	core.AssertContains(t, got, "content")
}

func TestAX7Frame_StaticModel_Bad(t *core.T) {
	m := StaticModel("")
	got := m.View(80, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7Frame_StaticModel_Ugly(t *core.T) {
	m := StaticModel(":check:")
	got := m.View(80, 1)

	core.AssertEqual(t, "✓", got)
	core.AssertContains(t, got, "✓")
}

func TestAX7Frame_Model_View_Good(t *core.T) {
	m := &staticModel{text: "content"}

	core.AssertEqual(t, "content", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "content")
}

func TestAX7Frame_Model_View_Bad(t *core.T) {
	m := &staticModel{}

	core.AssertEqual(t, "", m.View(80, 1))
	core.AssertEmpty(t, m.View(0, 0))
}

func TestAX7Frame_Model_View_Ugly(t *core.T) {
	m := &staticModel{text: ":warn:"}

	core.AssertEqual(t, "⚠", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "⚠")
}

func TestAX7Frame_Adapter_View_Good(t *core.T) {
	a := &modelAdapter{m: ax7FrameModel("content")}

	core.AssertEqual(t, "content", a.View(80, 1))
	core.AssertContains(t, a.View(80, 1), "content")
}

func TestAX7Frame_Adapter_View_Bad(t *core.T) {
	a := &modelAdapter{}

	core.AssertPanics(t, func() { _ = a.View(80, 1) })
	core.AssertNil(t, a.m)
}

func TestAX7Frame_Adapter_View_Ugly(t *core.T) {
	a := &modelAdapter{m: ax7FrameModel("")}

	core.AssertEqual(t, "", a.View(0, 0))
	core.AssertEmpty(t, a.View(0, 0))
}

func TestAX7Frame_Adapter_Init_Good(t *core.T) {
	a := &modelAdapter{m: ax7FrameModel("content")}

	core.AssertNil(t, a.Init())
	core.AssertNotNil(t, a.m)
}

func TestAX7Frame_Adapter_Init_Bad(t *core.T) {
	a := &modelAdapter{}

	core.AssertNil(t, a.Init())
	core.AssertNil(t, a.m)
}

func TestAX7Frame_Adapter_Init_Ugly(t *core.T) {
	var a *modelAdapter

	core.AssertNil(t, a.Init())
	core.AssertNil(t, a)
}

func TestAX7Frame_Adapter_Update_Good(t *core.T) {
	a := &modelAdapter{m: ax7FrameModel("content")}
	updated, cmd := a.Update(struct{}{})

	core.AssertEqual(t, a, updated)
	core.AssertNil(t, cmd)
}

func TestAX7Frame_Adapter_Update_Bad(t *core.T) {
	a := &modelAdapter{}
	updated, cmd := a.Update(nil)

	core.AssertEqual(t, a, updated)
	core.AssertNil(t, cmd)
}

func TestAX7Frame_Adapter_Update_Ugly(t *core.T) {
	var a *modelAdapter

	updated, cmd := a.Update(nil)
	core.AssertNil(t, updated)
	core.AssertNil(t, cmd)
}

func TestAX7Frame_DefaultKeyMap_Good(t *core.T) {
	km := DefaultKeyMap()

	core.AssertEqual(t, tea.KeyTab, km.FocusNext)
	core.AssertEqual(t, tea.KeyCtrlC, km.Quit)
}

func TestAX7Frame_DefaultKeyMap_Bad(t *core.T) {
	km := KeyMap{}

	core.AssertNotEqual(t, DefaultKeyMap(), km)
	core.AssertEqual(t, tea.KeyType(0), km.Quit)
}

func TestAX7Frame_DefaultKeyMap_Ugly(t *core.T) {
	km := DefaultKeyMap()

	core.AssertNotEqual(t, km.FocusNext, km.FocusPrev)
	core.AssertNotEqual(t, km.FocusLeft, km.FocusRight)
}

func TestAX7Frame_ColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestAX7Frame_ColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestAX7Frame_ColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(!original)
	defer SetColorEnabled(original)

	core.AssertEqual(t, !original, ColorEnabled())
}

func TestAX7Frame_SetColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestAX7Frame_SetColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestAX7Frame_SetColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestAX7Frame_NewStyle_Good(t *core.T) {
	s := NewStyle()

	core.AssertNotNil(t, s)
	core.AssertEqual(t, "text", s.Render("text"))
}

func TestAX7Frame_NewStyle_Bad(t *core.T) {
	s := NewStyle()

	core.AssertFalse(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestAX7Frame_NewStyle_Ugly(t *core.T) {
	s := NewStyle().Bold().Dim()

	core.AssertTrue(t, s.bold)
	core.AssertTrue(t, s.dim)
}

func TestAX7Frame_AnsiStyle_Bold_Good(t *core.T) {
	s := NewStyle().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertEqual(t, s, s.Bold())
}

func TestAX7Frame_AnsiStyle_Bold_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Bold() })
	core.AssertNil(t, s)
}

func TestAX7Frame_AnsiStyle_Bold_Ugly(t *core.T) {
	s := NewStyle().Bold().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestAX7Frame_AnsiStyle_Dim_Good(t *core.T) {
	s := NewStyle().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertEqual(t, s, s.Dim())
}

func TestAX7Frame_AnsiStyle_Dim_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Dim() })
	core.AssertNil(t, s)
}

func TestAX7Frame_AnsiStyle_Dim_Ugly(t *core.T) {
	s := NewStyle().Dim().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertFalse(t, s.bold)
}

func TestAX7Frame_AnsiStyle_Foreground_Good(t *core.T) {
	s := NewStyle().Foreground("#ff0000")

	core.AssertContains(t, s.fg, "38;2;255;0;0")
	core.AssertEqual(t, s, s.Foreground("#00ff00"))
}

func TestAX7Frame_AnsiStyle_Foreground_Bad(t *core.T) {
	s := NewStyle().Foreground("bad")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestAX7Frame_AnsiStyle_Foreground_Ugly(t *core.T) {
	s := NewStyle().Foreground("")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestAX7Frame_AnsiStyle_Render_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertContains(t, got, "text")
	core.AssertContains(t, got, "\033[1m")
}

func TestAX7Frame_AnsiStyle_Render_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestAX7Frame_AnsiStyle_Render_Ugly(t *core.T) {
	var s *AnsiStyle
	got := s.Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestAX7Frame_Truncate_Good(t *core.T) {
	got := Truncate("abcdef", 4)

	core.AssertEqual(t, "a...", got)
	core.AssertLen(t, got, 4)
}

func TestAX7Frame_Truncate_Bad(t *core.T) {
	got := Truncate("abcdef", 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7Frame_Truncate_Ugly(t *core.T) {
	got := Truncate("abcdef", 2)

	core.AssertEqual(t, "ab", got)
	core.AssertLen(t, got, 2)
}

func TestAX7Frame_Glyph_Good(t *core.T) {
	got := Glyph(":check:")

	core.AssertEqual(t, "✓", got)
	core.AssertNotEqual(t, ":check:", got)
}

func TestAX7Frame_Glyph_Bad(t *core.T) {
	got := Glyph(":missing:")

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestAX7Frame_Glyph_Ugly(t *core.T) {
	got := Glyph("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}
