# Agent Notes

This repository is the `dappco.re/go/cli` module. It provides the shared command-line runtime, terminal UI helpers, daemon helpers, package commands, and i18n utilities used by Core command tools.

Work in this repo should follow the v0.9 Core conventions:

- Import `dappco.re/go` as `core` when a file needs Core wrappers.
- Keep public-symbol tests in the sibling `<source>_test.go` file using `Test<File>_<Symbol>_{Good,Bad,Ugly}`.
- Keep public-symbol examples in the sibling `<source>_example_test.go` file.
- Prefer `core.Result` for fallible production functions so callers branch on `r.OK`.
- Do not put compliance tests in aggregate AX-7 files or versioned test files.

The CLI package mutates process-level state for stdio, theme, colors, and global runtime registration. Tests that touch those globals should restore them with `t.Cleanup` or the local reset helpers before returning.

Before handing work back, run:

```sh
GOWORK=off go mod tidy
GOWORK=off go vet ./...
GOWORK=off go test -count=1 ./...
gofmt -l .
bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh .
```
