package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrame_Good(t *testing.T) {
	t.Run("static render HCF", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("header"))
		f.Content(StaticModel("content"))
		f.Footer(StaticModel("footer"))

		out := f.String()
		assert.Contains(t, out, "header")
		assert.Contains(t, out, "content")
		assert.Contains(t, out, "footer")
	})

	t.Run("region order preserved", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("AAA"))
		f.Content(StaticModel("BBB"))
		f.Footer(StaticModel("CCC"))

		out := f.String()
		posA := indexOf(out, "AAA")
		posB := indexOf(out, "BBB")
		posC := indexOf(out, "CCC")
		assert.Less(t, posA, posB, "header before content")
		assert.Less(t, posB, posC, "content before footer")
	})

	t.Run("navigate and back", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("nav"))
		f.Content(StaticModel("page-1"))
		f.Footer(StaticModel("hints"))

		assert.Contains(t, f.String(), "page-1")

		// Navigate to page 2
		f.Navigate(StaticModel("page-2"))
		assert.Contains(t, f.String(), "page-2")
		assert.NotContains(t, f.String(), "page-1")

		// Navigate to page 3
		f.Navigate(StaticModel("page-3"))
		assert.Contains(t, f.String(), "page-3")

		// Back to page 2
		ok := f.Back()
		require.True(t, ok)
		assert.Contains(t, f.String(), "page-2")

		// Back to page 1
		ok = f.Back()
		require.True(t, ok)
		assert.Contains(t, f.String(), "page-1")

		// No more history
		ok = f.Back()
		assert.False(t, ok)
	})

	t.Run("empty regions skipped", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("only content"))

		out := f.String()
		assert.Equal(t, "only content\n", out)
	})

	t.Run("non-TTY run renders once", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		var buf bytes.Buffer
		f := NewFrame("HCF")
		f.out = &buf
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))

		f.Run() // non-TTY, should return immediately
		assert.Contains(t, buf.String(), "h")
		assert.Contains(t, buf.String(), "c")
		assert.Contains(t, buf.String(), "f")
	})

	t.Run("ModelFunc adapter", func(t *testing.T) {
		called := false
		m := ModelFunc(func(w, h int) string {
			called = true
			return "dynamic"
		})

		out := m.View(80, 24)
		assert.True(t, called)
		assert.Equal(t, "dynamic", out)
	})

	t.Run("RunFor exits after duration", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewFrame("C")
		f.out = &buf // non-TTY → RunFor renders once and returns
		f.Content(StaticModel("timed"))

		start := time.Now()
		f.RunFor(50 * time.Millisecond)
		elapsed := time.Since(start)

		assert.Less(t, elapsed, 200*time.Millisecond)
		assert.Contains(t, buf.String(), "timed")
	})
}

func TestFrame_Bad(t *testing.T) {
	t.Run("empty frame", func(t *testing.T) {
		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		assert.Equal(t, "", f.String())
	})

	t.Run("back on empty history", func(t *testing.T) {
		f := NewFrame("C")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("x"))
		assert.False(t, f.Back())
	})

	t.Run("invalid variant degrades gracefully", func(t *testing.T) {
		f := NewFrame("XYZ")
		f.out = &bytes.Buffer{}
		// No valid regions, so nothing renders
		assert.Equal(t, "", f.String())
	})
}

func TestStatusLine_Good(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := StatusLine("core dev", "18 repos", "main")
	out := m.View(80, 1)
	assert.Contains(t, out, "core dev")
	assert.Contains(t, out, "18 repos")
	assert.Contains(t, out, "main")
}

func TestKeyHints_Good(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := KeyHints("↑/↓ navigate", "q quit")
	out := m.View(80, 1)
	assert.Contains(t, out, "navigate")
	assert.Contains(t, out, "quit")
}

func TestBreadcrumb_Good(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := Breadcrumb("core", "dev", "health")
	out := m.View(80, 1)
	assert.Contains(t, out, "core")
	assert.Contains(t, out, "dev")
	assert.Contains(t, out, "health")
	assert.Contains(t, out, ">")
}

func TestStaticModel_Good(t *testing.T) {
	m := StaticModel("hello")
	assert.Equal(t, "hello", m.View(80, 24))
}

// indexOf returns the position of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := range len(s) - len(substr) + 1 {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
