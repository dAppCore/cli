#!/usr/bin/env bash

CORE_CLI_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CORE_CLI_ROOT="$(cd "$CORE_CLI_LIB_DIR/../../.." && pwd)"
CORE_CLI_BIN="${CORE_CLI_BIN:-$CORE_CLI_LIB_DIR/../bin/cli}"
CORE_CLI_GO_CACHE="${CORE_CLI_GO_CACHE:-/tmp/core-cli-ax10-go-build}"

build_core_binary() {
	mkdir -p "$(dirname "$CORE_CLI_BIN")"
	mkdir -p "$CORE_CLI_GO_CACHE"
	(
		cd "$CORE_CLI_ROOT/cmd/core"
		export GOCACHE="${GOCACHE:-$CORE_CLI_GO_CACHE}"
		GOWORK=off go build -trimpath -o "$CORE_CLI_BIN" .
	)
}

cli() {
	if [[ ! -x "$CORE_CLI_BIN" ]]; then
		build_core_binary
	fi

	"$CORE_CLI_BIN" "$@"
}

run() {
	if [[ ! -x "$CORE_CLI_BIN" ]]; then
		build_core_binary
	fi

	RUN_OUTPUT="$(mktemp)"

	set +e
	cli "$@" >"$RUN_OUTPUT" 2>&1
	RUN_EXIT_CODE=$?
	set -e

	export RUN_EXIT_CODE RUN_OUTPUT
}

assert_exit_code() {
	local expected="$1"

	if [[ "$RUN_EXIT_CODE" -ne "$expected" ]]; then
		printf 'expected exit %s, got %s\n' "$expected" "$RUN_EXIT_CODE" >&2
		if [[ -s "$RUN_OUTPUT" ]]; then
			printf 'output:\n' >&2
			cat "$RUN_OUTPUT" >&2
		fi
		return 1
	fi
}

assert_exit_code_nonzero() {
	if [[ "$RUN_EXIT_CODE" -eq 0 ]]; then
		printf 'expected non-zero exit, got 0\n' >&2
		if [[ -s "$RUN_OUTPUT" ]]; then
			printf 'output:\n' >&2
			cat "$RUN_OUTPUT" >&2
		fi
		return 1
	fi
}

assert_output_contains() {
	local needle="$1"

	if ! grep -Fq "$needle" "$RUN_OUTPUT"; then
		printf 'expected output to contain: %s\n' "$needle" >&2
		if [[ -s "$RUN_OUTPUT" ]]; then
			printf 'output:\n' >&2
			cat "$RUN_OUTPUT" >&2
		fi
		return 1
	fi
}

assert_output_matches() {
	local pattern="$1"

	if ! grep -Eq "$pattern" "$RUN_OUTPUT"; then
		printf 'expected output to match: %s\n' "$pattern" >&2
		if [[ -s "$RUN_OUTPUT" ]]; then
			printf 'output:\n' >&2
			cat "$RUN_OUTPUT" >&2
		fi
		return 1
	fi
}
