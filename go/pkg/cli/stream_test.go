package cli

import (
	core "dappco.re/go"
)

func TestStream_WithWordWrap_Good(t *core.T) {
	stream := NewStream(WithWordWrap(4), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 4, stream.wrap)
	core.AssertNotNil(t, stream.out)
}

func TestStream_WithWordWrap_Bad(t *core.T) {
	stream := NewStream(WithWordWrap(0), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 0, stream.wrap)
	core.AssertNotNil(t, stream.out)
}

func TestStream_WithWordWrap_Ugly(t *core.T) {
	stream := NewStream(WithWordWrap(-1), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, -1, stream.wrap)
	core.AssertNotNil(t, stream.out)
}

func TestStream_WithStreamOutput_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))

	core.AssertEqual(t, out, stream.out)
	core.AssertNotNil(t, stream.done)
}

func TestStream_WithStreamOutput_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(nil))

	core.AssertNil(t, stream.out)
	core.AssertNotNil(t, stream.done)
}

func TestStream_WithStreamOutput_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.Discard))

	core.AssertEqual(t, core.Discard, stream.out)
	core.AssertEqual(t, 0, stream.Column())
}

func TestStream_NewStream_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))

	core.AssertNotNil(t, stream)
	core.AssertEqual(t, out, stream.out)
}

func TestStream_NewStream_Bad(t *core.T) {
	stream := NewStream()

	core.AssertNotNil(t, stream)
	core.AssertNotNil(t, stream.out)
}

func TestStream_NewStream_Ugly(t *core.T) {
	stream := NewStream(WithWordWrap(3), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 3, stream.wrap)
	core.AssertEqual(t, 0, stream.Column())
}

func TestStream_Stream_Write_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("hello")

	core.AssertEqual(t, "hello", out.String())
	core.AssertEqual(t, 5, stream.Column())
}

func TestStream_Stream_Write_Bad(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("")

	core.AssertEqual(t, "", out.String())
	core.AssertEqual(t, 0, stream.Column())
}

func TestStream_Stream_Write_Ugly(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out), WithWordWrap(3))
	stream.Write("abcd")

	core.AssertContains(t, out.String(), "\n")
	core.AssertEqual(t, 1, stream.Column())
}

func TestStream_Stream_WriteFrom_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	err := cliResultError(stream.WriteFrom(core.NewReader("hello")))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "hello", out.String())
}

func TestStream_Stream_WriteFrom_Bad(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	err := cliResultError(stream.WriteFrom(core.NewReader("")))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "", out.String())
}

func TestStream_Stream_WriteFrom_Ugly(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out), WithWordWrap(3))
	err := cliResultError(stream.WriteFrom(core.NewReader("abcd")))

	core.AssertNoError(t, err)
	core.AssertContains(t, out.String(), "\n")
}

func TestStream_Stream_Done_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("hello")
	stream.Done()

	core.AssertContains(t, out.String(), "\n")
	core.AssertNotPanics(t, stream.Done)
}

func TestStream_Stream_Done_Bad(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Done()

	core.AssertEqual(t, "", out.String())
	core.AssertNotPanics(t, stream.Done)
}

func TestStream_Stream_Done_Ugly(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("hello\n")
	stream.Done()

	core.AssertEqual(t, "hello\n", out.String())
	core.AssertNotPanics(t, stream.Done)
}

func TestStream_Stream_Wait_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Done()

	core.AssertNotPanics(t, stream.Wait)
	core.AssertEqual(t, 0, stream.Column())
}

func TestStream_Stream_Wait_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	go stream.Done()

	core.AssertNotPanics(t, stream.Wait)
	core.AssertEqual(t, 0, stream.Column())
}

func TestStream_Stream_Wait_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("x")
	stream.Done()

	core.AssertNotPanics(t, stream.Wait)
	core.AssertEqual(t, 1, stream.Column())
}

func TestStream_Stream_Column_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("abc")

	core.AssertEqual(t, 3, stream.Column())
	core.AssertNotEqual(t, 0, stream.Column())
}

func TestStream_Stream_Column_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 0, stream.Column())
	core.AssertNotEqual(t, 1, stream.Column())
}

func TestStream_Stream_Column_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("a\nbc")

	core.AssertEqual(t, 2, stream.Column())
	core.AssertNotEqual(t, 4, stream.Column())
}

func TestStream_Stream_Captured_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("hello")

	core.AssertEqual(t, "hello", stream.Captured())
	core.AssertNotEmpty(t, stream.Captured())
}

func TestStream_Stream_Captured_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(core.Discard))
	stream.Write("hello")

	core.AssertEqual(t, "", stream.Captured())
	core.AssertEmpty(t, stream.Captured())
}

func TestStream_Stream_Captured_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("")

	core.AssertEqual(t, "", stream.Captured())
	core.AssertEmpty(t, stream.Captured())
}

func TestStream_Stream_CapturedOK_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("hello")
	got, ok := stream.CapturedOK()

	core.AssertTrue(t, ok)
	core.AssertEqual(t, "hello", got)
}

func TestStream_Stream_CapturedOK_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(core.Discard))
	got, ok := stream.CapturedOK()

	core.AssertFalse(t, ok)
	core.AssertEqual(t, "", got)
}

func TestStream_Stream_CapturedOK_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	got, ok := stream.CapturedOK()

	core.AssertTrue(t, ok)
	core.AssertEqual(t, "", got)
}
