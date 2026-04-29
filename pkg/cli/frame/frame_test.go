package frame

import (
	core "dappco.re/go"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type frameStringWriter interface {
	core.Writer
	String() string
}

type frameTestModel string

func (m frameTestModel) View(_, _ int) string { return string(m) }

type frameInteractiveModel struct {
	view    string
	updates int
}

func (m *frameInteractiveModel) View(_, _ int) string { return m.view }
func (m *frameInteractiveModel) Init() tea.Cmd        { return nil }
func (m *frameInteractiveModel) Update(_ tea.Msg) (FrameModel, tea.Cmd) {
	m.updates++
	return m, nil
}

func frameWithBuffer() (*Frame, frameStringWriter) {
	buf := core.NewBufferString("")
	f := NewFrame("HLCF").WithOutput(buf)
	return f, buf
}

func frameRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	out := core.NewBuilder()
	for range count {
		out.WriteString(s)
	}
	return out.String()
}

func TestFrame_ModelFunc_View_Good(t *core.T) {
	m := ModelFunc(func(width, height int) string { return core.Sprintf("%dx%d", width, height) })

	core.AssertEqual(t, "4x2", m.View(4, 2))
	core.AssertContains(t, m.View(3, 5), "3x5")
}

func TestFrame_ModelFunc_View_Bad(t *core.T) {
	m := ModelFunc(func(_, _ int) string { return "" })

	core.AssertEqual(t, "", m.View(0, 0))
	core.AssertEmpty(t, m.View(-1, -1))
}

func TestFrame_ModelFunc_View_Ugly(t *core.T) {
	m := ModelFunc(func(width, _ int) string { return frameRepeat("x", width) })

	core.AssertEqual(t, "xxx", m.View(3, 1))
	core.AssertNotPanics(t, func() { _ = m.View(0, 0) })
}

func TestFrame_NewFrame_Good(t *core.T) {
	f := NewFrame("HCF")

	core.AssertEqual(t, "HCF", f.variant)
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestFrame_NewFrame_Bad(t *core.T) {
	f := NewFrame("Z")

	core.AssertEqual(t, "Z", f.variant)
	core.AssertEmpty(t, f.String())
}

func TestFrame_NewFrame_Ugly(t *core.T) {
	f := NewFrame("")

	core.AssertEqual(t, "", f.variant)
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestFrame_Frame_WithOutput_Good(t *core.T) {
	f, buf := frameWithBuffer()

	core.AssertEqual(t, f, f.WithOutput(buf))
	core.AssertEqual(t, buf, f.out)
}

func TestFrame_Frame_WithOutput_Bad(t *core.T) {
	f := NewFrame("C")
	original := f.out

	core.AssertEqual(t, f, f.WithOutput(nil))
	core.AssertEqual(t, original, f.out)
}

func TestFrame_Frame_WithOutput_Ugly(t *core.T) {
	f := NewFrame("C")
	var buf bytes.Buffer

	f.WithOutput(&buf).Content(StaticModel("content")).Run()
	core.AssertContains(t, buf.String(), "content")
}

func TestFrame_Frame_Header_Good(t *core.T) {
	f := NewFrame("HC").Header(StaticModel("header"))

	core.AssertContains(t, f.String(), "header")
	core.AssertNotNil(t, f.models[RegionHeader])
}

func TestFrame_Frame_Header_Bad(t *core.T) {
	f := NewFrame("C").Header(StaticModel("ignored"))

	core.AssertNotContains(t, f.String(), "ignored")
	core.AssertNotNil(t, f.models[RegionHeader])
}

func TestFrame_Frame_Header_Ugly(t *core.T) {
	f := NewFrame("HC").Header(StaticModel(":check: header"))

	core.AssertContains(t, f.String(), "header")
	core.AssertContains(t, f.String(), "✓")
}

func TestFrame_Frame_Left_Good(t *core.T) {
	f := NewFrame("LC").Left(StaticModel("left")).Content(StaticModel("content"))

	core.AssertContains(t, f.String(), "left")
	core.AssertContains(t, f.String(), "content")
}

func TestFrame_Frame_Left_Bad(t *core.T) {
	f := NewFrame("C").Left(StaticModel("left"))

	core.AssertNotContains(t, f.String(), "left")
	core.AssertNotNil(t, f.models[RegionLeft])
}

func TestFrame_Frame_Left_Ugly(t *core.T) {
	f := NewFrame("LC").Left(StaticModel(""))

	core.AssertNotPanics(t, func() { _ = f.String() })
	core.AssertEqual(t, "", f.models[RegionLeft].View(1, 1))
}

func TestFrame_Frame_Content_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertContains(t, f.String(), "content")
	core.AssertNotNil(t, f.models[RegionContent])
}

func TestFrame_Frame_Content_Bad(t *core.T) {
	f := NewFrame("H").Content(StaticModel("content"))

	core.AssertContains(t, f.String(), "content")
	core.AssertNotNil(t, f.models[RegionContent])
}

func TestFrame_Frame_Content_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel(":warn:"))

	core.AssertContains(t, f.String(), "⚠")
	core.AssertNotEmpty(t, f.String())
}

