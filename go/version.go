// Package cli is the module root for dappco.re/go/cli. It holds the version
// the binary was built at, and nothing else: everything a caller reaches for
// lives under pkg/.
//
// The version sits here rather than in pkg/cli or in main so that one symbol
// name serves every binary built from this module — the release stamps
// dappco.re/go/cli.Version and each consumer reads it from the same place.
package cli

// Version is the version of the built binary, set at link time.
//
//	go build -ldflags="-X dappco.re/go/cli.Version=1.2.0"
//
// -X names a package-qualified symbol and silently discards one it cannot
// find, so a stamp aimed at a variable that does not exist leaves this at its
// default and the binary reports 0.0.0 with a green build.
var Version = "0.0.0"
