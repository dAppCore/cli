package workspace

import "forge.lthn.ai/core/go/pkg/cli"

func init() {
	cli.RegisterCommands(AddWorkspaceCommands)
}
