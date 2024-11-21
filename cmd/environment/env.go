/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func NewEnvironmentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environments",
		Short:   "Manage Environement data for the CLI",
		Long:    "Create, read, updated, and delete environments to use with the CLI. This allows an admin to execute commands against a specified NERM tenant",
		Example: "nerm environment create",
		Aliases: []string{"env"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newCreateCommand(),
		newListCommand(),
		newShowCommand(),
		newDeleteCommand(),
		newUpdateCommand(),
		newUseCommand(),
	)

	return cmd
}

func printEnvListTable(data []string) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()

	tbl := table.New("Environments", "","")
	tbl.WithHeaderFormatter(headerFmt)

	if len(data)>3{
		for i := 0; i < len(data); i += 3 {
			if (len(data) - i) == 2 {
				tbl.AddRow(data[len(data)-2], data[len(data)-1])
			} else if (len(data) - i) == 1 {
				tbl.AddRow(data[len(data)-1])
			} else {
				tbl.AddRow(data[i],data[i+1],data[i+2])
			}
		}
	} else {
		for k := range data {
			tbl.AddRow(data[k])
		}
	}

	tbl.Print()
}