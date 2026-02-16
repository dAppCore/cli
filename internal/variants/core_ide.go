//go:build ide

// core_ide.go imports packages for the Core IDE desktop application.
//
// Build with: go build -tags ide
//
// This is the Wails v3 GUI variant featuring:
//   - System tray with quick actions
//   - Tray panel for status/notifications
//   - Angular frontend
//   - All CLI commands available via IPC

package variants

import (
	// CLI commands available via IPC (IDE GUI is now in core/ide repo)
	_ "forge.lthn.ai/core/cli/internal/cmd/ai"
	_ "forge.lthn.ai/core/cli/internal/cmd/deploy"
	_ "forge.lthn.ai/core/cli/internal/cmd/dev"
	_ "forge.lthn.ai/core/cli/internal/cmd/php"
	_ "forge.lthn.ai/core/cli/internal/cmd/rag"
)
