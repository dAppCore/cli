# Build Variants

The Core CLI supports a modular build system that allows you to create different binary variants with different feature sets. This is useful for creating minimal builds, specialized builds, or builds tailored to specific use cases.

## Overview

The variant system works through Go's module system and selective import statements in `main.go`. By commenting out or including specific import statements, you can control which commands are included in the final binary.

## Available Variants

### Full Build (Default)

The default build includes all commands and features. This is the standard build used for development and general-purpose use.

**Size:** ~50-100 MB (depending on platform)

**Included:** All commands documented in [CLI Reference](cli-reference.md)

### Minimal Build

A minimal build includes only essential commands, suitable for resource-constrained environments or when you only need basic functionality.

**Recommended commands for minimal build:**
- `config` - Configuration management
- `help` - Help system
- `doctor` - Environment checks
- `test` - Basic testing

### Specialized Builds

You can create specialized builds for specific use cases:

#### Developer Build
Focus on development workflow commands:
- `dev` - Development workflow
- `git` - Git operations
- `go` - Go development
- `test` - Testing
- `qa` - Quality assurance
- `docs` - Documentation

#### DevOps Build
Focus on infrastructure and deployment:
- `deploy` - Deployment management
- `prod` - Production infrastructure
- `vm` - Virtual machine management
- `unifi` - Network management
- `monitor` - Security monitoring

#### AI/ML Build
Focus on AI and machine learning operations:
- `ai` - AI task management
- `ml` - ML pipeline
- `rag` - RAG system
- `mcp` - MCP server
- `collect` - Data collection

## How to Create a Variant

### Method 1: Comment Out Imports in main.go

Edit `main.go` and comment out the commands you don't want:

```go
package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	// Essential commands
	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/doctor"

	// Development commands
	_ "forge.lthn.ai/core/cli/cmd/dev"
	_ "forge.lthn.ai/core/cli/cmd/go"
	_ "forge.lthn.ai/core/cli/cmd/test"

	// Comment out to exclude:
	// _ "forge.lthn.ai/core/cli/cmd/ai"
	// _ "forge.lthn.ai/core/cli/cmd/ml"
	// _ "forge.lthn.ai/core/cli/cmd/deploy"
	// etc.
)

func main() {
	cli.Main()
}
```

Then build:

```bash
task cli:build
```

### Method 2: Create Multiple main.go Files

Create separate entry points for different variants:

```bash
# Directory structure
cmd/
  core/           # Full build
    main.go
  core-minimal/   # Minimal build
    main.go
  core-dev/       # Developer build
    main.go
  core-devops/    # DevOps build
    main.go
```

Each `main.go` includes only the relevant imports for that variant.

Build specific variant:

```bash
cd cmd/core-minimal
go build -o ../../bin/core-minimal .
```

### Method 3: Build Tags

Use Go build tags to conditionally include commands:

```go
//go:build !minimal
// +build !minimal

package ai

// Command registration...
```

Build with tags:

```bash
# Full build
go build -o bin/core .

# Minimal build
go build -tags minimal -o bin/core-minimal .
```

## External Variant Repositories

The Core CLI architecture supports external command repositories that can be optionally included. These are commented out in `main.go` by default:

```go
// Variant repos (optional — comment out to exclude)
// _ "forge.lthn.ai/core/php"
// _ "forge.lthn.ai/core/ci"
```

### Currently Available External Repos

#### core/php (Archived)
PHP and Laravel development commands. **Note:** This has been moved to its own repository and is no longer included by default.

#### core/ci (Archived)
CI/CD pipeline commands. **Note:** This has been moved to its own repository and is no longer included by default.

To include these in your build:

1. Add the module to your `go.mod`:
   ```bash
   go get forge.lthn.ai/core/php@latest
   ```

2. Uncomment the import in `main.go`:
   ```go
   _ "forge.lthn.ai/core/php"
   ```

3. Rebuild:
   ```bash
   task cli:build
   ```

