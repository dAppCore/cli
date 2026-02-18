package module

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/cli"
	"forge.lthn.ai/core/go/pkg/i18n"
	"forge.lthn.ai/core/go/pkg/marketplace"
)

var (
	installRepo    string
	installSignKey string
)

func addInstallCommand(parent *cli.Command) {
	installCmd := cli.NewCommand(
		"install <code>",
		i18n.T("Install a module from a Git repo"),
		i18n.T("Install a module by cloning its Git repository, verifying the manifest signature, and registering it.\n\nThe --repo flag is required and specifies the Git URL to clone from."),
		func(cmd *cli.Command, args []string) error {
			if installRepo == "" {
				return fmt.Errorf("--repo flag is required")
			}
			return runInstall(args[0], installRepo, installSignKey)
		},
	)
	installCmd.Args = cli.ExactArgs(1)
	installCmd.Example = "  core module install my-module --repo https://forge.lthn.ai/modules/my-module.git\n  core module install signed-mod --repo ssh://git@forge.lthn.ai:2223/modules/signed.git --sign-key abc123"

	cli.StringFlag(installCmd, &installRepo, "repo", "r", "", i18n.T("Git repository URL to clone"))
	cli.StringFlag(installCmd, &installSignKey, "sign-key", "k", "", i18n.T("Hex-encoded ed25519 public key for manifest verification"))

	parent.AddCommand(installCmd)
}

func runInstall(code, repo, signKey string) error {
	_, st, inst, err := moduleSetup()
	if err != nil {
		return err
	}
	defer st.Close()

	cli.Dim("Installing module " + code + " from " + repo + "...")

	mod := marketplace.Module{
		Code:    code,
		Repo:    repo,
		SignKey: signKey,
	}

	if err := inst.Install(context.Background(), mod); err != nil {
		return err
	}

	cli.Success("Module " + code + " installed successfully")
	return nil
}
