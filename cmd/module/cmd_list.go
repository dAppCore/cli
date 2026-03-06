package module

import (
	"fmt"

	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go-i18n"
)

func addListCommand(parent *cli.Command) {
	listCmd := cli.NewCommand(
		"list",
		i18n.T("List installed modules"),
		"",
		func(cmd *cli.Command, args []string) error {
			return runList()
		},
	)

	parent.AddCommand(listCmd)
}

func runList() error {
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

	table := cli.NewTable("Code", "Name", "Version", "Repo")
	for _, m := range installed {
		table.AddRow(m.Code, m.Name, m.Version, m.Repo)
	}

	fmt.Println()
	table.Render()
	fmt.Println()
	cli.Dim(fmt.Sprintf("%d module(s) installed", len(installed)))

	return nil
}