## Build Optimization

### Size Optimization

For smaller binaries, use the release build with stripped symbols:

```bash
task cli:build:release
```

This uses the following ldflags:
```
-s -w  # Strip debug info and symbol table
```

### Compression

Further reduce binary size with UPX:

```bash
# Install UPX
# macOS: brew install upx
# Linux: apt-get install upx-ucl

# Compress binary
upx --best --lzma bin/core
```

**Note:** Compressed binaries may trigger antivirus false positives and won't work with code signing on macOS.

## Version Information

All builds include embedded version information via ldflags:

```bash
# View version
core --version

# Build with custom version
go build -ldflags "-X forge.lthn.ai/core/go/pkg/cli.AppVersion=1.2.3" .
```

The Taskfile automatically sets version info from git tags:

```yaml
LDFLAGS_BASE: >-
  -X {{.PKG}}.AppVersion={{.SEMVER_VERSION}}
  -X {{.PKG}}.BuildCommit={{.SEMVER_COMMIT}}
  -X {{.PKG}}.BuildDate={{.SEMVER_DATE}}
  -X {{.PKG}}.BuildPreRelease={{.SEMVER_PRERELEASE}}
```

## Examples

### Create a Minimal Developer Build

1. Edit `main.go`:
```go
package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/dev"
	_ "forge.lthn.ai/core/cli/cmd/doctor"
	_ "forge.lthn.ai/core/cli/cmd/git"
	_ "forge.lthn.ai/core/cli/cmd/go"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/test"
)

func main() {
	cli.Main()
}
```

2. Build:
```bash
task cli:build:release
```

3. Result: ~20-30 MB binary with only development commands

### Create an AI/ML Specialist Build

1. Edit `main.go`:
```go
package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	_ "forge.lthn.ai/core/cli/cmd/ai"
	_ "forge.lthn.ai/core/cli/cmd/collect"
	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/daemon"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/mcp"
	_ "forge.lthn.ai/core/cli/cmd/ml"
	_ "forge.lthn.ai/core/cli/cmd/rag"
	_ "forge.lthn.ai/core/cli/cmd/session"
)

func main() {
	cli.Main()
}
```

2. Build:
```bash
task cli:build:release
```

3. Result: Binary optimized for AI/ML workflows

## Testing Variants

Test your variant to ensure all needed commands are present:

```bash
# Build variant
task cli:build

# Test available commands
./bin/core help

# Verify specific command
./bin/core dev --help
```

## Distribution

When distributing variant builds, use clear naming:

```bash
# Good naming examples
core-full-v1.2.3-linux-amd64
core-minimal-v1.2.3-darwin-arm64
core-dev-v1.2.3-windows-amd64.exe
```

Include a README describing what's included in each variant.

## Best Practices

1. **Document your variant** - Keep a record of which commands are included
2. **Test thoroughly** - Ensure all needed functionality is present
3. **Version consistently** - Use the same version numbering scheme
4. **Consider dependencies** - Some commands depend on others (e.g., `dev` uses `git`)
5. **Keep it simple** - Don't create too many variants unless necessary
6. **Use release builds** - Always use `-s -w` ldflags for distribution

## Troubleshooting

### Command Not Found

If a command is missing after building a variant:

1. Check that the import is uncommented in `main.go`
2. Verify the module is in `go.mod`
3. Run `go mod tidy`
4. Rebuild with `task cli:build`

### Build Errors

If you get import errors:

1. Ensure all required modules are in `go.mod`
2. Run `go mod download`
3. Check for version conflicts with `go mod why`

### Large Binary Size

If your minimal build is still too large:

1. Use release build flags: `task cli:build:release`
2. Remove unused commands from imports
3. Consider using UPX compression (with caveats noted above)
4. Check for large embedded assets

## See Also

- [CLI Reference](cli-reference.md) - Complete command documentation
- [Configuration](configuration.md) - Configuration system
- [Getting Started](getting-started.md) - Installation and setup
