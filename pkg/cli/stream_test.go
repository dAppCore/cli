package cli

import (
	"bytes"
	"dappco.re/go"
	"strings"
)

func TestStream_Good(t *core.T) {
	t.Run("basic write", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		s.Write("hello ")
		s.Write("world")
		s.Done()
		s.Wait()
		core.AssertEqual(t, "hello world\n", buf.String())
	})

	t.Run("write with newlines", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		s.Write("line1\nline2\n")
		s.Done()
		s.Wait()
		core.AssertEqual(t, "line1\nline2\n", buf.String())
	})

	t.Run("word wrap", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithWordWrap(10), WithStreamOutput(&buf))

		s.Write("1234567890ABCDE")
		s.Done()
		s.Wait()

		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		core.AssertEqual(t, 2, len(lines))
		core.AssertEqual(t, "1234567890", lines[0])
		core.AssertEqual(t, "ABCDE", lines[1])
	})

	t.Run("word wrap preserves explicit newlines", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithWordWrap(20), WithStreamOutput(&buf))

		s.Write("short\nanother")
		s.Done()
		s.Wait()

		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		core.AssertEqual(t, 2, len(lines))
		core.AssertEqual(t, "short", lines[0])
		core.AssertEqual(t, "another", lines[1])
	})

	t.Run("word wrap resets column on newline", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithWordWrap(5), WithStreamOutput(&buf))

		s.Write("12345\n67890ABCDE")
		s.Done()
		s.Wait()

		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		core.AssertEqual(t, 3, len(lines))
		core.AssertEqual(t, "12345", lines[0])
		core.AssertEqual(t, "67890", lines[1])
		core.AssertEqual(t, "ABCDE", lines[2])
	})

	t.Run("no wrap when disabled", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		s.Write(strings.Repeat("x", 200))
		s.Done()
		s.Wait()

		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		core.AssertEqual(t, 1, len(lines))
		core.AssertEqual(t, 200, len(lines[0]))
	})

	t.Run("column tracking", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		s.Write("hello")
		core.AssertEqual(t, 5, s.Column())

		s.Write(" world")
		core.AssertEqual(t, 11, s.Column())
	})

	t.Run("WriteFrom io.Reader", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		r := strings.NewReader("streamed content")
		err := s.WriteFrom(r)
		core.AssertNoError(t, err)

		s.Done()
		s.Wait()
		core.AssertEqual(t, "streamed content\n", buf.String())
	})

	t.Run("channel pattern", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		tokens := make(chan string, 3)
		tokens <- "one "
		tokens <- "two "
		tokens <- "three"
		close(tokens)

		go func() {
			for tok := range tokens {
				s.Write(tok)
			}
			s.Done()
		}()

		s.Wait()
		core.AssertEqual(t, "one two three\n", buf.String())
	})

	t.Run("Done adds trailing newline only if needed", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		s.Write("text\n") // ends with newline, col=0
		s.Done()
		s.Wait()
		core.AssertEqual(t, "text\n", buf.String()) // no double newline
	})
}

func TestStream_Bad(t *core.T) {
	t.Run("empty stream", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		s.Done()
		s.Wait()
		core.AssertEqual(t, "", buf.String())
	})
}

func TestStream_Ugly(t *core.T) {
	t.Run("Write after Done does not panic", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		s.Done()
		s.Wait()
		core.AssertNotPanics(t, func() {
			s.Write("late write")
		})
	})

	t.Run("word wrap width of 1 does not panic", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithWordWrap(1), WithStreamOutput(&buf))
		core.AssertNotPanics(t, func() {
			s.Write("hello")
			s.Done()
			s.Wait()
		})
	})

	t.Run("very large write does not panic", func(t *core.T) {
		var buf bytes.Buffer
		s := NewStream(WithStreamOutput(&buf))

		large := strings.Repeat("x", 100_000)
		core.AssertNotPanics(t, func() {
			s.Write(large)
			s.Done()
			s.Wait()
		})
		core.AssertEqual(t, 100_000, len(strings.TrimRight(buf.String(), "\n")))
	})
}