func TestFrame_Frame_Right_Good(t *core.T) {
	f := NewFrame("CR").Content(StaticModel("content")).Right(StaticModel("right"))

	core.AssertContains(t, f.String(), "right")
	core.AssertContains(t, f.String(), "content")
}

func TestFrame_Frame_Right_Bad(t *core.T) {
	f := NewFrame("C").Right(StaticModel("right"))

	core.AssertNotContains(t, f.String(), "right")
	core.AssertNotNil(t, f.models[RegionRight])
}

func TestFrame_Frame_Right_Ugly(t *core.T) {
	f := NewFrame("CR").Right(StaticModel(""))

	core.AssertNotPanics(t, func() { _ = f.String() })
	core.AssertEqual(t, "", f.models[RegionRight].View(1, 1))
}

func TestFrame_Frame_Footer_Good(t *core.T) {
	f := NewFrame("CF").Content(StaticModel("content")).Footer(StaticModel("footer"))

	core.AssertContains(t, f.String(), "footer")
	core.AssertContains(t, f.String(), "content")
}

func TestFrame_Frame_Footer_Bad(t *core.T) {
	f := NewFrame("C").Footer(StaticModel("footer"))

	core.AssertNotContains(t, f.String(), "footer")
	core.AssertNotNil(t, f.models[RegionFooter])
}

func TestFrame_Frame_Footer_Ugly(t *core.T) {
	f := NewFrame("CF").Footer(StaticModel(":info: footer"))

	core.AssertContains(t, f.String(), "footer")
	core.AssertContains(t, f.String(), "ℹ")
}

func TestFrame_Frame_Navigate_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel("second"))

	core.AssertContains(t, f.String(), "second")
	core.AssertLen(t, f.history, 1)
}

func TestFrame_Frame_Navigate_Bad(t *core.T) {
	f := NewFrame("C")
	f.Navigate(StaticModel("second"))

	core.AssertContains(t, f.String(), "second")
	core.AssertEmpty(t, f.history)
}

func TestFrame_Frame_Navigate_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel(""))

	core.AssertNotContains(t, f.String(), "first")
	core.AssertLen(t, f.history, 1)
}

func TestFrame_Frame_Back_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel("second"))

	core.AssertTrue(t, f.Back())
	core.AssertContains(t, f.String(), "first")
}

func TestFrame_Frame_Back_Bad(t *core.T) {
	f := NewFrame("C")

	core.AssertFalse(t, f.Back())
	core.AssertEmpty(t, f.history)
}

func TestFrame_Frame_Back_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("first"))
	f.Navigate(StaticModel("second"))
	f.Back()

	core.AssertFalse(t, f.Back())
	core.AssertContains(t, f.String(), "first")
}

