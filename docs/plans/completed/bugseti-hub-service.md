# BugSETI HubService — Completion Summary

**Completed:** 13 February 2026
**Module:** `forge.lthn.ai/core/cli` (extracted to `core/bugseti` repo on 16 Feb 2026)
**Status:** Complete — all Go-side tasks implemented and wired into app lifecycle

## What Was Built

Thin HTTP client service coordinating with the agentic portal's
`/api/bugseti/*` endpoints for issue claiming, stats sync, leaderboard,
and offline-first pending operations queue.

### Implementation (Tasks 1-8 from plan)

All 8 Go-side tasks were implemented across commits `a38ce05` through `177ce27`:

1. **Config fields** — HubURL, HubToken, ClientID, ClientName added to
   ConfigService with getters/setters (`a38ce05`)
2. **HubService types + constructor** — HubService, PendingOp, HubClaim,
   LeaderboardEntry, GlobalStats, ConflictError, NotFoundError (`a89acfa`)
3. **HTTP request helpers** — `doRequest()`, `doJSON()` with bearer auth,
   error classification (401/404/409), and connection tracking (`ab7ef52`)
4. **AutoRegister** — exchange forge token for ak_ hub token via
   `/auth/forge` endpoint (`21d5f5f`)
5. **Write operations** — Register, Heartbeat, ClaimIssue, UpdateStatus,
   ReleaseClaim, SyncStats (`a6456e2`)
6. **Read operations** — IsIssueClaimed, ListClaims, GetLeaderboard,
   GetGlobalStats (`7a92fe0`)
7. **Pending ops queue** — offline-first queue with disk persistence to
   `hub_pending.json`, drain-on-reconnect (`a567568`)
8. **main.go integration** — HubService wired as Wails service with
   auto-registration at startup (`177ce27`)

### Tests

All operations tested with `httptest.NewServer` mocks covering success,
network error, 409 conflict, 401 re-auth, and pending ops persist/reload
scenarios. Hub test file: `internal/bugseti/hub_test.go`.

### Key files (before extraction)

- `internal/bugseti/hub.go` — HubService implementation (25 exported methods)
- `internal/bugseti/hub_test.go` — comprehensive httptest-based test suite
- `internal/bugseti/config.go` — hub config fields and accessors
- `cmd/bugseti/main.go` — lifecycle wiring

### Task 9 (Laravel endpoint)

The portal-side `/api/bugseti/auth/forge` endpoint (Task 9) lives in the
`agentic` repo, not in `core/cli`. It was designed in this plan but
implemented separately.

### Extraction

BugSETI was extracted to its own repo on 16 Feb 2026 (`8167f66`):
`internal/bugseti/` moved to `core/bugseti`, `cmd/bugseti/` moved to
`core/bugseti/cmd/`.
