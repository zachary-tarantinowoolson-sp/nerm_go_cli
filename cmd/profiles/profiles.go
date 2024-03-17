/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"github.com/spf13/cobra"
)

func NewProfilesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profiles",
		Short:   "CRUD Profile data from an eviroment",
		Long:    "Create, read, update, and delete Profiles. This allows an admin to execute commands against a specified NERM tenant",
		Example: "nerm profiles count",
		Aliases: []string{"p"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newProfileCountCommand(),
	)

	return cmd
}
