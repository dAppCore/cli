# Core-IDE Job Runner — Completion Summary

**Completed:** 9 February 2026
**Module:** `forge.lthn.ai/core/cli` (extracted to `core/ide` repo during monorepo split)
**Status:** Complete — all components built, tested, and operational before extraction

## What Was Built

Autonomous job runner for core-ide that polls Forgejo for actionable pipeline
work, executes it via typed handler functions, captures JSONL training data,
and supports both headless (server) and desktop (Wails GUI) modes.

### Key components

- **`pkg/jobrunner/types.go`** — JobSource, JobHandler, PipelineSignal,
  ActionResult interfaces and structs
- **`pkg/jobrunner/poller.go`** — multi-source poller with configurable
  interval, ETag-based conditional requests, and idle backoff
- **`pkg/jobrunner/journal.go`** — append-only JSONL writer for training data
  capture (structural signals only, no content)
- **`pkg/jobrunner/forgejo/source.go`** — ForgejoSource adapter (evolved from
  original GitHubSource design to use pkg/forge SDK)
- **`pkg/jobrunner/forgejo/signals.go`** — PR/issue state extraction and
  signal building from Forgejo API responses

### Handlers

All six handlers from the design were implemented with tests:

- `publish_draft` — mark draft PRs as ready when checks pass
- `send_fix_command` — comment fix instructions for conflicts/reviews
- `resolve_threads` — resolve pre-commit review threads after fix
- `enable_auto_merge` — enable auto-merge when all checks pass
- `tick_parent` — update epic issue checklist when child PR merges
- `dispatch` — SCP ticket delivery to agent machines via SSH (added beyond
  original design)

### Headless / Desktop mode

- `hasDisplay()` detection for Linux/macOS/Windows
- `--headless` / `--desktop` CLI flag overrides
- Headless: poller + MCP bridge, signal handling, systemd-ready
- Desktop: Wails GUI with system tray, optional poller toggle

### Extraction

Code was fully operational and then extracted during the Feb 2026 monorepo
split (`abe74a1`). `pkg/jobrunner/` moved to `core/go`, `cmd/core-ide/` and
`internal/core-ide/` moved to `core/ide`. The agentci dispatch system
(`d9f3b72` through `886c67e`) built on top of the jobrunner before extraction.
