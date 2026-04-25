package cli

import (
	"dappco.re/go/core"
	"github.com/mattn/go-runewidth"
)

// StreamOption configures a Stream.
//
//	stream := cli.NewStream(cli.WithWordWrap(80))
//	stream.Wait()
type StreamOption func(*Stream)

// WithWordWrap sets the word-wrap column width.
func WithWordWrap(cols int) StreamOption {
	return func(s *Stream) { s.wrap = cols }
}

// WithStreamOutput sets the output writer (default: stdoutWriter()).
func WithStreamOutput(w Writer) StreamOption {
	return func(s *Stream) { s.out = w }
}

// Stream renders growing text as tokens arrive, with optional word-wrap.
// Safe for concurrent writes from a single producer goroutine.
//
//	stream := cli.NewStream(cli.WithWordWrap(80))
//	go func() {
//	    for token := range tokens {
//	        stream.Write(token)
//	    }
//	    stream.Done()
//	}()
//	stream.Wait()
type Stream struct {
	out  Writer
	wrap int
	col  int // current column position (visible characters)
	done chan struct{}
	once core.Once
	mu   core.Mutex
}

// NewStream creates a streaming text renderer.
func NewStream(opts ...StreamOption) *Stream {
	s := &Stream{
		out:  stdoutWriter(),
		done: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Write appends text to the stream. Handles word-wrap if configured.
func (s *Stream) Write(text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.wrap <= 0 {
		writeString(s.out, text)
		// Track visible width across newlines for Done() trailing-newline logic.
		if idx := LastIndex(text, "\n"); idx >= 0 {
			s.col = runewidth.StringWidth(text[idx+1:])
		} else {
			s.col += runewidth.StringWidth(text)
		}
		return
	}

	for _, r := range text {
		if r == '\n' {
			core.Print(s.out, "")
			s.col = 0
			continue
		}

		rw := runewidth.RuneWidth(r)
		if rw > 0 && s.col > 0 && s.col+rw > s.wrap {
			core.Print(s.out, "")
			s.col = 0
		}

		writeString(s.out, string(r))
		s.col += rw
	}
}

// WriteFrom reads from r and streams all content until EOF.
func (s *Stream) WriteFrom(r Reader) error {
	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			s.Write(string(buf[:n]))
		}
		if isEOF(err) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// Done signals that no more text will arrive.
func (s *Stream) Done() {
	s.once.Do(func() {
		s.mu.Lock()
		if s.col > 0 {
			core.Print(s.out, "") // ensure trailing newline
		}
		s.mu.Unlock()
		close(s.done)
	})
}

// Wait blocks until Done is called.
func (s *Stream) Wait() {
	<-s.done
}

// Content returns the current column position (for testing).
func (s *Stream) Column() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.col
}

// Captured returns the stream output as a string when the output writer is
// capture-capable. If the writer cannot be captured, it returns an empty string.
// Use CapturedOK when you need to distinguish that case.
func (s *Stream) Captured() string {
	out, _ := s.CapturedOK()
	return out
}

// CapturedOK returns the stream output and whether the configured writer
// supports capture. Any writer exposing a String() method is detected,
// which includes *strings.Builder, *bytes.Buffer, and test fixtures.
func (s *Stream) CapturedOK() (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if st, ok := s.out.(interface{ String() string }); ok {
		return st.String(), true
	}
	return "", false
}
