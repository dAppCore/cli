# Architecture

`dappco.re/go/cli` is a thin CLI layer on top of `dappco.re/go`. The module keeps terminal interaction, command registration, process lifecycle, and package-management commands out of application repositories while still using Core's Result and service patterns.

The main packages are:

- `pkg/cli`: public CLI runtime, terminal output, prompts, layout rendering, glyphs, task tracking, daemon helpers, and command registration.
- `pkg/cli/frame`: frame-oriented layout primitives for interactive terminal screens.
- `pkg/i18n`: translation loading and phrase helpers used by CLI output.
- `cmd/core/*`: command groups for the Core executable, including config, doctor, help, and package workflows.
- `internal/term`: small terminal capability helpers.

Global CLI state is concentrated in `pkg/cli`: runtime initialization, registered commands, stdio overrides, render style, glyph theme, and color settings. This makes command packages simple, but tests must reset state carefully because these values are process-wide.

Fallible production paths return `core.Result`. File IO, JSON parsing, path work, process execution, and formatting should use Core wrapper APIs so the CLI layer has one error-handling shape and one audit surface.
