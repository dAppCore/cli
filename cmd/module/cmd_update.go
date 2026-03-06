package module

import (
	"context"
	"errors"
	"fmt"

	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go-i18n"
)

var updateAll bool

func addUpdateCommand(parent *cli.Command) {
	updateCmd := cli.NewCommand(
		"update [code]",
		i18n.T("Update a module or all modules"),
		i18n.T("Update a specific module to the latest version, or use --all to update all installed modules."),
		func(cmd *cli.Command, args []string) error {
			if updateAll {
				return runUpdateAll()
			}
			if len(args) == 0 {
				return errors.New("module code required (or use --all)")
			}
			return runUpdate(args[0])
		},
	)

	cli.BoolFlag(updateCmd, &updateAll, "all", "a", false, i18n.T("Update all installed modules"))

	parent.AddCommand(updateCmd)
}

func runUpdate(code string) error {
	_, st, inst, err := moduleSetup()
	if err != nil {
		return err
	}
	defer st.Close()

	cli.Dim("Updating " + code + "...")

	if err := inst.Update(context.Background(), code); err != nil {
		return err
	}

	cli.Success("Module " + code + " updated successfully")
	return nil
}

func runUpdateAll() error {
	_, st, inst, err := moduleSetup()
	if err != nil {
		return err
	}
	defer st.Close()

	installed, err := inst.Installed()
	if err != nil {
		return err
	}

	if len(installed) == 0 {
		cli.Dim("No modules installed")
		return nil
	}

	ctx := context.Background()
	var updated, failed int
	for _, m := range installed {
		cli.Dim("Updating " + m.Code + "...")
		if err := inst.Update(ctx, m.Code); err != nil {
			cli.Errorf("Failed to update %s: %v", m.Code, err)
			failed++
			continue
		}
		cli.Success(m.Code + " updated")
		updated++
	}

	fmt.Println()
	cli.Dim(fmt.Sprintf("%d updated, %d failed", updated, failed))
	return nil
}
