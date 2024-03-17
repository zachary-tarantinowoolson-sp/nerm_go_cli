/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
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
