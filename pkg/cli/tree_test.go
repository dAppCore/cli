package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree_Good(t *testing.T) {
	t.Run("single root", func(t *testing.T) {
		tree := NewTree("root")
		assert.Equal(t, "root\n", tree.String())
	})

	t.Run("flat children", func(t *testing.T) {
		tree := NewTree("root")
		tree.Add("alpha")
		tree.Add("beta")
		tree.Add("gamma")

		expected := "root\n" +
			"├── alpha\n" +
			"├── beta\n" +
			"└── gamma\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("nested children", func(t *testing.T) {
		tree := NewTree("core-php")
		tree.Add("core-tenant").Add("core-bio")
		tree.Add("core-admin")
		tree.Add("core-api")

		expected := "core-php\n" +
			"├── core-tenant\n" +
			"│   └── core-bio\n" +
			"├── core-admin\n" +
			"└── core-api\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("deep nesting", func(t *testing.T) {
		tree := NewTree("a")
		tree.Add("b").Add("c").Add("d")

		expected := "a\n" +
			"└── b\n" +
			"    └── c\n" +
			"        └── d\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("mixed depth", func(t *testing.T) {
		tree := NewTree("root")
		a := tree.Add("a")
		a.Add("a1")
		a.Add("a2")
		tree.Add("b")

		expected := "root\n" +
			"├── a\n" +
			"│   ├── a1\n" +
			"│   └── a2\n" +
			"└── b\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("AddTree composes subtrees", func(t *testing.T) {
		sub := NewTree("sub-root")
		sub.Add("child")

		tree := NewTree("main")
		tree.AddTree(sub)

		expected := "main\n" +
			"└── sub-root\n" +
			"    └── child\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("styled nodes", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tree := NewTree("root")
		tree.AddStyled("green", SuccessStyle)
		tree.Add("plain")

		expected := "root\n" +
			"├── green\n" +
			"└── plain\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("WithStyle on root", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tree := NewTree("root").WithStyle(ErrorStyle)
		tree.Add("child")

		expected := "root\n" +
			"└── child\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("ASCII theme uses ASCII connectors", func(t *testing.T) {
		prevTheme := currentTheme
		prevColor := ColorEnabled()
		UseASCII()
		t.Cleanup(func() {
			currentTheme = prevTheme
			SetColorEnabled(prevColor)
		})

		tree := NewTree("core-php")
		tree.Add("core-tenant").Add("core-bio")
		tree.Add("core-admin")
		tree.Add("core-api")

		expected := "core-php\n" +
			"+-- core-tenant\n" +
			"|   `-- core-bio\n" +
			"+-- core-admin\n" +
			"`-- core-api\n"
		assert.Equal(t, expected, tree.String())
	})

	t.Run("glyph shortcodes render in labels", func(t *testing.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tree := NewTree(":check: root")
		tree.Add(":warn: child")

		out := tree.String()
		assert.Contains(t, out, "[OK] root")
		assert.Contains(t, out, "[WARN] child")
	})
}

func TestTree_Bad(t *testing.T) {
	t.Run("empty label", func(t *testing.T) {
		tree := NewTree("")
		assert.Equal(t, "\n", tree.String())
	})
}

func TestTree_Ugly(t *testing.T) {
	t.Run("nil style does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			tree := NewTree("root").WithStyle(nil)
			tree.Add("child")
			_ = tree.String()
		})
	})

	t.Run("AddStyled with nil style does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			tree := NewTree("root")
			tree.AddStyled("item", nil)
			_ = tree.String()
		})
	})

	t.Run("very deep nesting does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			node := NewTree("root")
			for range 100 {
				node = node.Add("child")
			}
			_ = NewTree("root").String()
		})
	})
}
