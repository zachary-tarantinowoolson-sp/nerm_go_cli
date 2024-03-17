/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"log"
	"nerm/cmd/utilities"

	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Create a new environment",
		Long:    "Create a new environment for use in the CLI",
		Example: "nerm env create | nerm env create env_name",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, environmentName := range args {
				utilities.CreateEnvironment(environmentName)
			}

			if len(args) == 0 {
				log.Print("Please specify a tenant_name to create")
			}

			return nil
		},
	}
}
