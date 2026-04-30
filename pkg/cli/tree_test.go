package cli

import (
	core "dappco.re/go"
)

func TestTree_NewTree_Good(t *core.T) {
	tree := NewTree("root")

	core.AssertNotNil(t, tree)
	core.AssertContains(t, tree.String(), "root")
}

func TestTree_NewTree_Bad(t *core.T) {
	tree := NewTree("")

	core.AssertNotNil(t, tree)
	core.AssertContains(t, tree.String(), "\n")
}

func TestTree_NewTree_Ugly(t *core.T) {
	tree := NewTree(":check:")

	core.AssertContains(t, tree.String(), "✓")
	core.AssertNotContains(t, tree.String(), ":check:")
}

func TestTree_TreeNode_Add_Good(t *core.T) {
	tree := NewTree("root")
	child := tree.Add("child")

	core.AssertNotNil(t, child)
	core.AssertContains(t, tree.String(), "child")
}

func TestTree_TreeNode_Add_Bad(t *core.T) {
	tree := NewTree("root")
	child := tree.Add("")

	core.AssertNotNil(t, child)
	core.AssertLen(t, tree.children, 1)
}

func TestTree_TreeNode_Add_Ugly(t *core.T) {
	tree := NewTree("root")
	tree.Add(":check:")

	core.AssertContains(t, tree.String(), "✓")
	core.AssertLen(t, tree.children, 1)
}

func TestTree_TreeNode_AddStyled_Good(t *core.T) {
	tree := NewTree("root")
	child := tree.AddStyled("child", NewStyle().Bold())

	core.AssertNotNil(t, child.style)
	core.AssertContains(t, tree.String(), "child")
}

func TestTree_TreeNode_AddStyled_Bad(t *core.T) {
	tree := NewTree("root")
	child := tree.AddStyled("child", nil)

	core.AssertNil(t, child.style)
	core.AssertContains(t, tree.String(), "child")
}

func TestTree_TreeNode_AddStyled_Ugly(t *core.T) {
	tree := NewTree("root")
	child := tree.AddStyled("", NewStyle().Dim())

	core.AssertNotNil(t, child.style)
	core.AssertLen(t, tree.children, 1)
}

func TestTree_TreeNode_AddTree_Good(t *core.T) {
	tree := NewTree("root")
	child := NewTree("child")

	core.AssertEqual(t, tree, tree.AddTree(child))
	core.AssertContains(t, tree.String(), "child")
}

func TestTree_TreeNode_AddTree_Bad(t *core.T) {
	tree := NewTree("root")
	tree.AddTree(nil)

	core.AssertLen(t, tree.children, 1)
	core.AssertPanics(t, func() { _ = tree.String() })
}

func TestTree_TreeNode_AddTree_Ugly(t *core.T) {
	tree := NewTree("root")
	child := NewTree("child").Add("leaf")
	tree.AddTree(child)

	core.AssertContains(t, tree.String(), "leaf")
	core.AssertLen(t, tree.children, 1)
}

func TestTree_TreeNode_WithStyle_Good(t *core.T) {
	tree := NewTree("root").WithStyle(NewStyle().Bold())

	core.AssertNotNil(t, tree.style)
	core.AssertContains(t, tree.String(), "root")
}

func TestTree_TreeNode_WithStyle_Bad(t *core.T) {
	tree := NewTree("root").WithStyle(nil)

	core.AssertNil(t, tree.style)
	core.AssertContains(t, tree.String(), "root")
}

func TestTree_TreeNode_WithStyle_Ugly(t *core.T) {
	tree := NewTree("").WithStyle(NewStyle().Dim())

	core.AssertNotNil(t, tree.style)
	core.AssertContains(t, tree.String(), "\n")
}

func TestTree_TreeNode_Children_Good(t *core.T) {
	tree := NewTree("root")
	tree.Add("child")
	var count int
	for range tree.Children() {
		count++
	}

	core.AssertEqual(t, 1, count)
}

func TestTree_TreeNode_Children_Bad(t *core.T) {
	tree := NewTree("root")
	var count int
	for range tree.Children() {
		count++
	}

	core.AssertEqual(t, 0, count)
}

func TestTree_TreeNode_Children_Ugly(t *core.T) {
	tree := NewTree("root")
	tree.Add("a")
	tree.Add("b")
	var names []string
	for child := range tree.Children() {
		names = append(names, child.label)
	}

	core.AssertEqual(t, []string{"a", "b"}, names)
}

func TestTree_TreeNode_String_Good(t *core.T) {
	tree := NewTree("root")
	tree.Add("child")
	got := tree.String()

	core.AssertContains(t, got, "root")
	core.AssertContains(t, got, "child")
}

func TestTree_TreeNode_String_Bad(t *core.T) {
	got := NewTree("").String()

	core.AssertContains(t, got, "\n")
	core.AssertEqual(t, 1, core.RuneCount(got))
}

func TestTree_TreeNode_String_Ugly(t *core.T) {
	tree := NewTree(":check:")
	tree.Add(":warn:")

	core.AssertContains(t, tree.String(), "✓")
	core.AssertContains(t, tree.String(), "⚠")
}

func TestTree_TreeNode_Render_Good(t *core.T) {
	tree := NewTree("root")
	tree.Add("child")
	out := cliCaptureStdout(t, func() { tree.Render() })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "child")
}

func TestTree_TreeNode_Render_Bad(t *core.T) {
	out := cliCaptureStdout(t, func() { NewTree("").Render() })

	core.AssertContains(t, out, "\n")
	core.AssertEqual(t, 1, core.RuneCount(out))
}

func TestTree_TreeNode_Render_Ugly(t *core.T) {
	tree := NewTree(":check:")
	out := cliCaptureStdout(t, func() { tree.Render() })

	core.AssertContains(t, out, "✓")
	core.AssertNotContains(t, out, ":check:")
}
