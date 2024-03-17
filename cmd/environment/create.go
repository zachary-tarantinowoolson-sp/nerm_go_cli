/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Create a new environment",
		Long:    "Create a new environment for use in the CLI",
		Example: "nerm env create | nerm env create tenant_name",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// for _, environmentName := range args {
			// 	util.CreateEnvironment(environmentName, false)
			// }

			// if len(args) == 0 {
			// 	err := util.CreateEnvironment("", false)
			// 	if err != nil {
			// 		return err
			// 	}
			// }

			return nil
		},
	}
}
