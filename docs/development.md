# Development

Use `GOWORK=off` when validating this module so the repository is tested against its declared module graph. In sandboxed runs, a temporary Go cache may be needed, but the release check should still use the standard commands from `AGENTS.md`.

Public APIs require three local tests in the source sibling test file:

- `Good` covers the normal path.
- `Bad` covers an invalid input or failure path.
- `Ugly` covers an edge case, repeated call, nil-ish state, or boundary value.

Examples live beside the source file as `<source>_example_test.go`. They should execute the symbol they document and print deterministic stdout with a matching `// Output:` block.

When adding CLI behavior, avoid leaking process state between tests. Use in-memory readers and writers through `SetStdin`, `SetStdout`, and `SetStderr`, and restore them with cleanup. For command execution, prefer small fake executables in a temporary directory over using the developer machine's installed tools.

Run the audit script as the final compliance gate. A passing unit test suite is necessary, but the repository is not compliant until the audit reports `verdict: COMPLIANT` with every counter at zero.
