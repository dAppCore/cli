package module

import (
	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go/pkg/i18n"
)

func addRemoveCommand(parent *cli.Command) {
	removeCmd := cli.NewCommand(
		"remove <code>",
		i18n.T("Remove an installed module"),
		"",
		func(cmd *cli.Command, args []string) error {
			return runRemove(args[0])
		},
	)
	removeCmd.Args = cli.ExactArgs(1)

	parent.AddCommand(removeCmd)
}

func runRemove(code string) error {
	_, st, inst, err := moduleSetup()
	if err != nil {
		return err
	}
	defer st.Close()

	if !cli.Confirm("Remove module " + code + "?") {
		cli.Dim("Cancelled")
		return nil
	}

	if err := inst.Remove(code); err != nil {
		return err
	}

	cli.Success("Module " + code + " removed")
	return nil
}