func TestFrame_Frame_Stop_Good(t *core.T) {
	f := NewFrame("C")
	f.Stop()

	_, closed := <-f.done
	core.AssertFalse(t, closed)
}

func TestFrame_Frame_Stop_Bad(t *core.T) {
	f := NewFrame("C")
	f.Stop()

	core.AssertNotPanics(t, func() { f.Stop() })
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestFrame_Frame_Stop_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))
	f.Stop()

	core.AssertNotPanics(t, func() { f.RunFor(time.Nanosecond) })
	core.AssertContains(t, f.String(), "content")
}

func TestFrame_Frame_Send_Good(t *core.T) {
	f := NewFrame("C")

	core.AssertNotPanics(t, func() { f.Send(struct{}{}) })
	core.AssertNil(t, f.program)
}

func TestFrame_Frame_Send_Bad(t *core.T) {
	var f *Frame

	core.AssertPanics(t, func() { f.Send(struct{}{}) })
	core.AssertNil(t, f)
}

func TestFrame_Frame_Send_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertNotPanics(t, func() { f.Send(nil) })
	core.AssertContains(t, f.String(), "content")
}

func TestFrame_Frame_WithKeyMap_Good(t *core.T) {
	f := NewFrame("C")
	km := DefaultKeyMap()

	core.AssertEqual(t, f, f.WithKeyMap(km))
	core.AssertEqual(t, km, f.keyMap)
}

func TestFrame_Frame_WithKeyMap_Bad(t *core.T) {
	f := NewFrame("C")
	f.WithKeyMap(KeyMap{})

	core.AssertEqual(t, tea.KeyType(0), f.keyMap.Quit)
	core.AssertEqual(t, RegionContent, f.Focused())
}

func TestFrame_Frame_WithKeyMap_Ugly(t *core.T) {
	f := NewFrame("HC")
	f.WithKeyMap(DefaultKeyMap())

	core.AssertNotPanics(t, func() { _, _ = f.Update(tea.KeyMsg{Type: f.keyMap.FocusNext}) })
	core.AssertEqual(t, RegionHeader, f.Focused())
}

func TestFrame_Frame_Focused_Good(t *core.T) {
	f := NewFrame("HC")
	f.Focus(RegionHeader)

	core.AssertEqual(t, RegionHeader, f.Focused())
	core.AssertNotEqual(t, RegionContent, f.Focused())
}

func TestFrame_Frame_Focused_Bad(t *core.T) {
	f := NewFrame("C")
	f.Focus(RegionHeader)

	core.AssertEqual(t, RegionContent, f.Focused())
	core.AssertNotEqual(t, RegionHeader, f.Focused())
}

func TestFrame_Frame_Focused_Ugly(t *core.T) {
	f := NewFrame("")
	f.Focus(RegionContent)

	core.AssertEqual(t, RegionContent, f.Focused())
	core.AssertEmpty(t, f.layout.regions)
}

func TestFrame_Frame_Focus_Good(t *core.T) {
	f := NewFrame("HCF")
	f.Focus(RegionFooter)

	core.AssertEqual(t, RegionFooter, f.Focused())
	core.AssertNotEqual(t, RegionContent, f.Focused())
}

func TestFrame_Frame_Focus_Bad(t *core.T) {
	f := NewFrame("C")
	f.Focus(RegionLeft)

	core.AssertEqual(t, RegionContent, f.Focused())
	core.AssertFalse(t, layoutHasRegion(f.layout, RegionLeft))
}

func TestFrame_Frame_Focus_Ugly(t *core.T) {
	f := NewFrame("HLCF")
	f.Focus(RegionLeft)

	core.AssertEqual(t, RegionLeft, f.Focused())
	core.AssertTrue(t, layoutHasRegion(f.layout, RegionLeft))
}

func TestFrame_Frame_Init_Good(t *core.T) {
	f := NewFrame("C").Content(&frameInteractiveModel{view: "content"})
	cmd := f.Init()

	core.AssertNil(t, cmd)
	core.AssertNotPanics(t, func() { _ = f.Init() })
}

