package cli

import (
	"fmt"
	"iter"
	"strings"
)

// TreeNode represents a node in a displayable tree structure.
// Use NewTree to create a root, then Add children.
//
//	tree := cli.NewTree("core-php")
//	tree.Add("core-tenant").Add("core-bio")
//	tree.Add("core-admin")
//	tree.Add("core-api")
//	fmt.Print(tree)
//	// core-php
//	// ├── core-tenant
//	// │   └── core-bio
//	// ├── core-admin
//	// └── core-api
type TreeNode struct {
	label    string
	style    *AnsiStyle
	children []*TreeNode
}

// NewTree creates a new tree with the given root label.
func NewTree(label string) *TreeNode {
	return &TreeNode{label: label}
}

// Add appends a child node and returns the child for chaining.
func (n *TreeNode) Add(label string) *TreeNode {
	child := &TreeNode{label: label}
	n.children = append(n.children, child)
	return child
}

// AddStyled appends a styled child node and returns the child for chaining.
func (n *TreeNode) AddStyled(label string, style *AnsiStyle) *TreeNode {
	child := &TreeNode{label: label, style: style}
	n.children = append(n.children, child)
	return child
}

// AddTree appends an existing tree as a child and returns the parent for chaining.
func (n *TreeNode) AddTree(child *TreeNode) *TreeNode {
	n.children = append(n.children, child)
	return n
}

// WithStyle sets the style on this node and returns it for chaining.
func (n *TreeNode) WithStyle(style *AnsiStyle) *TreeNode {
	n.style = style
	return n
}

// Children returns an iterator over the node's children.
func (n *TreeNode) Children() iter.Seq[*TreeNode] {
	return func(yield func(*TreeNode) bool) {
		for _, child := range n.children {
			if !yield(child) {
				return
			}
		}
	}
}

// String renders the tree with box-drawing characters.
// Implements fmt.Stringer.
func (n *TreeNode) String() string {
	var sb strings.Builder
	sb.WriteString(n.renderLabel())
	sb.WriteByte('\n')
	n.writeChildren(&sb, "")
	return sb.String()
}

// Render prints the tree to stdout.
func (n *TreeNode) Render() {
	fmt.Print(n.String())
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
