# AX-10 CLI Artifact Tests

The CLI integration suite validates the built binary as an artifact. Tasks build
`tests/cli/bin/cli`, run it with real command-line arguments, and assert on
stdout, stderr, and exit status from bash heredocs.

Taskfile path = command path. The Taskfile structure mirrors the CLI command
tree:

- `tests/cli/Taskfile.yaml` owns artifact build and suite orchestration.
- `tests/cli/config/Taskfile.yaml` covers `cli config ...`.
- `tests/cli/doctor/Taskfile.yaml` covers `cli doctor`.
- `tests/cli/pkg/search/Taskfile.yaml` covers `cli pkg search ...`.
- `tests/cli/version/Taskfile.yaml` covers `cli version`.

Run the full suite from this directory:

```sh
task all
```

Tests should isolate mutable state with temporary homes, registries, caches, and
tool shims. Commands that normally reach outside the sandbox, such as GitHub
lookups, should be driven from seeded fixtures instead of live network access.