func TestFrame_Frame_Init_Bad(t *core.T) {
	f := NewFrame("C")
	cmd := f.Init()

	core.AssertNil(t, cmd)
	core.AssertNotPanics(t, func() { _ = f.Init() })
}

func TestFrame_Frame_Init_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("plain"))
	cmd := f.Init()

	core.AssertNil(t, cmd)
	core.AssertNotPanics(t, func() { _ = f.Init() })
}

func TestFrame_Frame_Update_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))
	model, cmd := f.Update(tea.WindowSizeMsg{Width: 40, Height: 10})

	core.AssertEqual(t, f, model)
	core.AssertNil(t, cmd)
}

func TestFrame_Frame_Update_Bad(t *core.T) {
	f := NewFrame("C")
	model, cmd := f.Update(tea.KeyMsg{Type: DefaultKeyMap().Back})

	core.AssertEqual(t, f, model)
	core.AssertNil(t, cmd)
}

func TestFrame_Frame_Update_Ugly(t *core.T) {
	m := &frameInteractiveModel{view: "content"}
	f := NewFrame("C").Content(m)
	_, _ = f.Update(struct{}{})

	core.AssertEqual(t, 1, m.updates)
	core.AssertContains(t, f.String(), "content")
}

func TestFrame_Frame_View_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertContains(t, f.View(), "content")
	core.AssertNotEmpty(t, f.View())
}

func TestFrame_Frame_View_Bad(t *core.T) {
	f := NewFrame("C")

	core.AssertEqual(t, "", f.View())
	core.AssertEmpty(t, f.View())
}

func TestFrame_Frame_View_Ugly(t *core.T) {
	f := NewFrame("HCF").Header(StaticModel("h")).Content(StaticModel("c")).Footer(StaticModel("f"))

	core.AssertContains(t, f.View(), "h")
	core.AssertContains(t, f.View(), "f")
}

func TestFrame_Frame_Run_Good(t *core.T) {
	f, buf := frameWithBuffer()
	f.Content(StaticModel("content"))

	f.Run()
	core.AssertContains(t, buf.String(), "content")
}

func TestFrame_Frame_Run_Bad(t *core.T) {
	f, buf := frameWithBuffer()

	f.Run()
	core.AssertEqual(t, "", buf.String())
}

func TestFrame_Frame_Run_Ugly(t *core.T) {
	f, buf := frameWithBuffer()
	f.Content(StaticModel(":check:"))

	f.Run()
	core.AssertContains(t, buf.String(), "✓")
}

func TestFrame_Frame_RunFor_Good(t *core.T) {
	f, buf := frameWithBuffer()
	f.Content(StaticModel("content"))

	f.RunFor(time.Nanosecond)
	core.AssertContains(t, buf.String(), "content")
}

func TestFrame_Frame_RunFor_Bad(t *core.T) {
	f, buf := frameWithBuffer()

	f.RunFor(0)
	core.AssertEqual(t, "", buf.String())
}

func TestFrame_Frame_RunFor_Ugly(t *core.T) {
	f, buf := frameWithBuffer()
	f.Content(StaticModel(""))

	f.RunFor(time.Nanosecond)
	core.AssertEqual(t, "", buf.String())
}

func TestFrame_Frame_String_Good(t *core.T) {
	f := NewFrame("C").Content(StaticModel("content"))

	core.AssertEqual(t, "content\n", f.String())
	core.AssertContains(t, f.String(), "\n")
}

func TestFrame_Frame_String_Bad(t *core.T) {
	f := NewFrame("C")

	core.AssertEqual(t, "", f.String())
	core.AssertEmpty(t, f.String())
}

func TestFrame_Frame_String_Ugly(t *core.T) {
	f := NewFrame("C").Content(StaticModel("\033[31mred\033[0m"))

	core.AssertEqual(t, "red\n", f.String())
	core.AssertNotContains(t, f.String(), "\033")
}
