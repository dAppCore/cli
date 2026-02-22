// Package module provides CLI commands for managing marketplace modules.
//
// Commands:
//   - install: Install a module from a Git repo
//   - list: List installed modules
//   - update: Update a module or all modules
//   - remove: Remove an installed module
package module

import (
	"os"
	"path/filepath"

	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go/pkg/i18n"
	"forge.lthn.ai/core/go/pkg/marketplace"
	"forge.lthn.ai/core/go/pkg/store"
)

// AddModuleCommands registers the 'module' command and all subcommands.
func AddModuleCommands(root *cli.Command) {
	moduleCmd := &cli.Command{
		Use:   "module",
		Short: i18n.T("Manage marketplace modules"),
	}
	root.AddCommand(moduleCmd)

	addInstallCommand(moduleCmd)
	addListCommand(moduleCmd)
	addUpdateCommand(moduleCmd)
	addRemoveCommand(moduleCmd)
}

// moduleSetup returns the modules directory, store, and installer.
// The caller must defer st.Close().
func moduleSetup() (string, *store.Store, *marketplace.Installer, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil, nil, cli.Wrap(err, "failed to determine home directory")
	}

	modulesDir := filepath.Join(home, ".core", "modules")
	if err := os.MkdirAll(modulesDir, 0755); err != nil {
		return "", nil, nil, cli.Wrap(err, "failed to create modules directory")
	}

	dbPath := filepath.Join(modulesDir, "modules.db")
	st, err := store.New(dbPath)
	if err != nil {
		return "", nil, nil, cli.Wrap(err, "failed to open module store")
	}

	inst := marketplace.NewInstaller(modulesDir, st)
	return modulesDir, st, inst, nil
}
