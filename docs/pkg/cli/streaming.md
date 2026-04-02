---
title: Streaming Output
description: Real-time token-by-token text rendering with optional word-wrap.
---

# Streaming Output

The `Stream` type renders growing text as tokens arrive, with optional word-wrap. It is designed for displaying LLM output, log tails, or any content that arrives incrementally. Thread-safe for a single producer goroutine.

## Basic Usage

```go
stream := cli.NewStream()
go func() {
    for token := range tokens {
        stream.Write(token)
    }
    stream.Done()
}()
stream.Wait()
```

`Done()` ensures a trailing newline if the stream did not end with one. `Wait()` blocks until `Done()` is called.

## Word Wrap

Wrap text at a column boundary:

```go
stream := cli.NewStream(cli.WithWordWrap(80))
```

When word-wrap is enabled, the stream tracks the current column position and inserts line breaks when the column limit is reached.

## Custom Output Writer

By default, streams write to the CLI stdout writer (`stdoutWriter()`), so tests can
redirect output via `cli.SetStdout` and other callers can provide any `io.Writer`:

```go
var buf strings.Builder
stream := cli.NewStream(cli.WithStreamOutput(&buf))
// ... write tokens ...
stream.Done()
result, ok := stream.CapturedOK() // or buf.String()
```

`Captured()` returns the output as a string when using a `*strings.Builder` or any `fmt.Stringer`.
`CapturedOK()` reports whether capture is supported by the configured writer.

## Reading from `io.Reader`

Stream content from an HTTP response body, file, or any `io.Reader`:

```go
stream := cli.NewStream(cli.WithWordWrap(120))
err := stream.WriteFrom(resp.Body)
stream.Done()
```

`WriteFrom` reads in 256-byte chunks until EOF, calling `Write` for each chunk.

## API Reference

| Method | Description |
|--------|-------------|
| `NewStream(opts...)` | Create a stream with options |
| `Write(text)` | Append text (thread-safe) |
| `WriteFrom(r)` | Stream from `io.Reader` until EOF |
| `Done()` | Signal completion (adds trailing newline if needed) |
| `Wait()` | Block until `Done` is called |
| `Column()` | Current column position |
| `Captured()` | Get output as string (returns `""` if capture is unsupported) |
| `CapturedOK()` | Get output and support status |

## Options

| Option | Description |
|--------|-------------|
| `WithWordWrap(cols)` | Set the word-wrap column width |
| `WithStreamOutput(w)` | Set the output writer (default: `stdoutWriter()`) |

## Example: LLM Token Streaming

```go
func streamResponse(ctx context.Context, model *Model, prompt string) error {
    stream := cli.NewStream(cli.WithWordWrap(100))

    go func() {
        ch := model.Generate(ctx, prompt)
        for token := range ch {
            stream.Write(token)
        }
        stream.Done()
    }()

    stream.Wait()
    return nil
}
```
