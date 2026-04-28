package cli

import "dappco.re/go"

func TestTree_Good(t *core.T) {
	t.Run("single root", func(t *core.T) {
		tree := NewTree("root")
		core.AssertEqual(t, "root\n", tree.String())
	})

	t.Run("flat children", func(t *core.T) {
		tree := NewTree("root")
		tree.Add("alpha")
		tree.Add("beta")
		tree.Add("gamma")

		expected := "root\n" +
			"├── alpha\n" +
			"├── beta\n" +
			"└── gamma\n"
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("nested children", func(t *core.T) {
		tree := NewTree("core-php")
		tree.Add("core-tenant").Add("core-bio")
		tree.Add("core-admin")
		tree.Add("core-api")

		expected := "core-php\n" +
			"├── core-tenant\n" +
			"│   └── core-bio\n" +
			"├── core-admin\n" +
			"└── core-api\n"
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("deep nesting", func(t *core.T) {
		tree := NewTree("a")
		tree.Add("b").Add("c").Add("d")

		expected := "a\n" +
			"└── b\n" +
			"    └── c\n" +
			"        └── d\n"
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("mixed depth", func(t *core.T) {
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
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("AddTree composes subtrees", func(t *core.T) {
		sub := NewTree("sub-root")
		sub.Add("child")

		tree := NewTree("main")
		tree.AddTree(sub)

		expected := "main\n" +
			"└── sub-root\n" +
			"    └── child\n"
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("styled nodes", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tree := NewTree("root")
		tree.AddStyled("green", SuccessStyle)
		tree.Add("plain")

		expected := "root\n" +
			"├── green\n" +
			"└── plain\n"
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("WithStyle on root", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		tree := NewTree("root").WithStyle(ErrorStyle)
		tree.Add("child")

		expected := "root\n" +
			"└── child\n"
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("ASCII theme uses ASCII connectors", func(t *core.T) {
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
		core.AssertEqual(t, expected, tree.String())
	})

	t.Run("glyph shortcodes render in labels", func(t *core.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tree := NewTree(":check: root")
		tree.Add(":warn: child")

		out := tree.String()
		core.AssertContains(t, out, "[OK] root")
		core.AssertContains(t, out, "[WARN] child")
	})
}

func TestTree_Bad(t *core.T) {
	t.Run("empty label", func(t *core.T) {
		tree := NewTree("")
		core.AssertEqual(t, "\n", tree.String())
	})
}

func TestTree_Ugly(t *core.T) {
	t.Run("nil style does not panic", func(t *core.T) {
		core.AssertNotPanics(t, func() {
			tree := NewTree("root").WithStyle(nil)
			tree.Add("child")
			_ = tree.String()
		})
	})

	t.Run("AddStyled with nil style does not panic", func(t *core.T) {
		core.AssertNotPanics(t, func() {
			tree := NewTree("root")
			tree.AddStyled("item", nil)
			_ = tree.String()
		})
	})

	t.Run("very deep nesting does not panic", func(t *core.T) {
		core.AssertNotPanics(t, func() {
			node := NewTree("root")
			for range 100 {
				node = node.Add("child")
			}
			_ = NewTree("root").String()
		})
	})
}
