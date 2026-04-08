package cli

import (
	"fmt"
	"iter"
	"strings"
)

// TreeNode represents a node in a displayable tree structure.
type TreeNode struct {
	label    string
	style    *AnsiStyle
	children []*TreeNode
}

func NewTree(label string) *TreeNode {
	return &TreeNode{label: label}
}

func (n *TreeNode) Add(label string) *TreeNode {
	child := &TreeNode{label: label}
	n.children = append(n.children, child)
	return child
}

func (n *TreeNode) AddStyled(label string, style *AnsiStyle) *TreeNode {
	child := &TreeNode{label: label, style: style}
	n.children = append(n.children, child)
	return child
}

func (n *TreeNode) AddTree(child *TreeNode) *TreeNode {
	n.children = append(n.children, child)
	return n
}

func (n *TreeNode) WithStyle(style *AnsiStyle) *TreeNode {
	n.style = style
	return n
}

func (n *TreeNode) Children() iter.Seq[*TreeNode] {
	return func(yield func(*TreeNode) bool) {
		for _, child := range n.children {
			if !yield(child) {
				return
			}
		}
	}
}

func (n *TreeNode) String() string {
	var sb strings.Builder
	sb.WriteString(n.renderLabel())
	sb.WriteByte('\n')
	n.writeChildren(&sb, "")
	return sb.String()
}

func (n *TreeNode) Render() {
	fmt.Fprint(stdoutWriter(), n.String())
}

func (n *TreeNode) renderLabel() string {
	label := compileGlyphs(n.label)
	if n.style != nil {
		return n.style.Render(label)
	}
	return label
}

func (n *TreeNode) writeChildren(sb *strings.Builder, prefix string) {
	tee := Glyph(":tee:") + Glyph(":dash:") + Glyph(":dash:") + " "
	corner := Glyph(":corner:") + Glyph(":dash:") + Glyph(":dash:") + " "
	pipe := Glyph(":pipe:") + "   "

	for i, child := range n.children {
		last := i == len(n.children)-1
		connector := tee
		next := pipe
		if last {
			connector = corner
			next = "    "
		}
		sb.WriteString(prefix)
		sb.WriteString(connector)
		sb.WriteString(child.renderLabel())
		sb.WriteByte('\n')
		child.writeChildren(sb, prefix+next)
	}
}
