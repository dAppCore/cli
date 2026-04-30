package frame

import (
	core "dappco.re/go"
	tea "github.com/charmbracelet/bubbletea"
)

func TestFrameModel_Adapter_View_Good(t *core.T) {
	a := &modelAdapter{m: frameTestModel("content")}

	core.AssertEqual(t, "content", a.View(80, 1))
	core.AssertContains(t, a.View(80, 1), "content")
}

func TestFrameModel_Adapter_View_Bad(t *core.T) {
	a := &modelAdapter{}

	core.AssertPanics(t, func() { _ = a.View(80, 1) })
	core.AssertNil(t, a.m)
}

func TestFrameModel_Adapter_View_Ugly(t *core.T) {
	a := &modelAdapter{m: frameTestModel("")}

	core.AssertEqual(t, "", a.View(0, 0))
	core.AssertEmpty(t, a.View(0, 0))
}

func TestFrameModel_Adapter_Init_Good(t *core.T) {
	a := &modelAdapter{m: frameTestModel("content")}

	core.AssertNil(t, a.Init())
	core.AssertNotNil(t, a.m)
}

func TestFrameModel_Adapter_Init_Bad(t *core.T) {
	a := &modelAdapter{}

	core.AssertNil(t, a.Init())
	core.AssertNil(t, a.m)
}

func TestFrameModel_Adapter_Init_Ugly(t *core.T) {
	var a *modelAdapter

	core.AssertNil(t, a.Init())
	core.AssertNil(t, a)
}

func TestFrameModel_Adapter_Update_Good(t *core.T) {
	a := &modelAdapter{m: frameTestModel("content")}
	updated, cmd := a.Update(struct{}{})

	core.AssertEqual(t, a, updated)
	core.AssertNil(t, cmd)
}

func TestFrameModel_Adapter_Update_Bad(t *core.T) {
	a := &modelAdapter{}
	updated, cmd := a.Update(nil)

	core.AssertEqual(t, a, updated)
	core.AssertNil(t, cmd)
}

func TestFrameModel_Adapter_Update_Ugly(t *core.T) {
	var a *modelAdapter

	updated, cmd := a.Update(nil)
	core.AssertNil(t, updated)
	core.AssertNil(t, cmd)
}

func TestFrameModel_DefaultKeyMap_Good(t *core.T) {
	km := DefaultKeyMap()

	core.AssertEqual(t, tea.KeyTab, km.FocusNext)
	core.AssertEqual(t, tea.KeyCtrlC, km.Quit)
}

func TestFrameModel_DefaultKeyMap_Bad(t *core.T) {
	km := KeyMap{}

	core.AssertNotEqual(t, DefaultKeyMap(), km)
	core.AssertEqual(t, tea.KeyType(0), km.Quit)
}

func TestFrameModel_DefaultKeyMap_Ugly(t *core.T) {
	km := DefaultKeyMap()

	core.AssertNotEqual(t, km.FocusNext, km.FocusPrev)
	core.AssertNotEqual(t, km.FocusLeft, km.FocusRight)
}
